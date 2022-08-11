package keeper

import (
	"context"
	"fmt"
	"strings"

	govkeeper "github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/multistaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type msgServer struct {
	keeper     Keeper
	bankKeeper types.BankKeeper
	govKeeper  govkeeper.Keeper
	sk         types.StakingKeeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, bankKeeper types.BankKeeper, govKeeper govkeeper.Keeper, sk types.StakingKeeper) types.MsgServer {
	return &msgServer{
		keeper:     keeper,
		bankKeeper: bankKeeper,
		govKeeper:  govKeeper,
		sk:         sk,
	}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) UpsertStakingPool(goCtx context.Context, msg *types.MsgUpsertStakingPool) (*types.MsgUpsertStakingPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, err
	}
	validator, err := k.sk.GetValidator(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	// check sender is validator owner
	if sdk.AccAddress(validator.ValKey).String() != msg.Sender {
		return nil, types.ErrNotValidatorOwner
	}

	// check previous pool exists and if exists return error
	pool, found := k.keeper.GetStakingPoolByValidator(ctx, msg.Validator)
	if found {
		pool.Enabled = msg.Enabled
		if pool.Slashed > 0 {
			return nil, types.ErrActionNotSupportedForSlashedPool
		}
		k.keeper.SetStakingPool(ctx, pool)
	} else {
		// increase id when creating a new pool
		lastPoolId := k.keeper.GetLastPoolId(ctx) + 1
		k.keeper.SetLastPoolId(ctx, lastPoolId)

		k.keeper.SetStakingPool(ctx, types.StakingPool{
			Id:                 lastPoolId,
			Enabled:            msg.Enabled,
			Validator:          msg.Validator,
			TotalStakingTokens: []sdk.Coin{},
			TotalShareTokens:   []sdk.Coin{},
			TotalRewards:       []sdk.Coin{},
		})
	}

	return &types.MsgUpsertStakingPoolResponse{}, nil
}

func (k msgServer) Delegate(goCtx context.Context, msg *types.MsgDelegate) (*types.MsgDelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.keeper.Delegate(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgDelegateResponse{}, nil
}

func (k msgServer) Undelegate(goCtx context.Context, msg *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pool, found := k.keeper.GetStakingPoolByValidator(ctx, msg.ValidatorAddress)
	if !found {
		return nil, types.ErrStakingPoolNotFound
	}

	delegator, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	poolCoins := getPoolCoins(pool, msg.Amounts)

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, delegator, types.ModuleName, poolCoins)
	if err != nil {
		return nil, err
	}
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, poolCoins)
	if err != nil {
		return nil, err
	}

	pool.TotalStakingTokens = sdk.Coins(pool.TotalStakingTokens).Sub(msg.Amounts)
	pool.TotalShareTokens = sdk.Coins(pool.TotalShareTokens).Sub(poolCoins)
	k.keeper.SetStakingPool(ctx, pool)

	lastUndelegationId := k.keeper.GetLastUndelegationId(ctx) + 1
	k.keeper.SetLastUndelegationId(ctx, lastUndelegationId)
	properties := k.govKeeper.GetNetworkProperties(ctx)
	k.keeper.SetUndelegation(ctx, types.Undelegation{
		Id:      lastUndelegationId,
		Address: msg.DelegatorAddress,
		Expiry:  uint64(ctx.BlockTime().Unix()) + properties.UnstakingPeriod,
		Amount:  msg.Amounts,
	})

	balances := k.bankKeeper.GetAllBalances(ctx, delegator)
	prefix := fmt.Sprintf("v%d_", pool.Id)
	if !strings.Contains(balances.String(), prefix) {
		k.keeper.RemovePoolDelegator(ctx, pool.Id, delegator)
	}

	return &types.MsgUndelegateResponse{}, nil
}

func (k msgServer) ClaimRewards(goCtx context.Context, msg *types.MsgClaimRewards) (*types.MsgClaimRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	delegator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	rewards := k.keeper.GetDelegatorRewards(ctx, delegator)
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, delegator, rewards)
	if err != nil {
		return nil, err
	}

	k.keeper.RemoveDelegatorRewards(ctx, delegator)

	return &types.MsgClaimRewardsResponse{}, nil
}

func (k msgServer) ClaimUndelegation(goCtx context.Context, msg *types.MsgClaimUndelegation) (*types.MsgClaimUndelegationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	undelegation, found := k.keeper.GetUndelegationById(ctx, msg.UndelegationId)
	if !found {
		return nil, types.ErrUndelegationNotFound
	}

	if uint64(ctx.BlockTime().Unix()) < undelegation.Expiry {
		return nil, types.ErrNotEnoughTimePassed
	}

	delegator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delegator, undelegation.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgClaimUndelegationResponse{}, nil
}

func (k msgServer) SetCompoundInfo(goCtx context.Context, msg *types.MsgSetCompoundInfo) (*types.MsgSetCompoundInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	k.keeper.SetCompoundInfo(ctx, types.CompoundInfo{
		Delegator:      msg.Sender,
		AllDenom:       msg.AllDenom,
		CompoundDenoms: msg.CompoundDenoms,
	})

	return &types.MsgSetCompoundInfoResponse{}, nil
}

func (k msgServer) RegisterDelegator(goCtx context.Context, msg *types.MsgRegisterDelegator) (*types.MsgRegisterDelegatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	delegator, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return nil, err
	}

	k.keeper.RegisterDelegator(ctx, delegator)

	return &types.MsgRegisterDelegatorResponse{}, nil
}
