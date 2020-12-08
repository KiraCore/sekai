package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName is the name of the custom staking
	ModuleName = "customstaking"
	
	// QuerierRoute is the querier route for the staking module
	QuerierRoute = ModuleName
)

var (
	ValidatorsKey          = []byte{0x21} // Validators key prefix.
	ValidatorsByMonikerKey = []byte{0x22} // Validators by moniker prefix.
)

// GetValidatorKey gets the key for the validator with address
func GetValidatorKey(operatorAddr sdk.ValAddress) []byte {
	return append(ValidatorsKey, operatorAddr.Bytes()...)
}

func GetValidatorKeyAcc(address sdk.AccAddress) []byte {
	return append(ValidatorsKey, address.Bytes()...)
}

func GetValidatorByMonikerKey(moniker string) []byte {
	return append(ValidatorsByMonikerKey, []byte(moniker)...)
}
