package types

const (
	ProposalTypeSoftwareUpgrade          = "SoftwareUpgrade"
	ProposalTypeCancelSoftwareUpgrade    = "CancelSoftwareUpgrade"
	ProposalTypeUpsertTokenAlias         = "UpsertTokenAlias"
	ProposalTypeUpsertTokenRates         = "UpsertTokenRates"
	ProposalTypeTokensWhiteBlackChange   = "TokensWhiteBlackChange"
	ProposalTypeUnjailValidator          = "UnjailValidator"
	ProposalTypeResetWholeValidatorRank  = "ResetWholeValidatorRank"
	ProposalTypeUpdateSpendingPool       = "UpdateSpendingPoolProposal"
	ProposalTypeSpendingPoolDistribution = "SpendingPoolDistributionProposal"
	ProposalTypeSpendingPoolWithdraw     = "SpendingPoolWithdrawProposal"
	AssignPermissionProposalType         = "AssignPermission"
	SetNetworkPropertyProposalType       = "SetNetworkProperty"
	UpsertDataRegistryProposalType       = "UpsertDataRegistry"
	SetPoorNetworkMessagesProposalType   = "SetPoorNetworkMessages"
	CreateRoleProposalType               = "CreateRoleProposal"
	SetProposalDurationsProposalType     = "SetProposalDurationsProposal"
)

var AllProposalTypes []string = []string{
	ProposalTypeSoftwareUpgrade,
	ProposalTypeCancelSoftwareUpgrade,
	ProposalTypeUpsertTokenAlias,
	ProposalTypeUpsertTokenRates,
	ProposalTypeTokensWhiteBlackChange,
	ProposalTypeUnjailValidator,
	ProposalTypeResetWholeValidatorRank,
	AssignPermissionProposalType,
	SetNetworkPropertyProposalType,
	UpsertDataRegistryProposalType,
	SetPoorNetworkMessagesProposalType,
	CreateRoleProposalType,
	SetProposalDurationsProposalType,
}
