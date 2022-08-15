package ante

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	kiratypes "github.com/KiraCore/sekai/types"
	custodykeeper "github.com/KiraCore/sekai/x/custody/keeper"
	custodytypes "github.com/KiraCore/sekai/x/custody/types"
	feeprocessingkeeper "github.com/KiraCore/sekai/x/feeprocessing/keeper"
	feeprocessingtypes "github.com/KiraCore/sekai/x/feeprocessing/types"
	customgovkeeper "github.com/KiraCore/sekai/x/gov/keeper"
	customstakingkeeper "github.com/KiraCore/sekai/x/staking/keeper"
	tokenskeeper "github.com/KiraCore/sekai/x/tokens/keeper"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gogo/protobuf/proto"
	"time"
)

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(
	sk customstakingkeeper.Keeper,
	cgk customgovkeeper.Keeper,
	tk tokenskeeper.Keeper,
	fk feeprocessingkeeper.Keeper,
	ak keeper.AccountKeeper,
	bk types.BankKeeper,
	sigGasConsumer ante.SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
	ck custodykeeper.Keeper,
) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		NewCustodyDecorator(ck),
		NewZeroGasMeterDecorator(),
		ante.NewRejectExtensionOptionsDecorator(),
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.TxTimeoutHeightDecorator{},
		ante.NewValidateMemoDecorator(ak),
		ante.NewConsumeGasForTxSizeDecorator(ak),
		// custom fee range validator
		NewValidateFeeRangeDecorator(sk, cgk, tk, ak),
		ante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(ak),
		ante.NewDeductFeeDecorator(ak, bk, nil),
		// poor network management decorator
		NewPoorNetworkManagementDecorator(ak, cgk, sk),
		NewBlackWhiteTokensCheckDecorator(cgk, sk, tk),
		// custom execution fee consume decorator
		NewExecutionFeeRegistrationDecorator(ak, cgk, fk),
		ante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		ante.NewSigVerificationDecorator(ak, signModeHandler),
		ante.NewIncrementSequenceDecorator(ak),
	)
}

type CustodyDecorator struct {
	ck custodykeeper.Keeper
}

func NewCustodyDecorator(ck custodykeeper.Keeper) CustodyDecorator {
	return CustodyDecorator{
		ck: ck,
	}
}

func (cd CustodyDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	for _, msg := range feeTx.GetMsgs() {
		settings := cd.ck.GetCustodyInfoByAddress(ctx, msg.GetSigners()[0])

		if settings != nil && kiratypes.MsgType(msg) == kiratypes.MsgTypeCreateCustody {
			msg := msg.(*custodytypes.MsgCreteCustodyRecord)

			hash := sha256.Sum256([]byte(msg.OldKey))
			hashString := hex.EncodeToString(hash[:])

			if hashString != settings.Key {
				return ctx, sdkerrors.Wrap(custodytypes.ErrWrongKey, "Custody module")
			}
		}

		if kiratypes.MsgType(msg) == kiratypes.MsgApproveCustodyTransaction { //Check approve
			//Todo: Reward
		}

		if kiratypes.MsgType(msg) == kiratypes.MsgDeclineCustodyTransaction { //Check decline
			//Todo: Reward
		}

		if kiratypes.MsgType(msg) == bank.TypeMsgSend {
			protoTx := msg.(proto.Message)
			msg := msg.(*bank.MsgSend)

			if settings != nil && settings.CustodyEnabled {
				custodians := cd.ck.GetCustodyCustodiansByAddress(ctx, msg.GetSigners()[0])

				//Todo: Check fee exists

				if len(custodians.Addresses) == 0 {
					continue
				}

				record := custodytypes.CustodyPool{
					Address:      msg.GetSigners()[0],
					Transactions: new(custodytypes.TransactionPool),
				}

				hash := sha256.Sum256(ctx.TxBytes())
				hashString := hex.EncodeToString(hash[:])
				anyValue, anyErr := codectypes.NewAnyWithValue(protoTx)

				if anyErr != nil {
					return ctx, sdkerrors.Wrap(sdkerrors.ErrPackAny, anyErr.Error())
				}

				record.Transactions.Record[hashString] = &custodytypes.TransactionRecord{
					Transaction: anyValue,
					Votes:       0,
				}

				//Todo: if Exists in pool and votes anought do not block and remove from the list
				cd.ck.AddToCustodyPool(ctx, record)

				return ctx, sdkerrors.Wrap(sdkerrors.ErrConflict, "Waiting for custodians check.")
			}

			if settings != nil && settings.UseWhiteList {
				whiteList := cd.ck.GetCustodyWhiteListByAddress(ctx, msg.GetSigners()[0])

				if whiteList != nil && !whiteList.Addresses[msg.ToAddress] {
					return ctx, custodytypes.ErrNotInWhiteList
				}
			}

			if settings != nil && settings.UseLimits {
				limits := cd.ck.GetCustodyLimitsByAddress(ctx, msg.GetSigners()[0])

				logger := cd.ck.Logger(ctx)
				logger.Info("Limits", "msg", msg)
				logger.Info("Limits", "limits", limits)

				custodyLimitStatusRecord := custodytypes.CustodyLimitStatusRecord{
					Address:         msg.GetSigners()[0],
					CustodyStatuses: cd.ck.GetCustodyLimitsStatusByAddress(ctx, msg.GetSigners()[0]),
				}

				newAmount := msg.Amount.AmountOf(msg.Amount[0].Denom).Uint64()

				if custodyLimitStatusRecord.CustodyStatuses != nil && custodyLimitStatusRecord.CustodyStatuses.Statuses[msg.Amount[0].Denom] != nil {
					limit, _ := time.ParseDuration(limits.Limits[msg.Amount[0].Denom].Limit)
					rate := limits.Limits[msg.Amount[0].Denom].Amount / (uint64(limit.Milliseconds()))
					period := uint64(time.Now().Unix() - custodyLimitStatusRecord.CustodyStatuses.Statuses[msg.Amount[0].Denom].Time)
					newAmount = custodyLimitStatusRecord.CustodyStatuses.Statuses[msg.Amount[0].Denom].Amount + msg.Amount.AmountOf(msg.Amount[0].Denom).Uint64() - (period * rate)

					if newAmount <= 0 {
						return ctx, custodytypes.ErrNotInLimits
					}
				}

				custodyLimitStatusRecord.CustodyStatuses.Statuses[msg.Amount[0].Denom] = &custodytypes.CustodyStatus{
					Amount: newAmount,
				}

				cd.ck.AddToCustodyLimitsStatus(ctx, custodyLimitStatusRecord)
			}
		}
	}

	return next(ctx, tx, simulate)
}

// ValidateFeeRangeDecorator check if fee is within range defined as network properties
type ValidateFeeRangeDecorator struct {
	sk  customstakingkeeper.Keeper
	cgk customgovkeeper.Keeper
	tk  tokenskeeper.Keeper
	ak  keeper.AccountKeeper
}

// NewValidateFeeRangeDecorator check if fee is within range defined as network properties
func NewValidateFeeRangeDecorator(
	sk customstakingkeeper.Keeper,
	cgk customgovkeeper.Keeper,
	tk tokenskeeper.Keeper,
	ak keeper.AccountKeeper,
) ValidateFeeRangeDecorator {
	return ValidateFeeRangeDecorator{
		sk:  sk,
		cgk: cgk,
		ak:  ak,
		tk:  tk,
	}
}

// AnteHandle is a handler for ValidateFeeRangeDecorator
func (svd ValidateFeeRangeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	properties := svd.cgk.GetNetworkProperties(ctx)
	bondDenom := svd.sk.BondDenom(ctx)

	feeAmount := sdk.NewDec(0)
	feeCoins := feeTx.GetFee()
	tokensBlackWhite := svd.tk.GetTokenBlackWhites(ctx)
	for _, feeCoin := range feeCoins {
		rate := svd.tk.GetTokenRate(ctx, feeCoin.Denom)
		if !properties.EnableForeignFeePayments && feeCoin.Denom != bondDenom {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("foreign fee payments is disabled by governance"))
		}
		if rate == nil || !rate.FeePayments {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("currency you are trying to use was not whitelisted as fee payment"))
		}
		if tokensBlackWhite.IsFrozen(feeCoin.Denom, bondDenom, properties.EnableTokenBlacklist, properties.EnableTokenWhitelist) {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("currency you are trying to use as fee is frozen"))
		}
		feeAmount = feeAmount.Add(feeCoin.Amount.ToDec().Mul(rate.FeeRate))
	}

	// execution fee should be prepaid
	executionMaxFee := uint64(0)
	for _, msg := range feeTx.GetMsgs() {
		fee := svd.cgk.GetExecutionFee(ctx, kiratypes.MsgType(msg))
		if fee != nil { // execution fee exist
			maxFee := fee.FailureFee
			if fee.ExecutionFee > maxFee {
				maxFee = fee.ExecutionFee
			}
			executionMaxFee += maxFee
		}
	}

	if feeAmount.LT(sdk.NewDec(int64(properties.MinTxFee))) || feeAmount.GT(sdk.NewDec(int64(properties.MaxTxFee))) {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("fee %+v(%d) is out of range [%d, %d]%s", feeTx.GetFee(), feeAmount.RoundInt().Int64(), properties.MinTxFee, properties.MaxTxFee, bondDenom))
	}

	if feeAmount.LT(sdk.NewDec(int64(executionMaxFee))) {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("fee %+v(%d) is less than max execution fee %d%s", feeTx.GetFee(), feeAmount.RoundInt().Int64(), executionMaxFee, bondDenom))
	}

	return next(ctx, tx, simulate)
}

// ExecutionFeeRegistrationDecorator register paid execution fee
type ExecutionFeeRegistrationDecorator struct {
	ak  keeper.AccountKeeper
	cgk customgovkeeper.Keeper
	fk  feeprocessingkeeper.Keeper
}

// NewExecutionFeeRegistrationDecorator returns instance of CustomExecutionFeeConsumeDecorator
func NewExecutionFeeRegistrationDecorator(ak keeper.AccountKeeper, cgk customgovkeeper.Keeper, fk feeprocessingkeeper.Keeper) ExecutionFeeRegistrationDecorator {
	return ExecutionFeeRegistrationDecorator{
		ak,
		cgk,
		fk,
	}
}

// AnteHandle handle ExecutionFeeRegistrationDecorator
func (sgcd ExecutionFeeRegistrationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	sigTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	// execution fee consume gas
	for _, msg := range sigTx.GetMsgs() {
		fee := sgcd.cgk.GetExecutionFee(ctx, kiratypes.MsgType(msg))
		if fee != nil { // execution fee exist
			sgcd.fk.AddExecutionStart(ctx, msg)
		}
	}

	return next(ctx, tx, simulate)
}

// PoorNetworkManagementDecorator register poor network manager
type PoorNetworkManagementDecorator struct {
	ak  keeper.AccountKeeper
	cgk customgovkeeper.Keeper
	csk customstakingkeeper.Keeper
}

// NewPoorNetworkManagementDecorator returns instance of PoorNetworkManagementDecorator
func NewPoorNetworkManagementDecorator(ak keeper.AccountKeeper, cgk customgovkeeper.Keeper, csk customstakingkeeper.Keeper) PoorNetworkManagementDecorator {
	return PoorNetworkManagementDecorator{
		ak,
		cgk,
		csk,
	}
}

func findString(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return -1
}

// AnteHandle handle PoorNetworkManagementDecorator
func (pnmd PoorNetworkManagementDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	sigTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	// if not poor network, skip this process
	if pnmd.csk.IsNetworkActive(ctx) {
		return next(ctx, tx, simulate)
	}
	// handle messages on poor network
	pnmsgs := pnmd.cgk.GetPoorNetworkMessages(ctx)
	for _, msg := range sigTx.GetMsgs() {
		if kiratypes.MsgType(msg) == bank.TypeMsgSend {
			// on poor network, we introduce POOR_NETWORK_MAX_BANK_TX_SEND network property to limit transaction send amount
			msg := msg.(*bank.MsgSend)
			if len(msg.Amount) > 1 || msg.Amount[0].Denom != pnmd.csk.BondDenom(ctx) {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "only bond denom is allowed on poor network")
			}
			if msg.Amount[0].Amount.Uint64() > pnmd.cgk.GetNetworkProperties(ctx).PoorNetworkMaxBankSend {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "only restricted amount send is allowed on poor network")
			}
			// TODO: we could do restriction to send only when target account does not exist on chain yet for more restriction
			return next(ctx, tx, simulate)
		}
		if findString(pnmsgs.Messages, kiratypes.MsgType(msg)) >= 0 {
			return next(ctx, tx, simulate)
		}
		return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid transaction type on poor network")
	}

	return next(ctx, tx, simulate)
}

// BlackWhiteTokensCheckDecorator register black white tokens check decorator
type BlackWhiteTokensCheckDecorator struct {
	cgk customgovkeeper.Keeper
	csk customstakingkeeper.Keeper
	tk  tokenskeeper.Keeper
}

// NewBlackWhiteTokensCheckDecorator returns instance of BlackWhiteTokensCheckDecorator
func NewBlackWhiteTokensCheckDecorator(cgk customgovkeeper.Keeper, csk customstakingkeeper.Keeper, tk tokenskeeper.Keeper) BlackWhiteTokensCheckDecorator {
	return BlackWhiteTokensCheckDecorator{
		cgk,
		csk,
		tk,
	}
}

// AnteHandle handle NewPoorNetworkManagementDecorator
func (pnmd BlackWhiteTokensCheckDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	sigTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	bondDenom := pnmd.csk.BondDenom(ctx)
	tokensBlackWhite := pnmd.tk.GetTokenBlackWhites(ctx)
	properties := pnmd.cgk.GetNetworkProperties(ctx)
	for _, msg := range sigTx.GetMsgs() {
		if kiratypes.MsgType(msg) == bank.TypeMsgSend {
			msg := msg.(*bank.MsgSend)
			for _, amt := range msg.Amount {
				if tokensBlackWhite.IsFrozen(amt.Denom, bondDenom, properties.EnableTokenBlacklist, properties.EnableTokenWhitelist) {
					return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "token is frozen")
				}
			}
		}
	}

	return next(ctx, tx, simulate)
}

// ZeroGasMeterDecorator uses infinite gas decorator to avoid gas usage in the network
type ZeroGasMeterDecorator struct{}

// NewZeroGasMeterDecorator returns instance of ZeroGasMeterDecorator
func NewZeroGasMeterDecorator() ZeroGasMeterDecorator {
	return ZeroGasMeterDecorator{}
}

func (igm ZeroGasMeterDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	newCtx = ctx.WithGasMeter(feeprocessingtypes.NewZeroGasMeter())
	return next(newCtx, tx, simulate)
}
