package keeper

import (
	"github.com/KiraCore/sekai/x/layer2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Pow10(decimal uint64) sdk.Int {
	res := sdk.OneInt()
	for i := 0; i < int(decimal); i++ {
		res = res.Mul(sdk.NewInt(10))
	}
	return res
}

func (k Keeper) SetTokenInfo(ctx sdk.Context, info types.TokenInfo) {
	bz := k.cdc.MustMarshal(&info)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.TokenInfoKey(info.Denom), bz)
}

func (k Keeper) DeleteTokenInfo(ctx sdk.Context, denom string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.TokenInfoKey(denom))
}

func (k Keeper) GetTokenInfo(ctx sdk.Context, denom string) types.TokenInfo {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TokenInfoKey(denom))
	if bz == nil {
		return types.TokenInfo{}
	}

	info := types.TokenInfo{}
	k.cdc.MustUnmarshal(bz, &info)
	return info
}

func (k Keeper) GetTokenInfos(ctx sdk.Context) []types.TokenInfo {
	store := ctx.KVStore(k.storeKey)

	infos := []types.TokenInfo{}
	it := sdk.KVStorePrefixIterator(store, types.PrefixTokenInfoKey)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		info := types.TokenInfo{}
		k.cdc.MustUnmarshal(it.Value(), &info)
		infos = append(infos, info)
	}
	return infos
}
