package keeper

import (
	"context"

	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/collectives/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
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
		Contributers: q.keeper.GetCollectiveContributers(ctx, request.Name),
	}, nil
}

// Collectives query list of all staking collectives (output list of names),
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
	proposals, err := q.keeper.gk.GetProposals(ctx)
	if err != nil {
		return nil, err
	}

	collectiveProposals := []govtypes.Proposal{}
	for _, proposal := range proposals {
		switch proposal.GetContent().ProposalType() {
		case kiratypes.ProposalTypeCollectiveSendDonation:
			fallthrough
		case kiratypes.ProposalTypeCollectiveUpdate:
			fallthrough
		case kiratypes.ProposalTypeCollectiveRemove:
			collectiveProposals = append(collectiveProposals, proposal)
		}
	}

	return &types.CollectivesProposalsResponse{
		Proposals: collectiveProposals,
	}, nil
}

// query list of staking collectives by an individual KIRA address
func (q Querier) CollectivesByAccount(c context.Context, request *types.CollectivesByAccountRequest) (*types.CollectivesByAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	collectives := q.keeper.GetAllCollectives(ctx)
	accCollectives := []types.Collective{}
	contributions := []types.CollectiveContributor{}
	for _, collective := range collectives {
		cc := q.keeper.GetCollectiveContributer(ctx, collective.Name, request.Account)
		if cc.Name == "" {
			accCollectives = append(accCollectives, collective)
			contributions = append(contributions, cc)
		}
	}

	return &types.CollectivesByAccountResponse{
		Collectives:   accCollectives,
		Contributions: contributions,
	}, nil
}
