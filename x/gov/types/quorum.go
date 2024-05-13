package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func IsQuorum(percentage sdk.Dec, votes, totalVoters uint64) (bool, error) {
	if votes > totalVoters {
		return false, fmt.Errorf("there is more votes than voters: %d > %d", votes, totalVoters)
	}

	if percentage.GT(sdk.OneDec()) {
		return false, fmt.Errorf("quorum cannot be bigger than 1.00: %d", percentage)
	}

	necessaryApproval := sdk.NewDec(int64(totalVoters)).Mul(percentage)
	return sdk.NewDec(int64(votes)).GTE(necessaryApproval), nil
}
