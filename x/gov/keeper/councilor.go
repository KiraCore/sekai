package keeper

import (
	"fmt"

	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), KeyPrefixCouncilorIdentityRegistry)

	bz := k.cdc.MustMarshalBinaryBare(&councilor)

	councilorKey := GetCouncilorKey(councilor.Address)

	prefixStore.Set(councilorKey, bz)
	prefixStore.Set(GetCouncilorByMonikerKey(councilor.Moniker), councilorKey)
}

func (k Keeper) GetCouncilor(ctx sdk.Context, address sdk.AccAddress) (types.Councilor, error) {
	return k.getCouncilorByKey(ctx, GetCouncilorKey(address))
}

func (k Keeper) GetCouncilorByMoniker(ctx sdk.Context, moniker string) (types.Councilor, error) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), KeyPrefixCouncilorIdentityRegistry)

	councilorKey := prefixStore.Get(GetCouncilorByMonikerKey(moniker))
	if councilorKey == nil {
		return types.Councilor{}, fmt.Errorf("councilor with moniker %s not found", moniker)
	}

	return k.getCouncilorByKey(ctx, councilorKey)
}

func (k Keeper) getCouncilorByKey(ctx sdk.Context, key []byte) (types.Councilor, error) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), KeyPrefixCouncilorIdentityRegistry)

	bz := prefixStore.Get(key)
	if bz == nil {
		return types.Councilor{}, fmt.Errorf("councilor not found")
	}

	var councilor types.Councilor
	k.cdc.MustUnmarshalBinaryBare(bz, &councilor)

	return councilor, nil
}
