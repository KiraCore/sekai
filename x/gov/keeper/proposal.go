package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/gov/types"
)

func (k Keeper) GetNextProposalIDAndIncrement(ctx sdk.Context) uint64 {
	proposalID := k.GetNextProposalID(ctx)
	k.SetNextProposalID(ctx, proposalID+1)
	return proposalID
}

func (k Keeper) GetNextProposalID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(NextProposalIDPrefix)
	if bz == nil {
		return 1
	}

	proposalID := BytesToProposalID(bz)
	return proposalID
}

func (k Keeper) SetNextProposalID(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(NextProposalIDPrefix, ProposalIDToBytes(proposalID))
}

func (k Keeper) CreateAndSaveProposalWithContent(ctx sdk.Context, title, description string, content types.Content) (uint64, error) {
	blockTime := ctx.BlockTime()
	proposalID := k.GetNextProposalIDAndIncrement(ctx)

	properties := k.GetNetworkProperties(ctx)

	proposal, err := types.NewProposal(
		proposalID,
		title,
		description,
		content,
		blockTime,
		blockTime.Add(time.Second*time.Duration(properties.ProposalEndTime)),
		blockTime.Add(time.Second*time.Duration(properties.ProposalEndTime)+
			time.Second*time.Duration(properties.ProposalEnactmentTime),
		),
		ctx.BlockHeight()+int64(properties.MinProposalEndBlocks),
		ctx.BlockHeight()+int64(properties.MinProposalEndBlocks+properties.MinProposalEnactmentBlocks),
	)

	if err != nil {
		return proposalID, err
	}

	k.SaveProposal(ctx, proposal)
	k.AddToActiveProposals(ctx, proposal)

	return proposalID, nil
}

func (k Keeper) SaveProposal(ctx sdk.Context, proposal types.Proposal) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&proposal)
	store.Set(GetProposalKey(proposal.ProposalId), bz)
}

func (k Keeper) GetProposal(ctx sdk.Context, proposalID uint64) (types.Proposal, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetProposalKey(proposalID))
	if bz == nil {
		return types.Proposal{}, false
	}

	var prop types.Proposal
	k.cdc.MustUnmarshal(bz, &prop)

	return prop, true
}

func (k Keeper) GetProposals(ctx sdk.Context) ([]types.Proposal, error) {
	proposals := []types.Proposal{}
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), ProposalsPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var proposal types.Proposal
		bz := iterator.Value()
		k.cdc.MustUnmarshal(bz, &proposal)
		proposals = append(proposals, proposal)
	}

	return proposals, nil
}

func (k Keeper) SaveVote(ctx sdk.Context, vote types.Vote) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&vote)
	store.Set(VoteKey(vote.ProposalId, vote.Voter), bz)
}

func (k Keeper) GetVote(ctx sdk.Context, proposalID uint64, address sdk.AccAddress) (types.Vote, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(VoteKey(proposalID, address))
	if bz == nil {
		return types.Vote{}, false
	}

	var vote types.Vote
	k.cdc.MustUnmarshal(bz, &vote)

	return vote, true
}

func (k Keeper) GetVotes(ctx sdk.Context) []types.Vote {
	votes := []types.Vote{}
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), VotesPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var vote types.Vote
		bz := iterator.Value()
		k.cdc.MustUnmarshal(bz, &vote)
		votes = append(votes, vote)
	}

	return votes
}

func (k Keeper) GetProposalVotesIterator(ctx sdk.Context, proposalID uint64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, VotesKey(proposalID))
}

func (k Keeper) GetProposalVotes(ctx sdk.Context, proposalID uint64) types.Votes {
	var votes types.Votes

	iterator := k.GetProposalVotesIterator(ctx, proposalID)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var vote types.Vote
		k.cdc.MustUnmarshal(iterator.Value(), &vote)
		votes = append(votes, vote)
	}

	return votes
}

func (k Keeper) AddToActiveProposals(ctx sdk.Context, proposal types.Proposal) {
	store := ctx.KVStore(k.storeKey)
	store.Set(ActiveProposalKey(proposal), ProposalIDToBytes(proposal.ProposalId))
}

func (k Keeper) RemoveActiveProposal(ctx sdk.Context, proposal types.Proposal) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(ActiveProposalKey(proposal))
}

func (k Keeper) AddToEnactmentProposals(ctx sdk.Context, proposal types.Proposal) {
	store := ctx.KVStore(k.storeKey)
	store.Set(EnactmentProposalKey(proposal), ProposalIDToBytes(proposal.ProposalId))
}

func (k Keeper) RemoveEnactmentProposal(ctx sdk.Context, proposal types.Proposal) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(EnactmentProposalKey(proposal))
}

// GetActiveProposalsWithFinishedVotingEndTimeIterator returns the proposals that have endtime finished.
func (k Keeper) GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(ActiveProposalsPrefix, sdk.PrefixEndBytes(ActiveProposalByTimeKey(endTime)))
}

// GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator returns the proposals that have finished the enactment time.
func (k Keeper) GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(EnactmentProposalsPrefix, sdk.PrefixEndBytes(EnactmentProposalByTimeKey(endTime)))
}

func ActiveProposalByTimeKey(endTime time.Time) []byte {
	return append(ActiveProposalsPrefix, sdk.FormatTimeBytes(endTime)...)
}

func EnactmentProposalByTimeKey(endTime time.Time) []byte {
	return append(EnactmentProposalsPrefix, sdk.FormatTimeBytes(endTime)...)
}

func ActiveProposalKey(prop types.Proposal) []byte {
	return append(ActiveProposalByTimeKey(prop.VotingEndTime), ProposalIDToBytes(prop.ProposalId)...)
}

func EnactmentProposalKey(prop types.Proposal) []byte {
	return append(EnactmentProposalByTimeKey(prop.EnactmentEndTime), ProposalIDToBytes(prop.ProposalId)...)
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
