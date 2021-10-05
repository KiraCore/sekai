package types_test

import (
	"os"
	"testing"

	"github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	app.SetConfig()
	os.Exit(m.Run())
}

func TestPermissions_IsBlacklisted(t *testing.T) {
	perms := types.NewPermissions(
		[]types.PermValue{},
		[]types.PermValue{types.PermClaimValidator},
	)

	require.True(t, perms.IsBlacklisted(types.PermClaimValidator))
	require.False(t, perms.IsBlacklisted(types.PermSetPermissions))
}

func TestPermissions_IsWhitelisted(t *testing.T) {
	perms := types.NewPermissions([]types.PermValue{types.PermClaimValidator}, nil)

	require.True(t, perms.IsWhitelisted(types.PermClaimValidator))
	require.False(t, perms.IsWhitelisted(types.PermSetPermissions))
}

func TestPermissions_AddWhitelist(t *testing.T) {
	perms := types.NewPermissions(nil, nil)

	require.False(t, perms.IsWhitelisted(types.PermClaimValidator))

	err := perms.AddToWhitelist(types.PermSetPermissions)
	require.NoError(t, err)
	require.True(t, perms.IsWhitelisted(types.PermSetPermissions))

	// Add to whitelist value blacklisted gives error
	err = perms.AddToBlacklist(types.PermClaimValidator)
	require.NoError(t, err)

	err = perms.AddToWhitelist(types.PermClaimValidator)
	require.EqualError(t, err, "permission is already blacklisted")
}

func TestPermissions_AddBlacklist(t *testing.T) {
	perms := types.NewPermissions(nil, nil)

	require.False(t, perms.IsBlacklisted(types.PermSetPermissions))
	err := perms.AddToBlacklist(types.PermSetPermissions)
	require.NoError(t, err)
	require.True(t, perms.IsBlacklisted(types.PermSetPermissions))

	// Add to blacklist when is whitelisted gives error
	err = perms.AddToWhitelist(types.PermClaimValidator)
	require.NoError(t, err)

	err = perms.AddToBlacklist(types.PermClaimValidator)
	require.EqualError(t, err, "permission is already whitelisted")
}

func TestPermissions_RemoveFromWhitelist(t *testing.T) {
	perms := types.NewPermissions([]types.PermValue{
		types.PermSetPermissions,
	}, nil)

	// It fails if permission is not whitelisted.
	err := perms.RemoveFromWhitelist(types.PermClaimCouncilor)
	require.EqualError(t, err, "permission is not whitelisted")

	err = perms.RemoveFromWhitelist(types.PermSetPermissions)
	require.NoError(t, err)

	require.False(t, perms.IsWhitelisted(types.PermSetPermissions))
}

func TestPermissions_RemoveFromBlacklist(t *testing.T) {
	perms := types.NewPermissions(nil,
		[]types.PermValue{
			types.PermSetPermissions,
		},
	)

	// It fails if permission is not blacklisted.
	err := perms.RemoveFromBlacklist(types.PermClaimCouncilor)
	require.EqualError(t, err, "permission is not blacklisted")

	err = perms.RemoveFromBlacklist(types.PermSetPermissions)
	require.NoError(t, err)

	require.False(t, perms.IsBlacklisted(types.PermSetPermissions))
}
