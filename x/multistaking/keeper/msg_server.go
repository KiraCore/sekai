package keeper

import (
	"context"

	govkeeper "github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/multistaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		if pool.Slashed.IsPositive() {
			return nil, types.ErrActionNotSupportedForSlashedPool
		}
		k.keeper.SetStakingPool(ctx, pool)
	} else {
		// increase id when creating a new pool
		lastPoolId := k.keeper.GetLastPoolId(ctx) + 1
		k.keeper.SetLastPoolId(ctx, lastPoolId)
		pool = types.StakingPool{
			Id:                 lastPoolId,
			Enabled:            msg.Enabled,
			Validator:          msg.Validator,
			Commission:         msg.Commission,
			TotalStakingTokens: []sdk.Coin{},
			TotalShareTokens:   []sdk.Coin{},
			TotalRewards:       []sdk.Coin{},
		}
		k.keeper.SetStakingPool(ctx, pool)
	}

	if k.keeper.hooks != nil {
		k.keeper.hooks.AfterUpsertStakingPool(ctx, valAddr, pool)
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

	err := k.keeper.Undelegate(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgUndelegateResponse{}, nil
}

func (k msgServer) ClaimRewards(goCtx context.Context, msg *types.MsgClaimRewards) (*types.MsgClaimRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	delegator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	_ = k.keeper.ClaimRewards(ctx, delegator)
	return &types.MsgClaimRewardsResponse{}, nil
}

func (k msgServer) ClaimUndelegation(goCtx context.Context, msg *types.MsgClaimUndelegation) (*types.MsgClaimUndelegationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

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

	k.keeper.RemoveUndelegation(ctx, undelegation.Id)

	return &types.MsgClaimUndelegationResponse{}, nil
}

func (k msgServer) ClaimMaturedUndelegations(goCtx context.Context, msg *types.MsgClaimMaturedUndelegations) (*types.MsgClaimMaturedUndelegationsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	delegator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	allUndelegations := k.keeper.GetAllUndelegations(ctx)
	for _, undelegation := range allUndelegations {
		if undelegation.Address != msg.Sender {
			continue
		}

		if uint64(ctx.BlockTime().Unix()) < undelegation.Expiry {
			continue
		}

		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delegator, undelegation.Amount)
		if err != nil {
			return nil, err
		}
		k.keeper.RemoveUndelegation(ctx, undelegation.Id)
	}

	return &types.MsgClaimMaturedUndelegationsResponse{}, nil
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
