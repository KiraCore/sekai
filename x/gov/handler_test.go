package gov_test

import (
	"testing"

	types2 "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/x/gov"
	"github.com/KiraCore/sekai/x/gov/types"
)

func TestMain(m *testing.M) {
	app.SetConfig()
	m.Run()
}

// When a network actor has not been saved before, it creates one with default params
// and sets the permissions.
func TestNewHandler_SetPermissions_ActorWithoutPerms(t *testing.T) {
	addr, err := types2.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	proposerAddr, err := types2.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	// First we set Permissions to SetPermissions to proposerAddr.
	proposerActor := types.NewDefaultActor(proposerAddr)
	err = proposerActor.Permissions.AddToWhitelist(types.PermSetPermissions)
	require.NoError(t, err)
	app.CustomGovKeeper.SaveNetworkActor(ctx, proposerActor)

	handler := gov.NewHandler(app.CustomGovKeeper)
	_, err = handler(ctx, &types.MsgWhitelistPermissions{
		Proposer:   proposerAddr,
		Address:    addr,
		Permission: uint32(types.PermClaimValidator),
	})
	require.NoError(t, err)

	actor, err := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
	require.NoError(t, err)

	require.True(t, actor.Permissions.IsWhitelisted(types.PermClaimValidator))
}

// When a network actor has already permissions it just appends the permission.
func TestNewHandler_SetPermissions_ActorWithPerms(t *testing.T) {
	addr, err := types2.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	proposerAddr, err := types2.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	// First we set Permissions to SetPermissions to proposerAddr.
	proposerActor := types.NewDefaultActor(proposerAddr)
	err = proposerActor.Permissions.AddToWhitelist(types.PermSetPermissions)
	require.NoError(t, err)
	app.CustomGovKeeper.SaveNetworkActor(ctx, proposerActor)

	// Add some perms before to the actor.
	actor := types.NewDefaultActor(addr)
	err = actor.Permissions.AddToWhitelist(types.PermSetPermissions)
	require.NoError(t, err)
	app.CustomGovKeeper.SaveNetworkActor(ctx, actor)

	// Call the handler to add some permissions.
	handler := gov.NewHandler(app.CustomGovKeeper)
	_, err = handler(ctx, &types.MsgWhitelistPermissions{
		Proposer:   proposerAddr,
		Address:    addr,
		Permission: uint32(types.PermClaimValidator),
	})
	require.NoError(t, err)

	actor, err = app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
	require.NoError(t, err)

	require.True(t, actor.Permissions.IsWhitelisted(types.PermClaimValidator))
	require.True(t, actor.Permissions.IsWhitelisted(types.PermSetPermissions)) // This permission was already set before callid add permission.
}

func TestNewHandler_SetPermissionsWithoutSetPermissions(t *testing.T) {
	addr, err := types2.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	proposerAddr, err := types2.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	handler := gov.NewHandler(app.CustomGovKeeper)
	_, err = handler(ctx, &types.MsgWhitelistPermissions{
		Proposer:   proposerAddr,
		Address:    addr,
		Permission: uint32(types.PermClaimValidator),
	})
	require.EqualError(t, err, "SetPermissions: not enough permissions")
}

