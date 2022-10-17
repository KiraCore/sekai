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
		proposalID := keeper.BytesToProposalID(enactmentIterator.Value())
		slash := k.GetAverageVotesSlash(ctx, proposalID)
		processEnactmentProposal(ctx, k, proposalID, slash)
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

	// councilor rank update function on absent
	voteMap := make(map[string]bool)
	for _, voter := range availableVoters {
		voteMap[voter.Address.String()] = true
	}
	councilors := k.GetAllCouncilors(ctx)
	for _, councilor := range councilors {
		if !voteMap[councilor.Address.String()] {
			k.OnCouncilorAbsent(ctx, councilor.Address)
		}
	}

	// update to spending pool users if it's spending pool proposal
	content := proposal.GetContent()
	if content.VotePermission() == types.PermZero {
		router := k.GetProposalRouter()
		totalVoters = len(router.AllowedAddressesDynamicProposal(ctx, content))
		if totalVoters == 0 {
			totalVoters = 1
		}
	}
	numVotes := len(votes)

	properties := k.GetNetworkProperties(ctx)

	quorum := properties.VoteQuorum
	if content.VotePermission() == types.PermZero {
		router := k.GetProposalRouter()
		quorum = router.QuorumDynamicProposal(ctx, content)
	}

	isQuorum, err := types.IsQuorum(quorum, uint64(numVotes), uint64(totalVoters))
	if err != nil {
		panic(fmt.Sprintf("Invalid quorum on proposal: proposalID=%d, proposalType=%s, err=%+v", proposalID, proposal.GetContent().ProposalType(), err))
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

func processEnactmentProposal(ctx sdk.Context, k keeper.Keeper, proposalID uint64, slash uint64) {
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
		err := router.ApplyProposal(ctx, proposalID, proposal.GetContent(), slash)
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
