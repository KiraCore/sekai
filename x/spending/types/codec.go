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
	cdc.RegisterConcrete(&MsgUpsertTokenAlias{}, "kiraHub/MsgUpsertTokenAlias", nil)
	functionmeta.AddNewFunction((&MsgUpsertTokenAlias{}).Type(), `{
		"description": "MsgUpsertTokenAlias represents a message to register token alias.",
		"parameters": {
			"symbol": {
				"type":        "string",
				"description": "Ticker (eg. ATOM, KEX, BTC)."
			},
			"name": {
				"type":        "string",
				"description": "Token Name (e.g. Cosmos, Kira, Bitcoin)."
			},
			"icon": {
				"type":        "string",
				"description": "Graphical Symbol (url link to graphics)."
			},
			"decimals": {
				"type":        "uint32",
				"description": "Integer number of max decimals."
			},
			"denoms": {
				"type":        "array<string>",
				"description": "An array of token denoms to be aliased."
			},
			"proposer": {
				"type":        "string",
				"description": "proposer who propose this message."
			}
		}
	}`)

}

// RegisterInterfaces register Msg and structs
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpsertTokenRate{},
		&MsgUpsertTokenAlias{},
	)

	registry.RegisterInterface(
		"kira.gov.Content",
		(*govtypes.Content)(nil),
		&ProposalUpsertTokenAlias{},
		&ProposalUpsertTokenRates{},
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
