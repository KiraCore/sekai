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
	cdc.RegisterConcrete(&MsgUpsertTokenInfo{}, "kiraHub/MsgUpsertTokenInfo", nil)
	functionmeta.AddNewFunction((&MsgUpsertTokenInfo{}).Type(), `{
		"description": "MsgUpsertTokenInfo represents a message to register token rate.",
		"parameters": {
			"denom": {
				"type":        "string",
				"description": "denomination target."
			},
			"rate": {
				"type":        "float",
				"description": "Exchange rate in terms of KEX token. e.g. 0.1, 10.5"
			},
			"fee_payments": {
				"type":        "bool",
				"description": "defining if it is enabled or disabled as fee payment method."
			},
			"proposer": {
				"type":        "address",
				"description": "proposer who propose this message."
			}
		}
	}`)

	cdc.RegisterConcrete(&MsgEthereumTx{}, "kiraHub/MsgEthereumTx", nil)
	functionmeta.AddNewFunction((&MsgEthereumTx{}).Type(), `{
		"description": "MsgUpsertTokenInfo represents a message to register token rate.",
		"parameters": {
			"tx_type": {
				"type":        "string",
				"description": "ethereum tx type."
			},
			"sender": {
				"type":        "string",
				"description": "bech32 encoded address of ethereum tx"
			},
			"hash": {
				"type":        "string",
				"description": "ethereum tx hash in hex."
			},
			"data": {
				"type":        "bytes",
				"description": "rlp encoding of ethereum tx bytes."
			}
		}
	}`)
}

// RegisterInterfaces register Msg and structs
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpsertTokenInfo{},
		&MsgEthereumTx{},
	)

	registry.RegisterInterface(
		"kira.gov.Content",
		(*govtypes.Content)(nil),
		&ProposalUpsertTokenInfo{},
		&ProposalTokensWhiteBlackChange{},
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
