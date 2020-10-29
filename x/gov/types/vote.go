package types

type Votes []Vote
type VoteResult int64

const (
	Unknown VoteResult = iota
	Passed
	Rejected
	RejectedWithVeto
)

type CalculatedVotes struct {
	votes map[VoteOption]uint64
	total uint64
}

func CalculateVotes(votes Votes) CalculatedVotes {
	votesMap := make(map[VoteOption]uint64)
	for _, vote := range votes {
		votesMap[vote.Option]++
	}

	return CalculatedVotes{
		total: uint64(len(votes)),
		votes: votesMap,
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
