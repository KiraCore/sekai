package types

import (
	functionmeta "github.com/KiraCore/sekai/function_meta"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
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
	functionmeta.AddNewFunction((&MsgProposalAssignPermission{}).Type(), `{
		"description": "MsgProposalAssignPermission defines a proposal message to whitelist permission of an address.",
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
	cdc.RegisterConcrete(&MsgProposalSetNetworkProperty{}, "kiraHub/MsgProposalSetNetworkProperty", nil)
	functionmeta.AddNewFunction((&MsgProposalSetNetworkProperty{}).Type(), `{
		"description": "MsgProposalSetNetworkProperty defines a proposal message to modify a single network property.",
		"parameters": {
			"proposer": {
				"type":        "string",
				"description": "proposer who propose this message."
			},
			"network_property": {
				"type":        "enum<NetworkProperty>",
				"description": "network property to be modified."
			},
			"value": {
				"type":        "uint64",
				"description": "network property value to be set."
			}
		}
	}`)
	cdc.RegisterConcrete(&MsgProposalUpsertDataRegistry{}, "kiraHub/MsgProposalUpsertDataRegistry", nil)
	functionmeta.AddNewFunction((&MsgProposalUpsertDataRegistry{}).Type(), `{
		"description": "MsgProposalUpsertDataRegistry defines a proposal message to upsert data registry.",
		"parameters": {
			"proposer": {
				"type":        "string",
				"description": "proposer who propose this message."
			},
			"key": {
				"type":        "string",
				"description": "key of data registry."
			},
			"hash": {
				"type":        "string",
				"description": "hash value of data."
			},
			"reference": {
				"type":        "string",
				"description": "reference of data."
			},
			"encoding": {
				"type":        "string",
				"description": "encoding type of data."
			},
			"size": {
				"type":        "uint64",
				"description": "size of data."
			}
		}
	}`)
	cdc.RegisterConcrete(&MsgVoteProposal{}, "kiraHub/MsgVoteProposal", nil)
	functionmeta.AddNewFunction((&MsgVoteProposal{}).Type(), `{
		"description": "MsgVoteProposal defines a proposal message to vote on a submitted proposal.",
		"parameters": {
			"proposal_id": {
				"type":        "uint64",
				"description": "id of proposal to be voted."
			},
			"voter": {
				"type":        "address",
				"description": "the address of the voter who vote on the proposal."
			},
			"value": {
				"type":        "enum<VoteOption>",
				"description": "vote option: [yes, no, veto, abstain]"
			}
		}
	}`)
}

func registerCouncilorCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgClaimCouncilor{}, "kiraHub/MsgClaimCouncilor", nil)
	functionmeta.AddNewFunction((&MsgClaimCouncilor{}).Type(), `{
		"description": "MsgClaimCouncilor defines a message to claim councilor when the proposer.",
		"parameters": {
			"moniker": {
				"type":        "string",
				"description": "validator's name or nickname."
			},
			"website": {
				"type":        "string",
				"description": "validator's website."
			},
			"social": {
				"type":        "string",
				"description": "validator's social link."
			},
			"identity": {
				"type":        "string",
				"description": "validator's identity information."
			},
			"address": {
				"type":        "string",
				"description": "Address to be set as councilor. This address should be proposer address as well."
			}
		}
	}`)
}

func registerPermissionsCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgWhitelistPermissions{}, "kiraHub/MsgWhitelistPermissions", nil)
	functionmeta.AddNewFunction((&MsgWhitelistPermissions{}).Type(), `{
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
	functionmeta.AddNewFunction((&MsgCreateRole{}).Type(), `{
		"description": "MsgCreateRole defines a message to create a role.",
		"parameters": {
			"proposer": {
				"type":        "string",
				"description": "proposer who propose this message."
			},
			"role": {
				"type":        "uint32",
				"description": "Identifier of this role."
			}
		}
	}`)
	cdc.RegisterConcrete(&MsgAssignRole{}, "kiraHub/MsgAssignRole", nil)
	functionmeta.AddNewFunction((&MsgAssignRole{}).Type(), `{
		"description": "MsgAssignRole defines a message to assign a role to an address.",
		"parameters": {
			"proposer": {
				"type":        "string",
				"description": "proposer who propose this message."
			},
			"address": {
				"type":        "string",
				"description": "Address to set role to."
			},
			"role": {
				"type":        "uint32",
				"description": "role identifier."
			}
		}
	}`)
	cdc.RegisterConcrete(&MsgRemoveRole{}, "kiraHub/MsgRemoveRole", nil)
	functionmeta.AddNewFunction((&MsgRemoveRole{}).Type(), `{
		"description": "MsgRemoveRole defines a message to remove a role from an address.",
		"parameters": {
			"proposer": {
				"type":        "string",
				"description": "proposer who propose this message."
			},
			"address": {
				"type":        "string",
				"description": "Address to remove role from."
			},
			"role": {
				"type":        "uint32",
				"description": "role identifier."
			}
		}
	}`)

	cdc.RegisterConcrete(&MsgWhitelistRolePermission{}, "kiraHub/MsgWhitelistRolePermission", nil)
	functionmeta.AddNewFunction((&MsgWhitelistRolePermission{}).Type(), `{
		"description": "MsgWhitelistRolePermission defines a message to whitelist permission for a role.",
		"parameters": {
			"proposer": {
				"type":        "string",
				"description": "proposer who propose this message."
			},
			"role": {
				"type":        "uint32",
				"description": "role identifier."
			},
			"permission": {
				"type":        "uint32",
				"description": "Permission to be whitelisted."
			}
		}
	}`)
	cdc.RegisterConcrete(&MsgBlacklistRolePermission{}, "kiraHub/MsgBlacklistRolePermission", nil)
	functionmeta.AddNewFunction((&MsgBlacklistRolePermission{}).Type(), `{
		"description": "MsgBlacklistRolePermission defines a message to blacklist permission for a role.",
		"parameters": {
			"proposer": {
				"type":        "string",
				"description": "proposer who propose this message."
			},
			"role": {
				"type":        "uint32",
				"description": "role identifier."
			},
			"permission": {
				"type":        "uint32",
				"description": "Permission to be blacklisted."
			}
		}
	}`)
	cdc.RegisterConcrete(&MsgRemoveWhitelistRolePermission{}, "kiraHub/MsgRemoveWhitelistRolePermission", nil)
	functionmeta.AddNewFunction((&MsgRemoveWhitelistRolePermission{}).Type(), `{
		"description": "MsgRemoveWhitelistRolePermission defines a message to remove whitelisted permission for a role.",
		"parameters": {
			"proposer": {
				"type":        "string",
				"description": "proposer who propose this message."
			},
			"role": {
				"type":        "uint32",
				"description": "role identifier."
			},
			"permission": {
				"type":        "uint32",
				"description": "Permission to be removed from whitelisted listing."
			}
		}
	}`)
	cdc.RegisterConcrete(&MsgRemoveBlacklistRolePermission{}, "kiraHub/MsgRemoveBlacklistRolePermission", nil)
	functionmeta.AddNewFunction((&MsgRemoveBlacklistRolePermission{}).Type(), `{
		"description": "MsgRemoveBlacklistRolePermission defines a message to remove blacklisted permission for a role.",
		"parameters": {
			"proposer": {
				"type":        "string",
				"description": "proposer who propose this message."
			},
			"role": {
				"type":        "uint32",
				"description": "role identifier."
			},
			"permission": {
				"type":        "uint32",
				"description": "Permission to be removed from blacklisted listing."
			}
		}
	}`)
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
		&MsgProposalSetPoorNetworkMessages{},
		&MsgProposalCreateRole{},
		&MsgVoteProposal{},

		&MsgCreateIdentityRecord{},
		&MsgEditIdentityRecord{},
		&MsgRequestIdentityRecordsVerify{},
		&MsgApproveIdentityRecords{},
		&MsgCancelIdentityRecordsVerifyRequest{},
	)

	registry.RegisterInterface(
		"kira.gov.Content",
		(*Content)(nil),
		&AssignPermissionProposal{},
		&SetNetworkPropertyProposal{},
		&UpsertDataRegistryProposal{},
		&SetPoorNetworkMessagesProposal{},
		&CreateRoleProposal{},
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
