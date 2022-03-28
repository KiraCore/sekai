package types

// constants
var (
	ModuleName = "distributor"
	// RouterKey to be used for routing msgs
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	PrefixKeyFeesCollected = []byte("fees_collected")
	PrefixKeyFeesTreasury  = []byte("fees_treasury")
	PrefixKeySnapPeriod    = []byte("snap_period")
)
