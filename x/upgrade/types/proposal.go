package types

import (
	"github.com/KiraCore/sekai/x/gov/types"
)

const ProposalTypeSoftwareUpgrade = "SoftwareUpgrade"
const ProposalTypeCancelSoftwareUpgrade = "CancelSoftwareUpgrade"

func (m *ProposalSoftwareUpgrade) ProposalType() string {
	return ProposalTypeSoftwareUpgrade
}

func (m *ProposalSoftwareUpgrade) VotePermission() types.PermValue {
	return types.PermVoteSoftwareUpgradeProposal
}

func (m *ProposalCancelSoftwareUpgrade) ProposalType() string {
	return ProposalTypeCancelSoftwareUpgrade
}

func (m *ProposalCancelSoftwareUpgrade) VotePermission() types.PermValue {
	return types.PermVoteSoftwareUpgradeProposal
}
