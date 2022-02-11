package types

// constants
const (
	ModuleName = "spending"
	// RouterKey to be used for routing msgs
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	KeyPrefixSpendingPool = "spending"
	KeyPrefixClaimInfo    = "claim_info"
)

func SpendingPoolKey(name string) []byte {
	return append([]byte(KeyPrefixSpendingPool), name...)
}

func ClaimInfoKey(name string, address string) []byte {
	return append(append([]byte(KeyPrefixClaimInfo), name...), address...)
}
