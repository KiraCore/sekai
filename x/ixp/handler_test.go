package ixp_test

import (
	"bytes"
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/KiraCore/sekai/app"

	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/x/ixp"
	"github.com/KiraCore/sekai/x/ixp/handlers"
	ixptypes "github.com/KiraCore/sekai/x/ixp/types"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestMain(m *testing.M) {
	app.SetConfig()
	os.Exit(m.Run())
}

func ParseResponseID(result *sdk.Result, t *testing.T) string {
	resultParser := handlers.CreateOrderBookResponse{}
	err := json.Unmarshal(result.Data, &resultParser)
	require.NoError(t, err)
	return resultParser.ID
}

func NewAccountByIndex(accNum int) sdk.AccAddress {
	var buffer bytes.Buffer
	i := accNum + 100
	numString := strconv.Itoa(i)
	buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") //base address string

	buffer.WriteString(numString) //adding on final two digits to make addresses unique
	res, _ := sdk.AccAddressFromHex(buffer.String())
	bech := res.String()
	addr, _ := simapp.TestAddr(buffer.String(), bech)
	buffer.Reset()
	return addr
}

// GenSecp256k1PubKey is a function to generate a pubKey
func GenSecp256k1PubKey() (secp256k1.PrivKey, crypto.PubKey, error) {
	priv := secp256k1.GenPrivKey()
	return priv, priv.PubKey(), nil
}

// GenEd25519PubKey is a function to generate a pubKey
func GenEd25519PubKey() (ed25519.PrivKey, crypto.PubKey, error) {
	priv := ed25519.GenPrivKey()
	return priv, priv.PubKey(), nil
}

func TestNewHandler_MsgCreateOrderBook_HappyPath(t *testing.T) {
	kiraAddr1, err := types.AccAddressFromBech32("kira1da22wd7slpxpptasczs679mr5c8xtucqdzxc3n")
	emptyKiraAddr1 := types.AccAddress{}
	emptyKiraAddr2 := types.AccAddress(nil)
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	handler := ixp.NewHandler(app.IxpKeeper)

	tests := []struct {
		name        string
		constructor func() (*ixptypes.MsgCreateOrderBook, error)
	}{
		{
			name: "basic path test",
			constructor: func() (*ixptypes.MsgCreateOrderBook, error) {
				return ixptypes.NewMsgCreateOrderBook("base", "quote", "mnemonic", kiraAddr1)
			},
		},
		{
			name: "empty base case",
			// TODO: This shouldn't fail?
			constructor: func() (*ixptypes.MsgCreateOrderBook, error) {
				return ixptypes.NewMsgCreateOrderBook("", "quote", "mnemonic", kiraAddr1)
			},
		},
		{
			name: "empty quote case",
			// TODO: This shouldn't fail?
			constructor: func() (*ixptypes.MsgCreateOrderBook, error) {
				return ixptypes.NewMsgCreateOrderBook("base", "", "mnemonic", kiraAddr1)
			},
		},
		{
			name: "empty mnemonic case",
			// TODO: This shouldn't fail?
			constructor: func() (*ixptypes.MsgCreateOrderBook, error) {
				return ixptypes.NewMsgCreateOrderBook("base", "quote", "", kiraAddr1)
			},
		},
		{
			name: "empty curator case1",
			// TODO: This shouldn't fail?
			constructor: func() (*ixptypes.MsgCreateOrderBook, error) {
				return ixptypes.NewMsgCreateOrderBook("base", "quote", "mnemonic", emptyKiraAddr1)
			},
		},
		{
			name: "empty curator case2",
			// TODO: This shouldn't fail?
			constructor: func() (*ixptypes.MsgCreateOrderBook, error) {
				return ixptypes.NewMsgCreateOrderBook("base", "quote", "mnemonic", emptyKiraAddr2)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			theMsg, err := tt.constructor()
			require.NoError(t, err)

			result, err := handler(ctx, theMsg)
			require.NoError(t, err)

			orderbooks := app.IxpKeeper.GetOrderBookByID(ctx, ParseResponseID(result, t))
			require.Len(t, orderbooks, 1)

			orderbook := orderbooks[0]

			require.Equal(t, theMsg.Base, orderbook.Base)
			require.Equal(t, theMsg.Quote, orderbook.Quote)
			require.Equal(t, theMsg.Mnemonic, orderbook.Mnemonic)
			if theMsg.Curator != nil && theMsg.Curator.Empty() {
				require.Equal(t, emptyKiraAddr2, orderbook.Curator)
			} else {
				require.Equal(t, theMsg.Curator, orderbook.Curator)
			}
		})
	}
}

func TestNewHandler_MsgCreateOrder_HappyPath(t *testing.T) {
	kiraAddr1, err := types.AccAddressFromBech32("kira1da22wd7slpxpptasczs679mr5c8xtucqdzxc3n")
	require.NoError(t, err)
	emptyKiraAddr1 := types.AccAddress{}
	emptyKiraAddr2 := types.AccAddress(nil)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	handler := ixp.NewHandler(app.IxpKeeper)

	theMsg, err := ixptypes.NewMsgCreateOrderBook("base", "quote", "mnemonic", kiraAddr1)
	require.NoError(t, err)

	result, err := handler(ctx, theMsg)
	require.NoError(t, err)

	orderbooks := app.IxpKeeper.GetOrderBookByID(ctx, ParseResponseID(result, t))
	require.Len(t, orderbooks, 1)

	orderbook := orderbooks[0]
	bookID := orderbook.ID

	tests := []struct {
		name        string
		constructor func() (*ixptypes.MsgCreateOrder, error)
	}{
		{
			name: "buy order test",
			constructor: func() (*ixptypes.MsgCreateOrder, error) {
				return ixptypes.NewMsgCreateOrder(bookID, ixptypes.LimitOrderType_limitBuy, 10, 10, kiraAddr1)
			},
		},
		{
			name: "sell order test",
			constructor: func() (*ixptypes.MsgCreateOrder, error) {
				return ixptypes.NewMsgCreateOrder(bookID, ixptypes.LimitOrderType_limitSell, 10, 10, kiraAddr1)
			},
		},
		{
			name: "zero price test",
			// TODO: This shouldn't fail?
			constructor: func() (*ixptypes.MsgCreateOrder, error) {
				return ixptypes.NewMsgCreateOrder(bookID, ixptypes.LimitOrderType_limitBuy, 10, 0, kiraAddr1)
			},
		},
		{
			name: "zero amount test",
			// TODO: This shouldn't fail?
			constructor: func() (*ixptypes.MsgCreateOrder, error) {
				return ixptypes.NewMsgCreateOrder(bookID, ixptypes.LimitOrderType_limitBuy, 0, 10, kiraAddr1)
			},
		},
		{
			name: "empty curator test1",
			// TODO: This shouldn't fail?
			constructor: func() (*ixptypes.MsgCreateOrder, error) {
				return ixptypes.NewMsgCreateOrder(bookID, ixptypes.LimitOrderType_limitBuy, 10, 10, emptyKiraAddr1)
			},
		},
		{
			name: "empty curator test2",
			// TODO: This shouldn't fail?
			constructor: func() (*ixptypes.MsgCreateOrder, error) {
				return ixptypes.NewMsgCreateOrder(bookID, ixptypes.LimitOrderType_limitBuy, 10, 10, emptyKiraAddr2)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			createOrderMsg, err := tt.constructor()
			require.NoError(t, err)

			result, err = handler(ctx, createOrderMsg)
			require.NoError(t, err)

			order, err := app.IxpKeeper.GetOrderByID(ctx, ParseResponseID(result, t))
			require.NoError(t, err)

			require.Equal(t, createOrderMsg.OrderType, order.OrderType)
			require.Equal(t, createOrderMsg.OrderBookID, order.OrderBookID)
			require.Equal(t, createOrderMsg.Amount, order.Amount)
			require.Equal(t, createOrderMsg.LimitPrice, order.LimitPrice)
			require.Equal(t, createOrderMsg.ExpiryTime, order.ExpiryTime)
			if createOrderMsg.Curator != nil && createOrderMsg.Curator.Empty() {
				require.Equal(t, emptyKiraAddr2, order.Curator)
			} else {
				require.Equal(t, createOrderMsg.Curator, order.Curator)
			}
			require.Equal(t, order.IsCancelled, false)
		})
	}
}

func TestNewHandler_MsgCancelOrder_HappyPath(t *testing.T) {
	kiraAddr1, err := types.AccAddressFromBech32("kira1da22wd7slpxpptasczs679mr5c8xtucqdzxc3n")
	require.NoError(t, err)
	emptyKiraAddr1 := types.AccAddress{}
	emptyKiraAddr2 := types.AccAddress(nil)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	handler := ixp.NewHandler(app.IxpKeeper)

	theMsg, err := ixptypes.NewMsgCreateOrderBook("base", "quote", "mnemonic", kiraAddr1)

	require.NoError(t, err)

	result, err := handler(ctx, theMsg)
	require.NoError(t, err)

	orderbooks := app.IxpKeeper.GetOrderBookByID(ctx, ParseResponseID(result, t))
	require.Len(t, orderbooks, 1)

	orderbook := orderbooks[0]
	bookID := orderbook.ID

	tests := []struct {
		name        string
		constructor func() (*ixptypes.MsgCreateOrder, error)
	}{
		{
			name: "buy order test",
			constructor: func() (*ixptypes.MsgCreateOrder, error) {
				return ixptypes.NewMsgCreateOrder(bookID, ixptypes.LimitOrderType_limitBuy, 10, 10, kiraAddr1)
			},
		},
		{
			name: "sell order test",
			constructor: func() (*ixptypes.MsgCreateOrder, error) {
				return ixptypes.NewMsgCreateOrder(bookID, ixptypes.LimitOrderType_limitSell, 10, 10, kiraAddr1)
			},
		},
		{
			name: "zero price test",
			// TODO: This shouldn't fail?
			constructor: func() (*ixptypes.MsgCreateOrder, error) {
				return ixptypes.NewMsgCreateOrder(bookID, ixptypes.LimitOrderType_limitBuy, 10, 0, kiraAddr1)
			},
		},
		{
			name: "zero amount test",
			// TODO: This shouldn't fail?
			constructor: func() (*ixptypes.MsgCreateOrder, error) {
				return ixptypes.NewMsgCreateOrder(bookID, ixptypes.LimitOrderType_limitBuy, 0, 10, kiraAddr1)
			},
		},
		{
			name: "empty curator test1",
			// TODO: This shouldn't fail?
			constructor: func() (*ixptypes.MsgCreateOrder, error) {
				return ixptypes.NewMsgCreateOrder(bookID, ixptypes.LimitOrderType_limitBuy, 10, 10, emptyKiraAddr1)
			},
		},
		{
			name: "empty curator test2",
			// TODO: This shouldn't fail?
			constructor: func() (*ixptypes.MsgCreateOrder, error) {
				return ixptypes.NewMsgCreateOrder(bookID, ixptypes.LimitOrderType_limitBuy, 10, 10, emptyKiraAddr2)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			createOrderMsg, err := tt.constructor()
			require.NoError(t, err)

			result, err = handler(ctx, createOrderMsg)
			require.NoError(t, err)

			orderID := ParseResponseID(result, t)
			order, err := app.IxpKeeper.GetOrderByID(ctx, orderID)
			require.NoError(t, err)

			cancelOrderMsg, err := ixptypes.NewMsgCancelOrder(orderID, kiraAddr1)
			require.NoError(t, err)

			_, err = handler(ctx, cancelOrderMsg)
			require.NoError(t, err)
			order, err = app.IxpKeeper.GetOrderByID(ctx, orderID)
			require.NoError(t, err)

			require.Equal(t, createOrderMsg.OrderType, order.OrderType)
			require.Equal(t, createOrderMsg.OrderBookID, order.OrderBookID)
			require.Equal(t, createOrderMsg.Amount, order.Amount)
			require.Equal(t, createOrderMsg.LimitPrice, order.LimitPrice)
			require.Equal(t, createOrderMsg.ExpiryTime, order.ExpiryTime)
			if createOrderMsg.Curator != nil && createOrderMsg.Curator.Empty() {
				require.Equal(t, emptyKiraAddr2, order.Curator)
			} else {
				require.Equal(t, createOrderMsg.Curator, order.Curator)
			}
			require.Equal(t, order.IsCancelled, true)
		})
	}
}

func TestNewHandler_MsgUpsertSignerKey_HappyPath(t *testing.T) {

	pubKeyText := "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em"
	_, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, pubKeyText)
	require.NoError(t, err)
	emptyKiraAddr1 := types.AccAddress{}
	emptyKiraAddr2 := types.AccAddress(nil)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})
	handler := ixp.NewHandler(app.IxpKeeper)

	tests := []struct {
		name          string
		constructor   func(sdk.AccAddress) (*ixptypes.MsgUpsertSignerKey, error)
		handlerErr    string
		runAgain      bool
		againErr      string
		runDiffAddr   bool
		diffAddrErr   string
		runDiffPubKey bool
	}{
		{
			name: "one permission test",
			constructor: func(addr sdk.AccAddress) (*ixptypes.MsgUpsertSignerKey, error) {
				_, pubKey, err := GenSecp256k1PubKey()
				require.NoError(t, err)
				pubKeyText, err := types.Bech32ifyPubKey(types.Bech32PubKeyTypeAccPub, pubKey)
				require.NoError(t, err)
				return ixptypes.NewMsgUpsertSignerKey(pubKeyText, ixptypes.SignerKeyType_Secp256k1, 0, true, "", []int64{1}, addr)
			},
			runAgain:      true,
			againErr:      "",
			runDiffAddr:   true,
			diffAddrErr:   "this key is owned by another curator already",
			runDiffPubKey: true,
		},
		{
			name: "empty permission test1",
			constructor: func(addr sdk.AccAddress) (*ixptypes.MsgUpsertSignerKey, error) {
				_, pubKey, err := GenSecp256k1PubKey()
				require.NoError(t, err)
				pubKeyText, err := types.Bech32ifyPubKey(types.Bech32PubKeyTypeAccPub, pubKey)
				require.NoError(t, err)
				return ixptypes.NewMsgUpsertSignerKey(pubKeyText, ixptypes.SignerKeyType_Secp256k1, 0, true, "", []int64{}, addr)
			},
		},
		{
			name: "empty permission test2",
			constructor: func(addr sdk.AccAddress) (*ixptypes.MsgUpsertSignerKey, error) {
				_, pubKey, err := GenSecp256k1PubKey()
				require.NoError(t, err)
				pubKeyText, err := types.Bech32ifyPubKey(types.Bech32PubKeyTypeAccPub, pubKey)
				require.NoError(t, err)
				return ixptypes.NewMsgUpsertSignerKey(pubKeyText, ixptypes.SignerKeyType_Secp256k1, 0, true, "", nil, addr)
			},
		},
		{
			name: "empty curator test1",
			constructor: func(addr sdk.AccAddress) (*ixptypes.MsgUpsertSignerKey, error) {
				_, pubKey, err := GenSecp256k1PubKey()
				require.NoError(t, err)
				pubKeyText, err := types.Bech32ifyPubKey(types.Bech32PubKeyTypeAccPub, pubKey)
				require.NoError(t, err)
				return ixptypes.NewMsgUpsertSignerKey(pubKeyText, ixptypes.SignerKeyType_Secp256k1, 0, true, "", []int64{1}, emptyKiraAddr1)
			},
			handlerErr: "curator shouldn't be empty",
		},
		{
			name: "empty curator test2",
			constructor: func(addr sdk.AccAddress) (*ixptypes.MsgUpsertSignerKey, error) {
				_, pubKey, err := GenSecp256k1PubKey()
				require.NoError(t, err)
				pubKeyText, err := types.Bech32ifyPubKey(types.Bech32PubKeyTypeAccPub, pubKey)
				require.NoError(t, err)
				return ixptypes.NewMsgUpsertSignerKey(pubKeyText, ixptypes.SignerKeyType_Secp256k1, 0, true, "", []int64{1}, emptyKiraAddr2)
			},
			handlerErr: "curator shouldn't be empty",
		},
		{
			name: "ed25519 pubKey test",
			constructor: func(addr sdk.AccAddress) (*ixptypes.MsgUpsertSignerKey, error) {
				_, pubKey, err := GenEd25519PubKey()
				require.NoError(t, err)
				pubKeyText, err := types.Bech32ifyPubKey(types.Bech32PubKeyTypeAccPub, pubKey)
				require.NoError(t, err)
				return ixptypes.NewMsgUpsertSignerKey(pubKeyText, ixptypes.SignerKeyType_Secp256k1, 0, true, "", []int64{1}, addr)
			},
		},
	}
	for i, tt := range tests {
		addr := NewAccountByIndex(i)
		theMsg, err := tt.constructor(addr)
		require.NoError(t, err)

		_, err = handler(ctx, theMsg)
		if len(tt.handlerErr) != 0 {
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.handlerErr)
			continue
		} else {
			require.NoError(t, err)
		}

		signerkeys := app.IxpKeeper.GetSignerKeys(ctx, theMsg.Curator)
		require.Len(t, signerkeys, 1)

		signerkey := signerkeys[0]

		require.Equal(t, theMsg.PubKey, signerkey.PubKey)
		require.Equal(t, theMsg.KeyType, signerkey.KeyType)
		require.Equal(t, theMsg.ExpiryTime, signerkey.ExpiryTime)
		if theMsg.Permissions != nil && len(theMsg.Permissions) == 0 {
			require.Equal(t, []int64(nil), signerkey.Permissions)
		} else {
			require.Equal(t, theMsg.Permissions, signerkey.Permissions)
		}
		if theMsg.Curator != nil && theMsg.Curator.Empty() {
			require.Equal(t, emptyKiraAddr2, signerkey.Curator)
		} else {
			require.Equal(t, theMsg.Curator, signerkey.Curator)
		}
		if tt.runAgain {
			_, err = handler(ctx, theMsg)
			if tt.againErr == "" {
				require.NoError(t, err)
				signerkeys := app.IxpKeeper.GetSignerKeys(ctx, theMsg.Curator)
				require.Len(t, signerkeys, 1)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.againErr)
			}
		}
		if tt.runDiffAddr {
			diffMsg := *theMsg                                      // copy message
			diffMsg.Curator = NewAccountByIndex(i + len(tests) + 1) // update curator field
			_, err = handler(ctx, &diffMsg)                         // execute the message
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.diffAddrErr)
		}
		if tt.runDiffPubKey {
			diffMsg, err := tt.constructor(addr)
			require.True(t, diffMsg.Curator.String() == theMsg.Curator.String())
			require.NoError(t, err)
			_, err = handler(ctx, diffMsg)
			require.NoError(t, err)
			signerkeys := app.IxpKeeper.GetSignerKeys(ctx, theMsg.Curator)
			require.Len(t, signerkeys, 2)
		}
	}
}
