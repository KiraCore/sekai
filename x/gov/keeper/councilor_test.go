package keeper_test

import (
	"testing"

	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/x/gov/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_SaveCouncilor(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	addr := addrs[0]

	councilor := types.NewCouncilor(
		addr,
		types.CouncilorActive,
	)

	app.CustomGovKeeper.SaveCouncilor(ctx, councilor)

	savedCouncilor, found := app.CustomGovKeeper.GetCouncilor(ctx, councilor.Address)
	require.True(t, found)
	require.Equal(t, councilor, savedCouncilor)
}

func TestKeeper_TryCouncilorMonikerOrUsernameUpdate(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	addr := addrs[0]

	councilor := types.NewCouncilor(
		addr,
		types.CouncilorActive,
	)

	app.CustomGovKeeper.SaveCouncilor(ctx, councilor)

	savedCouncilor, found := app.CustomGovKeeper.GetCouncilor(ctx, councilor.Address)
	require.True(t, found)
	require.Equal(t, councilor, savedCouncilor)

	err := app.CustomGovKeeper.RegisterIdentityRecords(ctx, addr, []types.IdentityInfoEntry{
		{Key: "moniker", Info: "moniker1"},
	})
	require.NoError(t, err)
	err = app.CustomGovKeeper.RegisterIdentityRecords(ctx, addr, []types.IdentityInfoEntry{
		{Key: "moniker", Info: "xxxx"},
	})
	require.Error(t, err)

	err = app.CustomGovKeeper.RegisterIdentityRecords(ctx, addr, []types.IdentityInfoEntry{
		{Key: "username", Info: "username1"},
	})
	require.NoError(t, err)
	err = app.CustomGovKeeper.RegisterIdentityRecords(ctx, addr, []types.IdentityInfoEntry{
		{Key: "username", Info: "xxxx"},
	})
	require.Error(t, err)

	err = app.CustomGovKeeper.RegisterIdentityRecords(ctx, addr, []types.IdentityInfoEntry{
		{Key: "avatar", Info: "avatar1"},
	})
	require.NoError(t, err)
	err = app.CustomGovKeeper.RegisterIdentityRecords(ctx, addr, []types.IdentityInfoEntry{
		{Key: "avatar", Info: "xxxx"},
	})
	require.NoError(t, err)
}

func TestKeeper_OnCouncilorAct(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	addr := addrs[0]

	councilor := types.Councilor{
		Address:           addr,
		Status:            types.CouncilorActive,
		Rank:              1,
		AbstentionCounter: 1,
	}

	app.CustomGovKeeper.SaveCouncilor(ctx, councilor)

	savedCouncilor, found := app.CustomGovKeeper.GetCouncilor(ctx, councilor.Address)
	require.True(t, found)
	require.Equal(t, councilor, savedCouncilor)

	app.CustomGovKeeper.OnCouncilorAct(ctx, addr)
	savedCouncilor, found = app.CustomGovKeeper.GetCouncilor(ctx, councilor.Address)
	require.True(t, found)
	require.Equal(t, types.Councilor{
		Address:           addr,
		Status:            types.CouncilorActive,
		Rank:              2,
		AbstentionCounter: 0,
	}, savedCouncilor)
}

func TestKeeper_OnCouncilorAbsent(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	addr := addrs[0]

	councilor := types.Councilor{
		Address:           addr,
		Status:            types.CouncilorActive,
		Rank:              1,
		AbstentionCounter: 1,
	}

	app.CustomGovKeeper.SaveCouncilor(ctx, councilor)

	savedCouncilor, found := app.CustomGovKeeper.GetCouncilor(ctx, councilor.Address)
	require.True(t, found)
	require.Equal(t, councilor, savedCouncilor)

	app.CustomGovKeeper.OnCouncilorAbsent(ctx, addr)
	savedCouncilor, found = app.CustomGovKeeper.GetCouncilor(ctx, councilor.Address)
	require.True(t, found)
	require.Equal(t, types.Councilor{
		Address:           addr,
		Status:            types.CouncilorInactive,
		Rank:              0,
		AbstentionCounter: 2,
	}, savedCouncilor)
}

func TestKeeper_OnCouncilorJail(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	addr := addrs[0]

	councilor := types.Councilor{
		Address:           addr,
		Status:            types.CouncilorActive,
		Rank:              1,
		AbstentionCounter: 1,
	}

	app.CustomGovKeeper.SaveCouncilor(ctx, councilor)

	savedCouncilor, found := app.CustomGovKeeper.GetCouncilor(ctx, councilor.Address)
	require.True(t, found)
	require.Equal(t, councilor, savedCouncilor)

	app.CustomGovKeeper.OnCouncilorJail(ctx, addr)
	savedCouncilor, found = app.CustomGovKeeper.GetCouncilor(ctx, councilor.Address)
	require.True(t, found)
	require.Equal(t, types.Councilor{
		Address:           addr,
		Status:            types.CouncilorJailed,
		Rank:              0,
		AbstentionCounter: 1,
	}, savedCouncilor)
}

func TestKeeper_ResetWholeCouncilorRank(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	addr := addrs[0]

	councilor := types.Councilor{
		Address:           addr,
		Status:            types.CouncilorActive,
		Rank:              1,
		AbstentionCounter: 1,
	}

	app.CustomGovKeeper.SaveCouncilor(ctx, councilor)

	savedCouncilor, found := app.CustomGovKeeper.GetCouncilor(ctx, councilor.Address)
	require.True(t, found)
	require.Equal(t, councilor, savedCouncilor)

	app.CustomGovKeeper.ResetWholeCouncilorRank(ctx)
	savedCouncilor, found = app.CustomGovKeeper.GetCouncilor(ctx, councilor.Address)
	require.True(t, found)
	require.Equal(t, types.Councilor{
		Address:           addr,
		Status:            types.CouncilorActive,
		Rank:              0,
		AbstentionCounter: 0,
	}, savedCouncilor)
}
