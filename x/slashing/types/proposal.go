package types

import (
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const ProposalTypeResetWholeValidatorRank = "ResetWholeValidatorRank"

func NewResetWholeValidatorRankProposal(proposer sdk.AccAddress) *ProposalResetWholeValidatorRank {
	return &ProposalResetWholeValidatorRank{
		Proposer: proposer,
	}
}

func (m *ProposalResetWholeValidatorRank) ProposalType() string {
	return ProposalTypeResetWholeValidatorRank
}

func (m *ProposalResetWholeValidatorRank) ProposalPermission() types.PermValue {
	return types.PermCreateResetWholeValidatorRankProposal
}

func (m *ProposalResetWholeValidatorRank) VotePermission() types.PermValue {
	return types.PermVoteResetWholeValidatorRankProposal
}

// ValidateBasic returns basic validation
func (m *ProposalResetWholeValidatorRank) ValidateBasic() error {
	return nil
}
