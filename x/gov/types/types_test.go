package types_test

import (
	"os"
	"testing"

	customgovtypes "github.com/KiraCore/sekai/x/gov/types"

	"github.com/KiraCore/sekai/app"

	"github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	app.SetConfig()
	os.Exit(m.Run())
}

func TestPermissions_IsBlacklisted(t *testing.T) {
	perms := customgovtypes.NewPermissions(
		[]customgovtypes.PermValue{},
		[]customgovtypes.PermValue{customgovtypes.PermClaimValidator},
	)

	require.True(t, perms.IsBlacklisted(customgovtypes.PermClaimValidator))
	require.False(t, perms.IsBlacklisted(customgovtypes.PermSetPermissions))
}

func TestPermissions_IsWhitelisted(t *testing.T) {
	perms := customgovtypes.NewPermissions([]customgovtypes.PermValue{customgovtypes.PermClaimValidator}, nil)

	require.True(t, perms.IsWhitelisted(customgovtypes.PermClaimValidator))
	require.False(t, perms.IsWhitelisted(customgovtypes.PermSetPermissions))
}

func TestPermissions_AddWhitelist(t *testing.T) {
	perms := customgovtypes.NewPermissions(nil, nil)

	require.False(t, perms.IsWhitelisted(customgovtypes.PermClaimValidator))

	err := perms.AddToWhitelist(customgovtypes.PermSetPermissions)
	require.NoError(t, err)
	require.True(t, perms.IsWhitelisted(customgovtypes.PermSetPermissions))

	// Add to whitelist value blacklisted gives error
	err = perms.AddToBlacklist(customgovtypes.PermClaimValidator)
	require.NoError(t, err)

	err = perms.AddToWhitelist(customgovtypes.PermClaimValidator)
	require.EqualError(t, err, "permission is already blacklisted")
}

func TestPermissions_AddBlacklist(t *testing.T) {
	perms := customgovtypes.NewPermissions(nil, nil)

	require.False(t, perms.IsBlacklisted(customgovtypes.PermSetPermissions))
	err := perms.AddToBlacklist(customgovtypes.PermSetPermissions)
	require.NoError(t, err)
	require.True(t, perms.IsBlacklisted(customgovtypes.PermSetPermissions))

	// Add to blacklist when is whitelisted gives error
	err = perms.AddToWhitelist(customgovtypes.PermClaimValidator)
	require.NoError(t, err)

	err = perms.AddToBlacklist(customgovtypes.PermClaimValidator)
	require.EqualError(t, err, "permission is already whitelisted")
}

func TestPermissions_RemoveFromWhitelist(t *testing.T) {
	perms := customgovtypes.NewPermissions([]customgovtypes.PermValue{
		customgovtypes.PermSetPermissions,
	}, nil)

	// It fails if permission is not whitelisted.
	err := perms.RemoveFromWhitelist(customgovtypes.PermClaimCouncilor)
	require.EqualError(t, err, "permission is not whitelisted")

	err = perms.RemoveFromWhitelist(customgovtypes.PermSetPermissions)
	require.NoError(t, err)

	require.False(t, perms.IsWhitelisted(customgovtypes.PermSetPermissions))
}

//
// NETWORK ACTOR
//
func TestNewNetworkActor_SetRole(t *testing.T) {
	addr, err := types.AccAddressFromBech32("kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r")
	require.NoError(t, err)

	actor := customgovtypes.NewDefaultActor(addr)
	require.False(t, actor.HasRole(customgovtypes.RoleValidator))

	actor.SetRole(customgovtypes.RoleValidator)

	require.True(t, actor.HasRole(customgovtypes.RoleValidator))
}
