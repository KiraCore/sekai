package createOrder

import (
	"github.com/KiraCore/cosmos-sdk/codec"
)

func RegisterCodec(codec *codec.Codec) {
	codec.RegisterConcrete(Message{}, "kiraHub/createOrderBook", nil)
}

var PackageCodec = codec.New()

func init() {
	RegisterCodec(PackageCodec)
}