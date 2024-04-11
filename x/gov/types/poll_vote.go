package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type PollVotes []PollVote

type CalculatedPollVotes struct {
	votes          map[string]uint64
	actorsWithVeto uint64
	total          uint64
}

func CalculatePollVotes(votes PollVotes, actorsWithVeto uint64) CalculatedPollVotes {
	votesMap := make(map[string]uint64)
	for _, vote := range votes {
		var mapValue string

		if vote.Option == PollOptionCustom {
			mapValue = vote.CustomValue
		} else {
			mapValue = vote.Option.String()
		}

		votesMap[mapValue]++
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
	return c.votes[PollOptionAbstain.String()]
}

func (c CalculatedPollVotes) VetoVotes() uint64 {
	return c.votes[PollOptionNoWithVeto.String()]
}

func (c CalculatedPollVotes) ProcessResult(properties *NetworkProperties) PollResult {
	if c.actorsWithVeto != 0 {
		percentageActorsWithVeto := sdk.NewDec(int64(c.votes[PollOptionNoWithVeto.String()])).Quo(sdk.NewDec(int64(c.actorsWithVeto)))

		if properties.VetoThreshold.LTE(percentageActorsWithVeto) {
			return PollRejectedWithVeto
		}
	}

	for _, count := range c.votes {
		yesPercentage := (float32(count) / float32(c.total)) * 100
		if yesPercentage > 50.00 {
			return PollPassed
		}
	}

	var highestVoteCount uint64 = 0
	var highestList []string

	for i, votesCount := range c.votes {
		if votesCount > highestVoteCount {
			highestVoteCount = votesCount
			highestList = []string{i}
		}

		if votesCount == highestVoteCount {
			highestList = append(highestList, i)
		}
	}

	if len(highestList) >= 2 {
		return PollRejected
	}

	sumOtherThanYes := c.votes[PollOptionAbstain.String()] + c.votes[PollOptionNoWithVeto.String()]
	sumPercentage := (float32(sumOtherThanYes) / float32(c.total)) * 100

	if sumPercentage >= 50 {
		return PollRejected
	}

	return PollUnknown
}
