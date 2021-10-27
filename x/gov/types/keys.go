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
	KeyLastIdentityRecordId                   = []byte("last_identity_record_id")
	KeyLastIdRecordVerifyRequestId            = []byte("last_identity_record_verify_request_id")
	KeyPrefixIdentityRecord                   = []byte("identity_record_prefix")
	KeyPrefixIdentityRecordByAddress          = []byte("identity_record_by_address_prefix")
	KeyPrefixIdRecordVerifyRequest            = []byte("identity_record_verify_request_prefix")
	KeyPrefixIdRecordVerifyRequestByRequester = []byte("identity_record_verify_request_by_requester_prefix")
	KeyPrefixIdRecordVerifyRequestByApprover  = []byte("identity_record_verify_request_by_approver_prefix")

	// Roles
	RoleUndefined Role = 0x0
	RoleSudo      Role = 0x1
	RoleValidator Role = 0x2
)

// Role represents a Role in the registry.
type Role uint64
type Roles []uint64

func IdentityRecordByAddressPrefix(address string) []byte {
	return append(KeyPrefixIdentityRecordByAddress, address...)
}

func IdRecordVerifyRequestByRequesterPrefix(address string) []byte {
	return append(KeyPrefixIdRecordVerifyRequestByRequester, address...)
}

func IdRecordVerifyRequestByApproverPrefix(address string) []byte {
	return append(KeyPrefixIdRecordVerifyRequestByApprover, address...)
}
