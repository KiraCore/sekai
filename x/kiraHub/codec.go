package kiraHub

import (
	"github.com/KiraCore/cosmos-sdk/codec"

	"github.com/KiraCore/sekai/x/kiraHub/transactions/createOrderBook"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	createOrderBook.RegisterCodec(cdc)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
