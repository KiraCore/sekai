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
	cdc.RegisterConcrete(&MsgCreateSpendingPool{}, "kiraHub/MsgCreateSpendingPool", nil)
	functionmeta.AddNewFunction((&MsgCreateSpendingPool{}).Type(), `{
		"description": "MsgCreateSpendingPool represents a message to create a spending pool.",
		"parameters": {
			"name": {
				"type":        "string",
				"description": ""
			},
			"claim_start": {
				"type":        "time",
				"description": ""
			},
			"claim_end": {
				"type":        "time",
				"description": ""
			},
			"expire": {
				"type":        "uint64",
				"description": ""
			},
			"token": {
				"type":        "string",
				"description": ""
			},
			"rate": {
				"type":        "decimal",
				"description": ""
			},
			"vote_quorum": {
				"type":        "uint64",
				"description": ""
			},
			"vote_period": {
				"type":        "uint64",
				"description": ""
			},
			"vote_enactment": {
				"type":        "uint64",
				"description": ""
			},
			"owners": {
				"type":        "PermInfo",
				"description": ""
			},
			"beneficiaries": {
				"type":        "PermInfo",
				"description": ""
			},
			"sender": {
				"type":        "string",
				"description": ""
			}
		}
	}`)
	cdc.RegisterConcrete(&MsgDepositSpendingPool{}, "kiraHub/MsgDepositSpendingPool", nil)
	cdc.RegisterConcrete(&MsgRegisterSpendingPoolBeneficiary{}, "kiraHub/MsgRegisterSpendingPoolBeneficiary", nil)
	cdc.RegisterConcrete(&MsgClaimSpendingPool{}, "kiraHub/MsgClaimSpendingPool", nil)
}

// RegisterInterfaces register Msg and structs
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateSpendingPool{},
		&MsgDepositSpendingPool{},
		&MsgRegisterSpendingPoolBeneficiary{},
		&MsgClaimSpendingPool{},
	)

	registry.RegisterInterface(
		"kira.gov.Content",
		(*govtypes.Content)(nil),
		&UpdateSpendingPoolProposal{},
		&SpendingPoolDistributionProposal{},
		&SpendingPoolWithdrawProposal{},
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
