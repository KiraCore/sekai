package keeper

import (
	"context"

	kiratypes "github.com/KiraCore/sekai/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/spending/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Querier struct {
	keeper Keeper
	gk     types.CustomGovKeeper
}

func NewQuerier(keeper Keeper, gk types.CustomGovKeeper) types.QueryServer {
	return &Querier{
		keeper: keeper,
		gk:     gk,
	}
}

var _ types.QueryServer = Querier{}

// QueryPoolNames - query list of pool names
func (q Querier) QueryPoolNames(c context.Context, request *types.QueryPoolNamesRequest) (*types.QueryPoolNamesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	_ = ctx
	pools := q.keeper.GetAllSpendingPools(ctx)
	poolNames := []string{}
	for _, pool := range pools {
		poolNames = append(poolNames, pool.Name)
	}
	return &types.QueryPoolNamesResponse{
		Names: poolNames,
	}, nil
}

// QueryPoolByName - query pool by name
func (q Querier) QueryPoolByName(c context.Context, request *types.QueryPoolByNameRequest) (*types.QueryPoolByNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryPoolByNameResponse{
		Pool: q.keeper.GetSpendingPool(ctx, request.Name),
	}, nil
}

// QueryPoolProposals - query pool proposals by name
func (q Querier) QueryPoolProposals(c context.Context, request *types.QueryPoolProposalsRequest) (*types.QueryPoolProposalsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	proposals, err := q.gk.GetProposals(ctx)
	if err != nil {
		return nil, err
	}

	poolProposals := []govtypes.Proposal{}
	for _, proposal := range proposals {
		content := proposal.GetContent()
		switch content.ProposalType() {
		case kiratypes.ProposalTypeUpdateSpendingPool:
			exactContent, ok := content.(*types.UpdateSpendingPoolProposal)
			if !ok {
				return nil, types.ErrInvalidProposalExists
			}

			if exactContent.Name == request.PoolName {
				poolProposals = append(poolProposals, proposal)
			}

		case kiratypes.ProposalTypeSpendingPoolDistribution:

			exactContent, ok := content.(*types.SpendingPoolDistributionProposal)
			if !ok {
				return nil, types.ErrInvalidProposalExists
			}

			if exactContent.PoolName == request.PoolName {
				poolProposals = append(poolProposals, proposal)
			}

		case kiratypes.ProposalTypeSpendingPoolWithdraw:

			exactContent, ok := content.(*types.SpendingPoolWithdrawProposal)
			if !ok {
				return nil, types.ErrInvalidProposalExists
			}

			if exactContent.PoolName == request.PoolName {
				poolProposals = append(poolProposals, proposal)
			}
		}
	}
	return &types.QueryPoolProposalsResponse{
		Proposals: poolProposals,
	}, nil
}

// QueryPoolsByAccount - query pools where an account can claim rewards
func (q Querier) QueryPoolsByAccount(c context.Context, request *types.QueryPoolsByAccountRequest) (*types.QueryPoolsByAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	pools := q.keeper.GetAllSpendingPools(ctx)
	accPools := []types.SpendingPool{}

	acc, err := sdk.AccAddressFromBech32(request.Account)
	if err != nil {
		return nil, err
	}

	for _, pool := range pools {
		if q.keeper.IsAllowedBeneficiary(ctx, acc, *pool.Beneficiaries) {
			accPools = append(accPools, pool)
		}
	}

	return &types.QueryPoolsByAccountResponse{
		Pools: accPools,
	}, nil
}
