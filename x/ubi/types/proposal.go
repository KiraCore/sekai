package types

import (
	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/gov/types"
)

var _ types.Content = &UpsertUBIProposal{}

func NewUpsertUBIProposal(
	name string,
	distrStart uint64,
	distrEnd uint64,
	amount uint64,
	period uint64,
	poolName string,
) *UpsertUBIProposal {
	return &UpsertUBIProposal{
		Name:              name,
		DistributionStart: distrStart,
		DistributionEnd:   distrEnd,
		Amount:            amount,
		Period:            period,
		Pool:              poolName,
	}
}

func (m *UpsertUBIProposal) ProposalType() string {
	return kiratypes.ProposalTypeUpsertUBI
}

func (m *UpsertUBIProposal) ProposalPermission() types.PermValue {
	return types.PermCreateUpsertUBIProposal
}

func (m *UpsertUBIProposal) VotePermission() types.PermValue {
	return types.PermVoteUpsertUBIProposal
}

func (m *UpsertUBIProposal) ValidateBasic() error {
	return nil
}

var _ types.Content = &RemoveUBIProposal{}

func NewRemoveUBIProposal(
	name string,
) *UpsertUBIProposal {
	return &UpsertUBIProposal{
		Name: name,
	}
}

func (m *RemoveUBIProposal) ProposalType() string {
	return kiratypes.ProposalTypeRemoveUBI
}

func (m *RemoveUBIProposal) ProposalPermission() types.PermValue {
	return types.PermCreateRemoveUBIProposal
}

func (m *RemoveUBIProposal) VotePermission() types.PermValue {
	return types.PermVoteRemoveUBIProposal
}

func (m *RemoveUBIProposal) ValidateBasic() error {
	return nil
}
