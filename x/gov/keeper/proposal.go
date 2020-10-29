package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	sdktypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/KiraCore/sekai/x/gov/types"
)

func (k Keeper) GetNextProposalID(ctx sdk.Context) (uint64, error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(NextProposalIDPrefix)
	if bz == nil {
		return 0, errors.Wrap(sdktypes.ErrInvalidGenesis, "initial proposal ID hasn't been set")
	}

	proposalID := BytesToProposalID(bz)

	return proposalID, nil
}

func (k Keeper) SaveProposalID(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)

	store.Set(NextProposalIDPrefix, ProposalIDToBytes(proposalID))
}

func (k Keeper) SaveProposal(ctx sdk.Context, proposal types.ProposalAssignPermission) error {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&proposal)
	store.Set(GetProposalKey(proposal.ProposalId), bz)

	// Update NextProposal
	k.SaveProposalID(ctx, proposal.ProposalId+1)

	return nil
}

func (k Keeper) GetProposal(ctx sdk.Context, proposalID uint64) (types.ProposalAssignPermission, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetProposalKey(proposalID))
	if bz == nil {
		return types.ProposalAssignPermission{}, false
	}

	var prop types.ProposalAssignPermission
	k.cdc.MustUnmarshalBinaryBare(bz, &prop)

	return prop, true
}

func (k Keeper) SaveVote(ctx sdk.Context, vote types.Vote) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&vote)
	store.Set(VoteKey(vote.ProposalId, vote.Voter), bz)
}

func (k Keeper) GetVote(ctx sdk.Context, proposalID uint64, address sdk.AccAddress) (types.Vote, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(VoteKey(proposalID, address))
	if bz == nil {
		return types.Vote{}, false
	}

	var vote types.Vote
	k.cdc.MustUnmarshalBinaryBare(bz, &vote)

	return vote, true
}

func (k Keeper) GetProposalVotesIterator(ctx sdk.Context, proposalID uint64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, VotesKey(proposalID))
}

func (k Keeper) GetProposalVotes(ctx sdk.Context, proposalID uint64) types.Votes {
	var votes types.Votes

	iterator := k.GetProposalVotesIterator(ctx, proposalID)
	for ; iterator.Valid(); iterator.Next() {
		var vote types.Vote
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &vote)
		votes = append(votes, vote)
	}

	return votes
}

func (k Keeper) AddToActiveProposals(ctx sdk.Context, proposal types.ProposalAssignPermission) {
	store := ctx.KVStore(k.storeKey)
	store.Set(ActiveProposalKey(proposal), ProposalIDToBytes(proposal.ProposalId))
}

// GetActiveProposalsWithFinishedVotingEndTimeIterator returns the proposals that have endtime finished.
func (k Keeper) GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(ActiveProposalsPrefix, sdk.PrefixEndBytes(ActiveProposalByTimeKey(endTime)))
}

func ActiveProposalByTimeKey(endTime time.Time) []byte {
	return append(ActiveProposalsPrefix, sdk.FormatTimeBytes(endTime)...)
}

func ActiveProposalKey(prop types.ProposalAssignPermission) []byte {
	return append(ActiveProposalByTimeKey(prop.VotingEndTime), ProposalIDToBytes(prop.ProposalId)...)
}

func VotesKey(proposalID uint64) []byte {
	return append(VotesPrefix, ProposalIDToBytes(proposalID)...)
}

func VoteKey(proposalId uint64, address sdk.AccAddress) []byte {
	return append(VotesKey(proposalId), address.Bytes()...)
}

func GetProposalKey(proposalID uint64) []byte {
	return append(ProposalsPrefix, ProposalIDToBytes(proposalID)...)
}
