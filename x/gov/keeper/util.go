package keeper

import (
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CheckIfAllowedPermission
func CheckIfAllowedPermission(ctx sdk.Context, keeper Keeper, addr sdk.AccAddress, permValue types.PermValue) bool {
	actor, found := keeper.GetNetworkActorByAddress(ctx, addr)
	if !found {
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
		rolePerms, found := keeper.GetPermissionsForRole(ctx, types.Role(role))
		if found {
			roles[role] = &rolePerms
		}
	}

	return roles
}

// ProposalIDToBytes returns the byte representation of the proposalID
func ProposalIDToBytes(proposalID uint64) []byte {
	proposalIDBz := sdk.Uint64ToBigEndian(proposalID)
	return proposalIDBz
}

// BytesToProposalID returns proposalID in uint64 format from a byte array
func BytesToProposalID(bz []byte) uint64 {
	return sdk.BigEndianToUint64(bz)
}

func BoolToInt(v bool) uint64 {
	if v {
		return 1
	}
	return 0
}

func IntToBool(v uint64) bool {
	if v != 0 {
		return true
	}
	return false
}
