package types

// constants
const (
	ModuleName = "upgrade"
	// RouterKey to be used for routing msgs
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
	StoreKey     = ModuleName
)

var (
	KeyUpgradePlan = []byte{0x01}
)
