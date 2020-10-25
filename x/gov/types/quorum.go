package types

import "math"

func IsQuorum(percentage, votes, totalVoters uint64) bool {
	necessaryApproval := uint64(math.Ceil(float64(totalVoters*percentage) / 100.0))
	return votes >= necessaryApproval
}
