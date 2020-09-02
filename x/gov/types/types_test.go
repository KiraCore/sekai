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

	err := perms.AddToWhitelist(PermClaimGovernanceSeat)
	require.NoError(t, err)
	require.True(t, perms.IsWhitelisted(PermClaimGovernanceSeat))

	// Add to whitelist value blacklisted gives error
	perms.AddToBlacklist(PermClaimValidator)

	err = perms.AddToWhitelist(PermClaimValidator)
	require.EqualError(t, err, "permission is already blacklisted")
}

func TestPermissions_AddBlacklist(t *testing.T) {
	perms := NewPermissions(nil, nil)

	require.False(t, perms.IsBlacklisted(PermClaimGovernanceSeat))
	err := perms.AddToBlacklist(PermClaimGovernanceSeat)
	require.NoError(t, err)
	require.True(t, perms.IsBlacklisted(PermClaimGovernanceSeat))

	// Add to blacklist when is whitelisted gives error
	err = perms.AddToWhitelist(PermClaimValidator)
	require.NoError(t, err)

	err = perms.AddToBlacklist(PermClaimValidator)
	require.EqualError(t, err, "permission is already whitelisted")
}
