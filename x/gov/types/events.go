package types

// governance module events
const (
	EventTypeProposalVote                  = "proposal_vote"
	EventTypeSubmitProposal                = "submit_proposal"
	EventTypeCreateRole                    = "create_role"
	EventTypeAssignRole                    = "assign_role"
	EventTypeRemoveRole                    = "remove_role"
	EventTypeWhitelistPermisison           = "whitelist_permission"
	EventTypeBlacklistPermisison           = "blacklist_permission"
	EventTypeWhitelistRolePermisison       = "whitelist_role_permission"
	EventTypeBlacklistRolePermisison       = "blacklist_role_permission"
	EventTypeRemoveBlacklistRolePermisison = "remove_blacklist_role_permission"
	EventTypeRemoveWhitelistRolePermisison = "remove_whitelist_role_permission"
	EventTypeSetNetworkProperties          = "set_network_properties"
	EventTypeSetExecutionFee               = "set_execution_fee"
	EventTypeClaimCouncilor                = "claim_councilor"
	EventTypeAddToEnactment                = "add_to_enactments"
	EventTypeRemoveEnactment               = "remove_from_enactments"

	AttributeKeyProposalId          = "proposal_id"
	AttributeKeyProposalType        = "proposal_type"
	AttributeKeyProposalContent     = "proposal_content"
	AttributeKeyVoter               = "voter"
	AttributeKeyOption              = "option"
	AttributeKeyProposer            = "proposer"
	AttributeKeyAddress             = "address"
	AttributeKeyRoleId              = "role_id"
	AttributeKeyPermission          = "permission"
	AttributeKeyProperties          = "properties"
	AttributeKeyTransactionType     = "transaction_type"
	AttributeKeyExecutionFee        = "execution_fee"
	AttributeKeyFailureFee          = "failure_fee"
	AttributeKeyTimeout             = "time_out"
	AttributeKeyDefaultParameters   = "default_parameters"
	AttributeKeyProposalDescription = "description"

	// ---- Cosmos SDK gov native events ----
	// EventTypeProposalDeposit  = "proposal_deposit"
	// EventTypeInactiveProposal = "inactive_proposal"
	// EventTypeActiveProposal   = "active_proposal"

	// AttributeKeyProposalResult     = "proposal_result"
	// AttributeKeyProposalID         = "proposal_id"
	// AttributeKeyVotingPeriodStart  = "voting_period_start"
	// AttributeValueCategory         = "governance"
	// AttributeValueProposalDropped  = "proposal_dropped"  // didn't meet min deposit
	// AttributeValueProposalPassed   = "proposal_passed"   // met vote quorum
	// AttributeValueProposalRejected = "proposal_rejected" // didn't meet vote quorum
	// AttributeValueProposalFailed   = "proposal_failed"   // error on proposal handler
	// AttributeKeyProposalType       = "proposal_type"
)
