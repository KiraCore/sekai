package keeper_test

import (
	"testing"

	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestNewKeeper_SaveNetworkActor(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	addr := addrs[0]

	networkActor := types.NetworkActor{
		Address: addr,
	}

	app.CustomGovKeeper.SaveNetworkActor(ctx, networkActor)

	savedActor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, networkActor.Address)
	require.True(t, found)

	require.Equal(t, networkActor, savedActor)
}

func TestKeeper_GetNetworkActorByAddress_FailsIfItDoesNotExist(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	addr := addrs[0]

	_, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
	require.False(t, found)
}

func TestKeeper_AssignRoleToAddress(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	addr := addrs[0]

	actor := types.NewDefaultActor(addr)
	app.CustomGovKeeper.AssignRoleToActor(ctx, actor, types.RoleSudo)

	savedActor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
	require.True(t, found)
	require.True(t, savedActor.HasRole(types.RoleSudo))
}

func TestKeeper_AddPermissionToNetworkActor(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	addr := addrs[0]

	networkActor := types.NewNetworkActor(
		addr,
		nil,
		1,
		nil,
		types.NewPermissions(nil, nil),
		1,
	)

	app.CustomGovKeeper.SaveNetworkActor(ctx, networkActor)

	// We check he does not have permissions
	savedNetworkActor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
	require.True(t, found)
	require.False(t, savedNetworkActor.Permissions.IsWhitelisted(types.PermSetPermissions))

	// We add permissions and we save it again.
	err := app.CustomGovKeeper.AddWhitelistPermission(ctx, savedNetworkActor, types.PermSetPermissions)
	require.NoError(t, err)

	// And we check that now it has permissions
	savedNetworkActor, found = app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
	require.True(t, found)
	require.True(t, savedNetworkActor.Permissions.IsWhitelisted(types.PermSetPermissions))
}

func TestKeeper_RemoveWhitelistPermission(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))

	err := whitelistPermToMultipleAddrs(app, ctx, addrs, types.PermSetPermissions)
	require.NoError(t, err)

	iterator := app.CustomGovKeeper.GetNetworkActorsByWhitelistedPermission(ctx, types.PermSetPermissions)
	requireIteratorCount(t, iterator, 2)
	assertAddrsHaveWhitelistedPerm(t, app, ctx, addrs, types.PermSetPermissions)

	actor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addrs[0])
	require.True(t, found)
	err = app.CustomGovKeeper.RemoveWhitelistPermission(ctx, actor, types.PermSetPermissions)
	require.NoError(t, err)

	iterator = app.CustomGovKeeper.GetNetworkActorsByWhitelistedPermission(ctx, types.PermSetPermissions)
	requireIteratorCount(t, iterator, 1)

	assertAddrsDontHaveWhitelistedPerm(t, app, ctx, []sdk.AccAddress{addrs[0]}, types.PermSetPermissions)
	assertAddrsHaveWhitelistedPerm(t, app, ctx, []sdk.AccAddress{addrs[1]}, types.PermSetPermissions)
}

func TestKeeper_GetActorsByRole(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))

	addRoleToMultipleAddrs(app, ctx, addrs, types.RoleSudo)

	iterator := app.CustomGovKeeper.GetNetworkActorsByRole(ctx, types.RoleSudo)
	requireIteratorCount(t, iterator, 2)

	assertAddrsHaveRole(t, app, ctx, addrs, types.RoleSudo)
}

func TestKeeper_RemoveRole(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))

	addRoleToMultipleAddrs(app, ctx, addrs, types.RoleSudo)

	iterator := app.CustomGovKeeper.GetNetworkActorsByRole(ctx, types.RoleSudo)
	requireIteratorCount(t, iterator, 2)

	assertAddrsHaveRole(t, app, ctx, addrs, types.RoleSudo)

	actor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addrs[0])
	require.True(t, found)
	require.True(t, actor.HasRole(types.RoleSudo))

	app.CustomGovKeeper.RemoveRoleFromActor(ctx, actor, types.RoleSudo)

	actor, found = app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addrs[0])
	require.True(t, found)
	require.False(t, actor.HasRole(types.RoleSudo))

	iterator = app.CustomGovKeeper.GetNetworkActorsByRole(ctx, types.RoleSudo)
	requireIteratorCount(t, iterator, 1)

}

func TestKeeper_GetActorsByWhitelistedPerm(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))

	err := whitelistPermToMultipleAddrs(app, ctx, addrs, types.PermSetPermissions)
	require.NoError(t, err)

	iterator := app.CustomGovKeeper.GetNetworkActorsByWhitelistedPermission(ctx, types.PermSetPermissions)
	requireIteratorCount(t, iterator, 2)

	assertAddrsHaveWhitelistedPerm(t, app, ctx, addrs, types.PermSetPermissions)
}

func TestKeeper_GetNetworkActorsByAbsoluteWhitelistPermission(t *testing.T) {
	tests := []struct {
		name       string
		prepareApp func(app *simapp.SekaiApp, ctx sdk.Context) []types.NetworkActor
	}{
		{
			name: "some addresses whitelisted",
			prepareApp: func(app *simapp.SekaiApp, ctx sdk.Context) []types.NetworkActor {
				addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))

				err := whitelistPermToMultipleAddrs(app, ctx, addrs, types.PermSetPermissions)
				require.NoError(t, err)

				expectedActors := []types.NetworkActor{
					types.NewDefaultActor(addrs[0]),
					types.NewDefaultActor(addrs[1]),
				}

				return expectedActors
			},
		},
		{
			name: "some addresses whitelisted by role",
			prepareApp: func(app *simapp.SekaiApp, ctx sdk.Context) []types.NetworkActor {
				addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))

				// Create role
				app.CustomGovKeeper.CreateRole(ctx, types.Role(12345))
				err := app.CustomGovKeeper.WhitelistRolePermission(ctx, types.Role(12345), types.PermSetPermissions)
				require.NoError(t, err)

				for _, addr := range addrs {
					actor := types.NewDefaultActor(addr)
					app.CustomGovKeeper.AssignRoleToActor(ctx, actor, types.Role(12345))
				}

				expectedActors := []types.NetworkActor{
					types.NewDefaultActor(addrs[0]),
					types.NewDefaultActor(addrs[1]),
				}

				return expectedActors
			},
		},
		{
			name: "whitelisted address whitelisted by role and personal permission (case 1)",
			prepareApp: func(app *simapp.SekaiApp, ctx sdk.Context) []types.NetworkActor {
				addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))

				// Create role
				app.CustomGovKeeper.CreateRole(ctx, types.Role(12345))
				err := app.CustomGovKeeper.WhitelistRolePermission(ctx, types.Role(12345), types.PermSetPermissions)
				require.NoError(t, err)

				// We whitelist all by the role.
				for _, addr := range addrs {
					actor := types.NewDefaultActor(addr)
					app.CustomGovKeeper.AssignRoleToActor(ctx, actor, types.Role(12345))
				}

				err = app.CustomGovKeeper.AddWhitelistPermission(ctx, types.NewDefaultActor(addrs[0]), types.PermSetPermissions)
				require.NoError(t, err)

				expectedActors := []types.NetworkActor{
					types.NewDefaultActor(addrs[0]),
					types.NewDefaultActor(addrs[1]),
				}

				return expectedActors
			},
		},
		{
			name: "whitelisted address whitelisted by role and personal permission (case 2)",
			prepareApp: func(app *simapp.SekaiApp, ctx sdk.Context) []types.NetworkActor {
				addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))

				// Create role
				app.CustomGovKeeper.CreateRole(ctx, types.Role(12345))
				err := app.CustomGovKeeper.WhitelistRolePermission(ctx, types.Role(12345), types.PermSetPermissions)
				require.NoError(t, err)

				// We assign role to first actor.
				actor := types.NewDefaultActor(addrs[0])
				app.CustomGovKeeper.AssignRoleToActor(ctx, actor, types.Role(12345))

				err = whitelistPermToMultipleAddrs(app, ctx, addrs, types.PermSetPermissions)
				require.NoError(t, err)

				expectedActors := []types.NetworkActor{
					types.NewDefaultActor(addrs[0]),
					types.NewDefaultActor(addrs[1]),
				}

				return expectedActors
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			expectecActors := tt.prepareApp(app, ctx)
			savedPerms := app.CustomGovKeeper.GetNetworkActorsByAbsoluteWhitelistPermission(ctx, types.PermSetPermissions)

			require.Equal(t, len(expectecActors), len(savedPerms))
			for i, actor := range expectecActors {
				require.Equal(t, actor.Address, savedPerms[i].Address)
			}
		})
	}
}

func addRoleToMultipleAddrs(app *simapp.SekaiApp, ctx sdk.Context, addrs []sdk.AccAddress, role types.Role) {
	for _, addr := range addrs {
		app.CustomGovKeeper.AssignRoleToActor(ctx, types.NewDefaultActor(addr), role)
	}
}

func whitelistPermToMultipleAddrs(app *simapp.SekaiApp, ctx sdk.Context, addrs []sdk.AccAddress, permissions types.PermValue) error {
	for _, addr := range addrs {
		err := app.CustomGovKeeper.AddWhitelistPermission(ctx, types.NewDefaultActor(addr), permissions)
		if err != nil {
			return err
		}
	}

	return nil
}

func requireIteratorCount(t *testing.T, iterator sdk.Iterator, expectedCount int) {
	c := 0
	for ; iterator.Valid(); iterator.Next() {
		c++
	}

	require.Equal(t, expectedCount, c)
}

func assertAddrsHaveWhitelistedPerm(
	t *testing.T,
	app *simapp.SekaiApp,
	ctx sdk.Context,
	addrs []sdk.AccAddress,
	perm types.PermValue,
) {
	for _, addr := range addrs {
		actor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
		require.True(t, found)
		require.True(t, actor.Permissions.IsWhitelisted(perm))
	}
}

func assertAddrsHaveRole(
	t *testing.T,
	app *simapp.SekaiApp,
	ctx sdk.Context,
	addrs []sdk.AccAddress,
	role types.Role,
) {
	for _, addr := range addrs {
		actor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
		require.True(t, found)
		require.True(t, actor.HasRole(role))
	}
}

func assertAddrsDontHaveWhitelistedPerm(
	t *testing.T,
	app *simapp.SekaiApp,
	ctx sdk.Context,
	addrs []sdk.AccAddress,
	perm types.PermValue,
) {
	for _, addr := range addrs {
		actor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
		require.True(t, found)
		require.False(t, actor.Permissions.IsWhitelisted(perm))
	}
}
