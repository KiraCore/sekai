package types

import (
	"github.com/KiraCore/sekai/x/gov/types"
)

const ProposalTypeSoftwareUpgrade = "SoftwareUpgrade"
const ProposalTypeCancelSoftwareUpgrade = "CancelSoftwareUpgrade"

var _ types.Content = &ProposalSoftwareUpgrade{}

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
