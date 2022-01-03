package keeper

import (
	"strconv"

	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
)

func (k Keeper) GetNextRoleId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(NextRolePrefix)
	if bz == nil {
		return 1
	}
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetNextRoleId(ctx sdk.Context, nextRoleId uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(NextRolePrefix, sdk.Uint64ToBigEndian(nextRoleId))
}

func (k Keeper) SetRole(ctx sdk.Context, role types.Role) {
	bz, err := proto.Marshal(&role)
	if err != nil {
		panic(err)
	}
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), RoleIdToInfo)
	prefixStore.Set(roleToBytes(uint64(role.Id)), bz)

	prefixStore = prefix.NewStore(ctx.KVStore(k.storeKey), RoleSidToIdRegistry)
	prefixStore.Set([]byte(role.Sid), sdk.Uint64ToBigEndian(uint64(role.Id)))

	// set empty permissions
	perms := types.NewPermissions(nil, nil)
	k.savePermissionsForRole(ctx, uint64(role.Id), perms)
}

func (k Keeper) GetRole(ctx sdk.Context, roleId uint64) (types.Role, error) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), RoleIdToInfo)
	bz := prefixStore.Get(roleToBytes(uint64(roleId)))
	if bz == nil {
		return types.Role{}, types.ErrRoleDoesNotExist
	}
	role := types.Role{}
	err := proto.Unmarshal(bz, &role)
	return role, err
}

func (k Keeper) GetRoleBySid(ctx sdk.Context, sId string) (types.Role, error) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), RoleSidToIdRegistry)
	bz := prefixStore.Get([]byte(sId))
	if bz == nil {
		return types.Role{}, types.ErrRoleDoesNotExist
	}
	return k.GetRole(ctx, sdk.BigEndianToUint64(bz))
}

func (k Keeper) CreateRole(ctx sdk.Context, sid, description string) uint64 {
	newRoleId := k.GetNextRoleId(ctx)
	k.SetRole(ctx, types.Role{
		Id:          uint32(newRoleId),
		Sid:         sid,
		Description: description,
	})
	k.SetNextRoleId(ctx, newRoleId+1)
	return newRoleId
}

// savePermissionsForRole adds permissions to role in the  permission Registry.
func (k Keeper) savePermissionsForRole(ctx sdk.Context, role uint64, permissions *types.Permissions) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), RolePermissionRegistry)
	prefixStore.Set(roleToBytes(role), k.cdc.MustMarshal(permissions))
}

func (k Keeper) GetAllRoles(ctx sdk.Context) []types.Role {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), RoleIdToInfo)
	iterator := sdk.KVStorePrefixIterator(prefixStore, nil)
	defer iterator.Close()

	roles := []types.Role{}
	for ; iterator.Valid(); iterator.Next() {
		role := types.Role{}
		err := proto.Unmarshal(iterator.Value(), &role)
		if err != nil {
			panic(err)
		}
		roles = append(roles, role)
	}
	return roles
}

// GetPermissionsForRole returns the permissions assigned to the specific role.
func (k Keeper) GetPermissionsForRole(ctx sdk.Context, role uint64) (types.Permissions, bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), RolePermissionRegistry)
	bz := prefixStore.Get(roleToBytes(role))
	if bz == nil {
		return types.Permissions{}, false
	}

	var perm types.Permissions
	k.cdc.MustUnmarshal(bz, &perm)

	return perm, true
}

func (k Keeper) GetRoleIdFromIdentifierString(ctx sdk.Context, identifier string) (uint64, error) {
	if roleId, err := strconv.Atoi(identifier); err == nil {
		return uint64(roleId), nil
	}
	role, err := k.GetRoleBySid(ctx, identifier) // sid
	if err != nil {
		return 0, err
	}
	return uint64(role.Id), nil
}

func (k Keeper) SetWhiltelistPermRoleKey(ctx sdk.Context, role uint64, perm types.PermValue) {
	store := ctx.KVStore(k.storeKey)
	store.Set(prefixWhitelistRole(perm, role), roleToBytes(role))
}

func (k Keeper) WhitelistRolePermission(ctx sdk.Context, role uint64, perm types.PermValue) error {
	store := ctx.KVStore(k.storeKey)

	prefixStore := prefix.NewStore(store, RolePermissionRegistry)
	bz := prefixStore.Get(roleToBytes(role))
	if bz == nil {
		return types.ErrRoleDoesNotExist
	}

	var perms types.Permissions
	k.cdc.MustUnmarshal(bz, &perms)

	err := perms.AddToWhitelist(perm)
	if err != nil {
		return err
	}

	k.savePermissionsForRole(ctx, role, &perms)
	k.SetWhiltelistPermRoleKey(ctx, role, perm)

	return nil
}

func (k Keeper) BlacklistRolePermission(ctx sdk.Context, role uint64, perm types.PermValue) error {
	store := ctx.KVStore(k.storeKey)

	prefixStore := prefix.NewStore(store, RolePermissionRegistry)
	bz := prefixStore.Get(roleToBytes(role))
	if bz == nil {
		return types.ErrRoleDoesNotExist
	}

	var perms types.Permissions
	k.cdc.MustUnmarshal(bz, &perms)

	err := perms.AddToBlacklist(perm)
	if err != nil {
		return err
	}

	k.savePermissionsForRole(ctx, role, &perms)

	return nil
}

func (k Keeper) RemoveWhitelistRolePermission(ctx sdk.Context, role uint64, perm types.PermValue) error {
	store := ctx.KVStore(k.storeKey)

	prefixStore := prefix.NewStore(store, RolePermissionRegistry)
	bz := prefixStore.Get(roleToBytes(role))
	if bz == nil {
		return types.ErrRoleDoesNotExist
	}

	var perms types.Permissions
	k.cdc.MustUnmarshal(bz, &perms)

	err := perms.RemoveFromWhitelist(perm)
	if err != nil {
		return err
	}

	k.savePermissionsForRole(ctx, role, &perms)
	store.Delete(prefixWhitelistRole(perm, role))

	return nil
}

func (k Keeper) RemoveBlacklistRolePermission(ctx sdk.Context, role uint64, perm types.PermValue) error {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, RolePermissionRegistry)

	bz := prefixStore.Get(roleToBytes(role))
	if bz == nil {
		return types.ErrRoleDoesNotExist
	}

	var perms types.Permissions
	k.cdc.MustUnmarshal(bz, &perms)

	err := perms.RemoveFromBlacklist(perm)
	if err != nil {
		return err
	}

	k.savePermissionsForRole(ctx, role, &perms)

	return nil
}

func (k Keeper) IterateRoles(ctx sdk.Context) sdk.Iterator {
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), RolePermissionRegistry)
}

func (k Keeper) GetPermissionsFromIterator(iterator sdk.Iterator) types.Permissions {
	bz := iterator.Value()
	if bz == nil {
		return types.Permissions{}
	}

	var perms types.Permissions
	k.cdc.MustUnmarshal(bz, &perms)
	return perms
}

func (k Keeper) GetRolesByWhitelistedPerm(ctx sdk.Context, perm types.PermValue) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, prefixWhitelist(perm))
}

func (k Keeper) CheckIfAllowedPermission(ctx sdk.Context, addr sdk.AccAddress, permValue types.PermValue) bool {
	return CheckIfAllowedPermission(ctx, k, addr, permValue)
}

func prefixWhitelist(perm types.PermValue) []byte {
	return append(WhitelistRolePrefix, permToBytes(perm)...)
}

func prefixWhitelistRole(perm types.PermValue, role uint64) []byte {
	return append(prefixWhitelist(perm), roleToBytes(role)...)
}
