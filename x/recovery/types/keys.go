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
	RecoveryRecordKeyPrefix = []byte{0x01} // Prefix for recovery record
	RecoveryTokenKeyPrefix  = []byte{0x02} // Prefix for recovery token
)

func RecoveryRecordKey(address string) []byte {
	return append(RecoveryRecordKeyPrefix, address...)
}

func RecoveryTokenKey(address string) []byte {
	return append(RecoveryTokenKeyPrefix, address...)
}
