package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/x/gov/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestKeeper_UpsertDataRegistryEntry(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	entry := types.NewDataRegistryEntry(
		"someHAsh",
		"someURL",
		"someEncoding",
		1234,
	)

	app.CustomGovKeeper.UpsertDataRegistryEntry(ctx, "CodeOfConduct", entry)

	savedDataRegistry, found := app.CustomGovKeeper.GetDataRegistryEntry(ctx, "CodeOfConduct")
	require.True(t, found)

	require.Equal(t, entry, savedDataRegistry)

	_, found = app.CustomGovKeeper.GetDataRegistryEntry(ctx, "NonExistingKey")
	require.False(t, found)
}
