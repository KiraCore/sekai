package listOrders

import (
	"github.com/KiraCore/cosmos-sdk/codec"
)


func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(QueryListOrders{}, "kiraHub/query", nil)
}

var packageCodec = codec.New()

func init() {
	RegisterCodec(packageCodec)
	packageCodec.Seal()
}
