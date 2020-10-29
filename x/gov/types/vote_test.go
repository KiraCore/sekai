package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/types"
)

func TestVotes_GetResult(t *testing.T) {
	proposalID := uint64(12345)
	addr := types.AccAddress("some addr")

	votes := Votes{
		NewVote(proposalID, addr, OptionYes),
		NewVote(proposalID, addr, OptionYes),
		NewVote(proposalID, addr, OptionYes),
		NewVote(proposalID, addr, OptionYes),
		NewVote(proposalID, addr, OptionYes),
		NewVote(proposalID, addr, OptionNoWithVeto),
		NewVote(proposalID, addr, OptionNoWithVeto),
		NewVote(proposalID, addr, OptionNoWithVeto),
		NewVote(proposalID, addr, OptionAbstain),
		NewVote(proposalID, addr, OptionAbstain),
		NewVote(proposalID, addr, OptionNo),
	}

	calculatedVotes := CalculateVotes(votes)
	require.Equal(t, uint64(len(votes)), calculatedVotes.TotalVotes())
	require.Equal(t, uint64(5), calculatedVotes.YesVotes())
	require.Equal(t, uint64(1), calculatedVotes.NoVotes())
	require.Equal(t, uint64(2), calculatedVotes.AbstainVotes())
	require.Equal(t, uint64(3), calculatedVotes.VetoVotes())
}
