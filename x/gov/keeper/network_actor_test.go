package keeper_test

import (
	"testing"

	types3 "github.com/KiraCore/sekai/x/staking/types"

	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/x/gov/types"
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
	err := app.CustomGovKeeper.AddWhitelistPermission(ctx, addr, types.PermSetPermissions)
	require.NoError(t, err)

	// And we check that now it has permissions
	savedNetworkActor, found = app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
	require.True(t, found)
	require.True(t, savedNetworkActor.Permissions.IsWhitelisted(types.PermSetPermissions))
}

func TestKeeper_AddWhitelistPermission_Error(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, types2.TokensFromConsensusPower(10))
	addr := addrs[0]

	err := app.CustomGovKeeper.AddWhitelistPermission(ctx, addr, types.PermClaimCouncilor)
	require.EqualError(t, err, types3.ErrNetworkActorNotFound.Error())
}

func TestKeeper_GetActorsByWhitelistedPerm(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, types2.TokensFromConsensusPower(10))
	addr1 := addrs[0]
	addr2 := addrs[1]

	actor1 := types.NewDefaultActor(addr1)
	app.CustomGovKeeper.SaveNetworkActor(ctx, actor1)
	actor2 := types.NewDefaultActor(addr2)
	app.CustomGovKeeper.SaveNetworkActor(ctx, actor2)

	err := app.CustomGovKeeper.AddWhitelistPermission(ctx, addr1, types.PermSetPermissions)
	require.NoError(t, err)

	err = app.CustomGovKeeper.AddWhitelistPermission(ctx, addr2, types.PermSetPermissions)
	require.NoError(t, err)
}
