package keeper

import (
	"github.com/KiraCore/sekai/x/gov/types"
	types2 "github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SaveNetworkActor(ctx sdk.Context, actor types.NetworkActor) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), NetworkActorsPrefix)

	bz := k.cdc.MustMarshalBinaryBare(&actor)
	prefixStore.Set(actor.Address.Bytes(), bz)
}

func (k Keeper) GetNetworkActorByAddress(ctx sdk.Context, address sdk.AccAddress) (types.NetworkActor, bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), NetworkActorsPrefix)

	bz := prefixStore.Get(address.Bytes())
	if bz == nil {
		return types.NetworkActor{}, false
	}

	var na types.NetworkActor
	k.cdc.MustUnmarshalBinaryBare(bz, &na)

	return na, true
}

// AddWhitelistPermission checks if the actor exists, if not it fails.
func (k Keeper) AddWhitelistPermission(ctx sdk.Context, address sdk.AccAddress, perm types.PermValue) error {
	actor, found := k.GetNetworkActorByAddress(ctx, address)
	if !found {
		return types2.ErrNetworkActorNotFound
	}

	err := actor.Permissions.AddToWhitelist(perm)
	if err != nil {
		return err
	}

	k.SaveNetworkActor(ctx, actor)

	return nil
}
