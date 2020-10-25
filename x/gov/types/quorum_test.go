package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
			require.Equal(t, tt.reached, IsQuorum(tt.percentage, tt.votes, tt.voters))
		})
	}
}
