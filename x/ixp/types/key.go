package types

// keys
const (
	// ModuleName is the name of the module
	ModuleName = "ixp"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	QuerierRoute     = ModuleName
	TransactionRoute = ModuleName

	CreateOrderBookTransaction = "createorderbook"
	CreateOrderTransaction     = "createorder"
	CancelOrderTransaction     = "cancelorder"
	UpsertSignerKeyTransaction = "upsertsignerkey"

	ListOrderBooksQuery              = "listOrderBooks"
	ListOrderBooksQueryByTradingPair = "listOrderBooksByTradingPair"
	ListOrders                       = "listOrders"
)
