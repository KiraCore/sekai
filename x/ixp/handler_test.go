package ixp_test

import (
	"testing"

	"github.com/KiraCore/sekai/app"

	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/x/ixp"
	ixptypes "github.com/KiraCore/sekai/x/ixp/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestMain(m *testing.M) {
	app.SetConfig()
}

func TestNewHandler_MsgCreateOrderBook_HappyPath(t *testing.T) {
	kiraAddr1, err := types.AccAddressFromBech32("kira1da22wd7slpxpptasczs679mr5c8xtucqdzxc3n")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	handler := ixp.NewHandler(app.IxpKeeper)

	theMsg, err := ixptypes.NewMsgCreateOrderBook("base", "quote", "mnemonic", kiraAddr1)

	require.NoError(t, err)

	_, err = handler(ctx, theMsg) // TODO: should parse ID from handler response
	require.NoError(t, err)

	orderbooks := app.IxpKeeper.GetOrderBookByQuote(ctx, "quote") // TODO replace this to by handler getter
	require.Len(t, orderbooks, 1)

	orderbook := orderbooks[0]

	require.Equal(t, theMsg.Base, orderbook.Base)
	require.Equal(t, theMsg.Quote, orderbook.Quote)
	require.Equal(t, theMsg.Mnemonic, orderbook.Mnemonic)
	require.Equal(t, theMsg.Curator, orderbook.Curator)
}

func TestNewHandler_MsgCreateOrder_HappyPath(t *testing.T) {
	kiraAddr1, err := types.AccAddressFromBech32("kira1da22wd7slpxpptasczs679mr5c8xtucqdzxc3n")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	handler := ixp.NewHandler(app.IxpKeeper)

	theMsg, err := ixptypes.NewMsgCreateOrderBook("base", "quote", "mnemonic", kiraAddr1)

	require.NoError(t, err)

	_, err = handler(ctx, theMsg) // TODO: should parse ID from handler response
	require.NoError(t, err)

	orderbooks := app.IxpKeeper.GetOrderBookByQuote(ctx, "quote") // TODO replace this to by handler getter
	require.Len(t, orderbooks, 1)

	orderbook := orderbooks[0]
	bookID := orderbook.ID

	createOrderMsg, err := ixptypes.NewMsgCreateOrder(bookID, ixptypes.LimitOrderType_limitBuy, 10, 10, kiraAddr1)
	require.NoError(t, err)

	_, err = handler(ctx, createOrderMsg) // TODO: should parse ID from handler response
	require.NoError(t, err)

	orders := app.IxpKeeper.GetOrders(ctx, bookID, 0, 0)
	require.Len(t, orders, 1)

	order := orders[0]
	require.Equal(t, createOrderMsg.OrderType, order.OrderType)
	require.Equal(t, createOrderMsg.OrderBookID, order.OrderBookID)
	require.Equal(t, createOrderMsg.Amount, order.Amount)
	require.Equal(t, createOrderMsg.LimitPrice, order.LimitPrice)
	require.Equal(t, createOrderMsg.ExpiryTime, order.ExpiryTime)
	require.Equal(t, createOrderMsg.Curator, order.Curator)
	require.Equal(t, order.IsCancelled, false)
}

func TestNewHandler_MsgCancelOrder_HappyPath(t *testing.T) {
	kiraAddr1, err := types.AccAddressFromBech32("kira1da22wd7slpxpptasczs679mr5c8xtucqdzxc3n")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	handler := ixp.NewHandler(app.IxpKeeper)

	theMsg, err := ixptypes.NewMsgCreateOrderBook("base", "quote", "mnemonic", kiraAddr1)

	require.NoError(t, err)

	_, err = handler(ctx, theMsg) // TODO: should parse ID from handler response
	require.NoError(t, err)

	orderbooks := app.IxpKeeper.GetOrderBookByQuote(ctx, "quote") // TODO replace this to by handler getter
	require.Len(t, orderbooks, 1)

	orderbook := orderbooks[0]
	bookID := orderbook.ID

	createOrderMsg, err := ixptypes.NewMsgCreateOrder(bookID, ixptypes.LimitOrderType_limitBuy, 10, 10, kiraAddr1)
	require.NoError(t, err)

	_, err = handler(ctx, createOrderMsg) // TODO: should parse ID from handler response
	require.NoError(t, err)

	orders := app.IxpKeeper.GetOrders(ctx, bookID, 0, 0)
	require.Len(t, orders, 1)

	order := orders[0]
	orderID := order.ID

	cancelOrderMsg, err := ixptypes.NewMsgCancelOrder(orderID, kiraAddr1)
	require.NoError(t, err)

	_, err = handler(ctx, cancelOrderMsg)
	require.NoError(t, err)
	orders = app.IxpKeeper.GetOrders(ctx, bookID, 0, 0)
	require.Len(t, orders, 1)

	order = orders[0]
	require.Equal(t, createOrderMsg.OrderType, order.OrderType)
	require.Equal(t, createOrderMsg.OrderBookID, order.OrderBookID)
	require.Equal(t, createOrderMsg.Amount, order.Amount)
	require.Equal(t, createOrderMsg.LimitPrice, order.LimitPrice)
	require.Equal(t, createOrderMsg.ExpiryTime, order.ExpiryTime)
	require.Equal(t, createOrderMsg.Curator, order.Curator)
	require.Equal(t, order.IsCancelled, true)
}

func TestNewHandler_MsgUpsertSignerKey_HappyPath(t *testing.T) {
	kiraAddr1, err := types.AccAddressFromBech32("kira1da22wd7slpxpptasczs679mr5c8xtucqdzxc3n")
	require.NoError(t, err)

	pubKeyText := "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em"
	_, err = types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, pubKeyText)
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	handler := ixp.NewHandler(app.IxpKeeper)

	theMsg, err := ixptypes.NewMsgUpsertSignerKey(pubKeyText, ixptypes.SignerKeyType_Secp256k1, 0, true, []int64{}, kiraAddr1)
	require.NoError(t, err)

	_, err = handler(ctx, theMsg)
	require.NoError(t, err)

	signerkeys := app.IxpKeeper.GetSignerKeys(ctx, kiraAddr1)
	require.Len(t, signerkeys, 1)

	signerkey := signerkeys[0]

	require.Equal(t, theMsg.PubKey, signerkey.PubKey)
	require.Equal(t, theMsg.KeyType, signerkey.KeyType)
	require.Equal(t, theMsg.ExpiryTime, signerkey.ExpiryTime)
	require.Equal(t, theMsg.Permissions, signerkey.Permissions)
	require.Equal(t, theMsg.Curator, signerkey.Curator)
}
