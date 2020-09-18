package gov_test

import (
	"testing"

	"github.com/KiraCore/sekai/x/gov"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/x/gov/types"
)

func TestQuerier_PermissionsByAddress(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.TokensFromConsensusPower(10))
	addr1 := addrs[0]
	addr2 := addrs[1]

	permissions := types.NewPermissions(
		[]types.PermValue{
			types.PermClaimValidator,
		},
		nil,
	)

	networkActor := types.NewNetworkActor(
		addr1,
		types.Roles{},
		1,
		nil,
		permissions,
		123,
	)

	app.CustomGovKeeper.SaveNetworkActor(ctx, networkActor)

	querier := gov.NewQuerier(app.CustomGovKeeper)

	resp, err := querier.PermissionsByAddress(sdk.WrapSDKContext(ctx), &types.PermissionsByAddressRequest{ValAddr: addr1})
	require.NoError(t, err)

	require.Equal(t, permissions, resp.Permissions)

	// Get permissions by address that is not saved.
	_, err = querier.PermissionsByAddress(sdk.WrapSDKContext(ctx), &types.PermissionsByAddressRequest{ValAddr: addr2})
	require.EqualError(t, err, "network actor not found: key not found")
}

func TestQuerier_CouncilorByAddress(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.TokensFromConsensusPower(10))
	addr1 := addrs[0]
	addr2 := addrs[1]

	councilor := types.NewCouncilor(
		"TheMoniker",
		"TheWebsite",
		"TheSocial",
		"TheIdentity",
		addr1,
	)

	app.CustomGovKeeper.SaveCouncilor(ctx, councilor)

	querier := gov.NewQuerier(app.CustomGovKeeper)

	resp, err := querier.CouncilorByAddress(
		sdk.WrapSDKContext(ctx),
		&types.CouncilorByAddressRequest{ValAddr: addr1},
	)
	require.NoError(t, err)
	require.Equal(t, councilor, resp.Councilor)

	// Councilor by Moniker
	resp, err = querier.CouncilorByMoniker(
		sdk.WrapSDKContext(ctx),
		&types.CouncilorByMonikerRequest{
			Moniker: councilor.Moniker,
		},
	)
	require.NoError(t, err)
	require.Equal(t, councilor, resp.Councilor)

	// Non existing Councilor
	resp, err = querier.CouncilorByAddress(
		sdk.WrapSDKContext(ctx),
		&types.CouncilorByAddressRequest{ValAddr: addr2},
	)
	require.EqualError(t, err, "councilor not found: key not found")
}
