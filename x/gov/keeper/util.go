package keeper

import (
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CheckIfAllowedPermission
func CheckIfAllowedPermission(ctx sdk.Context, keeper Keeper, addr sdk.AccAddress, permValue types.PermValue) bool {
	actor, err := keeper.GetNetworkActorByAddress(ctx, addr)
	if err != nil {
		return false
	}

	// Get All Roles for actor
	roles := getRolePermissions(ctx, keeper, actor)

	permMap := map[uint32]bool{}

	// Add whitelist perms into list
	// First roles
	for _, rolePerm := range roles {
		for _, rp := range rolePerm.Whitelist {
			permMap[rp] = true
		}
	}

	// Second individual roles
	for _, ap := range actor.Permissions.Whitelist {
		permMap[ap] = true
	}

	// Remove All Blacklisted perms
	for _, rolePerm := range roles {
		for _, rp := range rolePerm.Blacklist {
			permMap[rp] = false
		}
	}

	// Remove personal Permissions
	for _, pp := range actor.Permissions.Blacklist {
		permMap[pp] = false
	}

	isAllowed, ok := permMap[uint32(permValue)]
	if !ok {
		return false
	}

	return isAllowed
}

func getRolePermissions(ctx sdk.Context, keeper Keeper, actor types.NetworkActor) map[uint64]*types.Permissions {
	roles := map[uint64]*types.Permissions{}
	for _, role := range actor.Roles {
		roles[role], _ = keeper.GetPermissionsForRole(ctx, types.Role(role)) // TODO take care of roles.
	}

	return roles
}
