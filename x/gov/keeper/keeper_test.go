package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/x/gov/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
)

func TestKeeper_SetNetworkProperty(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	app.CustomGovKeeper.SetNetworkProperties(ctx, &types.NetworkProperties{
		MinTxFee: 100,
		MaxTxFee: 50000,
	})

	err := app.CustomGovKeeper.SetNetworkProperty(ctx, types.MinTxFee, types.NetworkPropertyValue{Value: 300})
	require.Nil(t, err)

	savedMinTxFee, err := app.CustomGovKeeper.GetNetworkProperty(ctx, types.MinTxFee)
	require.Nil(t, err)
	require.Equal(t, uint64(300), savedMinTxFee.Value)
}

func TestKeeper_EnsureUniqueKeys(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	app.CustomGovKeeper.SetIdentityRecord(ctx, types.IdentityRecord{
		Id:      1,
		Address: "addr1",
		Key:     "nickname",
		Value:   "jack",
	})
	notUniqueKey := app.CustomGovKeeper.EnsureUniqueKeys(ctx, "", "nickname")
	require.Equal(t, notUniqueKey, "")
	app.CustomGovKeeper.SetIdentityRecord(ctx, types.IdentityRecord{
		Id:      2,
		Address: "addr2",
		Key:     "nickname",
		Value:   "jack",
	})
	notUniqueKey = app.CustomGovKeeper.EnsureUniqueKeys(ctx, "", "nickname")
	require.Equal(t, notUniqueKey, "nickname")
}

func TestKeeper_EnsureOldUniqueKeysNotRemoved(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	removedOldKey := app.CustomGovKeeper.EnsureOldUniqueKeysNotRemoved(ctx, "", "nickname")
	require.Equal(t, removedOldKey, "")

	removedOldKey = app.CustomGovKeeper.EnsureOldUniqueKeysNotRemoved(ctx, "nickname", "")
	require.Equal(t, removedOldKey, "nickname")
}
