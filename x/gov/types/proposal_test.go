package types

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestNewProposalAssignPermission_ProposalIsPendingByDefault_Deprecated(t *testing.T) { // TODO Remove this type
	proposal := NewProposalAssignPermission(
		1234,
		types.AccAddress{0x12},
		PermSetPermissions,
		time.Now(),
		time.Now(),
		time.Now(),
	)

	require.Equal(t, Pending, proposal.Result)
}

func TestNewProposalAssignPermission_ProposalIsPendingByDefault(t *testing.T) {
	proposal, err := NewAssignPermissionProposal(
		1234,
		types.AccAddress{0x12},
		PermSetPermissions,
		time.Now(),
		time.Now(),
		time.Now(),
	)

	require.NoError(t, err)
	require.Equal(t, Pending, proposal.Result)
}
