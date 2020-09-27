package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KiraCore/sekai/simapp"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestDefaultProposalIdAtDefaultGenesis(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	proposalID, err := app.CustomGovKeeper.GetProposalID(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(1), proposalID)
}
