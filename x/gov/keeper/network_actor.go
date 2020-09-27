package keeper

import (
	"fmt"

	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SaveNetworkActor(ctx sdk.Context, actor types.NetworkActor) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), NetworkActorsPrefix)

	bz := k.cdc.MustMarshalBinaryBare(&actor)
	prefixStore.Set(actor.Address.Bytes(), bz)
}

func (k Keeper) GetNetworkActorByAddress(ctx sdk.Context, address sdk.AccAddress) (types.NetworkActor, error) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), NetworkActorsPrefix)

	bz := prefixStore.Get(address.Bytes())
	if bz == nil {
		return types.NetworkActor{}, fmt.Errorf("network actor not found")
	}

	var na types.NetworkActor
	k.cdc.MustUnmarshalBinaryBare(bz, &na)

	return na, nil
}
