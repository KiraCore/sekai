package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPermissions_IsBlacklisted(t *testing.T) {
	perms := NewPermissions(
		[]PermValue{},
		[]PermValue{PermClaimValidator},
	)

	require.True(t, perms.IsBlacklisted(PermClaimValidator))
	require.False(t, perms.IsBlacklisted(PermClaimGovernanceSeat))
}

func TestPermissions_IsWhitelisted(t *testing.T) {
	perms := NewPermissions([]PermValue{PermClaimGovernanceSeat}, nil)

	require.True(t, perms.IsWhitelisted(PermClaimGovernanceSeat))
	require.False(t, perms.IsWhitelisted(PermClaimValidator))
}

func TestPermissions_AddWhitelist(t *testing.T) {
	perms := NewPermissions(nil, nil)

	require.False(t, perms.IsWhitelisted(PermClaimGovernanceSeat))
	perms.AddToWhitelist(PermClaimGovernanceSeat)
	require.True(t, perms.IsWhitelisted(PermClaimGovernanceSeat))
}

func TestPermissions_AddBlacklist(t *testing.T) {
	perms := NewPermissions(nil, nil)

	require.False(t, perms.IsBlacklisted(PermClaimGovernanceSeat))
	perms.AddToBlacklist(PermClaimGovernanceSeat)
	require.True(t, perms.IsBlacklisted(PermClaimGovernanceSeat))
}
