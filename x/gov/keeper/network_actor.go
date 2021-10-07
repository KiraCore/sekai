package keeper

import (
	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SaveNetworkActor(ctx sdk.Context, actor types.NetworkActor) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), NetworkActorsPrefix)

	bz := k.cdc.MustMarshal(&actor)
	prefixStore.Set(actor.Address.Bytes(), bz)
}

func (k Keeper) GetNetworkActorByAddress(ctx sdk.Context, address sdk.AccAddress) (types.NetworkActor, bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), NetworkActorsPrefix)

	bz := prefixStore.Get(address.Bytes())
	if bz == nil {
		return types.NetworkActor{}, false
	}

	var na types.NetworkActor
	k.cdc.MustUnmarshal(bz, &na)

	return na, true
}

func (k Keeper) GetNetworkActorFromIterator(iterator sdk.Iterator) *types.NetworkActor {
	bz := iterator.Value()
	if bz == nil {
		return nil
	}

	var na types.NetworkActor
	k.cdc.MustUnmarshal(bz, &na)

	return &na
}

// AddWhitelistPermission whitelist a permission to an address. It saves the actor after it.
func (k Keeper) AddWhitelistPermission(ctx sdk.Context, actor types.NetworkActor, perm types.PermValue) error {
	err := actor.Permissions.AddToWhitelist(perm)
	if err != nil {
		return err
	}

	k.SaveNetworkActor(ctx, actor)

	store := ctx.KVStore(k.storeKey)
	store.Set(WhitelistAddressPermKey(actor.Address, perm), actor.Address.Bytes())

	return nil
}

// RemoveWhitelistPermission removes the whitelisted permission that an address has
func (k Keeper) RemoveWhitelistPermission(ctx sdk.Context, actor types.NetworkActor, perm types.PermValue) error {
	err := actor.Permissions.RemoveFromWhitelist(perm)
	if err != nil {
		return err
	}

	k.SaveNetworkActor(ctx, actor)

	store := ctx.KVStore(k.storeKey)
	store.Delete(WhitelistAddressPermKey(actor.Address, perm))

	return nil
}

func (k Keeper) AssignRoleToActor(ctx sdk.Context, actor types.NetworkActor, role types.Role) {
	actor.SetRole(role)
	k.SaveNetworkActor(ctx, actor)

	store := ctx.KVStore(k.storeKey)
	store.Set(roleAddressKey(role, actor.Address), actor.Address.Bytes())
}

func (k Keeper) RemoveRoleFromActor(ctx sdk.Context, actor types.NetworkActor, role types.Role) {
	actor.RemoveRole(role)
	k.SaveNetworkActor(ctx, actor)

	store := ctx.KVStore(k.storeKey)
	store.Delete(roleAddressKey(role, actor.Address))
}

// GetNetworkActorsIterator returns all the actors iterator
func (k Keeper) GetNetworkActorsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, NetworkActorsPrefix)
}

// GetNetworkActorsByWhitelistedPermission returns all the actors that have Perm in individual whitelist.
func (k Keeper) GetNetworkActorsByWhitelistedPermission(ctx sdk.Context, perm types.PermValue) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, WhitelistPermKey(perm))
}

// GetNetworkActorsByRole returns all network actors that have role assigned.
func (k Keeper) GetNetworkActorsByRole(ctx sdk.Context, role types.Role) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, roleKey(role))
}

// GetNetworkActorsByAbsoluteWhitelistPermission returns all actors that have a specific whitelist permission,
// it does not matter if it is by role or by individual permission.
func (k Keeper) GetNetworkActorsByAbsoluteWhitelistPermission(ctx sdk.Context, perm types.PermValue) []types.NetworkActor {
	duplicateMap := map[string]bool{}

	var actors []types.NetworkActor
	iterator := k.GetNetworkActorsByWhitelistedPermission(ctx, perm)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		if _, ok := duplicateMap[sdk.AccAddress(iterator.Value()).String()]; !ok {
			duplicateMap[sdk.AccAddress(iterator.Value()).String()] = true
			actors = append(actors, k.getNetworkActorOrFail(ctx, iterator.Value()))
		}
	}

	rolesIter := k.GetRolesByWhitelistedPerm(ctx, perm)
	defer rolesIter.Close()

	for ; rolesIter.Valid(); rolesIter.Next() {
		actorIter := k.GetNetworkActorsByRole(ctx, bytesToRole(rolesIter.Value()))

		for ; actorIter.Valid(); actorIter.Next() {
			if _, ok := duplicateMap[sdk.AccAddress(actorIter.Value()).String()]; !ok {
				duplicateMap[sdk.AccAddress(actorIter.Value()).String()] = true
				actors = append(actors, k.getNetworkActorOrFail(ctx, actorIter.Value()))
			}
		}
	}

	return actors
}

func (k Keeper) getNetworkActorOrFail(ctx sdk.Context, addr sdk.AccAddress) types.NetworkActor {
	actor, found := k.GetNetworkActorByAddress(ctx, addr)
	if !found {
		panic("expected network actor not found")
	}

	return actor
}

// WhitelistAddressPermKey returns the prefix key in format <0x31 + Perm_Bytes + address_bytes>
func WhitelistAddressPermKey(address sdk.AccAddress, perm types.PermValue) []byte {
	return append(WhitelistPermKey(perm), address.Bytes()...)
}

// WhitelistPermKey returns the prefix key in format <0x31 + Perm_Bytes>
func WhitelistPermKey(perm types.PermValue) []byte {
	return append(WhitelistActorPrefix, permToBytes(perm)...)
}

// roleAddressKey returns the prefix key in format <0x33 + Role_Bytes + address_bytes>
func roleAddressKey(role types.Role, address sdk.AccAddress) []byte {
	return append(roleKey(role), address.Bytes()...)
}

// roleKey returns a prefix key in format <0x32 + Role_Bytes>
func roleKey(role types.Role) []byte {
	return append(RoleActorPrefix, roleToBytes(role)...)
}

// permToBytes returns a PermValue in bytes representation.
func permToBytes(perm types.PermValue) []byte {
	return sdk.Uint64ToBigEndian(uint64(perm))
}

// roleToBytes returns a Role in bytes representation.
func roleToBytes(role types.Role) []byte {
	return sdk.Uint64ToBigEndian(uint64(role))
}

// bytesToRole converts byte representation of a role to Role type.
func bytesToRole(bz []byte) types.Role {
	return types.Role(sdk.BigEndianToUint64(bz))
}
