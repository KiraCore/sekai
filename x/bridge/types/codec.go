package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterCodec register codec and metadata
func RegisterCodec(cdc *codec.LegacyAmino) {
	//cdc.RegisterConcrete(&MsgCreateCustodyRecord{}, "kiraHub/MsgCreateCustodyRecord", nil)
}

// RegisterInterfaces register Msg and structs
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgChangeCosmosEthereum{},
		&MsgChangeEthereumCosmos{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterCodec(amino)
	amino.Seal()
}
