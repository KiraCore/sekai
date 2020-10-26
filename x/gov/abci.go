package gov

import (
	"github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	iterator := k.GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx, ctx.BlockTime())
	for ; iterator.Valid(); iterator.Next() {
		processProposal(ctx, k, keeper.BytesToProposalID(iterator.Value()))
	}
}

func processProposal(ctx sdk.Context, k keeper.Keeper, proposalID uint64) {
	votes := k.GetProposalVotes(ctx, proposalID)

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
}
