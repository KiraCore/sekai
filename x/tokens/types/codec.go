package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/staking module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterCodec(amino)
	amino.Seal()
}
