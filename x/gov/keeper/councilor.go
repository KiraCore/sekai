package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/gov/types"
)

var (
	CouncilorsKey          = []byte{0x21} // Councilors key prefix.
	CouncilorsByMonikerKey = []byte{0x22} // Councilors by moniker prefix.
)

func GetCouncilorKey(address sdk.AccAddress) []byte {
	return append(CouncilorsKey, address.Bytes()...)
}

func GetCouncilorByMonikerKey(moniker string) []byte {
	return append(CouncilorsByMonikerKey, []byte(moniker)...)
}

func (k Keeper) SaveCouncilor(ctx sdk.Context, councilor types.Councilor) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), CouncilorIdentityRegistryPrefix)

	bz := k.cdc.MustMarshal(&councilor)

	councilorKey := GetCouncilorKey(councilor.Address)

	prefixStore.Set(councilorKey, bz)
	prefixStore.Set(GetCouncilorByMonikerKey(councilor.Moniker), councilorKey)
}

func (k Keeper) GetCouncilor(ctx sdk.Context, address sdk.AccAddress) (types.Councilor, bool) {
	return k.getCouncilorByKey(ctx, GetCouncilorKey(address))
}

func (k Keeper) GetCouncilorByMoniker(ctx sdk.Context, moniker string) (types.Councilor, bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), CouncilorIdentityRegistryPrefix)

	councilorKey := prefixStore.Get(GetCouncilorByMonikerKey(moniker))
	if councilorKey == nil {
		return types.Councilor{}, false
	}

	return k.getCouncilorByKey(ctx, councilorKey)
}

func (k Keeper) getCouncilorByKey(ctx sdk.Context, key []byte) (types.Councilor, bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), CouncilorIdentityRegistryPrefix)

	bz := prefixStore.Get(key)
	if bz == nil {
		return types.Councilor{}, false
	}

	var councilor types.Councilor
	k.cdc.MustUnmarshal(bz, &councilor)

	return councilor, true
}
