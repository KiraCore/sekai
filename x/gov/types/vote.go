package types

type Votes []Vote

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

func (v CalculatedVotes) TotalVotes() uint64 {
	return v.total
}

func (v CalculatedVotes) YesVotes() uint64 {
	return v.votes[OptionYes]
}

func (v CalculatedVotes) NoVotes() uint64 {
	return v.votes[OptionNo]
}

func (v CalculatedVotes) AbstainVotes() uint64 {
	return v.votes[OptionAbstain]
}

func (v CalculatedVotes) VetoVotes() uint64 {
	return v.votes[OptionNoWithVeto]
}
