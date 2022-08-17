package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/types"
)

func TestVotes_Getters(t *testing.T) {
	proposalID := uint64(12345)
	addr := types.AccAddress("some addr")

	votes := Votes{
		NewVote(proposalID, addr, OptionYes, 0),
		NewVote(proposalID, addr, OptionYes, 0),
		NewVote(proposalID, addr, OptionYes, 0),
		NewVote(proposalID, addr, OptionYes, 0),
		NewVote(proposalID, addr, OptionYes, 0),
		NewVote(proposalID, addr, OptionNoWithVeto, 0),
		NewVote(proposalID, addr, OptionNoWithVeto, 0),
		NewVote(proposalID, addr, OptionNoWithVeto, 0),
		NewVote(proposalID, addr, OptionAbstain, 0),
		NewVote(proposalID, addr, OptionAbstain, 0),
		NewVote(proposalID, addr, OptionNo, 0),
	}

	calculatedVotes := CalculateVotes(votes, 3)
	require.Equal(t, uint64(len(votes)), calculatedVotes.TotalVotes())
	require.Equal(t, uint64(5), calculatedVotes.YesVotes())
	require.Equal(t, uint64(1), calculatedVotes.NoVotes())
	require.Equal(t, uint64(2), calculatedVotes.AbstainVotes())
	require.Equal(t, uint64(3), calculatedVotes.VetoVotes())
}

func TestCalculatedVotes_ProcessResult(t *testing.T) {
	proposalID := uint64(12345)
	addr := types.AccAddress("some addr")

	tests := []struct {
		name           string
		votes          Votes
		actorsWithVeto uint64
		result         VoteResult
	}{
		{
			name: "more than 50% Yes",
			votes: Votes{
				NewVote(proposalID, addr, OptionYes, 0),
				NewVote(proposalID, addr, OptionYes, 0),
				NewVote(proposalID, addr, OptionYes, 0),
				NewVote(proposalID, addr, OptionYes, 0),
				NewVote(proposalID, addr, OptionYes, 0),
				NewVote(proposalID, addr, OptionYes, 0),
				NewVote(proposalID, addr, OptionNoWithVeto, 0),
				NewVote(proposalID, addr, OptionNo, 0),
				NewVote(proposalID, addr, OptionAbstain, 0),
				NewVote(proposalID, addr, OptionAbstain, 0),
			},
			actorsWithVeto: 3,
			result:         Passed,
		},
		{
			name: "different votes than yes equal or more than 50 : equal 50%",
			votes: Votes{
				NewVote(proposalID, addr, OptionYes, 0),
				NewVote(proposalID, addr, OptionYes, 0),
				NewVote(proposalID, addr, OptionYes, 0),
				NewVote(proposalID, addr, OptionYes, 0),
				NewVote(proposalID, addr, OptionYes, 0),
				NewVote(proposalID, addr, OptionNoWithVeto, 0),
				NewVote(proposalID, addr, OptionNo, 0),
				NewVote(proposalID, addr, OptionNo, 0),
				NewVote(proposalID, addr, OptionAbstain, 0),
				NewVote(proposalID, addr, OptionAbstain, 0),
			},
			result: Rejected,
		},
		{
			name: "50% or more of actors with Veto reject by voting No With Veto",
			votes: Votes{
				NewVote(proposalID, addr, OptionYes, 0),
				NewVote(proposalID, addr, OptionYes, 0),
				NewVote(proposalID, addr, OptionYes, 0),
				NewVote(proposalID, addr, OptionYes, 0),
				NewVote(proposalID, addr, OptionYes, 0),
				NewVote(proposalID, addr, OptionNoWithVeto, 0),
				NewVote(proposalID, addr, OptionNoWithVeto, 0),
				NewVote(proposalID, addr, OptionNo, 0),
				NewVote(proposalID, addr, OptionAbstain, 0),
				NewVote(proposalID, addr, OptionAbstain, 0),
			},
			actorsWithVeto: 3,
			result:         RejectedWithVeto,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			calc := CalculateVotes(tt.votes, tt.actorsWithVeto)
			require.Equal(t, tt.result, calc.ProcessResult())
		})
	}
}
