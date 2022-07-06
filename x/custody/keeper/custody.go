package keeper

import (
	"github.com/KiraCore/sekai/x/custody/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetCustodyInfoByAddress(ctx sdk.Context, address sdk.AccAddress) *types.CustodySettings {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PrefixKeyCustodyRecord))
	bz := prefixStore.Get(address)

	if bz == nil {
		return nil
	}

	info := new(types.CustodySettings)
	k.cdc.MustUnmarshal(bz, info)

	return info
}

func (k Keeper) SetCustodyRecord(ctx sdk.Context, record types.CustodyRecord) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyCustodyRecord), record.Address...)

	store.Set(key, k.cdc.MustMarshal(&record.CustodySettings))
}

func (k Keeper) GetCustodyWhiteListByAddress(ctx sdk.Context, address sdk.AccAddress) *types.CustodyWhiteList {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PrefixKeyCustodyWhiteList))
	bz := prefixStore.Get(address)

	if bz == nil {
		return nil
	}

	info := new(types.CustodyWhiteList)
	k.cdc.MustUnmarshal(bz, info)

	return info
}

func (k Keeper) DropCustodyWhiteListByAddress(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyCustodyWhiteList), address...)

	store.Delete(key)
}

func (k Keeper) AddToCustodyWhiteList(ctx sdk.Context, record types.CustodyWhiteListRecord) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyCustodyWhiteList), record.Address...)

	store.Set(key, k.cdc.MustMarshal(record.CustodyWhiteList))
}
