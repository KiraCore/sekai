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

	// identity registrar
	KeyLastIdentityRecordId        = []byte("last_identity_record_id")
	KeyLastIdRecordVerifyRequestId = []byte("last_identity_record_verify_request_id")
	KeyPrefixIdentityRecord        = []byte("identity_record")
	KeyPrefixIdRecordVerifyRequest = []byte("identity_record_verify_request")

	// Roles
	RoleUndefined Role = 0x0
	RoleSudo      Role = 0x1
	RoleValidator Role = 0x2
)

// Role represents a Role in the registry.
type Role uint64
type Roles []uint64
