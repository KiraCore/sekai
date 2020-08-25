package keeper

import (
	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	cdc      codec.BinaryMarshaler
	storeKey sdk.StoreKey
}

func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryMarshaler) Keeper {
	return Keeper{cdc: cdc, storeKey: storeKey}
}

func (k Keeper) SetPermissionsForRole(ctx sdk.Context, role types.Role, permissions *types.Permissions) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPermissionsRegistry)

	prefixStore.Set(types.RoleToKey(role), k.cdc.MustMarshalBinaryBare(permissions))
}

func (k Keeper) GetPermissionsForRole(ctx sdk.Context, councilor types.Role) *types.Permissions {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPermissionsRegistry)
	bz := prefixStore.Get(types.RoleToKey(councilor))

	perm := new(types.Permissions)
	k.cdc.MustUnmarshalBinaryBare(bz, perm)

	return perm
}

func (k Keeper) SetNetworkActor(ctx sdk.Context, actor types.NetworkActor) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixActors)

	bz := k.cdc.MustMarshalBinaryBare(&actor)
	prefixStore.Set(actor.Address.Bytes(), bz)
}

func (k Keeper) GetNetworkActorByAddress(ctx sdk.Context, address sdk.AccAddress) types.NetworkActor {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixActors)

	bz := prefixStore.Get(address.Bytes())

	var na types.NetworkActor
	k.cdc.MustUnmarshalBinaryBare(bz, &na)

	return na
}
