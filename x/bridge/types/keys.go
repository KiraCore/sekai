package types

var (
	ModuleName = "bridge"
	RouterKey  = ModuleName
	QueryRoute = ModuleName
	StoreKey   = ModuleName

	PrefixKeyBridgeCosmosEthereumRecord = "bridge_cosmos_ethereum_record_prefix_"
	PrefixKeyBridgeEthereumCosmosRecord = "bridge_ethereum_cosmos_record_prefix_"

	BridgeAddressKey                    = []byte("bridge_address")
	BridgeCosmosEthereumExchangeRateKey = []byte("bridge_cosmos_ethereum_rate")
	BridgeEthereumCosmosExchangeRateKey = []byte("bridge_ethereum_cosmos_rate")
)
