package types

// constants
const (
	ModuleName = "customgov"
	// RouterKey to be used for routing msgs
	RouterKey = ModuleName
	// QuerierRoute is the querier route for the staking module
	QuerierRoute = ModuleName
)

// constants
var (
	KeyPrefixNetworkProperties = []byte("network_properties")
	KeyPrefixExecutionFee      = []byte("execution_fee")

	// Roles
	RoleUndefined Role = 0x0
	RoleSudo      Role = 0x1
	RoleValidator Role = 0x2
)

// Role represents a Role in the registry.
type Role uint64
type Roles []uint64
