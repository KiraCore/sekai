package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// custom msg types
const (
	//evidence
	TypeMsgSubmitEvidence = "submit_evidence"

	// governance
	MsgTypeSubmitProposal = "submit-proposal"
	MsgTypeVoteProposal   = "vote-proposal"

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

	MsgTypeRegisterIdentityRecords            = "register-identity-records"
	MsgTypeEditIdentityRecord                 = "edit-identity-record"
	MsgTypeRequestIdentityRecordsVerify       = "request-identity-records-verify"
	MsgTypeHandleIdentityRecordsVerifyRequest = "handle-identity-records-verify-request"
	MsgTypeCancelIdentityRecordsVerifyRequest = "cancel-identity-records-verify-request"

	// staking module
	MsgTypeClaimValidator = "claim-validator"

	// multistaking module
	MsgTypeUpsertStakingPool = "upsert_staking_pool"
	MsgTypeDelegate          = "delegate"
	MsgTypeUndelegate        = "undelegate"
	MsgTypeClaimRewards      = "claim_rewards"
	MsgTypeClaimUndelegation = "claim_undelegation"
	MsgTypeSetCompoundInfo   = "set_compound_info"
	MsgTypeRegisterDelegator = "register_delegator"

	// basket module
	MsgTypeDisableBasketDeposits  = "disable-basket-deposits"
	MsgTypeDisableBasketWithdraws = "disable-basket-withdraws"
	MsgTypeDisableBasketSwaps     = "disable-basket-swaps"
	MsgTypeBasketTokenMint        = "basket-token-mint"
	MsgTypeBasketTokenBurn        = "basket-token-burn"
	MsgTypeBasketTokenSwap        = "basket-token-swap"
	MsgTypeBasketClaimRewards     = "basket-claim-rewards"

	// tokens module
	MsgTypeUpsertTokenAlias = "upsert-token-alias"
	MsgTypeUpsertTokenRate  = "upsert-token-rate"

	// slashing module
	MsgTypeActivate = "activate"
	MsgTypePause    = "pause"
	MsgTypeUnpause  = "unpause"

	//upgrade module

	// spending module
	MsgTypeCreateSpendingPool              = "create-spending-pool"
	MsgTypeDepositSpendingPool             = "deposit-spending-pool"
	MsgTypeRegisterSpendingPoolBeneficiary = "register-spending-pool-beneficiary"
	MsgTypeClaimSpendingPool               = "claim-spending-pool"

	// custody module
	MsgTypeCreateCustody               = "create-custody"
	MsgTypeAddToCustodyWhiteList       = "add-to-custody-whitelist"
	MsgTypeAddToCustodyCustodians      = "add-to-custody-custodians"
	MsgTypeRemoveFromCustodyCustodians = "remove-from-custody-custodians"
	MsgTypeDropCustodyCustodians       = "drop-custody-custodians"
	MsgTypeRemoveFromCustodyWhiteList  = "remove-from-custody-whitelist"
	MsgTypeDropCustodyWhiteList        = "drop-custody-whitelist"
	MsgApproveCustodyTransaction       = "approve-custody-transaction"
	MsgDeclineCustodyTransaction       = "decline-custody-transaction"
	MsgPasswordConfirmTransaction      = "password-confirm-transaction"
	MsgTypeSend                        = "custody-send"
)

// Msg defines the interface a transaction message must fulfill.
type Msg interface {
	sdk.Msg

	// Type returns type of message
	Type() string
}

// MsgFuncIDMapping defines function_id mapping
var MsgFuncIDMapping = map[string]int64{
	bank.TypeMsgSend:      1,
	bank.TypeMsgMultiSend: 2,

	TypeMsgSubmitEvidence: 3,

	MsgTypeSubmitProposal:                     10,
	MsgTypeVoteProposal:                       11,
	MsgTypeRegisterIdentityRecords:            12,
	MsgTypeEditIdentityRecord:                 13,
	MsgTypeRequestIdentityRecordsVerify:       14,
	MsgTypeHandleIdentityRecordsVerifyRequest: 15,
	MsgTypeCancelIdentityRecordsVerifyRequest: 16,

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

	MsgTypeCreateSpendingPool:              41,
	MsgTypeDepositSpendingPool:             42,
	MsgTypeRegisterSpendingPoolBeneficiary: 43,
	MsgTypeClaimSpendingPool:               44,

	MsgTypeCreateCustody:               45,
	MsgTypeAddToCustodyWhiteList:       46,
	MsgTypeAddToCustodyCustodians:      47,
	MsgTypeRemoveFromCustodyCustodians: 48,
	MsgTypeDropCustodyCustodians:       49,
	MsgTypeRemoveFromCustodyWhiteList:  50,
	MsgTypeDropCustodyWhiteList:        51,
	MsgApproveCustodyTransaction:       52,
	MsgDeclineCustodyTransaction:       53,
}

func MsgType(msg sdk.Msg) string {
	kiraMsg, ok := msg.(Msg)
	if !ok {
		return ""
	}
	return kiraMsg.Type()
}
