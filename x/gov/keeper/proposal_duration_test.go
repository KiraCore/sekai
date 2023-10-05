package keeper_test

import (
	"testing"

	simapp "github.com/KiraCore/sekai/app"
	kiratypes "github.com/KiraCore/sekai/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/require"
)

func TestProposalDurationSetGet(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	// get not specifically define type
	duration := app.CustomGovKeeper.GetProposalDuration(ctx, kiratypes.ProposalTypeSetProposalDurations)
	require.Equal(t, duration, uint64(0))

	// try to set correct value
	err := app.CustomGovKeeper.SetProposalDuration(ctx, kiratypes.ProposalTypeSetProposalDurations, 2400)
	require.NoError(t, err)

	duration = app.CustomGovKeeper.GetProposalDuration(ctx, kiratypes.ProposalTypeSetProposalDurations)
	require.Equal(t, duration, uint64(2400))

	// try setting again with lower than minimum value
	err = app.CustomGovKeeper.SetProposalDuration(ctx, kiratypes.ProposalTypeSetProposalDurations, 1)
	require.Error(t, err)

	duration = app.CustomGovKeeper.GetProposalDuration(ctx, kiratypes.ProposalTypeSetProposalDurations)
	require.Equal(t, duration, uint64(2400))

	// check get all functionality
	allDurations := app.CustomGovKeeper.GetAllProposalDurations(ctx)
	require.Equal(t, len(allDurations), 1)
}
