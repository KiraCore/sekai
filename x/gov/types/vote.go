package types

import "github.com/cosmos/cosmos-sdk/types"

type Votes []Vote

func NewVote(proposalID uint64, addr types.AccAddress, option VoteOption) Vote {
	return Vote{
		ProposalId: proposalID,
		Voter:      addr,
		Option:     option,
	}
}

type CalculatedVotes struct {
	votes          map[VoteOption]uint64
	actorsWithVeto uint64
	total          uint64
}

func CalculateVotes(votes Votes, actorsWithVeto uint64) CalculatedVotes {
	votesMap := make(map[VoteOption]uint64)
	for _, vote := range votes {
		votesMap[vote.Option]++
	}

	return CalculatedVotes{
		total:          uint64(len(votes)),
		actorsWithVeto: actorsWithVeto,
		votes:          votesMap,
	}
}

func (c CalculatedVotes) TotalVotes() uint64 {
	return c.total
}

func (c CalculatedVotes) YesVotes() uint64 {
	return c.votes[OptionYes]
}

func (c CalculatedVotes) NoVotes() uint64 {
	return c.votes[OptionNo]
}

func (c CalculatedVotes) AbstainVotes() uint64 {
	return c.votes[OptionAbstain]
}

func (c CalculatedVotes) VetoVotes() uint64 {
	return c.votes[OptionNoWithVeto]
}

func (c CalculatedVotes) ProcessResult() VoteResult {
	if c.actorsWithVeto != 0 {
		percentageActorsWithVeto := (float32(c.votes[OptionNoWithVeto]) / float32(c.actorsWithVeto)) * 100
		if percentageActorsWithVeto >= 50 {
			return RejectedWithVeto
		}
	}

	yesPercentage := (float32(c.votes[OptionYes]) / float32(c.total)) * 100
	if yesPercentage > 50.00 {
		return Passed
	}

	sumOtherThanYes := c.votes[OptionNo] + c.votes[OptionAbstain] + c.votes[OptionNoWithVeto]
	sumPercentage := (float32(sumOtherThanYes) / float32(c.total)) * 100

	if sumPercentage >= 50 {
		return Rejected
	}

	return Unknown
}
