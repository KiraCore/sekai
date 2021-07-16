package gov

import (
	"fmt"

	"github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	enactmentIterator := k.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, ctx.BlockTime())
	defer enactmentIterator.Close()
	for ; enactmentIterator.Valid(); enactmentIterator.Next() {
		processEnactmentProposal(ctx, k, keeper.BytesToProposalID(enactmentIterator.Value()))
	}

	activeIterator := k.GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx, ctx.BlockTime())
	defer activeIterator.Close()
	for ; activeIterator.Valid(); activeIterator.Next() {
		processProposal(ctx, k, keeper.BytesToProposalID(activeIterator.Value()))
	}
}

func processProposal(ctx sdk.Context, k keeper.Keeper, proposalID uint64) {
	proposal, found := k.GetProposal(ctx, proposalID)
	if !found {
		panic("proposal was expected to exist")
	}

	// skip execution if block height condition does not meet
	if proposal.MinVotingEndBlockHeight > ctx.BlockHeight() {
		return
	}

	votes := k.GetProposalVotes(ctx, proposalID)

	availableVoters := k.GetNetworkActorsByAbsoluteWhitelistPermission(ctx, proposal.GetContent().VotePermission())
	totalVoters := len(availableVoters)
	numVotes := len(votes)

	properties := k.GetNetworkProperties(ctx)
	quorum := properties.VoteQuorum

	isQuorum, err := types.IsQuorum(quorum, uint64(numVotes), uint64(totalVoters))
	if err != nil {
		panic(err)
	}

	if isQuorum {
		numActorsWithVeto := len(types.GetActorsWithVoteWithVeto(availableVoters))
		calculatedVote := types.CalculateVotes(votes, uint64(numActorsWithVeto))

		proposal.Result = calculatedVote.ProcessResult()
		if proposal.Result == types.Passed { // This is done in order to show that proposal is in enactment, but after enactment passes it will be passed.
			proposal.Result = types.Enactment
		}
	} else {
		proposal.Result = types.QuorumNotReached
	}

	// enactment time should be at least 1 block from voting period finish
	proposal.MinEnactmentEndBlockHeight = ctx.BlockHeight() + int64(properties.MinProposalEnactmentBlocks)
	k.SaveProposal(ctx, proposal)
	k.RemoveActiveProposal(ctx, proposal)
	k.AddToEnactmentProposals(ctx, proposal)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddToEnactment,
			sdk.NewAttribute(types.AttributeKeyProposalId, fmt.Sprintf("%d", proposal.ProposalId)),
			sdk.NewAttribute(types.AttributeKeyProposalDescription, proposal.Description),
		),
	)
}

func processEnactmentProposal(ctx sdk.Context, k keeper.Keeper, proposalID uint64) {
	router := k.GetProposalRouter()
	proposal, found := k.GetProposal(ctx, proposalID)
	if !found {
		panic("proposal was expected to exist")
	}

	// skip execution if block height condition does not meet
	if proposal.MinEnactmentEndBlockHeight > ctx.BlockHeight() {
		return
	}

	if proposal.Result == types.Enactment {
		err := router.ApplyProposal(ctx, proposal.GetContent())
		if err != nil {
			proposal.ExecResult = "execution failed"
		} else {
			proposal.ExecResult = "executed successfully"
		}
		proposal.Result = types.Passed
		k.SaveProposal(ctx, proposal)
	}

	k.RemoveEnactmentProposal(ctx, proposal)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveEnactment,
			sdk.NewAttribute(types.AttributeKeyProposalId, fmt.Sprintf("%d", proposal.ProposalId)),
			sdk.NewAttribute(types.AttributeKeyProposalDescription, proposal.Description),
		),
	)
}
