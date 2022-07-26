package types

var (
	ModuleName   = "custody"
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
	StoreKey     = ModuleName

	PrefixKeyCustodyRecord       = "custody_record_prefix_"
	PrefixKeyCustodyCustodians   = "custody_custodians_prefix_"
	PrefixKeyCustodyWhiteList    = "custody_white_list_prefix_"
	PrefixKeyCustodyLimits       = "custody_limits_prefix_"
	PrefixKeyCustodyLimitsStatus = "custody_limits_status_prefix_"
	CustodyBufferSizeKey         = []byte("custody_buffer_size")
	CustodyTxSizeKey             = []byte("custody_tx_size")
)
