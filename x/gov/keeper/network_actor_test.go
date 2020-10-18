package keeper_test

import (
	"testing"

	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestNewKeeper_SaveNetworkActor(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, types2.TokensFromConsensusPower(10))
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

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, types2.TokensFromConsensusPower(10))
	addr := addrs[0]

	_, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
	require.False(t, found)
}

func TestKeeper_AssignRoleToAddress(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, types2.TokensFromConsensusPower(10))
	addr := addrs[0]

	actor := types.NewDefaultActor(addr)
	app.CustomGovKeeper.AssignRoleToAddress(ctx, actor, types.RoleSudo)

	savedActor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
	require.True(t, found)
	require.True(t, savedActor.HasRole(types.RoleSudo))
}

func TestKeeper_AddPermissionToNetworkActor(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, types2.TokensFromConsensusPower(10))
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

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, types2.TokensFromConsensusPower(10))

	err := whitelistPermToMultipleAddrs(app, ctx, addrs, types.PermSetPermissions)
	require.NoError(t, err)

	iterator := app.CustomGovKeeper.GetNetworkActorByWhitelistedPermission(ctx, types.PermSetPermissions)
	requireIteratorCount(t, iterator, 2)
	assertAddrsHaveWhitelistedPerm(t, app, ctx, addrs, types.PermSetPermissions)

	actor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addrs[0])
	require.True(t, found)
	err = app.CustomGovKeeper.RemoveWhitelistPermission(ctx, actor, types.PermSetPermissions)
	require.NoError(t, err)

	iterator = app.CustomGovKeeper.GetNetworkActorByWhitelistedPermission(ctx, types.PermSetPermissions)
	requireIteratorCount(t, iterator, 1)

	assertAddrsDontHaveWhitelistedPerm(t, app, ctx, []types2.AccAddress{addrs[0]}, types.PermSetPermissions)
	assertAddrsHaveWhitelistedPerm(t, app, ctx, []types2.AccAddress{addrs[1]}, types.PermSetPermissions)
}

func TestKeeper_GetActorsByRole(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, types2.TokensFromConsensusPower(10))

	addRoleToMultipleAddrs(app, ctx, addrs, types.RoleSudo)

	iterator := app.CustomGovKeeper.GetNetworkActorsByRole(ctx, types.RoleSudo)
	requireIteratorCount(t, iterator, 2)

	assertAddrsHaveRole(t, app, ctx, addrs, types.RoleSudo)
}

func TestKeeper_RemoveRole(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, types2.TokensFromConsensusPower(10))

	addRoleToMultipleAddrs(app, ctx, addrs, types.RoleSudo)

	iterator := app.CustomGovKeeper.GetNetworkActorsByRole(ctx, types.RoleSudo)
	requireIteratorCount(t, iterator, 2)

	assertAddrsHaveRole(t, app, ctx, addrs, types.RoleSudo)
}

func TestKeeper_GetActorsByWhitelistedPerm(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, types2.TokensFromConsensusPower(10))

	err := whitelistPermToMultipleAddrs(app, ctx, addrs, types.PermSetPermissions)
	require.NoError(t, err)

	iterator := app.CustomGovKeeper.GetNetworkActorByWhitelistedPermission(ctx, types.PermSetPermissions)
	requireIteratorCount(t, iterator, 2)

	assertAddrsHaveWhitelistedPerm(t, app, ctx, addrs, types.PermSetPermissions)
}

func addRoleToMultipleAddrs(app *simapp.SimApp, ctx types2.Context, addrs []types2.AccAddress, role types.Role) {
	for _, addr := range addrs {
		app.CustomGovKeeper.AssignRoleToAddress(ctx, types.NewDefaultActor(addr), role)
	}
}

func whitelistPermToMultipleAddrs(app *simapp.SimApp, ctx types2.Context, addrs []types2.AccAddress, permissions types.PermValue) error {
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
	app *simapp.SimApp,
	ctx types2.Context,
	addrs []types2.AccAddress,
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
	app *simapp.SimApp,
	ctx types2.Context,
	addrs []types2.AccAddress,
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
	app *simapp.SimApp,
	ctx types2.Context,
	addrs []types2.AccAddress,
	perm types.PermValue,
) {
	for _, addr := range addrs {
		actor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
		require.True(t, found)
		require.False(t, actor.Permissions.IsWhitelisted(perm))
	}
}
