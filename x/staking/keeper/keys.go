package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

// Keys for Staking store.
// 0x00<ValAddress> : The validator
// 0x01<moniker_bytes> : The Proposal
var (
	ValidatorsKey          = []byte{0x00} // Validators key prefix.
	ValidatorsByMonikerKey = []byte{0x01} // Validators by moniker prefix.
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
