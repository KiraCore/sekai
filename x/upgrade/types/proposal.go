package types

import (
	"github.com/KiraCore/sekai/x/gov/types"
)

const ProposalTypeSoftwareUpgrade = "SoftwareUpgrade"
const ProposalTypeCancelSoftwareUpgrade = "CancelSoftwareUpgrade"

var _ types.Content = &ProposalSoftwareUpgrade{}

func NewSoftwareUpgradeProposal(name string, resources []Resource,
	upgradeTime int64, oldChainId, newChainId, rollBackMemo string,
	maxEnrollmentDuration int64, upgradeMemo string,
	instateUpgrade, rebootRequired, skipHandler bool,
) *ProposalSoftwareUpgrade {
	return &ProposalSoftwareUpgrade{
		Name:                 name,
		Resources:            resources,
		UpgradeTime:          upgradeTime,
		OldChainId:           oldChainId,
		NewChainId:           newChainId,
		RollbackChecksum:     rollBackMemo,
		MaxEnrolmentDuration: maxEnrollmentDuration,
		Memo:                 upgradeMemo,
		InstateUpgrade:       instateUpgrade,
		RebootRequired:       rebootRequired,
		SkipHandler:          skipHandler,
	}
}

func (m *ProposalSoftwareUpgrade) ProposalType() string {
	return ProposalTypeSoftwareUpgrade
}

func (m *ProposalSoftwareUpgrade) ProposalPermission() types.PermValue {
	return types.PermCreateSoftwareUpgradeProposal
}

func (m *ProposalSoftwareUpgrade) VotePermission() types.PermValue {
	return types.PermVoteSoftwareUpgradeProposal
}

// ValidateBasic returns basic validation
func (m *ProposalSoftwareUpgrade) ValidateBasic() error {
	return nil
}

var _ types.Content = &ProposalCancelSoftwareUpgrade{}

func NewCancelSoftwareUpgradeProposal(name string) *ProposalCancelSoftwareUpgrade {
	return &ProposalCancelSoftwareUpgrade{
		Name: name,
	}
}

func (m *ProposalCancelSoftwareUpgrade) ProposalType() string {
	return ProposalTypeCancelSoftwareUpgrade
}

func (m *ProposalCancelSoftwareUpgrade) ProposalPermission() types.PermValue {
	return types.PermCreateSoftwareUpgradeProposal
}
func (m *ProposalCancelSoftwareUpgrade) VotePermission() types.PermValue {
	return types.PermVoteSoftwareUpgradeProposal
}

// ValidateBasic returns basic validation
func (m *ProposalCancelSoftwareUpgrade) ValidateBasic() error {
	return nil
}
