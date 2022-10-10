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

func (k Keeper) SaveCouncilor(ctx sdk.Context, councilor types.Councilor) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), CouncilorIdentityRegistryPrefix)

	bz := k.cdc.MustMarshal(&councilor)

	councilorKey := GetCouncilorKey(councilor.Address)

	prefixStore.Set(councilorKey, bz)
}

func (k Keeper) GetCouncilor(ctx sdk.Context, address sdk.AccAddress) (types.Councilor, bool) {
	return k.getCouncilorByKey(ctx, GetCouncilorKey(address))
}

func (k Keeper) GetCouncilorByMoniker(ctx sdk.Context, moniker string) (types.Councilor, bool) {
	addresses := k.GetAddressesByIdRecordKey(ctx, "moniker", moniker)

	if len(addresses) != 1 {
		return types.Councilor{}, false
	}

	return k.GetCouncilor(ctx, addresses[0])
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

func (k Keeper) GetAllCouncilors(ctx sdk.Context) []types.Councilor {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), CouncilorIdentityRegistryPrefix)

	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	councilors := []types.Councilor{}
	for ; iterator.Valid(); iterator.Next() {
		var councilor types.Councilor
		k.cdc.MustUnmarshal(iterator.Value(), &councilor)
		councilors = append(councilors, councilor)
	}

	return councilors
}

func (k Keeper) ResetWholeCouncilorRank(ctx sdk.Context) {
	councilors := k.GetAllCouncilors(ctx)
	for _, councilor := range councilors {
		councilor.Status = types.CouncilorActive
		councilor.Rank = 0
		k.SaveCouncilor(ctx, councilor)
	}
}
