package signerkey

import (
	"github.com/KiraCore/cosmos-sdk/codec"
)

// RegisterCodec register codec for this message type
func RegisterCodec(codec *codec.Codec) {
	codec.RegisterConcrete(Message{}, "kiraHub/createOrder", nil)
}

// PackageCodec returns new codec
var PackageCodec = codec.New()

func init() {
	RegisterCodec(PackageCodec)
}
