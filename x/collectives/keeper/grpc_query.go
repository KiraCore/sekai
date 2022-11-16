package keeper

import (
	"context"

	"github.com/KiraCore/sekai/x/collectives/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Querier struct {
	keeper Keeper
}

func NewQuerier(keeper Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}

var _ types.QueryServer = Querier{}

// Collective queries a collective
func (q Querier) Collective(c context.Context, request *types.CollectiveRequest) (*types.CollectiveResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return &types.CollectiveResponse{
		Collective:   q.keeper.GetCollective(ctx, request.Name),
		Contributers: q.keeper.GetAllCollectiveContributers(ctx, request.Name),
	}, nil
}

// Collectives query list of all staking collectives (output list of names),
// if `name` / `id` is specified then output full details of a single collective.
func (q Querier) Collectives(c context.Context, request *types.CollectivesRequest) (*types.CollectivesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return &types.CollectivesResponse{
		Collectives: q.keeper.GetAllCollectives(ctx),
	}, nil
}

// list id of all proposals in regards to staking collectives,
// (or proposals in regards to a specific collective if `name` / `id` is specified in the query)
func (q Querier) CollectivesProposals(c context.Context, request *types.CollectivesProposalsRequest) (*types.CollectivesProposalsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_ = ctx
	// TODO:
	return &types.CollectivesProposalsResponse{}, nil
}

// query list of staking collectives by an individual KIRA address
func (q Querier) CollectivesByAccount(c context.Context, request *types.CollectivesByAccountRequest) (*types.CollectivesByAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_ = ctx
	// TODO:
	return &types.CollectivesByAccountResponse{}, nil
}
