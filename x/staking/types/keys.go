package types

import sdk "github.com/KiraCore/cosmos-sdk/types"

const (
	// ModuleName is the name of the custom staking
	ModuleName = "customstaking"

	ClaimValidator = "claim-validator"
)

var ValidatorsKey = []byte{0x21} // Validators key

// GetValidatorKey gets the key for the validator with address
func GetValidatorKey(operatorAddr sdk.ValAddress) []byte {
	return append(ValidatorsKey, operatorAddr.Bytes()...)
}
