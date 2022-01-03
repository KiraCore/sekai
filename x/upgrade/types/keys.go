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
	KeyCurrentPlan = []byte{0x01}
	KeyNextPlan    = []byte{0x02}
)
