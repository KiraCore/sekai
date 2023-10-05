package keeper_test

import (
	"testing"

	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/gov/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServerClaimCouncilor(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	// try claim with not allowed user
	msgServer := keeper.NewMsgServerImpl(app.CustomGovKeeper)
	_, err := msgServer.ClaimCouncilor(sdk.WrapSDKContext(ctx), types.NewMsgClaimCouncilor(
		addr, "moniker", "username", "description", "social", "contact", "avatar",
	))
	require.Error(t, err)

	// try claim with allowed user
	actor := types.NewDefaultActor(addr)
	err = app.CustomGovKeeper.AddWhitelistPermission(ctx, actor, types.PermClaimCouncilor)
	require.NoError(t, err)
	_, err = msgServer.ClaimCouncilor(sdk.WrapSDKContext(ctx), types.NewMsgClaimCouncilor(
		addr, "moniker", "username", "description", "social", "contact", "avatar",
	))
	require.NoError(t, err)

	// check councilor created
	councilor, found := app.CustomGovKeeper.GetCouncilor(ctx, addr)
	require.True(t, found)
	require.Equal(t, councilor, types.Councilor{
		Address:           addr,
		Status:            types.CouncilorActive,
		Rank:              0,
		AbstentionCounter: 0,
	})

	// check moniker, username, description, social, contact, avatar registered correctly
	recordId := app.CustomGovKeeper.GetIdentityRecordIdByAddressKey(ctx, addr, "moniker")
	record := app.CustomGovKeeper.GetIdentityRecordById(ctx, recordId)
	require.Equal(t, record.Key, "moniker")
	require.Equal(t, record.Value, "moniker")

	recordId = app.CustomGovKeeper.GetIdentityRecordIdByAddressKey(ctx, addr, "username")
	record = app.CustomGovKeeper.GetIdentityRecordById(ctx, recordId)
	require.Equal(t, record.Key, "username")
	require.Equal(t, record.Value, "username")

	recordId = app.CustomGovKeeper.GetIdentityRecordIdByAddressKey(ctx, addr, "description")
	record = app.CustomGovKeeper.GetIdentityRecordById(ctx, recordId)
	require.Equal(t, record.Key, "description")
	require.Equal(t, record.Value, "description")

	recordId = app.CustomGovKeeper.GetIdentityRecordIdByAddressKey(ctx, addr, "social")
	record = app.CustomGovKeeper.GetIdentityRecordById(ctx, recordId)
	require.Equal(t, record.Key, "social")
	require.Equal(t, record.Value, "social")

	recordId = app.CustomGovKeeper.GetIdentityRecordIdByAddressKey(ctx, addr, "contact")
	record = app.CustomGovKeeper.GetIdentityRecordById(ctx, recordId)
	require.Equal(t, record.Key, "contact")
	require.Equal(t, record.Value, "contact")

	recordId = app.CustomGovKeeper.GetIdentityRecordIdByAddressKey(ctx, addr, "avatar")
	record = app.CustomGovKeeper.GetIdentityRecordById(ctx, recordId)
	require.Equal(t, record.Key, "avatar")
	require.Equal(t, record.Value, "avatar")
}

func TestMsgServerCouncilorPause(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	// try pause with not available councilor
	msgServer := keeper.NewMsgServerImpl(app.CustomGovKeeper)
	_, err := msgServer.CouncilorPause(sdk.WrapSDKContext(ctx), types.NewMsgCouncilorPause(
		addr,
	))
	require.Error(t, err)

	// test jailed councilor pause
	app.CustomGovKeeper.SaveCouncilor(ctx, types.Councilor{
		Address:           addr,
		Status:            types.CouncilorJailed,
		Rank:              0,
		AbstentionCounter: 0,
	})
	_, err = msgServer.CouncilorPause(sdk.WrapSDKContext(ctx), types.NewMsgCouncilorPause(
		addr,
	))
	require.Error(t, err)

	// test inactive councilor
	app.CustomGovKeeper.SaveCouncilor(ctx, types.Councilor{
		Address:           addr,
		Status:            types.CouncilorInactive,
		Rank:              0,
		AbstentionCounter: 0,
	})
	_, err = msgServer.CouncilorPause(sdk.WrapSDKContext(ctx), types.NewMsgCouncilorPause(
		addr,
	))
	require.Error(t, err)

	// test paused councilor pause
	app.CustomGovKeeper.SaveCouncilor(ctx, types.Councilor{
		Address:           addr,
		Status:            types.CouncilorPaused,
		Rank:              0,
		AbstentionCounter: 0,
	})
	_, err = msgServer.CouncilorPause(sdk.WrapSDKContext(ctx), types.NewMsgCouncilorPause(
		addr,
	))
	require.Error(t, err)

	// active councilor pause
	app.CustomGovKeeper.SaveCouncilor(ctx, types.Councilor{
		Address:           addr,
		Status:            types.CouncilorActive,
		Rank:              0,
		AbstentionCounter: 0,
	})
	_, err = msgServer.CouncilorPause(sdk.WrapSDKContext(ctx), types.NewMsgCouncilorPause(
		addr,
	))
	require.NoError(t, err)
	councilor, found := app.CustomGovKeeper.GetCouncilor(ctx, addr)
	require.True(t, found)
	require.Equal(t, councilor.Status, types.CouncilorPaused)
}

func TestMsgServerCouncilorUnpause(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	// try unpause with not available councilor
	msgServer := keeper.NewMsgServerImpl(app.CustomGovKeeper)
	_, err := msgServer.CouncilorUnpause(sdk.WrapSDKContext(ctx), types.NewMsgCouncilorUnpause(
		addr,
	))
	require.Error(t, err)

	// test active councilor unpause
	app.CustomGovKeeper.SaveCouncilor(ctx, types.Councilor{
		Address:           addr,
		Status:            types.CouncilorActive,
		Rank:              0,
		AbstentionCounter: 0,
	})
	_, err = msgServer.CouncilorUnpause(sdk.WrapSDKContext(ctx), types.NewMsgCouncilorUnpause(
		addr,
	))
	require.Error(t, err)

	// test paused councilor pause
	app.CustomGovKeeper.SaveCouncilor(ctx, types.Councilor{
		Address:           addr,
		Status:            types.CouncilorPaused,
		Rank:              0,
		AbstentionCounter: 0,
	})
	_, err = msgServer.CouncilorUnpause(sdk.WrapSDKContext(ctx), types.NewMsgCouncilorUnpause(
		addr,
	))
	require.NoError(t, err)

	councilor, found := app.CustomGovKeeper.GetCouncilor(ctx, addr)
	require.True(t, found)
	require.Equal(t, councilor.Status, types.CouncilorActive)
}

func TestMsgServerCouncilorActivate(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	// try activate with not available councilor
	msgServer := keeper.NewMsgServerImpl(app.CustomGovKeeper)
	_, err := msgServer.CouncilorActivate(sdk.WrapSDKContext(ctx), types.NewMsgCouncilorActivate(
		addr,
	))
	require.Error(t, err)

	// test active councilor activate
	app.CustomGovKeeper.SaveCouncilor(ctx, types.Councilor{
		Address:           addr,
		Status:            types.CouncilorActive,
		Rank:              0,
		AbstentionCounter: 100,
	})
	_, err = msgServer.CouncilorActivate(sdk.WrapSDKContext(ctx), types.NewMsgCouncilorActivate(
		addr,
	))
	require.Error(t, err)

	// test inactive councilor activate
	app.CustomGovKeeper.SaveCouncilor(ctx, types.Councilor{
		Address:           addr,
		Status:            types.CouncilorInactive,
		Rank:              0,
		AbstentionCounter: 100,
	})
	_, err = msgServer.CouncilorActivate(sdk.WrapSDKContext(ctx), types.NewMsgCouncilorActivate(
		addr,
	))
	require.NoError(t, err)

	councilor, found := app.CustomGovKeeper.GetCouncilor(ctx, addr)
	require.True(t, found)
	// check status change
	require.Equal(t, councilor.Status, types.CouncilorActive)
	// check absention counter change
	require.Equal(t, councilor.AbstentionCounter, int64(0))
}
