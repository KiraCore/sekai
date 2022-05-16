package types

const (
	// ModuleName is the name of the multistaking
	ModuleName = "multistaking"

	// QuerierRoute is the querier route for the multistaking module
	QuerierRoute = ModuleName
)

var (
	KeyPrefixStakingPool  = []byte{0x1}
	KeyPrefixUndelegation = []byte{0x2}
	KeyLastPoolId         = []byte{0x3}
	KeyLastUndelegationId = []byte{0x4}
)
