package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/KiraCore/sekai/x/gov/types"
)

func (k Keeper) CreateRole(ctx sdk.Context, role types.Role) {
	perms := types.NewPermissions(nil, nil)
	k.SavePermissionsForRole(ctx, role, perms)
}

// SavePermissionsForRole adds permissions to role in the  permission Registry.
func (k Keeper) SavePermissionsForRole(ctx sdk.Context, role types.Role, permissions *types.Permissions) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), RolePermissionRegistry)
	prefixStore.Set(roleToBytes(role), k.cdc.MustMarshalBinaryBare(permissions))
}

// GetPermissionsForRole returns the permissions assigned to the specific role.
func (k Keeper) GetPermissionsForRole(ctx sdk.Context, role types.Role) (types.Permissions, bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), RolePermissionRegistry)
	bz := prefixStore.Get(roleToBytes(role))
	if bz == nil {
		return types.Permissions{}, false
	}

	var perm types.Permissions
	k.cdc.MustUnmarshalBinaryBare(bz, &perm)

	return perm, true
}

func (k Keeper) WhitelistRolePermission(ctx sdk.Context, role types.Role, perm types.PermValue) error {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), RolePermissionRegistry)
	bz := prefixStore.Get(roleToBytes(role))
	if bz == nil {
		return types.ErrRoleDoesNotExist
	}

	var perms types.Permissions
	k.cdc.MustUnmarshalBinaryBare(bz, &perms)

	err := perms.AddToWhitelist(perm)
	if err != nil {
		return errors.Wrap(types.ErrWhitelisting, err.Error())
	}

	k.SavePermissionsForRole(ctx, role, &perms)
	return nil
}

func (k Keeper) CheckIfAllowedPermission(ctx sdk.Context, addr sdk.AccAddress, permValue types.PermValue) bool {
	return CheckIfAllowedPermission(ctx, k, addr, permValue)
}

func prefixWhitelist(role types.Role) []byte {
	return append(WhitelistRolePrefix, roleToBytes(role)...)
}
