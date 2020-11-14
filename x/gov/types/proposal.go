package types

import (
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types"
)

const (
	AssignPermissionProposalType   = "AssignPermission"
	UpsertDataRegistryProposalType = "UpsertDataRegistry"
)

var _ Content = &AssignPermissionProposal{}

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

func NewAssignPermissionProposal(
	address types.AccAddress,
	permission PermValue,
) Content {
	return &AssignPermissionProposal{
		Address:    address,
		Permission: uint32(permission),
	}
}

func (m *AssignPermissionProposal) ProposalType() string {
	return AssignPermissionProposalType
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
