package types

import (
	functionmeta "github.com/KiraCore/sekai/function_meta"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers concrete types on LegacyAmino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgActivate{}, "cosmos-sdk/MsgActivate", nil)
	functionmeta.AddNewFunction((&MsgActivate{}).Type(), `{
		"description": "MsgActivate defines a message to activate an inactive validator.",
		"parameters": {
			"validator_addr": {
				"type":        "string",
				"description": "bech32 format of validator address. e.g. kiravaloper1ewgq8gtsefakhal687t8hnsw5zl4y8eksup39w"
			}
		}
	}`)
	cdc.RegisterConcrete(&MsgPause{}, "cosmos-sdk/MsgPause", nil)
	functionmeta.AddNewFunction((&MsgPause{}).Type(), `{
		"description": "MsgPause defines a message to pause an active validator.",
		"parameters": {
			"validator_addr": {
				"type":        "string",
				"description": "bech32 format of validator address. e.g. kiravaloper1ewgq8gtsefakhal687t8hnsw5zl4y8eksup39w"
			}
		}
	}`)
	cdc.RegisterConcrete(&MsgUnpause{}, "cosmos-sdk/MsgUnpause", nil)
	functionmeta.AddNewFunction((&MsgUnpause{}).Type(), `{
		"description": "MsgUnpause defines a message to unpause a paused validator.",
		"parameters": {
			"validator_addr": {
				"type":        "string",
				"description": "bech32 format of validator address. e.g. kiravaloper1ewgq8gtsefakhal687t8hnsw5zl4y8eksup39w"
			}
		}
	}`)
}

// RegisterInterfaces register interfaces on registry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgActivate{},
		&MsgPause{},
		&MsgUnpause{},
	)

	registry.RegisterInterface(
		"kira.gov.Content",
		(*govtypes.Content)(nil),
		&ProposalResetWholeValidatorRank{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/slashing module codec. Note, the codec
	// should ONLY be used in certain instances of tests and for JSON encoding as Amino
	// is still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/slashing and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
