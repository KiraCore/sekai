package types

import (
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// custom msg types
const (
	// governance
	MsgTypeProposalSetNetworkProperty     = "proposal-set-network-property"
	MsgTypeProposalAssignPermission       = "proposal-assign-permission"
	MsgTypeProposalUpsertDataRegistry     = "proposal-upsert-data-registry"
	MsgTypeProposalSetPoorNetworkMessages = "proposal-set-poor-network-messages"
	MsgTypeProposalCreateRole             = "proposal-create-role"
	MsgTypeVoteProposal                   = "vote-proposal"

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

	MsgTypeCreateIdentityRecord               = "create-identity-record"
	MsgTypeEditIdentityRecord                 = "edit-identity-record"
	MsgTypeRequestIdentityRecordsVerify       = "request-identity-records-verify"
	MsgTypeApproveIdentityRecords             = "approve-identity-records"
	MsgTypeCancelIdentityRecordsVerifyRequest = "cancel-identity-records-verify-request"

	// staking module
	MsgTypeClaimValidator          = "claim-validator"
	MsgTypeProposalUnjailValidator = "proposal-unjail-validator"

	// tokens module
	MsgTypeUpsertTokenAlias               = "upsert-token-alias"
	MsgTypeUpsertTokenRate                = "upsert-token-rate"
	MsgTypeProposalUpsertTokenAlias       = "propose-upsert-token-alias"
	MsgTypeProposalUpsertTokenRates       = "propose-upsert-token-rates"
	MsgTypeProposalTokensWhiteBlackChange = "propose-tokens-white-black-change"

	// slashing module
	MsgTypeActivate                        = "activate"
	MsgTypePause                           = "pause"
	MsgTypeUnpause                         = "unpause"
	MsgTypeProposalResetWholeValidatorRank = "proposal-reset-whole-validator-rank"

	//upgrade module
	MsgProposalSoftwareUpgrade       = "propose-software-upgrade"
	MsgProposalCancelSoftwareUpgrade = "propose-cancel-software-upgrade"
)

// MsgFuncIDMapping defines function_id mapping
var MsgFuncIDMapping = map[string]int64{
	bank.TypeMsgSend:      1,
	bank.TypeMsgMultiSend: 2,

	MsgTypeProposalAssignPermission:       3,
	MsgTypeProposalSetNetworkProperty:     4,
	MsgTypeProposalUpsertDataRegistry:     5,
	MsgTypeProposalSetPoorNetworkMessages: 6,
	MsgTypeProposalUpsertTokenAlias:       7,
	MsgTypeProposalUpsertTokenRates:       8,
	MsgTypeProposalTokensWhiteBlackChange: 9,
	MsgTypeProposalUnjailValidator:        10,
	MsgProposalSoftwareUpgrade:            11,
	MsgTypeVoteProposal:                   12,

	MsgTypeSetNetworkProperties:          20,
	MsgTypeSetExecutionFee:               21,
	MsgTypeClaimCouncilor:                22,
	MsgTypeWhitelistPermissions:          23,
	MsgTypeBlacklistPermissions:          24,
	MsgTypeCreateRole:                    25,
	MsgTypeAssignRole:                    26,
	MsgTypeRemoveRole:                    27,
	MsgTypeWhitelistRolePermission:       28,
	MsgTypeBlacklistRolePermission:       29,
	MsgTypeRemoveWhitelistRolePermission: 30,
	MsgTypeRemoveBlacklistRolePermission: 31,
	MsgTypeClaimValidator:                32,
	MsgTypeUpsertTokenAlias:              33,
	MsgTypeUpsertTokenRate:               34,
	MsgTypeActivate:                      35,
	MsgTypePause:                         36,
	MsgTypeUnpause:                       37,
}
