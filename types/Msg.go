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
	MsgTypeSubmitProposal = "submit_proposal"
	MsgTypeVoteProposal   = "vote_proposal"
	MsgTypeCreatePoll     = "create_poll"
	MsgTypeVotePoll       = "vote_poll"
	MsgTypeAddressPoll    = "address_poll"

	MsgTypeWhitelistPermissions = "whitelist_permissions"
	MsgTypeBlacklistPermissions = "blacklist_permissions"

	MsgTypeClaimCouncilor       = "claim_councilor"
	MsgTypeSetNetworkProperties = "set_network_properties"
	MsgTypeSetExecutionFee      = "set_execution_fee"

	MsgTypeCreateRole   = "create_role"
	MsgTypeAssignRole   = "assign_role"
	MsgTypeUnassignRole = "unassign_role"

	MsgTypeWhitelistRolePermission       = "whitelist_role_permission"
	MsgTypeBlacklistRolePermission       = "blacklist_role_permission"
	MsgTypeRemoveWhitelistRolePermission = "remove_whitelist_role_permission"
	MsgTypeRemoveBlacklistRolePermission = "remove_blacklist_role_permission"

	MsgTypeRegisterIdentityRecords            = "register_identity_records"
	MsgTypeDeleteIdentityRecords              = "delete_identity_records"
	MsgTypeRequestIdentityRecordsVerify       = "request_identity_records_verify"
	MsgTypeHandleIdentityRecordsVerifyRequest = "handle_identity_records_verify_request"
	MsgTypeCancelIdentityRecordsVerifyRequest = "cancel_identity_records_verify_request"

	// staking module
	MsgTypeClaimValidator = "claim_validator"

	// multistaking module
	MsgTypeUpsertStakingPool         = "upsert_staking_pool"
	MsgTypeDelegate                  = "delegate"
	MsgTypeUndelegate                = "undelegate"
	MsgTypeClaimRewards              = "claim_rewards"
	MsgTypeClaimUndelegation         = "claim_undelegation"
	MsgTypeClaimMaturedUndelegations = "claim_matured_undelegations"
	MsgTypeSetCompoundInfo           = "set_compound_info"
	MsgTypeRegisterDelegator         = "register_delegator"

	// basket module
	MsgTypeDisableBasketDeposits  = "disable_basket_deposits"
	MsgTypeDisableBasketWithdraws = "disable_basket_withdraws"
	MsgTypeDisableBasketSwaps     = "disable_basket_swaps"
	MsgTypeBasketTokenMint        = "basket_token_mint"
	MsgTypeBasketTokenBurn        = "basket_token_burn"
	MsgTypeBasketTokenSwap        = "basket_token_swap"
	MsgTypeBasketClaimRewards     = "basket_claim_rewards"

	// tokens module
	MsgTypeUpsertTokenInfo = "upsert_token_info"
	MsgTypeEthereumTx      = "ethereum_tx"

	// slashing module
	MsgTypeActivate = "activate"
	MsgTypePause    = "pause"
	MsgTypeUnpause  = "unpause"

	// recovery module
	MsgTypeRegisterRecoverySecret             = "register_recovery_secret"
	MsgTypeRotateRecoveryAddress              = "rotate_recovery_address"
	MsgTypeIssueRecoveryTokens                = "issue_recovery_tokens"
	MsgTypeBurnRecoveryTokens                 = "burn_recovery_tokens"
	MsgTypeRegisterRRTokenHolder              = "register_rrtoken_holder"
	MsgTypeClaimRRHolderRewards               = "claim_rrholder_rewards"
	MsgTypeRotateValidatorByHalfRRTokenHolder = "rotate_validator_by_half_rr_token_holder"

	//upgrade module

	// spending module
	MsgTypeCreateSpendingPool              = "create_spending_pool"
	MsgTypeDepositSpendingPool             = "deposit_spending_pool"
	MsgTypeRegisterSpendingPoolBeneficiary = "register_spending_pool_beneficiary"
	MsgTypeClaimSpendingPool               = "claim_spending_pool"

	// custody module
	MsgTypeCreateCustody               = "create_custody"
	MsgTypeDisableCustody              = "disable_custody"
	MsgTypeDropCustody                 = "drop_custody"
	MsgTypeAddToCustodyWhiteList       = "add_to_custody_whitelist"
	MsgTypeAddToCustodyCustodians      = "add_to_custody_custodians"
	MsgTypeRemoveFromCustodyCustodians = "remove_from_custody_custodians"
	MsgTypeDropCustodyCustodians       = "drop_custody_custodians"
	MsgTypeRemoveFromCustodyWhiteList  = "remove_from_custody_whitelist"
	MsgTypeDropCustodyWhiteList        = "drop_custody_whitelist"
	MsgApproveCustodyTransaction       = "approve_custody_transaction"
	MsgDeclineCustodyTransaction       = "decline_custody_transaction"
	MsgPasswordConfirmTransaction      = "password_confirm_transaction"
	MsgTypeSend                        = "custody_send"

	// bridge module
	MsgTypeChangeCosmosEthereum = "change-cosmos-ethereum"
	MsgTypeChangeEthereumCosmos = "change-ethereum-cosmos"

	// collectives module
	MsgTypeCreateCollective   = "create_collective"
	MsgTypeBondCollective     = "bond_collective"
	MsgTypeDonateCollective   = "donate_collective"
	MsgTypeWithdrawCollective = "withdraw_collective"

	// layer2 module
	MsgTypeCreateDappProposal       = "create_dapp_proposal"
	MsgTypeBondDappProposal         = "bond_dapp_proposal"
	MsgTypeReclaimDappBondProposal  = "reclaim_dapp_bond_proposal"
	MsgTypeJoinDappVerifierWithBond = "join_dapp_verifier_with_bond"
	MsgTypeExitDapp                 = "exit_dapp"
	MsgTypeVoteDappOperatorTx       = "vote_dapp_operator_tx"
	MsgTypeRedeemDappPoolTx         = "redeem_dapp_pool_tx"
	MsgTypeSwapDappPoolTx           = "swap_dapp_pool_tx"
	MsgTypeConvertDappPoolTx        = "convert_dapp_pool_tx"
	MsgTypePauseDappTx              = "pause_dapp_tx"
	MsgTypeUnPauseDappTx            = "unpause_dapp_tx"
	MsgTypeReactivateDappTx         = "reactivate_dapp_tx"
	MsgTypeExecuteDappTx            = "execute_dapp_tx"
	MsgTypeDenounceLeaderTx         = "denounce_leader_tx"
	MsgTypeTransitionDappTx         = "transition_dapp_tx"
	MsgTypeApproveDappTransitionTx  = "approve_dapp_transition_tx"
	MsgTypeRejectDappTransitionTx   = "reject_dapp_transition_tx"
	MsgTypeUpsertDappProposalTx     = "upsert_dapp_proposal_tx"
	MsgTypeVoteUpsertDappProposalTx = "vote_upsert_dapp_proposal_tx"
	MsgTypeTransferDappTx           = "transfer_dapp_tx"
	MsgTypeAckTransferDappTx        = "ack_transfer_dapp_tx"
	MsgTypeMintCreateFtTx           = "mint_create_ft_tx"
	MsgTypeMintCreateNftTx          = "mint_create_nft_tx"
	MsgTypeMintIssueTx              = "mint_issue_tx"
	MsgTypeMintBurnTx               = "mint_burn_tx"
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
	MsgTypeDeleteIdentityRecords:              13,
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
	MsgTypeUnassignRole:                  27,
	MsgTypeWhitelistRolePermission:       28,
	MsgTypeBlacklistRolePermission:       29,
	MsgTypeRemoveWhitelistRolePermission: 30,
	MsgTypeRemoveBlacklistRolePermission: 31,
	MsgTypeClaimValidator:                32,
	MsgTypeUpsertTokenInfo:               34,
	MsgTypeActivate:                      35,
	MsgTypePause:                         36,
	MsgTypeUnpause:                       37,
	MsgTypeEthereumTx:                    38,

	MsgTypeCreateSpendingPool:              41,
	MsgTypeDepositSpendingPool:             42,
	MsgTypeRegisterSpendingPoolBeneficiary: 43,
	MsgTypeClaimSpendingPool:               44,

	MsgTypeUpsertStakingPool: 51,
	MsgTypeDelegate:          52,
	MsgTypeUndelegate:        53,
	MsgTypeClaimRewards:      54,
	MsgTypeClaimUndelegation: 55,
	MsgTypeSetCompoundInfo:   56,
	MsgTypeRegisterDelegator: 57,

	MsgTypeCreateCustody:               61,
	MsgTypeAddToCustodyWhiteList:       62,
	MsgTypeAddToCustodyCustodians:      63,
	MsgTypeRemoveFromCustodyCustodians: 64,
	MsgTypeDropCustodyCustodians:       65,
	MsgTypeRemoveFromCustodyWhiteList:  66,
	MsgTypeDropCustodyWhiteList:        67,
	MsgApproveCustodyTransaction:       68,
	MsgDeclineCustodyTransaction:       69,
}

func MsgType(msg sdk.Msg) string {
	kiraMsg, ok := msg.(Msg)
	if !ok {
		return ""
	}
	return kiraMsg.Type()
}
