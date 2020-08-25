package keeper_test

import (
	"testing"

	types2 "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/stretchr/testify/require"

	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/x/gov/types"
)

func TestKeeper_SaveGetPermissionsForRole(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	perm := types.NewPermissions(
		nil, []types.PermValue{types.PermClaimValidator},
	)

	app.CustomGovKeeper.SetPermissionsForRole(ctx, types.RoleCouncilor, perm)

	savedPerms := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleCouncilor)
	require.Equal(t, perm, savedPerms)
}

func TestNewKeeper_SaveNetworkActor(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, types2.TokensFromConsensusPower(10))
	addr := addrs[0]

	networkActor := types.NetworkActor{
		Address:     addr,
		Roles:       nil,
		Status:      0,
		Votes:       nil,
		Permissions: nil,
		Skin:        0,
	}

	app.CustomGovKeeper.SetNetworkActor(ctx, networkActor)

	savedActor, err := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, networkActor.Address)
	require.NoError(t, err)

	require.Equal(t, networkActor, savedActor)
}

func TestKeeper_GetNetworkActorByAddress_FailsIfItDoesNotExist(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, types2.TokensFromConsensusPower(10))
	addr := addrs[0]

	_, err := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
	require.EqualError(t, err, "network actor not found")
}
