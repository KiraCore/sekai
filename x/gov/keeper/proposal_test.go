package keeper_test

import (
	"testing"

	"github.com/KiraCore/sekai/x/gov/types"
	types2 "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"

	"github.com/KiraCore/sekai/simapp"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestDefaultProposalIdAtDefaultGenesis(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	proposalID, err := app.CustomGovKeeper.GetNextProposalID(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(1), proposalID)
}

func TestSaveProposalReturnsTheProposalID_AndIncreasesLast(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	proposalID, err := app.CustomGovKeeper.GetNextProposalID(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(1), proposalID)

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, types2.TokensFromConsensusPower(10))
	addr := addrs[0]

	proposal := types.NewProposalAssignPermission(addr, types.PermClaimValidator)
	saveProposal, err := app.CustomGovKeeper.SaveProposal(ctx, proposal)
	require.NoError(t, err)

	require.Equal(t, uint64(1), saveProposal)

	// nextProposalID should be 2
	proposalID, err = app.CustomGovKeeper.GetNextProposalID(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(2), proposalID)
}
