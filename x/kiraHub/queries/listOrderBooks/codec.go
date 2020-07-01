package listOrderBooks

import (
	"github.com/KiraCore/cosmos-sdk/codec"
)


func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(QueryListOrderBooks{}, "kiraHub/queryOrderBooks", nil)
	cdc.RegisterConcrete(QueryListOrderBooksByTP{}, "kiraHub/queryOrderBooksByTP", nil)
}

var packageCodec = codec.New()

func init() {
	RegisterCodec(packageCodec)
	packageCodec.Seal()
}