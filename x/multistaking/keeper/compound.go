package keeper

import (
	"github.com/KiraCore/sekai/x/multistaking/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetCompoundInfoByAddress(ctx sdk.Context, addr string) types.CompoundInfo {
	compoundInfo := types.CompoundInfo{}
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.KeyPrefixCompoundInfo), []byte(addr)...)
	bz := store.Get(key)
	if bz == nil {
		return compoundInfo
	}
	k.cdc.MustUnmarshal(bz, &compoundInfo)
	return compoundInfo
}

func (k Keeper) GetAllCompoundInfo(ctx sdk.Context) []types.CompoundInfo {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCompoundInfo)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	compoundInfos := []types.CompoundInfo{}
	for ; iterator.Valid(); iterator.Next() {
		info := types.CompoundInfo{}
		k.cdc.MustUnmarshal(iterator.Value(), &info)
		compoundInfos = append(compoundInfos, info)
	}
	return compoundInfos
}

func (k Keeper) SetCompoundInfo(ctx sdk.Context, info types.CompoundInfo) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixCompoundInfo, []byte(info.Delegator)...)
	store.Set(key, k.cdc.MustMarshal(&info))
}

func (k Keeper) RemoveCompoundInfo(ctx sdk.Context, info types.CompoundInfo) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixCompoundInfo, []byte(info.Delegator)...)
	store.Delete(key)
}
