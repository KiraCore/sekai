package keeper

import (
	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: implement
func (k Keeper) GetNextRoleId(ctx sdk.Context) uint64 {
	return 1
}

// TODO: implement
func (k Keeper) SetRole(ctx sdk.Context, role types.Role) {
	perms := types.NewPermissions(nil, nil)
	k.savePermissionsForRole(ctx, uint64(role.Id), perms)
}

// TODO: implement
func (k Keeper) GetRole(ctx sdk.Context, roleId uint64) types.Role {
	return types.Role{}
}

// TODO: implement
func (k Keeper) GetRoleBySid(ctx sdk.Context, sId string) (types.Role, bool) {
	return types.Role{}, true
}

// TODO: implement
func (k Keeper) CreateRole(ctx sdk.Context, sid, description string) uint64 {
	perms := types.NewPermissions(nil, nil)
	k.savePermissionsForRole(ctx, uint64(role.Id), perms)
	return 0
}

// savePermissionsForRole adds permissions to role in the  permission Registry.
func (k Keeper) savePermissionsForRole(ctx sdk.Context, role uint64, permissions *types.Permissions) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), RolePermissionRegistry)
	prefixStore.Set(roleToBytes(role), k.cdc.MustMarshal(permissions))
}

func (k Keeper) GetAllRoles(ctx sdk.Context) []types.Role {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), RolePermissionRegistry)
	iterator := sdk.KVStorePrefixIterator(prefixStore, nil)
	defer iterator.Close()

	roles := []types.Role{}
	for ; iterator.Valid(); iterator.Next() {
		// role := bytesToRole(iterator.Key())
		role := types.Role{}
		// TODO: unmarshal from iterator.Value()
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
	store.Set(prefixWhitelistRole(perm, role), roleToBytes(role))

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
