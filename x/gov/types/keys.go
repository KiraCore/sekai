package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const ModuleName = "customgov"

const WhitelistPermissions = "whitelist-permissions"
const BlacklistPermissions = "blacklist-permissions"

// constants
var (
	KeyPrefixPermissionsRegistry = []byte("permissions_registry")
	KeyPrefixActors              = []byte("network_actors")
	KeyPrefixNetworkProperties   = []byte("network_properties")

	// Roles
	RoleUndefined Role = 0x0
	RoleCouncilor Role = 0x1
	RoleValidator Role = 0x2
	RoleGovLeader Role = 0x3

	PermClaimValidator      PermValue = 1
	PermClaimGovernanceSeat PermValue = 2
	PermAddPermissions      PermValue = 3
)

// Role represents a Role in the registry.
type Role uint64
type Roles []uint64

// PermValue represents a single permission value, like claim-role-validator.
type PermValue uint32

// RoleToKey returns bytes to be used as a key for a given capability index.
func RoleToKey(index Role) []byte {
	return sdk.Uint64ToBigEndian(uint64(index))
}

// RoleFromKey returns an index from a call to IndexToKey for a given capability
// index.
func RoleFromKey(key []byte) Role {
	return Role(sdk.BigEndianToUint64(key))
}
