package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

// Keys for Staking store.
// 0x00<ValAddress> : The validator
// 0x01<moniker_bytes> : The Key of the Validator.
// 0x02<ValAddress> : The Validator
// 0x03<ValAddress> : The Validator
// 0x04<ValAddress> : The Validator Address
// 0x05<ValAddress> : The Validator Address
// 0x06<ValAddress> : Validator Jail Info
var (
	ValidatorsKey              = []byte{0x00} // Validators key prefix.
	ValidatorsByConsAddressKey = []byte{0x02} // Validators by consensus addres (PubKey).
	PendingValidatorQueue      = []byte{0x03} // Validators that are pending to join into the end blocker.
	RemovingValidatorQueue     = []byte{0x04} // Validators that are pending to be removed from the validator set.
	ReactivatingValidatorQueue = []byte{0x05} // Validators that are pending to be reactivated in the set.
	ValidatorJailInfo          = []byte{0x06} // Validator Jail Info (JailTime, etc)
	LastValidatorPowerKey      = []byte{0x07}
)

// GetValidatorKey gets the key for the validator with address
func GetValidatorKey(operatorAddr sdk.ValAddress) []byte {
	return append(ValidatorsKey, operatorAddr.Bytes()...)
}

func GetValidatorKeyAcc(address sdk.AccAddress) []byte {
	return append(ValidatorsKey, address.Bytes()...)
}

func GetValidatorByConsAddrKey(addr sdk.ConsAddress) []byte {
	return append(ValidatorsByConsAddressKey, addr.Bytes()...)
}

func GetPendingValidatorKey(operatorAddress sdk.ValAddress) []byte {
	return append(PendingValidatorQueue, operatorAddress.Bytes()...)
}

func GetRemovingValidatorKey(operatorAddress sdk.ValAddress) []byte {
	return append(RemovingValidatorQueue, operatorAddress.Bytes()...)
}

func GetReactivatingValidatorKey(operatorAddress sdk.ValAddress) []byte {
	return append(ReactivatingValidatorQueue, operatorAddress.Bytes()...)
}

func GetValidatorJailInfoKey(operatorAddress sdk.ValAddress) []byte {
	return append(ValidatorJailInfo, operatorAddress.Bytes()...)
}
