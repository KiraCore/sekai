package types

const (
	// ModuleName is the name of the module
	ModuleName = "kiraHub"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

  // RouterKey to be used for routing msgs
  RouterKey = ModuleName

	QuerierRoute = ModuleName
	TransactionRoute = ModuleName

	CreateOrderBookTransaction = "createorderbook"
	CreateOrderTransaction = "createorder"

	ListOrderBooksQuery = "listOrderBooks"
	ListOrderBooksQueryByTP = "listOrderBooksByTP"
)
