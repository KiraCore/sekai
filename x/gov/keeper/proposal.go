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

	proposalID := getProposalIDFromBytes(bz)

	return proposalID, nil
}

func (k Keeper) SaveProposalID(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)

	store.Set(NextProposalIDPrefix, getProposalIDBytes(proposalID))
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

	return proposalID, err
}

func GetProposalKey(proposalID uint64) []byte {
	return append(ProposalsPrefix, getProposalIDBytes(proposalID)...)
}
