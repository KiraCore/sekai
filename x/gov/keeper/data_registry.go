package keeper

import (
	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UpsertDataRegistryEntry insert/updates a data registry entry.
func (k Keeper) UpsertDataRegistryEntry(ctx sdk.Context, key string, entry types.DataRegistryEntry) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), DataRegistryPrefix)

	bz := k.cdc.MustMarshalBinaryBare(&entry)
	prefixStore.Set([]byte(key), bz)
}

// GetDataRegistryEntry returns the Entry of the registry by key.
func (k Keeper) GetDataRegistryEntry(ctx sdk.Context, key string) (types.DataRegistryEntry, bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), DataRegistryPrefix)

	bz := prefixStore.Get([]byte(key))
	if bz == nil {
		return types.DataRegistryEntry{}, false
	}

	var na types.DataRegistryEntry
	k.cdc.MustUnmarshalBinaryBare(bz, &na)

	return na, true
}
