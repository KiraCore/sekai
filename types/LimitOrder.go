package types

import (
	"fmt"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"strings"
)

type LimitOrder struct {
	ID          string `json:"id"`
	Index       uint32 `json:"index"`
	OrderBookID string `json:"order_book_id"`
	OrderType   uint8  `json:"order_type"`
	Amount      int64  `json:"amount"`
	LimitPrice  int64  `json:"limit_price"`
	ExpiryTime  int64  `json:"curator"`
	IsCancelled bool   `json:"is_cancelled"`
	Curator 	sdk.AccAddress `json:"curator"`
}

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
		Curator: nil,
	}
}

func (o LimitOrder) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ID: %s, OrderBookID: %s, OrderType: %d, Amount: %d, LimitPrice: %d, ExpiryTime: %d, Curator: %d`,
		o.ID, o.OrderBookID, o.OrderType, o.Amount, o.LimitPrice, o.ExpiryTime, o.Curator))
}
