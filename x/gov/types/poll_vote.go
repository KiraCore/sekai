package types

type PollVotes []PollVote

type CalculatedPollVotes struct {
	votes          map[PollVoteOption]uint64
	actorsWithVeto uint64
	total          uint64
}

func CalculatePollVotes(votes PollVotes, actorsWithVeto uint64) CalculatedPollVotes {
	votesMap := make(map[PollVoteOption]uint64)
	for _, vote := range votes {
		votesMap[vote.Option]++
	}

	return CalculatedPollVotes{
		total:          uint64(len(votes)),
		actorsWithVeto: actorsWithVeto,
		votes:          votesMap,
	}
}

func (c CalculatedPollVotes) TotalVotes() uint64 {
	return c.total
}

func (c CalculatedPollVotes) AbstainVotes() uint64 {
	return c.votes[PollOptionAbstain]
}

func (c CalculatedPollVotes) VetoVotes() uint64 {
	return c.votes[PollOptionNoWithVeto]
}

func (c CalculatedPollVotes) ProcessResult() PollResult {
	if c.actorsWithVeto != 0 {
		percentageActorsWithVeto := (float32(c.votes[PollOptionNoWithVeto]) / float32(c.actorsWithVeto)) * 100
		if percentageActorsWithVeto >= 50 {
			return PollRejectedWithVeto
		}
	}

	yesPercentage := (float32(c.votes[PollOptionCustom]) / float32(c.total)) * 100
	if yesPercentage > 50.00 {
		return PollPassed
	}

	sumOtherThanYes := c.votes[PollOptionAbstain] + c.votes[PollOptionNoWithVeto]
	sumPercentage := (float32(sumOtherThanYes) / float32(c.total)) * 100

	if sumPercentage >= 50 {
		return PollRejected
	}

	return PollUnknown
}
