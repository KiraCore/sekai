package keeper

import (
	"context"

	"github.com/KiraCore/sekai/x/spending/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Querier struct {
	keeper Keeper
}

func NewQuerier(keeper Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}

var _ types.QueryServer = Querier{}

// query-pools - query list of pool names
func (q Querier) QueryPoolNames(c context.Context, request *types.QueryPoolNamesRequest) (*types.QueryPoolNamesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	_ = ctx
	return &types.QueryPoolNamesResponse{}, nil
}

// query-pool - query pool by name
func (q Querier) QueryPoolByName(c context.Context, request *types.QueryPoolByNameRequest) (*types.QueryPoolByNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	_ = ctx
	return &types.QueryPoolByNameResponse{}, nil
}

// query-pool-proposals - query pool proposals by name
func (q Querier) QueryPoolProposals(c context.Context, request *types.QueryPoolProposalsRequest) (*types.QueryPoolProposalsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	_ = ctx
	return &types.QueryPoolProposalsResponse{}, nil
}

// query-pool-proposals - query pool proposals by name
func (q Querier) QueryPoolsByAccount(c context.Context, request *types.QueryPoolsByAccountRequest) (*types.QueryPoolsByAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	_ = ctx
	return &types.QueryPoolsByAccountResponse{}, nil
}
