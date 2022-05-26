package keeper

import (
	"context"

	"github.com/KiraCore/sekai/x/multistaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Querier struct {
	keeper Keeper
}

func NewQuerier(keeper Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}

var _ types.QueryServer = Querier{}

func (q Querier) StakingPools(c context.Context, request *types.QueryStakingPoolsRequest) (*types.QueryStakingPoolsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryStakingPoolsResponse{
		Pools: q.keeper.GetAllStakingPools(ctx),
	}, nil
}

func (q Querier) OutstandingRewards(c context.Context, request *types.QueryOutstandingRewardsRequest) (*types.QueryOutstandingRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	delegator, err := sdk.AccAddressFromBech32(request.Delegator)
	if err != nil {
		return nil, err
	}

	return &types.QueryOutstandingRewardsResponse{
		Rewards: q.keeper.GetDelegatorRewards(ctx, delegator),
	}, nil
}

func (q Querier) Undelegations(c context.Context, request *types.QueryUndelegationsRequest) (*types.QueryUndelegationsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryUndelegationsResponse{
		Undelegations: q.keeper.GetAllUndelegations(ctx),
	}, nil
}
