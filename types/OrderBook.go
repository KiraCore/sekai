package types

import (
	"fmt"
	"strings"

	sdk "github.com/KiraCore/cosmos-sdk/types"
)

type OrderBook struct {
	Index int32            `json:"index"`
	Base sdk.Coins		   `json:"base"`
	Quote sdk.Coins		   `json:"quote"`
	Mnemonic string 	   `json:"mnemonic"`
	Curator sdk.AccAddress `json:"curator"`
}

func NewOrderBook() OrderBook {
	return OrderBook{
		Index: nil,
		Base: nil,
		Quote: nil,
		Mnemonic: "",
		Curator: nil,
	}
}

func (o OrderBook) String() string {
	if o.Mnemonic == "" {
		return strings.TrimSpace(fmt.Sprintf(`Index: %s, Base: %s, Quote: %s, Curator: %s`, o.Index, o.Base, o.Quote, o.Curator))
	}
	return strings.TrimSpace(fmt.Sprintf(`Index: %s, Base: %s, Quote: %s, Mnemonic: %s, Curator: %s`, o.Index, o.Base, o.Quote, o.Mnemonic, o.Curator))
}