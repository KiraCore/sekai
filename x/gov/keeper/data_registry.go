package keeper

import (
	"strings"

	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UpsertDataRegistryEntry insert/updates a data registry entry.
func (k Keeper) UpsertDataRegistryEntry(ctx sdk.Context, key string, entry types.DataRegistryEntry) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), DataRegistryPrefix)

	bz := k.cdc.MustMarshal(&entry)
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
	k.cdc.MustUnmarshal(bz, &na)

	return na, true
}

// ListDataRegistryEntryKeys returns all keys of data registry
func (k Keeper) ListDataRegistryEntryKeys(ctx sdk.Context) []string {
	var keys []string

	// get iterator for token aliases
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, DataRegistryPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := strings.TrimPrefix(string(iterator.Key()), string(DataRegistryPrefix))
		keys = append(keys, key)
	}

	return keys
}

// AllDataRegistry returns all of data registry
func (k Keeper) AllDataRegistry(ctx sdk.Context) map[string]*types.DataRegistryEntry {
	var dataRegistry map[string]*types.DataRegistryEntry

	keys := k.ListDataRegistryEntryKeys(ctx)

	// get iterator for token aliases
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, DataRegistryPrefix)
	defer iterator.Close()

	for _, key := range keys {
		entry, ok := k.GetDataRegistryEntry(ctx, key)
		if ok {
			dataRegistry[key] = &entry
		} else {
			dataRegistry[key] = nil
		}
	}

	return dataRegistry
}
