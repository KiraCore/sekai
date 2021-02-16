package gov

import (
	"github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper, router ProposalRouter) {
	iterator := k.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, ctx.BlockTime())
	for ; iterator.Valid(); iterator.Next() {
		processEnactmentProposal(ctx, k, router, keeper.BytesToProposalID(iterator.Value()))
	}

	iterator = k.GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx, ctx.BlockTime())
	for ; iterator.Valid(); iterator.Next() {
		processProposal(ctx, k, keeper.BytesToProposalID(iterator.Value()))
	}
}

func processProposal(ctx sdk.Context, k keeper.Keeper, proposalID uint64) {
	proposal, found := k.GetProposal(ctx, proposalID)
	if !found {
		panic("proposal was expected to exist")
	}

	votes := k.GetProposalVotes(ctx, proposalID)

	availableVoters := k.GetNetworkActorsByAbsoluteWhitelistPermission(ctx, proposal.GetContent().VotePermission())
	totalVoters := len(availableVoters)
	numVotes := len(votes)

	quorum := k.GetNetworkProperties(ctx).VoteQuorum

	isQuorum, err := types.IsQuorum(quorum, uint64(numVotes), uint64(totalVoters))
	if err != nil {
		panic(err)
	}

	if isQuorum {
		numActorsWithVeto := len(types.GetActorsWithVoteWithVeto(availableVoters))
		calculatedVote := types.CalculateVotes(votes, uint64(numActorsWithVeto))

		proposal.Result = calculatedVote.ProcessResult()
	} else {
		proposal.Result = types.QuorumNotReached
	}

	k.SaveProposal(ctx, proposal)
	k.RemoveActiveProposal(ctx, proposal)
	k.AddToEnactmentProposals(ctx, proposal)
}

func processEnactmentProposal(ctx sdk.Context, k keeper.Keeper, router ProposalRouter, proposalID uint64) {
	proposal, found := k.GetProposal(ctx, proposalID)
	if !found {
		panic("proposal was expected to exist")
	}

	if proposal.Result == types.Passed {
		router.ApplyProposal(ctx, proposal.GetContent())
	}

	k.RemoveEnactmentProposal(ctx, proposal)
}
