package types

// constants
const (
	ModuleName = "custody"
	// RouterKey to be used for routing msgs
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
	StoreKey     = ModuleName

	PrefixKeyCustodyRecord       = "custody_record_prefix_"
	PrefixKeyCustodyWhiteList    = "custody_white_list_prefix_"
	PrefixKeyCustodyLimits       = "custody_limits_prefix_"
	PrefixKeyCustodyLimitsStatus = "custody_limits_status_prefix_"
)
