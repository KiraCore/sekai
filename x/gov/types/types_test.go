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
	require.False(t, perms.IsBlacklisted(PermSetPermissions))
}

func TestPermissions_IsWhitelisted(t *testing.T) {
	perms := NewPermissions([]PermValue{PermClaimValidator}, nil)

	require.True(t, perms.IsWhitelisted(PermClaimValidator))
	require.False(t, perms.IsWhitelisted(PermSetPermissions))
}

func TestPermissions_AddWhitelist(t *testing.T) {
	perms := NewPermissions(nil, nil)

	require.False(t, perms.IsWhitelisted(PermClaimValidator))

	err := perms.AddToWhitelist(PermSetPermissions)
	require.NoError(t, err)
	require.True(t, perms.IsWhitelisted(PermSetPermissions))

	// Add to whitelist value blacklisted gives error
	err = perms.AddToBlacklist(PermClaimValidator)
	require.NoError(t, err)

	err = perms.AddToWhitelist(PermClaimValidator)
	require.EqualError(t, err, "permission is already blacklisted")
}

func TestPermissions_AddBlacklist(t *testing.T) {
	perms := NewPermissions(nil, nil)

	require.False(t, perms.IsBlacklisted(PermSetPermissions))
	err := perms.AddToBlacklist(PermSetPermissions)
	require.NoError(t, err)
	require.True(t, perms.IsBlacklisted(PermSetPermissions))

	// Add to blacklist when is whitelisted gives error
	err = perms.AddToWhitelist(PermClaimValidator)
	require.NoError(t, err)

	err = perms.AddToBlacklist(PermClaimValidator)
	require.EqualError(t, err, "permission is already whitelisted")
}
