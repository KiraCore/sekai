package keeper

import (
	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) CreateRole(ctx sdk.Context, role types.Role) {
	perms := types.NewPermissions(nil, nil)
	k.savePermissionsForRole(ctx, role, perms)
}

// savePermissionsForRole adds permissions to role in the  permission Registry.
func (k Keeper) savePermissionsForRole(ctx sdk.Context, role types.Role, permissions *types.Permissions) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), RolePermissionRegistry)
	prefixStore.Set(roleToBytes(role), k.cdc.MustMarshal(permissions))
}

// GetPermissionsForRole returns the permissions assigned to the specific role.
func (k Keeper) GetPermissionsForRole(ctx sdk.Context, role types.Role) (types.Permissions, bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), RolePermissionRegistry)
	bz := prefixStore.Get(roleToBytes(role))
	if bz == nil {
		return types.Permissions{}, false
	}

	var perm types.Permissions
	k.cdc.MustUnmarshal(bz, &perm)

	return perm, true
}

func (k Keeper) WhitelistRolePermission(ctx sdk.Context, role types.Role, perm types.PermValue) error {
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

func (k Keeper) BlacklistRolePermission(ctx sdk.Context, role types.Role, perm types.PermValue) error {
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

func (k Keeper) RemoveWhitelistRolePermission(ctx sdk.Context, role types.Role, perm types.PermValue) error {
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

func (k Keeper) RemoveBlacklistRolePermission(ctx sdk.Context, role types.Role, perm types.PermValue) error {
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

func prefixWhitelistRole(perm types.PermValue, role types.Role) []byte {
	return append(prefixWhitelist(perm), roleToBytes(role)...)
}
