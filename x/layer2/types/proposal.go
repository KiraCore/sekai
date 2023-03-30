package types

import (
	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/gov/types"
)

func (m *ProposalJoinDapp) ProposalType() string {
	return kiratypes.ProposalTypeJoinDapp
}

func (m *ProposalJoinDapp) ProposalPermission() types.PermValue {
	return types.PermZero
}

func (m *ProposalJoinDapp) VotePermission() types.PermValue {
	return types.PermZero
}

// ValidateBasic returns basic validation
func (m *ProposalJoinDapp) ValidateBasic() error {
	return nil
}

func (m *ProposalUpsertDapp) ProposalType() string {
	return kiratypes.ProposalTypeUpsertDapp
}

func (m *ProposalUpsertDapp) ProposalPermission() types.PermValue {
	return types.PermZero
}

func (m *ProposalUpsertDapp) VotePermission() types.PermValue {
	return types.PermZero
}

// ValidateBasic returns basic validation
func (m *ProposalUpsertDapp) ValidateBasic() error {
	return nil
}
