package types

const (
	ProposalTypeSoftwareUpgrade                    = "SoftwareUpgrade"
	ProposalTypeCancelSoftwareUpgrade              = "CancelSoftwareUpgrade"
	ProposalTypeUpsertTokenAlias                   = "UpsertTokenAlias"
	ProposalTypeUpsertTokenRates                   = "UpsertTokenRates"
	ProposalTypeTokensWhiteBlackChange             = "TokensWhiteBlackChange"
	ProposalTypeUnjailValidator                    = "UnjailValidator"
	ProposalTypeResetWholeValidatorRank            = "ResetWholeValidatorRank"
	ProposalTypeUpdateSpendingPool                 = "UpdateSpendingPoolProposal"
	ProposalTypeSpendingPoolDistribution           = "SpendingPoolDistributionProposal"
	ProposalTypeSpendingPoolWithdraw               = "SpendingPoolWithdrawProposal"
	ProposalTypeUpsertUBI                          = "UpsertUBIProposal"
	ProposalTypeRemoveUBI                          = "RemoveUBIProposal"
	AssignPermissionProposalType                   = "AssignPermission"
	WhitelistAccountPermissionProposalType         = "WhitelistAccountPermission"
	BlacklistAccountPermissionProposalType         = "BlacklistAccountPermission"
	RemoveWhitelistedAccountPermissionProposalType = "RemoveWhitelistedAccountPermission"
	RemoveBlacklistedAccountPermissionProposalType = "RemoveBlacklistedAccountPermission"
	AssignRoleToAccountProposalType                = "AssignRoleToAccount"
	UnassignRoleFromAccountProposalType            = "UnassignRoleFromAccount"
	SetNetworkPropertyProposalType                 = "SetNetworkProperty"
	UpsertDataRegistryProposalType                 = "UpsertDataRegistry"
	SetPoorNetworkMessagesProposalType             = "SetPoorNetworkMessages"
	CreateRoleProposalType                         = "CreateRoleProposal"
	RemoveRoleProposalType                         = "RemoveRoleProposal"
	WhitelistRolePermissionProposalType            = "WhitelistRolePermission"
	BlacklistRolePermissionProposalType            = "BlacklistRolePermission"
	RemoveWhitelistedRolePermissionProposalType    = "RemoveWhitelistedRolePermission"
	RemoveBlacklistedRolePermissionProposalType    = "RemoveBlacklistedRolePermission"
	SetProposalDurationsProposalType               = "SetProposalDurationsProposal"
)

var AllProposalTypes []string = []string{
	ProposalTypeSoftwareUpgrade,
	ProposalTypeCancelSoftwareUpgrade,
	ProposalTypeUpsertTokenAlias,
	ProposalTypeUpsertTokenRates,
	ProposalTypeTokensWhiteBlackChange,
	ProposalTypeUnjailValidator,
	ProposalTypeResetWholeValidatorRank,
	WhitelistAccountPermissionProposalType,
	BlacklistAccountPermissionProposalType,
	RemoveWhitelistedAccountPermissionProposalType,
	RemoveBlacklistedAccountPermissionProposalType,
	AssignRoleToAccountProposalType,
	UnassignRoleFromAccountProposalType,
	SetNetworkPropertyProposalType,
	UpsertDataRegistryProposalType,
	SetPoorNetworkMessagesProposalType,
	CreateRoleProposalType,
	RemoveRoleProposalType,
	WhitelistRolePermissionProposalType,
	BlacklistRolePermissionProposalType,
	RemoveWhitelistedRolePermissionProposalType,
	RemoveBlacklistedRolePermissionProposalType,
	SetProposalDurationsProposalType,
}
