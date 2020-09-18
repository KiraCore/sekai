package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/gov/types"
)

// SetPermissionsForRole adds permissions to role in the  permission Registry.
func (k Keeper) SetPermissionsForRole(ctx sdk.Context, role types.Role, permissions *types.Permissions) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), KeyPrefixPermissionsRegistry)

	prefixStore.Set(types.RoleToKey(role), k.cdc.MustMarshalBinaryBare(permissions))
}

// GetPermissionsForRole returns the permissions assigned to the specific role.
func (k Keeper) GetPermissionsForRole(ctx sdk.Context, role types.Role) (*types.Permissions, error) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), KeyPrefixPermissionsRegistry)
	bz := prefixStore.Get(types.RoleToKey(role))
	if bz == nil {
		return nil, fmt.Errorf("role not found")
	}

	perm := new(types.Permissions)
	k.cdc.MustUnmarshalBinaryBare(bz, perm)

	return perm, nil
}
