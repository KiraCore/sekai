package ixp_test

import (
	"testing"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/KiraCore/sekai/x/ixp"

	ixptypes "github.com/KiraCore/sekai/x/ixp/types"

	"github.com/KiraCore/sekai/simapp"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestQuerier_GetOrderBooks(t *testing.T) {
	kiraAddr1, err := types.AccAddressFromBech32("kira1da22wd7slpxpptasczs679mr5c8xtucqdzxc3n")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	app.IxpKeeper.CreateOrderBook(ctx, "quote", "base", kiraAddr1, "mnemonic")

	querier := ixp.NewQuerier(app.IxpKeeper)

	qOrderBooksResp, err := querier.GetOrderBooks(types.WrapSDKContext(ctx), &ixptypes.GetOrderBooksRequest{
		QueryType:  "Quote",
		QueryValue: "quote",
	})
	require.NoError(t, err)

	require.Len(t, qOrderBooksResp.Orderbooks, 1)
	orderbook := qOrderBooksResp.Orderbooks[0]

	require.Equal(t, orderbook.Base, "base")
	require.Equal(t, orderbook.Quote, "quote")
	require.Equal(t, orderbook.Mnemonic, "mnemonic")
	require.Equal(t, orderbook.Curator, kiraAddr1)
}

// TODO should add tests for GetOrderBooksByTradingPair
// TODO should add tests for GetOrders
// TODO should add tests for GetSignerKeys
