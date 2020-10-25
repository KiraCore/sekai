package types

import (
	"fmt"
	"math"
)

func IsQuorum(percentage, votes, totalVoters uint64) (bool, error) {
	if votes > totalVoters {
		return false, fmt.Errorf("there is more votes than voters")
	}

	if percentage > 100 {
		return false, fmt.Errorf("quorum cannot be bigger than 100")
	}

	necessaryApproval := uint64(math.Ceil(float64(totalVoters*percentage) / 100.0))
	return votes >= necessaryApproval, nil
}
