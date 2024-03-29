package types

import (
	"encoding/binary"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "customslashing"

	// StoreKey is the store key string for slashing
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute is the querier route for slashing
	QuerierRoute = ModuleName
)

// Keys for slashing store
// Items are stored with the following key: values
//
// - 0x01<consAddress_Bytes>: ValidatorSigningInfo
//
// - 0x02<consAddress_Bytes><period_Bytes>: bool
//
// - 0x03<accAddr_Bytes>: crypto.PubKey
var (
	ValidatorSigningInfoKeyPrefix         = []byte{0x01} // Prefix for signing info
	ValidatorMissedBlockBitArrayKeyPrefix = []byte{0x02} // Prefix for missed block bit array
	AddrPubkeyRelationKeyPrefix           = []byte{0x03} // Prefix for address-pubkey relation
	SlashedValidatorsByTimeKeyPrefix      = []byte{0x04} // Prefix for slashed validators
)

// ValidatorSigningInfoKey - stored by *Consensus* address (not operator address)
func ValidatorSigningInfoKey(v sdk.ConsAddress) []byte {
	return append(ValidatorSigningInfoKeyPrefix, v.Bytes()...)
}

// ValidatorMissedBlockBitArrayPrefixKey - stored by *Consensus* address (not operator address)
func ValidatorMissedBlockBitArrayPrefixKey(v sdk.ConsAddress) []byte {
	return append(ValidatorMissedBlockBitArrayKeyPrefix, v.Bytes()...)
}

// ValidatorMissedBlockBitArrayKey - stored by *Consensus* address (not operator address)
func ValidatorMissedBlockBitArrayKey(v sdk.ConsAddress, i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return append(ValidatorMissedBlockBitArrayPrefixKey(v), b...)
}

// AddrPubkeyRelationKey gets pubkey relation key used to get the pubkey from the address
func AddrPubkeyRelationKey(address []byte) []byte {
	return append(AddrPubkeyRelationKeyPrefix, address...)
}

func SlashedValidatorByTimeKey(timestamp time.Time, val sdk.ValAddress) []byte {
	return append(SlashedValidatorByTimePrefixKey(timestamp), val...)
}

func SlashedValidatorByTimePrefixKey(timestamp time.Time) []byte {
	return append(SlashedValidatorsByTimeKeyPrefix, sdk.FormatTimeBytes(timestamp)...)
}
