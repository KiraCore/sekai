package types

import (
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// constants
const (
	AssignPermissionProposalType   = "AssignPermission"
	SetNetworkPropertyProposalType = "SetNetworkProperty"
	UpsertDataRegistryProposalType = "UpsertDataRegistry"
	SetPoorNetworkMsgsProposalType = "SetPoorNetworkMsgs"
	CreateRoleProposalType         = "CreateRoleProposal"
)

var _ Content = &AssignPermissionProposal{}

// NewProposal creates a new proposal
func NewProposal(
	proposalID uint64,
	content Content,
	submitTime time.Time,
	votingEndTime time.Time,
	enactmentEndTime time.Time,
) (Proposal, error) {
	msg, ok := content.(proto.Message)
	if !ok {
		return Proposal{}, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("%T does not implement proto.Message", content))
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return Proposal{}, err
	}

	return Proposal{
		ProposalId:       proposalID,
		SubmitTime:       submitTime,
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

// VotePermission returns permission to vote on this proposal
func (m *SetNetworkPropertyProposal) VotePermission() PermValue {
	return PermVoteSetNetworkPropertyProposal
}

func (m *AssignPermissionProposal) VotePermission() PermValue {
	return PermVoteSetPermissionProposal
}

func NewUpsertDataRegistryProposal(key, hash, reference, encoding string, size uint64) Content {
	return &UpsertDataRegistryProposal{
		Key:       key,
		Hash:      hash,
		Reference: reference,
		Encoding:  encoding,
		Size_:     size,
	}
}

func (m *UpsertDataRegistryProposal) ProposalType() string {
	return UpsertDataRegistryProposalType
}

func (m *UpsertDataRegistryProposal) VotePermission() PermValue {
	return PermVoteUpsertDataRegistryProposal
}

func NewSetPoorNetworkMessagesProposal(msgs []string) Content {
	return &SetPoorNetworkMessagesProposal{
		Messages: msgs,
	}
}

func (m *SetPoorNetworkMessagesProposal) ProposalType() string {
	return SetPoorNetworkMsgsProposalType
}

func (m *SetPoorNetworkMessagesProposal) VotePermission() PermValue {
	return PermVoteSetPoorNetworkMessagesProposal
}

func NewCreateRoleProposal(role Role) Content {
	return &CreateRoleProposal{
		Role: uint32(role),
	}
}

func (m *CreateRoleProposal) ProposalType() string {
	return CreateRoleProposalType
}

func (m *CreateRoleProposal) VotePermission() PermValue {
	return PermVoteCreateRoleProposal
}
