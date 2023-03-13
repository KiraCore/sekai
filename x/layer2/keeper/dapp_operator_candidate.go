package keeper

import (
	"github.com/KiraCore/sekai/x/layer2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetDappOperatorCandidate(ctx sdk.Context, candidate types.DappOperatorCandidate) {
	bz := k.cdc.MustMarshal(&candidate)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DappOperatorCandidateKey(candidate.DappName, candidate.Candidate), bz)
}

func (k Keeper) DeleteDappOperatorCandidate(ctx sdk.Context, name, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.DappOperatorCandidateKey(name, address))
}

func (k Keeper) GetDappOperatorCandidate(ctx sdk.Context, name string, user string) types.DappOperatorCandidate {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DappOperatorCandidateKey(name, user))
	if bz == nil {
		return types.DappOperatorCandidate{}
	}

	candidate := types.DappOperatorCandidate{}
	k.cdc.MustUnmarshal(bz, &candidate)
	return candidate
}

func (k Keeper) GetDappOperatorCandidates(ctx sdk.Context, name string) []types.DappOperatorCandidate {
	store := ctx.KVStore(k.storeKey)

	candidates := []types.DappOperatorCandidate{}
	it := sdk.KVStorePrefixIterator(store, append([]byte(types.PrefixDappOperatorCandidateKey), name...))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		candidate := types.DappOperatorCandidate{}
		k.cdc.MustUnmarshal(it.Value(), &candidate)
		candidates = append(candidates, candidate)
	}
	return candidates
}

func (k Keeper) GetAllDappOperatorCandidates(ctx sdk.Context) []types.DappOperatorCandidate {
	store := ctx.KVStore(k.storeKey)

	candidates := []types.DappOperatorCandidate{}
	it := sdk.KVStorePrefixIterator(store, []byte(types.PrefixDappOperatorCandidateKey))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		candidate := types.DappOperatorCandidate{}
		k.cdc.MustUnmarshal(it.Value(), &candidate)
		candidates = append(candidates, candidate)
	}
	return candidates
}
