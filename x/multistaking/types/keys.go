package types

import "strings"

const (
	// ModuleName is the name of the multistaking
	ModuleName = "multistaking"

	// QuerierRoute is the querier route for the multistaking module
	QuerierRoute = ModuleName
)

var (
	KeyPrefixStakingPool   = []byte{0x1}
	KeyPrefixUndelegation  = []byte{0x2}
	KeyLastPoolId          = []byte{0x3}
	KeyLastUndelegationId  = []byte{0x4}
	KeyPrefixPoolDelegator = []byte{0x5}
	KeyPrefixRewards       = []byte{0x6}
	KeyPrefixCompoundInfo  = []byte{0x7}
)

func GetOriginalDenom(denom string) string {
	if denom == "" || denom[0] != 'v' {
		return denom
	}
	split := strings.Split(denom, "/")
	return strings.TrimPrefix(denom, split[0]+"/")
}
