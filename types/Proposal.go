package types

const (
	ProposalTypeSoftwareUpgrade                    = "SoftwareUpgrade"
	ProposalTypeCancelSoftwareUpgrade              = "CancelSoftwareUpgrade"
	ProposalTypeUpsertTokenAlias                   = "UpsertTokenAlias"
	ProposalTypeUpsertTokenRates                   = "UpsertTokenRates"
	ProposalTypeTokensWhiteBlackChange             = "TokensWhiteBlackChange"
	ProposalTypeUnjailValidator                    = "UnjailValidator"
	ProposalTypeResetWholeValidatorRank            = "ResetWholeValidatorRank"
	ProposalTypeSlashValidator                     = "SlashValidator"
	ProposalTypeUpdateSpendingPool                 = "UpdateSpendingPoolProposal"
	ProposalTypeSpendingPoolDistribution           = "SpendingPoolDistributionProposal"
	ProposalTypeSpendingPoolWithdraw               = "SpendingPoolWithdrawProposal"
	ProposalTypeUpsertUBI                          = "UpsertUBIProposal"
	ProposalTypeRemoveUBI                          = "RemoveUBIProposal"
	ProposalTypeResetWholeCouncilorRank            = "ResetWholeCouncilorRank"
	ProposalTypeJailCouncilor                      = "JailCouncilor"
	ProposalTypeWhitelistAccountPermission         = "WhitelistAccountPermission"
	ProposalTypeBlacklistAccountPermission         = "BlacklistAccountPermission"
	ProposalTypeRemoveWhitelistedAccountPermission = "RemoveWhitelistedAccountPermission"
	ProposalTypeRemoveBlacklistedAccountPermission = "RemoveBlacklistedAccountPermission"
	ProposalTypeAssignRoleToAccount                = "AssignRoleToAccount"
	ProposalTypeUnassignRoleFromAccount            = "UnassignRoleFromAccount"
	ProposalTypeSetNetworkProperty                 = "SetNetworkProperty"
	ProposalTypeUpsertDataRegistry                 = "UpsertDataRegistry"
	ProposalTypeSetPoorNetworkMessages             = "SetPoorNetworkMessages"
	ProposalTypeCreateRole                         = "CreateRoleProposal"
	ProposalTypeRemoveRole                         = "RemoveRoleProposal"
	ProposalTypeWhitelistRolePermission            = "WhitelistRolePermission"
	ProposalTypeBlacklistRolePermission            = "BlacklistRolePermission"
	ProposalTypeRemoveWhitelistedRolePermission    = "RemoveWhitelistedRolePermission"
	ProposalTypeRemoveBlacklistedRolePermission    = "RemoveBlacklistedRolePermission"
	ProposalTypeSetProposalDurations               = "SetProposalDurationsProposal"

	ProposalTypeCreateBasket          = "CreateBasket"
	ProposalTypeEditBasket            = "EditBasket"
	ProposalTypeBasketWithdrawSurplus = "BasketWithdrawSurplus"

	ProposalTypeCollectiveSendDonation = "CollectiveSendDonation"
	ProposalTypeCollectiveUpdate       = "CollectiveUpdate"
	ProposalTypeCollectiveRemove       = "CollectiveRemove"

	ProposalTypeJoinDapp       = "JoinDapp"
	ProposalTypeTransitionDapp = "TransitionDapp"
	ProposalTypeUpsertDapp     = "UpsertDapp"
)

var AllProposalTypes []string = []string{
	ProposalTypeSoftwareUpgrade,
	ProposalTypeCancelSoftwareUpgrade,
	ProposalTypeUpsertTokenAlias,
	ProposalTypeUpsertTokenRates,
	ProposalTypeTokensWhiteBlackChange,
	ProposalTypeUnjailValidator,
	ProposalTypeResetWholeValidatorRank,
	ProposalTypeSlashValidator,
	ProposalTypeResetWholeCouncilorRank,
	ProposalTypeWhitelistAccountPermission,
	ProposalTypeBlacklistAccountPermission,
	ProposalTypeRemoveWhitelistedAccountPermission,
	ProposalTypeRemoveBlacklistedAccountPermission,
	ProposalTypeAssignRoleToAccount,
	ProposalTypeUnassignRoleFromAccount,
	ProposalTypeSetNetworkProperty,
	ProposalTypeUpsertDataRegistry,
	ProposalTypeSetPoorNetworkMessages,
	ProposalTypeCreateRole,
	ProposalTypeRemoveRole,
	ProposalTypeWhitelistRolePermission,
	ProposalTypeBlacklistRolePermission,
	ProposalTypeRemoveWhitelistedRolePermission,
	ProposalTypeRemoveBlacklistedRolePermission,
	ProposalTypeSetProposalDurations,
}
