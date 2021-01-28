package types

import (
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewProposalUnjailValidator(proposer sdk.AccAddress) *ProposalUnjailValidator {
	return &ProposalUnjailValidator{Proposer: proposer}
}

func (m *ProposalUnjailValidator) ProposalType() string {
	panic("implement me")
}

func (m *ProposalUnjailValidator) VotePermission() types.PermValue {
	panic("implement me")
}
