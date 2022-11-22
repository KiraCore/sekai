package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/KiraCore/sekai/x/custody/types"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
)

type msgServer struct {
	keeper Keeper
	cgk    types.CustomGovKeeper
	bk     types.BankKeeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, cgk types.CustomGovKeeper, bk types.BankKeeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
		cgk:    cgk,
		bk:     bk,
	}
}

var _ types.MsgServer = msgServer{}

func (s msgServer) CreateCustody(goCtx context.Context, msg *types.MsgCreteCustodyRecord) (*types.MsgCreteCustodyRecordResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	record := types.CustodyRecord{
		Address:         msg.Address,
		CustodySettings: &msg.CustodySettings,
	}

	record.CustodySettings.Key = msg.NewKey

	s.keeper.SetCustodyRecord(ctx, record)

	return &types.MsgCreteCustodyRecordResponse{}, nil
}

func (s msgServer) AddToCustodians(goCtx context.Context, msg *types.MsgAddToCustodyCustodians) (*types.MsgAddToCustodyCustodiansResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	record := types.CustodyCustodiansRecord{
		Address:           msg.Address,
		CustodyCustodians: s.keeper.GetCustodyCustodiansByAddress(ctx, msg.Address),
	}

	if record.CustodyCustodians == nil {
		record.CustodyCustodians = new(types.CustodyCustodianList)
		record.CustodyCustodians.Addresses = map[string]bool{}
	}

	for _, address := range msg.AddAddress {
		record.CustodyCustodians.Addresses[address.String()] = true
	}

	keyRecord := types.CustodyKeyRecord{
		Address: msg.Address,
		Key:     msg.NewKey,
	}

	s.keeper.SetCustodyRecordKey(ctx, keyRecord)
	s.keeper.AddToCustodyCustodians(ctx, record)

	return &types.MsgAddToCustodyCustodiansResponse{}, nil
}

func (s msgServer) RemoveFromCustodians(goCtx context.Context, msg *types.MsgRemoveFromCustodyCustodians) (*types.MsgRemoveFromCustodyCustodiansResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	record := types.CustodyCustodiansRecord{
		Address:           msg.Address,
		CustodyCustodians: s.keeper.GetCustodyCustodiansByAddress(ctx, msg.Address),
	}

	if record.CustodyCustodians == nil {
		return nil, errors.Wrap(types.ErrNoWhiteLists, "Can not remove from the empty whitelist")
	}

	if !record.CustodyCustodians.Addresses[msg.RemoveAddress.String()] {
		return nil, errors.Wrap(types.ErrNoWhiteListsElement, "Can not remove missing element from the whitelist")
	}

	record.CustodyCustodians.Addresses[msg.RemoveAddress.String()] = false
	s.keeper.AddToCustodyCustodians(ctx, record)

	return &types.MsgRemoveFromCustodyCustodiansResponse{}, nil
}

func (s msgServer) ApproveTransaction(goCtx context.Context, msg *types.MsgApproveCustodyTransaction) (*types.MsgApproveCustodyTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	hash := strings.ToLower(msg.Hash)
	allowCustodians := true
	allowPassword := true
	settings := s.keeper.GetCustodyInfoByAddress(ctx, msg.TargetAddress)
	custodians := s.keeper.GetCustodyCustodiansByAddress(ctx, msg.TargetAddress)

	//todo: add fee to msg.FromAddress
	//todo: add ValidateBasic and make possibility to vote only 1 time

	record := types.CustodyPool{
		Address:      msg.TargetAddress,
		Transactions: s.keeper.GetCustodyPoolByAddress(ctx, msg.TargetAddress),
	}

	if record.Transactions != nil && record.Transactions.Record[hash] != nil {
		record.Transactions.Record[hash].Votes += 1
	}

	if settings != nil && settings.CustodyEnabled && len(custodians.Addresses) > 0 {
		allowCustodians = (record.Transactions.Record[hash].Votes)*100/uint64(len(custodians.Addresses)) >= settings.CustodyMode
	}

	if settings != nil && settings.UsePassword {
		allowPassword = record.Transactions.Record[hash].Confirmed
	}

	if allowCustodians && allowPassword {
		tx := record.Transactions.Record[hash].Transaction

		to, err := sdk.AccAddressFromBech32(tx.ToAddress)
		if err != nil {
			return nil, err
		}

		if err := s.bk.IsSendEnabledCoins(ctx, tx.Amount...); err != nil {
			return nil, err
		}

		from, err := sdk.AccAddressFromBech32(tx.FromAddress)
		if err != nil {
			return nil, err
		}

		err = s.bk.SendCoins(ctx, from, to, tx.Amount)
		if err != nil {
			return nil, err
		}

		if s.bk.BlockedAddr(to) {
			return nil, errors.Wrapf(errors.ErrUnauthorized, "%s is not allowed to receive funds", tx.ToAddress)
		}

		defer func() {
			for _, a := range tx.Amount {
				if a.Amount.IsInt64() {
					telemetry.SetGaugeWithLabels(
						[]string{"tx", "msg", "send"},
						float32(a.Amount.Int64()),
						[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
					)
				}
			}
		}()

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			),
		)

		delete(record.Transactions.Record, hash)
	}

	s.keeper.AddToCustodyPool(ctx, record)

	return &types.MsgApproveCustodyTransactionResponse{}, nil
}

func (s msgServer) DeclineTransaction(goCtx context.Context, msg *types.MsgDeclineCustodyTransaction) (*types.MsgDeclineCustodyTransactionResponse, error) {
	//todo: add fee to msg.FromAddress
	//todo: make possibility to vote only 1 time

	return &types.MsgDeclineCustodyTransactionResponse{}, nil
}

func (s msgServer) DropCustodians(goCtx context.Context, msg *types.MsgDropCustodyCustodians) (*types.MsgDropCustodyCustodiansResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	s.keeper.DropCustodyCustodiansByAddress(ctx, msg.Address)

	return &types.MsgDropCustodyCustodiansResponse{}, nil
}

func (s msgServer) AddToWhiteList(goCtx context.Context, msg *types.MsgAddToCustodyWhiteList) (*types.MsgAddToCustodyWhiteListResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	record := types.CustodyWhiteListRecord{
		Address:          msg.Address,
		CustodyWhiteList: s.keeper.GetCustodyWhiteListByAddress(ctx, msg.Address),
	}

	if record.CustodyWhiteList == nil {
		record.CustodyWhiteList = new(types.CustodyWhiteList)
		record.CustodyWhiteList.Addresses = map[string]bool{}
	}

	for _, address := range msg.AddAddress {
		record.CustodyWhiteList.Addresses[address.String()] = true
	}

	s.keeper.AddToCustodyWhiteList(ctx, record)

	return &types.MsgAddToCustodyWhiteListResponse{}, nil
}

func (s msgServer) RemoveFromWhiteList(goCtx context.Context, msg *types.MsgRemoveFromCustodyWhiteList) (*types.MsgRemoveFromCustodyWhiteListResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	record := types.CustodyWhiteListRecord{
		Address:          msg.Address,
		CustodyWhiteList: s.keeper.GetCustodyWhiteListByAddress(ctx, msg.Address),
	}

	if record.CustodyWhiteList == nil {
		return nil, errors.Wrap(types.ErrNoWhiteLists, "Can not remove from the empty whitelist")
	}

	if !record.CustodyWhiteList.Addresses[msg.RemoveAddress.String()] {
		return nil, errors.Wrap(types.ErrNoWhiteListsElement, "Can not remove missing element from the whitelist")
	}

	record.CustodyWhiteList.Addresses[msg.RemoveAddress.String()] = false
	s.keeper.AddToCustodyWhiteList(ctx, record)

	return &types.MsgRemoveFromCustodyWhiteListResponse{}, nil
}

func (s msgServer) DropWhiteList(goCtx context.Context, msg *types.MsgDropCustodyWhiteList) (*types.MsgDropCustodyWhiteListResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	s.keeper.DropCustodyWhiteListByAddress(ctx, msg.Address)

	return &types.MsgDropCustodyWhiteListResponse{}, nil
}

func (s msgServer) AddToLimits(goCtx context.Context, msg *types.MsgAddToCustodyLimits) (*types.MsgAddToCustodyLimitsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	record := types.CustodyLimitRecord{
		Address:       msg.Address,
		CustodyLimits: s.keeper.GetCustodyLimitsByAddress(ctx, msg.Address),
	}

	if record.CustodyLimits == nil {
		record.CustodyLimits = new(types.CustodyLimits)
		record.CustodyLimits.Limits = map[string]*types.CustodyLimit{}
	}

	custodyLimit := types.CustodyLimit{
		Amount: msg.Amount,
		Limit:  msg.Limit,
	}

	record.CustodyLimits.Limits[msg.Denom] = &custodyLimit
	s.keeper.AddToCustodyLimits(ctx, record)

	return &types.MsgAddToCustodyLimitsResponse{}, nil
}

func (s msgServer) RemoveFromLimits(goCtx context.Context, msg *types.MsgRemoveFromCustodyLimits) (*types.MsgRemoveFromCustodyLimitsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	record := types.CustodyLimitRecord{
		Address:       msg.Address,
		CustodyLimits: s.keeper.GetCustodyLimitsByAddress(ctx, msg.Address),
	}

	if record.CustodyLimits == nil {
		return nil, errors.Wrap(types.ErrNoWhiteLists, "Can not remove from the empty limits")
	}

	if record.CustodyLimits.Limits[msg.Denom] == nil {
		return nil, errors.Wrap(types.ErrNoWhiteListsElement, "Can not remove missing element from the limits")
	}

	record.CustodyLimits.Limits[msg.Denom] = new(types.CustodyLimit)
	s.keeper.AddToCustodyLimits(ctx, record)

	return &types.MsgRemoveFromCustodyLimitsResponse{}, nil
}

func (s msgServer) DropLimits(goCtx context.Context, msg *types.MsgDropCustodyLimits) (*types.MsgDropCustodyLimitsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	s.keeper.DropCustodyLimitsByAddress(ctx, msg.Address)

	return &types.MsgDropCustodyLimitsResponse{}, nil
}

func (s msgServer) Send(goCtx context.Context, msg *types.MsgSend) (*types.MsgSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	to, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return nil, err
	}

	if s.bk.BlockedAddr(to) {
		return nil, errors.Wrapf(errors.ErrUnauthorized, "%s is not allowed to receive funds", msg.ToAddress)
	}

	if err := s.bk.IsSendEnabledCoins(ctx, msg.Amount...); err != nil {
		return nil, err
	}

	settings := s.keeper.GetCustodyInfoByAddress(ctx, msg.GetSigners()[0])
	custodians := s.keeper.GetCustodyCustodiansByAddress(ctx, msg.GetSigners()[0])

	record := types.CustodyPool{
		Address: msg.GetSigners()[0],
		Transactions: &types.TransactionPool{
			Record: map[string]*types.TransactionRecord{},
		},
	}

	hash := sha256.Sum256(ctx.TxBytes())
	hashString := hex.EncodeToString(hash[:])

	record.Transactions.Record[hashString] = &types.TransactionRecord{
		Transaction: msg,
		Votes:       0,
		Confirmed:   false,
	}

	if settings != nil && (settings.CustodyEnabled && len(custodians.Addresses) > 0 || settings.UsePassword) {
		s.keeper.AddToCustodyPool(ctx, record)
	} else {
		from, err := sdk.AccAddressFromBech32(msg.FromAddress)
		if err != nil {
			return nil, err
		}

		err = s.bk.SendCoins(ctx, from, to, msg.Amount)
		if err != nil {
			return nil, err
		}

		defer func() {
			for _, a := range msg.Amount {
				if a.Amount.IsInt64() {
					telemetry.SetGaugeWithLabels(
						[]string{"tx", "msg", "send"},
						float32(a.Amount.Int64()),
						[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
					)
				}
			}
		}()

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			),
		)
	}

	return &types.MsgSendResponse{}, nil
}

func (s msgServer) PasswordConfirm(goCtx context.Context, msg *types.MsgPasswordConfirmTransaction) (*types.MsgPasswordConfirmTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	hash := strings.ToLower(msg.Hash)
	allowCustodians := true
	allowPassword := true
	settings := s.keeper.GetCustodyInfoByAddress(ctx, msg.SenderAddress)
	custodians := s.keeper.GetCustodyCustodiansByAddress(ctx, msg.SenderAddress)

	record := types.CustodyPool{
		Address:      msg.SenderAddress,
		Transactions: s.keeper.GetCustodyPoolByAddress(ctx, msg.SenderAddress),
	}

	if record.Transactions != nil && record.Transactions.Record[hash] != nil {
		record.Transactions.Record[hash].Confirmed = true
	}

	if settings != nil && settings.CustodyEnabled && len(custodians.Addresses) > 0 {
		allowCustodians = (record.Transactions.Record[hash].Votes)*100/uint64(len(custodians.Addresses)) >= settings.CustodyMode
	}

	if settings != nil && settings.UsePassword {
		allowPassword = record.Transactions.Record[hash].Confirmed
	}

	if allowCustodians && allowPassword {
		tx := record.Transactions.Record[hash].Transaction

		to, err := sdk.AccAddressFromBech32(tx.ToAddress)
		if err != nil {
			return nil, err
		}

		if err := s.bk.IsSendEnabledCoins(ctx, tx.Amount...); err != nil {
			return nil, err
		}

		from, err := sdk.AccAddressFromBech32(tx.FromAddress)
		if err != nil {
			return nil, err
		}

		err = s.bk.SendCoins(ctx, from, to, tx.Amount)
		if err != nil {
			return nil, err
		}

		if s.bk.BlockedAddr(to) {
			return nil, errors.Wrapf(errors.ErrUnauthorized, "%s is not allowed to receive funds", tx.ToAddress)
		}

		defer func() {
			for _, a := range tx.Amount {
				if a.Amount.IsInt64() {
					telemetry.SetGaugeWithLabels(
						[]string{"tx", "msg", "send"},
						float32(a.Amount.Int64()),
						[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
					)
				}
			}
		}()

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			),
		)

		delete(record.Transactions.Record, hash)
	}

	s.keeper.AddToCustodyPool(ctx, record)

	return &types.MsgPasswordConfirmTransactionResponse{}, nil
}