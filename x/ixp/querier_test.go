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

func TestQuerier_GetOrderBooksByTradingPair(t *testing.T) {
	kiraAddr1, err := types.AccAddressFromBech32("kira1da22wd7slpxpptasczs679mr5c8xtucqdzxc3n")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	app.IxpKeeper.CreateOrderBook(ctx, "quote", "base", kiraAddr1, "mnemonic")

	querier := ixp.NewQuerier(app.IxpKeeper)

	qOrderBooksResp, err := querier.GetOrderBooksByTradingPair(types.WrapSDKContext(ctx), &ixptypes.GetOrderBooksByTradingPairRequest{
		Base:  "base",
		Quote: "quote",
	})
	require.NoError(t, err)

	require.Len(t, qOrderBooksResp.Orderbooks, 1)
	orderbook := qOrderBooksResp.Orderbooks[0]

	require.Equal(t, orderbook.Base, "base")
	require.Equal(t, orderbook.Quote, "quote")
	require.Equal(t, orderbook.Mnemonic, "mnemonic")
	require.Equal(t, orderbook.Curator, kiraAddr1)
}

func TestQuerier_GetOrders(t *testing.T) {
	kiraAddr1, err := types.AccAddressFromBech32("kira1da22wd7slpxpptasczs679mr5c8xtucqdzxc3n")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	bookID, err := app.IxpKeeper.CreateOrderBook(ctx, "quote", "base", kiraAddr1, "mnemonic")
	require.NoError(t, err)

	orderID, err := app.IxpKeeper.CreateOrder(ctx, bookID, ixptypes.LimitOrderType_limitBuy, 10, 10, 10000, kiraAddr1)
	require.NoError(t, err)

	querier := ixp.NewQuerier(app.IxpKeeper)

	qOrdersResp, err := querier.GetOrders(types.WrapSDKContext(ctx), &ixptypes.GetOrdersRequest{
		OrderBookID: bookID,
		MaxOrders:   0,
		MinAmount:   0,
	})
	require.NoError(t, err)

	require.Len(t, qOrdersResp.Orders, 1)
	order := qOrdersResp.Orders[0]

	require.Equal(t, order.ID, orderID)
	require.Equal(t, order.OrderBookID, bookID)
	require.Equal(t, order.OrderType, ixptypes.LimitOrderType_limitBuy)
	require.Equal(t, order.Amount, int64(10))
	require.Equal(t, order.LimitPrice, int64(10))
	require.Equal(t, order.ExpiryTime, int64(10000))
}

// TODO should add tests for GetSignerKeys
