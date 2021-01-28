package types

import (
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// custom msg types
const (
	// governance
	MsgTypeProposalSetNetworkProperty = "proposal-set-network-property"
	MsgTypeProposalAssignPermission   = "proposal-assign-permission"
	MsgTypeProposalUpsertDataRegistry = "proposal-upsert-data-registry"
	MsgTypeProposalUpsertTokenAlias   = "proposal-upsert-token-alias"
	MsgTypeProposalSetPoorNetworkMsgs = "proposal-set-poor-network-messages"
	MsgTypeVoteProposal               = "vote-proposal"

	MsgTypeWhitelistPermissions = "whitelist-permissions"
	MsgTypeBlacklistPermissions = "blacklist-permissions"

	MsgTypeClaimCouncilor       = "claim-councilor"
	MsgTypeSetNetworkProperties = "set-network-properties"
	MsgTypeSetExecutionFee      = "set-execution-fee"

	MsgTypeCreateRole = "create-role"
	MsgTypeAssignRole = "assign-role"
	MsgTypeRemoveRole = "remove-role"

	MsgTypeWhitelistRolePermission       = "whitelist-role-permission"
	MsgTypeBlacklistRolePermission       = "blacklist-role-permission"
	MsgTypeRemoveWhitelistRolePermission = "remove-whitelist-role-permission"
	MsgTypeRemoveBlacklistRolePermission = "remove-blacklist-role-permission"

	// staking module
	MsgTypeClaimValidator = "claim-validator"

	// tokens module
	MsgTypeUpsertTokenAlias         = "upsert-token-alias"
	MsgTypeUpsertTokenRate          = "upsert-token-rate"
	MsgProposalUpsertTokenAliasType = "propose-upsert-token-alias"
	MsgProposalUpsertTokenRatesType = "propose-upsert-token-rates"

	// slashing module
	MsgTypeActivate = "activate"
	MsgTypePause    = "pause"
	MsgTypeUnpause  = "unpause"
)

// MsgFuncIDMapping defines function_id mapping
var MsgFuncIDMapping = map[string]int64{
	bank.TypeMsgSend:                     1,
	bank.TypeMsgMultiSend:                2,
	MsgTypeSetNetworkProperties:          3,
	MsgTypeSetExecutionFee:               4,
	MsgTypeProposalAssignPermission:      5,
	MsgTypeProposalSetNetworkProperty:    6,
	MsgTypeProposalUpsertDataRegistry:    7,
	MsgTypeVoteProposal:                  8,
	MsgTypeClaimCouncilor:                9,
	MsgTypeWhitelistPermissions:          10,
	MsgTypeBlacklistPermissions:          11,
	MsgTypeCreateRole:                    12,
	MsgTypeAssignRole:                    13,
	MsgTypeRemoveRole:                    14,
	MsgTypeWhitelistRolePermission:       15,
	MsgTypeBlacklistRolePermission:       16,
	MsgTypeRemoveWhitelistRolePermission: 17,
	MsgTypeRemoveBlacklistRolePermission: 18,
	MsgTypeClaimValidator:                19,
	MsgTypeUpsertTokenAlias:              20,
	MsgTypeUpsertTokenRate:               21,
	MsgTypeProposalUpsertTokenAlias:      22,
	MsgProposalUpsertTokenAliasType:      23,
	MsgProposalUpsertTokenRatesType:      24,
	MsgTypeActivate:                      25,
	MsgTypePause:                         26,
	MsgTypeUnpause:                       27,
}
