package types

import (
	"regexp"
)

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

func ValidateSpendingPoolName(name string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z][_0-9a-zA-Z]*$`)
	return regex.MatchString(name)
}
