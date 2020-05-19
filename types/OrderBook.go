package types

import (
	"fmt"
	"strings"

	sdk "github.com/KiraCore/cosmos-sdk/types"
)

type OrderBook struct {
	ID string              `json:"id"`
	Index string           `json:"index"`
	Base string		       `json:"base"`
	Quote string		   `json:"quote"`
	Mnemonic string 	   `json:"mnemonic"`
	Curator sdk.AccAddress `json:"curator"`
}

func NewOrderBook() OrderBook {
	return OrderBook{
		Index: "",
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