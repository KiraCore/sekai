package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

// Keys for Staking store.
// 0x00<ValAddress> : The validator
// 0x01<moniker_bytes> : The Key of the Validator.
// 0x02<ValAddress> : The Validator
var (
	ValidatorsKey          = []byte{0x00} // Validators key prefix.
	ValidatorsByMonikerKey = []byte{0x01} // Validators by moniker prefix.
	PendingValidatorQueue  = []byte{0x02} // Validators that are pending to join into the end blocker.
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

func GetPendingValidatorKey(operatorAddress sdk.ValAddress) []byte {
	return append(PendingValidatorQueue, operatorAddress.Bytes()...)
}
