package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (k Keeper) GetProposalID(ctx sdk.Context) (uint64, error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(KeyProposalKeyIDPrefix)
	if bz == nil {
		return 0, errors.Wrap(types.ErrInvalidGenesis, "initial proposal ID hasn't been set")
	}

	proposalID := getProposalIDFromBytes(bz)

	return proposalID, nil
}

func (k Keeper) SaveProposalID(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)

	store.Set(KeyProposalKeyIDPrefix, getProposalIDBytes(proposalID))
}
