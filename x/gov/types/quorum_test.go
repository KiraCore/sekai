package types

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestIsQuorum_Errors(t *testing.T) {
	tests := []struct {
		name                      string
		percentage, votes, voters uint64
		expectedErr               error
	}{
		{
			name:        "more votes than total voters",
			percentage:  66,
			votes:       11,
			voters:      10,
			expectedErr: fmt.Errorf("there is more votes than voters: 11 > 10"),
		},
		{
			name:        "invalid quorum",
			percentage:  101,
			votes:       7,
			voters:      10,
			expectedErr: fmt.Errorf("quorum cannot be bigger than 1.00: 1.010000000000000000"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := IsQuorum(sdk.NewDecWithPrec(int64(tt.percentage), 2), tt.votes, tt.voters)
			require.EqualError(t, err, tt.expectedErr.Error())
		})
	}
}

func TestIsQuorum(t *testing.T) {
	tests := []struct {
		name                      string
		percentage, votes, voters uint64
		reached                   bool
	}{
		{
			name:       "quorum not reached",
			percentage: 66,
			votes:      6,
			voters:     10,
			reached:    false,
		},
		{
			name:       "quorum reached",
			percentage: 66,
			votes:      7,
			voters:     10,
			reached:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			quorum, err := IsQuorum(sdk.NewDecWithPrec(int64(tt.percentage), 2), tt.votes, tt.voters)
			require.NoError(t, err)
			require.Equal(t, tt.reached, quorum)
		})
	}
}
