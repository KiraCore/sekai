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

	_ = ctx
	return &types.QueryStakingPoolsResponse{}, nil
}

func (q Querier) Delegations(c context.Context, request *types.QueryDelegationsRequest) (*types.QueryDelegationsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	_ = ctx
	return &types.QueryDelegationsResponse{}, nil
}

func (q Querier) Undelegations(c context.Context, request *types.QueryUndelegationsRequest) (*types.QueryUndelegationsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	_ = ctx
	return &types.QueryUndelegationsResponse{}, nil
}
