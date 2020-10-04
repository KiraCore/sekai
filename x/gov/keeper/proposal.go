package keeper

import (
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

	proposalID := GetProposalIDFromBytes(bz)

	return proposalID, nil
}

func (k Keeper) SaveProposalID(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)

	store.Set(NextProposalIDPrefix, GetProposalIDBytes(proposalID))
}

func (k Keeper) SaveProposal(ctx sdk.Context, proposal types.ProposalAssignPermission) (uint64, error) {
	store := ctx.KVStore(k.storeKey)

	proposalID, err := k.GetNextProposalID(ctx)
	if err != nil {
		return 0, err
	}

	bz := k.cdc.MustMarshalBinaryBare(&proposal)
	store.Set(GetProposalKey(proposalID), bz)

	// Update NextProposal
	k.SaveProposalID(ctx, proposalID+1)

	return proposalID, nil
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

func VotesKey(proposalID uint64) []byte {
	return append(VotesPrefix, GetProposalIDBytes(proposalID)...)
}

func VoteKey(proposalId uint64, address sdk.AccAddress) []byte {
	return append(VotesKey(proposalId), address.Bytes()...)
}

func GetProposalKey(proposalID uint64) []byte {
	return append(ProposalsPrefix, GetProposalIDBytes(proposalID)...)
}
