package types

import (
	"github.com/KiraCore/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(OrderBook{}, "kiraHub/OrderBook", nil)
}

var PackageCodec = codec.New()

func init() {
	RegisterCodec(PackageCodec)
}