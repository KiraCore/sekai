package types

import (
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types"
)

const AssignPermissionProposalType = "AssignPermission"

var _ Content = &AssignPermissionProposal{}

func NewAssignPermissionProposal(
	proposalID uint64,
	address types.AccAddress,
	permission PermValue,
	votingStartTime time.Time,
	votingEndTime time.Time,
	enactmentEndTime time.Time,
) (Proposal, error) {
	var content Content = &AssignPermissionProposal{
		Address:    address,
		Permission: uint32(permission),
	}

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
		Result: Pending,
	}, nil
}

func (m *AssignPermissionProposal) ProposalType() string {
	return AssignPermissionProposalType
}

func NewProposalAssignPermission(
	proposalID uint64,
	address types.AccAddress,
	permission PermValue,
	votingStartTime time.Time,
	votingEndTime time.Time,
	enactmentEndTime time.Time,
) ProposalAssignPermission {
	return ProposalAssignPermission{
		ProposalId:       proposalID,
		Address:          address,
		Permission:       uint32(permission),
		VotingStartTime:  votingStartTime,
		VotingEndTime:    votingEndTime,
		EnactmentEndTime: enactmentEndTime,
		Result:           Pending,
	}
}
