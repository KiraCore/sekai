package types

import (
	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewProposalCollectiveSendDonation(
	name string,
	address string,
	amounts sdk.Coins,
) *ProposalCollectiveSendDonation {
	return &ProposalCollectiveSendDonation{
		Name:    name,
		Address: address,
		Amounts: amounts,
	}
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

func NewProposalCollectiveUpdate(
	name, description string,
	status CollectiveStatus,
	depositWhitelist DepositWhitelist,
	ownersWhitelist OwnersWhitelist,
	weightedSpendingPool []WeightedSpendingPool,
	claimStart, claimPeriod, claimEnd uint64,
	voteQuorum sdk.Dec,
	votePeriod, voteEnactment uint64,
) *ProposalCollectiveUpdate {
	return &ProposalCollectiveUpdate{
		Name:             name,
		Description:      description,
		Status:           status,
		DepositWhitelist: depositWhitelist,
		OwnersWhitelist:  ownersWhitelist,
		SpendingPools:    weightedSpendingPool,
		ClaimStart:       claimStart,
		ClaimPeriod:      claimPeriod,
		ClaimEnd:         claimEnd,
		VoteQuorum:       voteQuorum,
		VotePeriod:       votePeriod,
		VoteEnactment:    voteEnactment,
	}
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

func NewProposalCollectiveRemove(name string) *ProposalCollectiveRemove {
	return &ProposalCollectiveRemove{
		Name: name,
	}
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
