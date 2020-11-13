package types

import (
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types"
)

// constants
const (
	AssignPermissionProposalType   = "AssignPermission"
	SetNetworkPropertyProposalType = "SetNetworkProperty"
)

var _ Content = &AssignPermissionProposal{}

// NewProposal creates a new proposal
func NewProposal(
	proposalID uint64,
	content Content,
	votingStartTime time.Time,
	votingEndTime time.Time,
	enactmentEndTime time.Time,
) (Proposal, error) {
	msg, ok := content.(proto.Message)
	if !ok {
		return Proposal{}, fmt.Errorf("%T does not implement proto.Message", content)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return Proposal{}, err
	}

	return Proposal{
		ProposalId:       proposalID,
		VotingStartTime:  votingStartTime,
		VotingEndTime:    votingEndTime,
		EnactmentEndTime: enactmentEndTime,
		Content:          any,
		Result:           Pending,
	}, nil
}

// GetContent returns the proposal Content
func (p Proposal) GetContent() Content {
	content, ok := p.Content.GetCachedValue().(Content)
	if !ok {
		return nil
	}
	return content
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (p Proposal) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var content Content
	return unpacker.UnpackAny(p.Content, &content)
}

// NewAssignPermissionProposal creates a new assign permission proposal
func NewAssignPermissionProposal(
	address types.AccAddress,
	permission PermValue,
) Content {
	return &AssignPermissionProposal{
		Address:    address,
		Permission: uint32(permission),
	}
}

// ProposalType returns proposal's type
func (m *AssignPermissionProposal) ProposalType() string {
	return AssignPermissionProposalType
}

// NewSetNetworkPropertyProposal creates a new set network property proposal
func NewSetNetworkPropertyProposal(
	property NetworkProperty,
	value uint64,
) Content {
	return &SetNetworkPropertyProposal{
		NetworkProperty: property,
		Value:           value,
	}
}

// ProposalType returns proposal's type
func (m *SetNetworkPropertyProposal) ProposalType() string {
	return SetNetworkPropertyProposalType
}

// VotePermissionByProposal returns required permission to vote on a proposal type
func VotePermissionByProposal(proposalType string) PermValue {
	switch proposalType {
	case SetNetworkPropertyProposalType:
		return PermVoteSetNetworkPropertyProposal
	case AssignPermissionProposalType:
		return PermVoteSetPermissionProposal
	default:
		return PermZero
	}
}
