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
	MsgTypeSetNetworkProperties:          19,
	MsgTypeSetExecutionFee:               20,
	MsgTypeProposalAssignPermission:      21,
	MsgTypeProposalSetNetworkProperty:    22,
	MsgTypeProposalUpsertDataRegistry:    23,
	MsgTypeVoteProposal:                  24,
	MsgTypeClaimCouncilor:                25,
	MsgTypeWhitelistPermissions:          26,
	MsgTypeBlacklistPermissions:          27,
	MsgTypeCreateRole:                    28,
	MsgTypeAssignRole:                    29,
	MsgTypeRemoveRole:                    30,
	MsgTypeWhitelistRolePermission:       31,
	MsgTypeBlacklistRolePermission:       32,
	MsgTypeRemoveWhitelistRolePermission: 33,
	MsgTypeRemoveBlacklistRolePermission: 34,
	MsgTypeClaimValidator:                35,
	MsgTypeUpsertTokenAlias:              36,
	MsgTypeUpsertTokenRate:               37,
	MsgTypeProposalUpsertTokenAlias:      38,
	MsgProposalUpsertTokenAliasType:      39,
	MsgProposalUpsertTokenRatesType:      40,
	MsgTypeActivate:                      41,
	MsgTypePause:                         42,
	MsgTypeUnpause:                       43,
}
