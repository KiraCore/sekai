package types

import (
	"fmt"
	"strings"
)

type OrderBook struct {
	ID string              `json:"id"`
	Index string           `json:"index"`
	Base string		       `json:"base"`
	Quote string		   `json:"quote"`
	Mnemonic string 	   `json:"mnemonic"`
	Curator string         `json:"curator"`
}

func NewOrderBook() OrderBook {
	return OrderBook{
		Index: "",
		Base: "",
		Quote: "",
		Mnemonic: "",
		Curator: "",
	}
}

func (o OrderBook) String() string {
	if o.Mnemonic == "" {
		return strings.TrimSpace(fmt.Sprintf(`Index: %s, Base: %s, Quote: %s, Curator: %s`, o.Index, o.Base, o.Quote, o.Curator))
	}
	return strings.TrimSpace(fmt.Sprintf(`Index: %s, Base: %s, Quote: %s, Mnemonic: %s, Curator: %s`, o.Index, o.Base, o.Quote, o.Mnemonic, o.Curator))
}