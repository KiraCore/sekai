package keeper_test

import (
	"testing"

	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestKeeper_SaveGetPermissionsForRole(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	perm := types.NewPermissions(
		nil, []types.PermValue{types.PermClaimValidator},
	)

	app.CustomGovKeeper.SetPermissionsForRole(ctx, types.RoleSudo, perm)

	savedPerms, err := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleSudo)
	require.NoError(t, err)
	require.Equal(t, perm, savedPerms)
}

func TestKeeper_HasGenesisDefaultRoles(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	roleSudo, err := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleSudo)
	require.NoError(t, err)
	require.True(t, roleSudo.IsWhitelisted(types.PermSetPermissions))

	roleValidator, err := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
	require.NoError(t, err)
	require.True(t, roleValidator.IsWhitelisted(types.PermClaimValidator))
}

func TestKeeper_SetPermissionsOverwritesOldPerms(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	roleValidator, err := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
	require.NoError(t, err)
	require.True(t, roleValidator.IsWhitelisted(types.PermClaimValidator))
	require.False(t, roleValidator.IsWhitelisted(types.PermSetPermissions))

	// We set whitelisted PermSetPermissions and Blacklisted PermClaimValidator
	newPerms := types.NewPermissions([]types.PermValue{types.PermSetPermissions}, []types.PermValue{types.PermClaimValidator})
	app.CustomGovKeeper.SetPermissionsForRole(ctx, types.RoleValidator, newPerms)

	newRoleValidatorPerms, err := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
	require.NoError(t, err)
	require.True(t, newRoleValidatorPerms.IsWhitelisted(types.PermSetPermissions))
	require.False(t, newRoleValidatorPerms.IsWhitelisted(types.PermClaimValidator))
}

func TestKeeper_GetPermissionsForRole_ReturnsErrorWhenDoesNotExist(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	_, err := app.CustomGovKeeper.GetPermissionsForRole(ctx, 12345)
	require.Error(t, err)
}
