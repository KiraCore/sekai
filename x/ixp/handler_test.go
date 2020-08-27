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
