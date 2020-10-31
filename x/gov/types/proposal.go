package types

import (
	"time"

	"github.com/cosmos/cosmos-sdk/types"
)

func NewProposalAssignPermission(
	proposalID uint64,
	address types.AccAddress,
	permission PermValue,
	votingStartTime time.Time,
	votingEndTime time.Time,
	enactmentEndTime time.Time,
) ProposalAssignPermission {
	return ProposalAssignPermission{
		ProposalId:      proposalID,
		Address:         address,
		Permission:      uint32(permission),
		VotingStartTime: votingStartTime,
		VotingEndTime:   votingEndTime,
		EnactmentEndTime: enactmentEndTime,
		Result:          Pending,
	}
}
