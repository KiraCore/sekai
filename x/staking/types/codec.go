package types

import (
	functionmeta "github.com/KiraCore/sekai/function_meta"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterCodec register codec and metadata
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgClaimValidator{}, "kiraHub/MsgClaimValidator", nil)

	functionmeta.AddNewFunction((&MsgClaimValidator{}).Type(), `{
		"description": "MsgClaimValidator defines a message for claiming a new validator.",
		"parameters": {
			"moniker": {
				"type":        "string",
				"description": "validator's name or nickname."
			},
			"val_key": {
				"type":        "val_address",
				"description": "validator operator address"
			},
			"pub_key": {
				"type":        "string",
				"description": "validator bech32 public key"
			}
		}
	}`)
}

// RegisterInterfaces register Msg and structs
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgClaimValidator{},
	)

	registry.RegisterInterface(
		"kira.gov.Content",
		(*govtypes.Content)(nil),
		&ProposalUnjailValidator{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
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
