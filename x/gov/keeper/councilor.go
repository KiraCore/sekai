package keeper

import (
	"fmt"

	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SaveCouncilor(ctx sdk.Context, councilor types.Councilor) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), KeyPrefixCouncilorIdentityRegistry)

	bz := k.cdc.MustMarshalBinaryBare(&councilor)
	prefixStore.Set(councilor.Address.Bytes(), bz)
}

func (k Keeper) GetCouncilor(ctx sdk.Context, address sdk.AccAddress) (types.Councilor, error) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), KeyPrefixCouncilorIdentityRegistry)

	bz := prefixStore.Get(address.Bytes())
	if bz == nil {
		return types.Councilor{}, fmt.Errorf("councilor not found")
	}

	var co types.Councilor
	k.cdc.MustUnmarshalBinaryBare(bz, &co)

	return co, nil
}
