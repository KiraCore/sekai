package types

import (
	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/gov/types"
)

func NewProposalCollectiveSendDonation() *ProposalCollectiveSendDonation {
	return &ProposalCollectiveSendDonation{}
}

func (m *ProposalCollectiveSendDonation) ProposalType() string {
	return kiratypes.ProposalTypeCollectiveSendDonation
}

func (m *ProposalCollectiveSendDonation) ProposalPermission() types.PermValue {
	return types.PermZero
}

func (m *ProposalCollectiveSendDonation) VotePermission() types.PermValue {
	return types.PermZero
}

// ValidateBasic returns basic validation
func (m *ProposalCollectiveSendDonation) ValidateBasic() error {
	return nil
}

func NewProposalCollectiveUpdate() *ProposalCollectiveUpdate {
	return &ProposalCollectiveUpdate{}
}

func (m *ProposalCollectiveUpdate) ProposalType() string {
	return kiratypes.ProposalTypeCollectiveUpdate
}

func (m *ProposalCollectiveUpdate) ProposalPermission() types.PermValue {
	return types.PermZero
}

func (m *ProposalCollectiveUpdate) VotePermission() types.PermValue {
	return types.PermZero
}

// ValidateBasic returns basic validation
func (m *ProposalCollectiveUpdate) ValidateBasic() error {
	return nil
}

func NewProposalCollectiveRemove() *ProposalCollectiveRemove {
	return &ProposalCollectiveRemove{}
}

func (m *ProposalCollectiveRemove) ProposalType() string {
	return kiratypes.ProposalTypeCollectiveRemove
}

func (m *ProposalCollectiveRemove) ProposalPermission() types.PermValue {
	return types.PermZero
}

func (m *ProposalCollectiveRemove) VotePermission() types.PermValue {
	return types.PermZero
}

// ValidateBasic returns basic validation
func (m *ProposalCollectiveRemove) ValidateBasic() error {
	return nil
}
