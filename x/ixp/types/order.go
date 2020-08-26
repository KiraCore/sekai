package types

// NewLimitOrder generates new LimitOrder.
func NewLimitOrder() LimitOrder {
	return LimitOrder{
		ID:          "",
		Index:       0,
		OrderBookID: "",
		OrderType:   0,
		Amount:      0,
		LimitPrice:  0,
		ExpiryTime:  0,
		IsCancelled: false,
		Curator:     nil,
	}
}

// NewOrderBook generates new OrderBook.
func NewOrderBook() OrderBook {
	return OrderBook{
		Index:    0,
		Base:     "",
		Quote:    "",
		Mnemonic: "",
		Curator:  nil,
	}
}
