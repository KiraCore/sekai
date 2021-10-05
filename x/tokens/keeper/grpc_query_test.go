package keeper_test

import (
	"testing"

	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/x/tokens/keeper"
	"github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestQuerier_GetTokenAlias(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	querier := keeper.NewQuerier(app.TokensKeeper)

	resp, err := querier.GetTokenAlias(
		sdk.WrapSDKContext(ctx),
		&types.TokenAliasRequest{Symbol: "KEX"},
	)
	require.NoError(t, err)
	require.Equal(t, "KEX", resp.Data.Symbol)
	require.Equal(t, "Kira", resp.Data.Name)
	require.Equal(t, "", resp.Data.Icon)
	require.Equal(t, uint32(0x6), resp.Data.Decimals)
	require.Equal(t, []string{"ukex", "mkex"}, resp.Data.Denoms)
}

func TestQuerier_GetTokenAliasesByDenom(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	querier := keeper.NewQuerier(app.TokensKeeper)

	resp, err := querier.GetTokenAliasesByDenom(
		sdk.WrapSDKContext(ctx),
		&types.TokenAliasesByDenomRequest{Denoms: []string{"ukex"}},
	)
	require.NoError(t, err)
	require.Equal(t, len(resp.Data), 1)
	require.Equal(t, "KEX", resp.Data["ukex"].Symbol)
	require.Equal(t, "Kira", resp.Data["ukex"].Name)
	require.Equal(t, "", resp.Data["ukex"].Icon)
	require.Equal(t, uint32(0x6), resp.Data["ukex"].Decimals)
	require.Equal(t, []string{"ukex", "mkex"}, resp.Data["ukex"].Denoms)
}

func TestQuerier_GetAllTokenAliases(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	querier := keeper.NewQuerier(app.TokensKeeper)

	resp, err := querier.GetAllTokenAliases(
		sdk.WrapSDKContext(ctx),
		&types.AllTokenAliasesRequest{},
	)
	require.NoError(t, err)
	require.Equal(t, len(resp.Data), 1)
	require.Equal(t, "KEX", resp.Data[0].Symbol)
	require.Equal(t, "Kira", resp.Data[0].Name)
	require.Equal(t, "", resp.Data[0].Icon)
	require.Equal(t, uint32(0x6), resp.Data[0].Decimals)
	require.Equal(t, []string{"ukex", "mkex"}, resp.Data[0].Denoms)
}

func TestQuerier_GetTokenRate(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	querier := keeper.NewQuerier(app.TokensKeeper)

	resp, err := querier.GetTokenRate(
		sdk.WrapSDKContext(ctx),
		&types.TokenRateRequest{Denom: "ukex"},
	)
	require.NoError(t, err)
	require.Equal(t, "ukex", resp.Data.Denom)
	require.Equal(t, sdk.NewDec(1), resp.Data.Rate)
	require.Equal(t, true, resp.Data.FeePayments)
}

func TestQuerier_GetTokenRatesByDenom(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	querier := keeper.NewQuerier(app.TokensKeeper)

	resp, err := querier.GetTokenRatesByDenom(
		sdk.WrapSDKContext(ctx),
		&types.TokenRatesByDenomRequest{Denoms: []string{"ukex"}},
	)
	require.NoError(t, err)
	require.Equal(t, len(resp.Data), 1)
	require.Equal(t, "ukex", resp.Data["ukex"].Denom)
	require.Equal(t, sdk.NewDec(1), resp.Data["ukex"].Rate)
	require.Equal(t, true, resp.Data["ukex"].FeePayments)
}

func TestQuerier_GetAllTokenRates(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	querier := keeper.NewQuerier(app.TokensKeeper)

	resp, err := querier.GetAllTokenRates(
		sdk.WrapSDKContext(ctx),
		&types.AllTokenRatesRequest{},
	)
	require.NoError(t, err)
	require.Equal(t, len(resp.Data), 4)
	require.Equal(t, "frozen", resp.Data[0].Denom)
	require.Equal(t, sdk.NewDecWithPrec(1, 1), resp.Data[0].Rate)
	require.Equal(t, true, resp.Data[0].FeePayments)
}

func TestQuerier_GetTokenBlackWhites(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	querier := keeper.NewQuerier(app.TokensKeeper)

	resp, err := querier.GetTokenBlackWhites(
		sdk.WrapSDKContext(ctx),
		&types.TokenBlackWhitesRequest{},
	)
	require.NoError(t, err)
	require.Equal(t, resp.Data.Blacklisted, []string{"frozen"})
	require.Equal(t, resp.Data.Whitelisted, []string{"ukex"})
}
