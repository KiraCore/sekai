package gov

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KiraCore/sekai/x/gov/types"

	"github.com/KiraCore/sekai/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestQuerier_PermissionsByAddress(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(10))
	addr1 := addrs[0]

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

	app.CustomGovKeeper.SetNetworkActor(ctx, networkActor)

	querier := NewQuerier(app.CustomGovKeeper)

	resp, err := querier.PermissionsByAddress(sdk.WrapSDKContext(ctx), &types.PermissionsByAddressRequest{ValAddr: addr1})
	require.NoError(t, err)

	require.Equal(t, permissions, resp.Permissions)
}
