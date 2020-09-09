package keeper

import (
	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetPermissionsForRole adds permissions to role in the  permission Registry.
func (k Keeper) SetPermissionsForRole(ctx sdk.Context, role types.Role, permissions *types.Permissions) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), KeyPrefixPermissionsRegistry)

	prefixStore.Set(types.RoleToKey(role), k.cdc.MustMarshalBinaryBare(permissions))
}

// GetPermissionsForRole returns the permissions assigned to the specific role.
func (k Keeper) GetPermissionsForRole(ctx sdk.Context, role types.Role) *types.Permissions {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), KeyPrefixPermissionsRegistry)
	bz := prefixStore.Get(types.RoleToKey(role))

	perm := new(types.Permissions)
	k.cdc.MustUnmarshalBinaryBare(bz, perm)

	return perm
}
