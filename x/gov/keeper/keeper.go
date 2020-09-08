package keeper

import (
	"fmt"

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

func (k Keeper) GetPermissionsForRole(ctx sdk.Context, role types.Role) *types.Permissions {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPermissionsRegistry)
	bz := prefixStore.Get(types.RoleToKey(role))

	perm := new(types.Permissions)
	k.cdc.MustUnmarshalBinaryBare(bz, perm)

	return perm
}

func (k Keeper) SaveNetworkActor(ctx sdk.Context, actor types.NetworkActor) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixActors)

	bz := k.cdc.MustMarshalBinaryBare(&actor)
	prefixStore.Set(actor.Address.Bytes(), bz)
}

func (k Keeper) GetNetworkActorByAddress(ctx sdk.Context, address sdk.AccAddress) (types.NetworkActor, error) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixActors)

	bz := prefixStore.Get(address.Bytes())
	if bz == nil {
		return types.NetworkActor{}, fmt.Errorf("network actor not found")
	}

	var na types.NetworkActor
	k.cdc.MustUnmarshalBinaryBare(bz, &na)

	return na, nil
}

func (k Keeper) SaveCouncilor(ctx sdk.Context, councilor types.Councilor) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCouncilorIdentityRegistry)

	bz := k.cdc.MustMarshalBinaryBare(&councilor)
	prefixStore.Set(councilor.Address.Bytes(), bz)
}

func (k Keeper) GetCouncilor(ctx sdk.Context, address sdk.AccAddress) (types.Councilor, error) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCouncilorIdentityRegistry)

	bz := prefixStore.Get(address.Bytes())
	if bz == nil {
		return types.Councilor{}, fmt.Errorf("councilor not found")
	}

	var co types.Councilor
	k.cdc.MustUnmarshalBinaryBare(bz, &co)

	return co, nil
}
