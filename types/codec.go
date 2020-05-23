package types

import (
	"github.com/KiraCore/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(OrderBook{}, "kiraHub/OrderBook", nil)
	cdc.RegisterConcrete(LimitOrder{}, "kiraHub/LimitOrder", nil)
}

var PackageCodec = codec.New()

func init() {
	RegisterCodec(PackageCodec)
}
