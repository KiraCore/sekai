package types

import (
	"github.com/KiraCore/sekai/x/gov/types"
)

const ProposalTypeSoftwareUpgrade = "SoftwareUpgrade"

func (m *ProposalSoftwareUpgrade) ProposalType() string {
	return ProposalTypeSoftwareUpgrade
}

func (m *ProposalSoftwareUpgrade) VotePermission() types.PermValue {
	return types.PermVoteSoftwareUpgradeProposal
}
