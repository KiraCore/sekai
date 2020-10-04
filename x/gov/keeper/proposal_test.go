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

	proposal := types.NewProposalAssignPermission(1, addr, types.PermClaimValidator, ctx.BlockTime())
	err = app.CustomGovKeeper.SaveProposal(ctx, proposal)
	require.NoError(t, err)

	// nextProposalID should be 2
	proposalID, err = app.CustomGovKeeper.GetNextProposalID(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(2), proposalID)

	// Get proposal
	savedProposal, found := app.CustomGovKeeper.GetProposal(ctx, proposal.ProposalId)
	require.True(t, found)
	require.Equal(t, proposal, savedProposal)
}

func TestKeeper_SaveVote(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, types2.TokensFromConsensusPower(10))
	addr := addrs[0]

	// Vote not saved yet
	_, found := app.CustomGovKeeper.GetVote(ctx, 1, addr)
	require.False(t, found)

	vote := types.NewVote(1, addr, types.OptionAbstain)

	app.CustomGovKeeper.SaveVote(ctx, vote)

	savedVote, found := app.CustomGovKeeper.GetVote(ctx, 1, addr)
	require.True(t, found)
	require.Equal(t, vote, savedVote)
}
