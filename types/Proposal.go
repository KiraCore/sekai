package types

const (
	ProposalTypeSoftwareUpgrade                    = "SoftwareUpgrade"
	ProposalTypeCancelSoftwareUpgrade              = "CancelSoftwareUpgrade"
	ProposalTypeUpsertTokenInfos                   = "UpsertTokenInfos"
	ProposalTypeTokensWhiteBlackChange             = "TokensWhiteBlackChange"
	ProposalTypeUnjailValidator                    = "UnjailValidator"
	ProposalTypeResetWholeValidatorRank            = "ResetWholeValidatorRank"
	ProposalTypeSlashValidator                     = "SlashValidator"
	ProposalTypeUpdateSpendingPool                 = "UpdateSpendingPool"
	ProposalTypeSpendingPoolDistribution           = "SpendingPoolDistribution"
	ProposalTypeSpendingPoolWithdraw               = "SpendingPoolWithdraw"
	ProposalTypeUpsertUBI                          = "UpsertUBI"
	ProposalTypeRemoveUBI                          = "RemoveUBI"
	ProposalTypeResetWholeCouncilorRank            = "ResetWholeCouncilorRank"
	ProposalTypeJailCouncilor                      = "JailCouncilor"
	ProposalTypeSetExecutionFees                   = "SetExecutionFees"
	ProposalTypeWhitelistAccountPermission         = "WhitelistAccountPermission"
	ProposalTypeBlacklistAccountPermission         = "BlacklistAccountPermission"
	ProposalTypeRemoveWhitelistedAccountPermission = "RemoveWhitelistedAccountPermission"
	ProposalTypeRemoveBlacklistedAccountPermission = "RemoveBlacklistedAccountPermission"
	ProposalTypeAssignRoleToAccount                = "AssignRoleToAccount"
	ProposalTypeUnassignRoleFromAccount            = "UnassignRoleFromAccount"
	ProposalTypeSetNetworkProperty                 = "SetNetworkProperty"
	ProposalTypeUpsertDataRegistry                 = "UpsertDataRegistry"
	ProposalTypeSetPoorNetworkMessages             = "SetPoorNetworkMessages"
	ProposalTypeCreateRole                         = "CreateRole"
	ProposalTypeRemoveRole                         = "RemoveRole"
	ProposalTypeWhitelistRolePermission            = "WhitelistRolePermission"
	ProposalTypeBlacklistRolePermission            = "BlacklistRolePermission"
	ProposalTypeRemoveWhitelistedRolePermission    = "RemoveWhitelistedRolePermission"
	ProposalTypeRemoveBlacklistedRolePermission    = "RemoveBlacklistedRolePermission"
	ProposalTypeSetProposalDurations               = "SetProposalDurations"

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
	ProposalTypeUpsertTokenInfos,
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
