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
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), KeyPrefixCouncilorIdentityRegistry)

	bz := prefixStore.Get(GetCouncilorKey(address))
	if bz == nil {
		return types.Councilor{}, fmt.Errorf("councilor not found")
	}

	var co types.Councilor
	k.cdc.MustUnmarshalBinaryBare(bz, &co)

	return co, nil
}

func (k Keeper) GetCouncilorByMoniker(ctx sdk.Context, moniker string) (types.Councilor, error) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), KeyPrefixCouncilorIdentityRegistry)

	councilorKey := prefixStore.Get(GetCouncilorByMonikerKey(moniker))

	bz := prefixStore.Get(councilorKey)
	var councilor types.Councilor

	k.cdc.MustUnmarshalBinaryBare(bz, &councilor)

	return councilor, nil
}
