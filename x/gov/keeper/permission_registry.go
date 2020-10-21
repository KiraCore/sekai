package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/gov/types"
)

// SetPermissionsForRole adds permissions to role in the  permission Registry.
func (k Keeper) SetPermissionsForRole(ctx sdk.Context, role types.Role, permissions *types.Permissions) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), RolePermissionRegistry)

	prefixStore.Set(types.RoleToKey(role), k.cdc.MustMarshalBinaryBare(permissions))
}

// GetPermissionsForRole returns the permissions assigned to the specific role.
func (k Keeper) GetPermissionsForRole(ctx sdk.Context, role types.Role) (types.Permissions, bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), RolePermissionRegistry)
	bz := prefixStore.Get(types.RoleToKey(role))
	if bz == nil {
		return types.Permissions{}, false
	}

	var perm types.Permissions
	k.cdc.MustUnmarshalBinaryBare(bz, &perm)

	return perm, true
}

func (k Keeper) CheckIfAllowedPermission(ctx sdk.Context, addr sdk.AccAddress, permValue types.PermValue) bool {
	return CheckIfAllowedPermission(ctx, k, addr, permValue)
}
