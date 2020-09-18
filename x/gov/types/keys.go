package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const ModuleName = "customgov"

// RouterKey to be used for routing msgs
const RouterKey = ModuleName

// constants
var (
	KeyPrefixPermissionsRegistry = []byte("permissions_registry")
	KeyPrefixActors              = []byte("network_actors")
	KeyPrefixNetworkProperties   = []byte("network_properties")
	KeyPrefixExecutionFee        = []byte("execution_fee")

	// Roles
	RoleUndefined Role = 0x0
	RoleSudo      Role = 0x1
	RoleValidator Role = 0x2

	PermClaimValidator PermValue = 1
	PermSetPermissions PermValue = 2
	PermChangeTxFee    PermValue = 3
	PermClaimCouncilor PermValue = 3
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
