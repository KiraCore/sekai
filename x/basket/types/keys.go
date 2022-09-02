package types

// constants
var (
	ModuleName = "basket"
	// RouterKey to be used for routing msgs
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	KeyLastBasketId = []byte("last_basket_id")
	PrefixBasketKey = []byte("basket")
)
