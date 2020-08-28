package keeper_test

import (
	"testing"

	"github.com/KiraCore/sekai/app"
	"github.com/cosmos/cosmos-sdk/types"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/KiraCore/sekai/simapp"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	app.SetConfig()
}

func TestKeeper_CreateOrderBook(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	kiraAddr1, err := types.AccAddressFromBech32("kira1da22wd7slpxpptasczs679mr5c8xtucqdzxc3n")
	require.NoError(t, err)

	obID, err := app.IxpKeeper.CreateOrderBook(ctx, "quote", "base", kiraAddr1, "mnemonic")
	require.NoError(t, err)

	t.Log("obID", obID)
}
