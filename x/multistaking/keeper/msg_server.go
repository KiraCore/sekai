package keeper

import (
	"context"

	govkeeper "github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/multistaking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	keeper    Keeper
	govKeeper govkeeper.Keeper
	sk        types.StakingKeeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, govKeeper govkeeper.Keeper, sk types.StakingKeeper) types.MsgServer {
	return &msgServer{
		keeper:    keeper,
		govKeeper: govKeeper,
		sk:        sk,
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
		k.keeper.SetStakingPool(ctx, pool)
	} else {
		k.keeper.SetStakingPool(ctx, types.StakingPool{
			Enabled:            msg.Enabled,
			Validator:          msg.Validator,
			TotalStakingTokens: []sdk.Coin{},
			TotalShare:         []sdk.Coin{},
			TotalRewards:       []sdk.Coin{},
		})
	}

	return &types.MsgUpsertStakingPoolResponse{}, nil
}

func (k msgServer) Delegate(goCtx context.Context, msg *types.MsgDelegate) (*types.MsgDelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return &types.MsgDelegateResponse{}, nil
}

func (k msgServer) Undelegate(goCtx context.Context, msg *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return &types.MsgUndelegateResponse{}, nil
}

func (k msgServer) ClaimRewards(goCtx context.Context, msg *types.MsgClaimRewards) (*types.MsgClaimRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return &types.MsgClaimRewardsResponse{}, nil
}

func (k msgServer) ClaimUndelegation(goCtx context.Context, msg *types.MsgClaimUndelegation) (*types.MsgClaimUndelegationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return &types.MsgClaimUndelegationResponse{}, nil
}
