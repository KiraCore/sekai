package types

import (
	functionmeta "github.com/KiraCore/sekai/function_meta"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterCodec register codec and metadata
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*Content)(nil), nil)

	registerPermissionsCodec(cdc)
	registerRolesCodec(cdc)
	registerCouncilorCodec(cdc)
	registerProposalCodec(cdc)

	cdc.RegisterConcrete(&MsgSetNetworkProperties{}, "kiraHub/MsgSetNetworkProperties", nil)
	functionmeta.AddNewFunction((&MsgSetNetworkProperties{}).Type(), `{
		"function_id": 0,
		"description": "MsgSetNetworkProperties defines a message to set network properties with specific permission.",
		"parameters": {
			"network_properties": {
				"type":        "<Object>",
				"description": "network properties to be set.",
				"fields": {
					"min_tx_fee": {
						"type":        "uint64",
						"description": "minimum transaction fee"
					},
					"max_tx_fee": {
						"type":        "uint64",
						"description": "maximum transaction fee"
					},
					"vote_quorum": {
						"type":        "uint64",
						"description": "vote quorum"
					},
					"proposal_end_time": {
						"type":        "uint64",
						"description": "proposal end time"
					},
					"proposal_enactment_time": {
						"type":        "uint64",
						"description": "proposal enactment time"
					},
					"enable_foreign_fee_payments": {
						"type":        "bool",
						"description": "flag to show if foreign fee payment is enabled"
					}
				}
			},
			"proposer": {
				"type":        "address",
				"description": "proposer who propose this message."
			}
		}
	}`)
	cdc.RegisterConcrete(&MsgSetExecutionFee{}, "kiraHub/MsgSetExecutionFee", nil)
	functionmeta.AddNewFunction((&MsgSetExecutionFee{}).Type(), `{
		"function_id": 0,
		"description": "MsgSetExecutionFee defines a message to set execution fee with specific permission.",
		"parameters": {
			"name": {
				"type":        "string",
				"description": "Friendly Name of the Function (max 128 characters)"
			},
			"transaction_type": {
				"type":        "string",
				"description": "Type of the transaction that given permission allows to execute"
			},
			"execution_fee": {
				"type":        "uint64",
				"description": "How much user should pay for executing this specific function"
			},
			"failure_fee": {
				"type":        "uint64",
				"description": "How much user should pay if function fails to execute"
			},
			"timeout": {
				"type":        "uint64",
				"description": "After what time function execution should fail"
			},
			"default_parameters": {
				"type":        "bool",
				"description": "Default values that the function in question will consume as input parameters before execution"
			},
			"proposer": {
				"type":        "address",
				"description": "proposer who propose this message."
			}
		}
	}`)
}

func registerProposalCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgProposalAssignPermission{}, "kiraHub/MsgProposalAssignPermission", nil)
	cdc.RegisterConcrete(&MsgProposalSetNetworkProperty{}, "kiraHub/MsgProposalSetNetworkProperty", nil)
	cdc.RegisterConcrete(&MsgProposalUpsertDataRegistry{}, "kiraHub/MsgProposalUpsertDataRegistry", nil)
	cdc.RegisterConcrete(&MsgVoteProposal{}, "kiraHub/MsgVoteProposal", nil)
}

func registerCouncilorCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgClaimCouncilor{}, "kiraHub/MsgClaimCouncilor", nil)
}

func registerPermissionsCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgWhitelistPermissions{}, "kiraHub/MsgWhitelistPermissions", nil)
	functionmeta.AddNewFunction((&MsgWhitelistPermissions{}).Type(), `{
		"function_id": 0,
		"description": "MsgWhitelistPermissions defines a message to whitelist permission of an address.",
		"parameters": {
			"proposer": {
				"type":        "string",
				"description": "proposer who propose this message."
			},
			"address": {
				"type":        "string",
				"description": "Address to whitelist permission to."
			},
			"permission": {
				"type":        "uint32",
				"description": "Permission to be whitelisted."
			}
		}
	}`)
	cdc.RegisterConcrete(&MsgBlacklistPermissions{}, "kiraHub/MsgBlacklistPermissions", nil)
	functionmeta.AddNewFunction((&MsgBlacklistPermissions{}).Type(), `{
		"function_id": 0,
		"description": "MsgBlacklistPermissions defines a message to blacklist permission of an address.",
		"parameters": {
			"proposer": {
				"type":        "string",
				"description": "proposer who propose this message."
			},
			"address": {
				"type":        "string",
				"description": "Address to blacklist permission to."
			},
			"permission": {
				"type":        "uint32",
				"description": "Permission to be blacklisted."
			}
		}
	}`)
}

func registerRolesCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateRole{}, "kiraHub/MsgCreateRole", nil)
	cdc.RegisterConcrete(&MsgAssignRole{}, "kiraHub/MsgAssignRole", nil)
	cdc.RegisterConcrete(&MsgRemoveRole{}, "kiraHub/MsgRemoveRole", nil)

	cdc.RegisterConcrete(&MsgWhitelistRolePermission{}, "kiraHub/MsgWhitelistRolePermission", nil)
	cdc.RegisterConcrete(&MsgBlacklistRolePermission{}, "kiraHub/MsgBlacklistRolePermission", nil)
	cdc.RegisterConcrete(&MsgRemoveWhitelistRolePermission{}, "kiraHub/MsgRemoveWhitelistRolePermission", nil)
	cdc.RegisterConcrete(&MsgRemoveBlacklistRolePermission{}, "kiraHub/MsgRemoveBlacklistRolePermission", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgWhitelistPermissions{},
		&MsgBlacklistPermissions{},

		&MsgSetNetworkProperties{},
		&MsgSetExecutionFee{},

		&MsgClaimCouncilor{},

		&MsgAssignRole{},
		&MsgCreateRole{},
		&MsgRemoveRole{},

		&MsgWhitelistRolePermission{},
		&MsgBlacklistRolePermission{},
		&MsgRemoveWhitelistRolePermission{},
		&MsgRemoveBlacklistRolePermission{},

		&MsgProposalAssignPermission{},
		&MsgProposalSetNetworkProperty{},
		&MsgProposalUpsertDataRegistry{},
		&MsgVoteProposal{},
	)

	registry.RegisterInterface(
		"kira.gov.Content",
		(*Content)(nil),
		&AssignPermissionProposal{},
		&SetNetworkPropertyProposal{},
		&UpsertDataRegistryProposal{},
	)

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
