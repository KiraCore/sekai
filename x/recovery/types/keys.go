package types

const (
	// ModuleName is the name of the module
	ModuleName = "recovery"

	// StoreKey is the store key string for recovery
	StoreKey = ModuleName

	// RouterKey is the message route for recovery
	RouterKey = ModuleName

	// QuerierRoute is the querier route for recovery
	QuerierRoute = ModuleName
)

var (
	RecoveryChallengeKeyPrefix    = []byte{0x01} // Prefix for recovery challenge
	RecoveryRecordKeyPrefix       = []byte{0x02} // Prefix for recovery record
	RecoveryTokenByDenomKeyPrefix = []byte{0x03} // Prefix for recovery token by denom
	RecoveryTokenKeyPrefix        = []byte{0x04} // Prefix for recovery token
	RotationHistoryKeyPrefix      = []byte{0x05} // Prefix for rotation history
	KeyPrefixRewards              = []byte{0x06}
	KeyPrefixRRTokenHolder        = []byte{0x07}
)

func RecoveryChallengeKey(challenge string) []byte {
	return append(RecoveryChallengeKeyPrefix, challenge...)
}

func RecoveryRecordKey(address string) []byte {
	return append(RecoveryRecordKeyPrefix, address...)
}

func RecoveryTokenKey(address string) []byte {
	return append(RecoveryTokenKeyPrefix, address...)
}

func RecoveryTokenByDenomKey(denom string) []byte {
	return append(RecoveryTokenByDenomKeyPrefix, denom...)
}

func RotationHistoryKey(address string) []byte {
	return append(RotationHistoryKeyPrefix, address...)
}
