package types

import (
	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewUnjailValidatorProposal(proposer sdk.AccAddress, valAddr sdk.ValAddress, reference string) *ProposalUnjailValidator {
	return &ProposalUnjailValidator{
		Proposer:  proposer.String(),
		ValAddr:   valAddr.String(),
		Reference: reference,
	}
}

func (m *ProposalUnjailValidator) ProposalType() string {
	return kiratypes.ProposalTypeUnjailValidator
}

func (m *ProposalUnjailValidator) ProposalPermission() types.PermValue {
	return types.PermCreateUnjailValidatorProposal
}

func (m *ProposalUnjailValidator) VotePermission() types.PermValue {
	return types.PermVoteUnjailValidatorProposal
}

// ValidateBasic returns basic validation
func (m *ProposalUnjailValidator) ValidateBasic() error {
	return nil
}
