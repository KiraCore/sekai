package types

// constants
const (
	ModuleName = "custody"
	// RouterKey to be used for routing msgs
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
	StoreKey     = ModuleName

	PrefixKeyCustodyRecord    = "custody_record_prefix_"
	PrefixKeyCustodyWhiteList = "custody_white_list_prefix_"
)
