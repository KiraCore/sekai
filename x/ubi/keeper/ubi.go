package keeper

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/ubi/types"
)

func (k Keeper) GetUBIRecordByName(ctx sdk.Context, name string) *types.UBIRecord {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PrefixKeyUBIRecord))
	bz := prefixStore.Get([]byte(name))
	if bz == nil {
		return nil
	}

	rate := new(types.UBIRecord)
	k.cdc.MustUnmarshal(bz, rate)

	return rate
}

func (k Keeper) GetUBIRecords(ctx sdk.Context) []types.UBIRecord {
	var records []types.UBIRecord

	// get iterator for token rates
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.PrefixKeyUBIRecord))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		name := strings.TrimPrefix(string(iterator.Key()), string(types.PrefixKeyUBIRecord))
		record := k.GetUBIRecordByName(ctx, name)
		if record != nil {
			records = append(records, *record)
		}
	}
	return records
}

func (k Keeper) SetUBIRecord(ctx sdk.Context, record types.UBIRecord) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyUBIRecord), []byte(record.Name)...)
	store.Set(key, k.cdc.MustMarshal(&record))
}

func (k Keeper) DeleteUBIRecord(ctx sdk.Context, name string) error {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyUBIRecord), []byte(name)...)
	if !store.Has(key) {
		return fmt.Errorf("ubi record does not exist: %s", name)
	}

	store.Delete(key)
	return nil
}
