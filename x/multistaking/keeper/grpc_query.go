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
	allUndelegations := q.keeper.GetAllUndelegations(ctx)
	undelegations := []types.Undelegation{}
	for _, del := range allUndelegations {
		if request.Delegator != "" && del.Address != request.Delegator {
			continue
		}
		if request.ValAddress != "" && del.ValAddress != request.ValAddress {
			continue
		}
		undelegations = append(undelegations, del)
	}

	return &types.QueryUndelegationsResponse{
		Undelegations: undelegations,
	}, nil
}

func (q Querier) CompoundInfo(c context.Context, request *types.QueryCompoundInfoRequest) (*types.QueryCompoundInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryCompoundInfoResponse{
		Info: q.keeper.GetCompoundInfoByAddress(ctx, request.Delegator),
	}, nil
}

func (q Querier) StakingPoolDelegators(c context.Context, request *types.QueryStakingPoolDelegatorsRequest) (*types.QueryStakingPoolDelegatorsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	pool, found := q.keeper.GetStakingPoolByValidator(ctx, request.Validator)
	if !found {
		return nil, types.ErrStakingPoolNotFound
	}
	delegators := q.keeper.GetPoolDelegators(ctx, pool.Id)
	delegatorAddrs := []string{}
	for _, delegator := range delegators {
		delegatorAddrs = append(delegatorAddrs, delegator.String())
	}
	return &types.QueryStakingPoolDelegatorsResponse{
		Pool:       pool,
		Delegators: delegatorAddrs,
	}, nil
}
