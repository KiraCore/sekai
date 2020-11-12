package gov

import (
	"github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	iterator := k.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, ctx.BlockTime())
	for ; iterator.Valid(); iterator.Next() {
		processEnactmentProposal(ctx, k, keeper.BytesToProposalID(iterator.Value()))
	}

	iterator = k.GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx, ctx.BlockTime())
	for ; iterator.Valid(); iterator.Next() {
		processProposal(ctx, k, keeper.BytesToProposalID(iterator.Value()))
	}
}

func processProposal(ctx sdk.Context, k keeper.Keeper, proposalID uint64) {
	votes := k.GetProposalVotes(ctx, proposalID)

	// TODO: this should get availableVoters by proposal type
	availableVoters := k.GetNetworkActorsByAbsoluteWhitelistPermission(ctx, types.PermVoteSetPermissionProposal)
	totalVoters := len(availableVoters)
	numVotes := len(votes)

	quorum := k.GetNetworkProperties(ctx).VoteQuorum

	isQuorum, err := types.IsQuorum(quorum, uint64(numVotes), uint64(totalVoters))
	if err != nil {
		panic(err)
	}

	if !isQuorum {
		return
	}

	numActorsWithVeto := len(types.GetActorsWithVoteWithVeto(availableVoters))
	calculatedVote := types.CalculateVotes(votes, uint64(numActorsWithVeto))

	proposal, found := k.GetProposal(ctx, proposalID)
	if !found {
		panic("proposal was expected to exist")
	}

	proposal.Result = calculatedVote.ProcessResult()

	k.SaveProposal(ctx, proposal)
	k.RemoveActiveProposal(ctx, proposal)
	k.AddToEnactmentProposals(ctx, proposal)
}

func processEnactmentProposal(ctx sdk.Context, k keeper.Keeper, proposalID uint64) {
	proposal, found := k.GetProposal(ctx, proposalID)
	if !found {
		panic("proposal was expected to exist")
	}

	if proposal.Result == types.Passed {
		switch proposal.GetContent().ProposalType() {
		case types.AssignPermissionProposalType:
			applyAssignPermissionProposal(ctx, k, proposal)
		default:
			panic("invalid proposal type")
		}

	}

	k.RemoveEnactmentProposal(ctx, proposal)
}

func applyAssignPermissionProposal(ctx sdk.Context, k keeper.Keeper, proposal types.Proposal) {
	p := proposal.GetContent().(*types.AssignPermissionProposal)

	actor, found := k.GetNetworkActorByAddress(ctx, p.Address)
	if !found {
		actor = types.NewDefaultActor(p.Address)
	}

	err := k.AddWhitelistPermission(ctx, actor, types.PermValue(p.Permission))
	if err != nil {
		panic("network actor has this permission")
	}
}
