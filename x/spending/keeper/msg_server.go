package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/spending/types"
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

// CreateSpendingPool is a function to allow creating a new spending pool.
// This function can be sent by any account. The person sending the transaction automatically becomes the pool owner.
// The original owner should provide a unique pool name when sending create tx.
func (k msgServer) CreateSpendingPool(
	goCtx context.Context,
	msg *types.MsgCreateSpendingPool,
) (*types.MsgCreateSpendingPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !types.ValidateSpendingPoolName(msg.Name) {
		return nil, types.ErrInvalidSpendingPoolName
	}

	err := k.keeper.CreateSpendingPool(ctx, types.SpendingPool{
		Name:                    msg.Name,
		ClaimStart:              msg.ClaimStart,
		ClaimEnd:                msg.ClaimEnd,
		ClaimExpiry:             msg.ClaimExpiry,
		Rates:                   msg.Rates,
		VoteQuorum:              msg.VoteQuorum,
		VotePeriod:              msg.VotePeriod,
		VoteEnactment:           msg.VoteEnactment,
		Owners:                  &msg.Owners,
		Beneficiaries:           &msg.Beneficiaries,
		Balances:                sdk.Coins{},
		DynamicRate:             msg.DynamicRate,
		DynamicRatePeriod:       msg.DynamicRatePeriod,
		LastDynamicRateCalcTime: uint64(ctx.BlockTime().Unix()),
	})

	if err != nil {
		return nil, err
	}

	return &types.MsgCreateSpendingPoolResponse{}, nil
}

// DepositSpendingPool is a function to allow depositing tokens to the pool address (name).
// Any KIRA address should be able to call this function and deposit tokens.
func (k msgServer) DepositSpendingPool(
	goCtx context.Context,
	msg *types.MsgDepositSpendingPool,
) (*types.MsgDepositSpendingPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	err = k.bk.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, msg.Amount)
	if err != nil {
		return nil, err
	}

	pool := k.keeper.GetSpendingPool(ctx, msg.PoolName)
	if pool == nil {
		return nil, types.ErrPoolDoesNotExist
	}

	pool.Balances = sdk.Coins(pool.Balances).Add(msg.Amount...)
	k.keeper.SetSpendingPool(ctx, *pool)

	return &types.MsgDepositSpendingPoolResponse{}, nil
}

// RegisterSpendingPoolBeneficiary is a function to register beneficiary account to be
// eligible for claims
func (k msgServer) RegisterSpendingPoolBeneficiary(
	goCtx context.Context,
	msg *types.MsgRegisterSpendingPoolBeneficiary,
) (*types.MsgRegisterSpendingPoolBeneficiaryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pool := k.keeper.GetSpendingPool(ctx, msg.PoolName)
	if pool == nil {
		return nil, types.ErrPoolDoesNotExist
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	if !k.keeper.IsAllowedBeneficiary(ctx, sender, *pool.Beneficiaries) {
		return nil, types.ErrNotPoolOwner
	}

	k.keeper.SetClaimInfo(ctx, types.ClaimInfo{
		PoolName:  pool.Name,
		Account:   sender.String(),
		LastClaim: uint64(ctx.BlockTime().Unix()),
	})

	return &types.MsgRegisterSpendingPoolBeneficiaryResponse{}, nil
}

// ClaimSpendingPool is a function to allow claiming tokens from the pool.
// Only beneficiaries should be able to send this transaction.
// Funds can be claimed only for the period between current bloct time and value set in the claims property in accordance to the current distribution rate. If the pool doesn't have a sufficient balance of a specific token as defined by tokens property then that specific token should NOT be sent in any amount. If the pool has sufficient funds as defined by the amount in the tokens property then exact amount owed should be sent to the beneficiary. All tokens that can be sent should be sent all at once to the account that is claiming them. If the claim expiration period elapsed and funds were NOT claimed by the beneficiary then the funds will NOT be sent. Beneficiary will only receive tokens if he already registered and his account is present in the claims array. Claiming of specific token should be only possible if and only if the spending pool has sufficient funds to distribute funds to ALL accounts eligible for claiming them (either all eligible accounts can claim a specific token or no one).
func (k msgServer) ClaimSpendingPool(
	goCtx context.Context,
	msg *types.MsgClaimSpendingPool,
) (*types.MsgClaimSpendingPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	err = k.keeper.ClaimSpendingPool(ctx, msg.PoolName, sender)
	if err != nil {
		return nil, err
	}

	return &types.MsgClaimSpendingPoolResponse{}, nil
}
