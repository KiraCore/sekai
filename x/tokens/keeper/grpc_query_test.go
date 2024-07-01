package keeper_test

import (
	"testing"

	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/x/tokens/keeper"
	"github.com/KiraCore/sekai/x/tokens/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestQuerier_GetTokenInfo(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	querier := keeper.NewQuerier(app.TokensKeeper)

	resp, err := querier.GetTokenInfo(
		sdk.WrapSDKContext(ctx),
		&types.TokenInfoRequest{Denom: "ukex"},
	)
	require.NoError(t, err)
	require.Equal(t, "ukex", resp.Data.Denom)
	require.Equal(t, sdk.NewDec(1), resp.Data.FeeRate)
	require.Equal(t, true, resp.Data.FeeEnabled)
}

func TestQuerier_GetTokenInfosByDenom(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	querier := keeper.NewQuerier(app.TokensKeeper)

	resp, err := querier.GetTokenInfosByDenom(
		sdk.WrapSDKContext(ctx),
		&types.TokenInfosByDenomRequest{Denoms: []string{"ukex"}},
	)
	require.NoError(t, err)
	require.Equal(t, len(resp.Data), 1)
	require.Equal(t, "ukex", resp.Data["ukex"].Data.Denom)
	require.Equal(t, sdk.NewDec(1), resp.Data["ukex"].Data.FeeRate)
	require.Equal(t, true, resp.Data["ukex"].Data.FeeEnabled)
}

func TestQuerier_GetAllTokenInfos(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	querier := keeper.NewQuerier(app.TokensKeeper)

	resp, err := querier.GetAllTokenInfos(
		sdk.WrapSDKContext(ctx),
		&types.AllTokenInfosRequest{},
	)
	require.NoError(t, err)
	require.Equal(t, len(resp.Data), 4)
	require.Equal(t, "xeth", resp.Data[0].Data.Denom)
	require.Equal(t, sdk.NewDecWithPrec(1, 1), resp.Data[0].Data.FeeRate)
	require.Equal(t, true, resp.Data[0].Data.FeeEnabled)
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
