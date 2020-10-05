package keeper_test

import (
	"testing"
	"time"

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

	proposal := types.NewProposalAssignPermission(1, addr, types.PermClaimValidator, ctx.BlockTime(), ctx.BlockTime().Add(10*time.Minute))
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

func TestKeeper_AddProposalToActiveQueue(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, types2.TokensFromConsensusPower(10))
	addr := addrs[0]

	baseEndTime := time.Now()
	for _, i := range []uint64{1, 2, 3} {
		endTime := baseEndTime.Add(time.Second * time.Duration(i))

		proposal := types.NewProposalAssignPermission(
			i,
			addr,
			types.PermSetPermissions,
			baseEndTime,
			endTime,
		)

		err := app.CustomGovKeeper.SaveProposal(ctx, proposal)
		require.NoError(t, err)
		app.CustomGovKeeper.AddToActiveProposals(ctx, proposal)
	}

	// We only get until endtime of the second proposal.
	iterator := app.CustomGovKeeper.GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx, baseEndTime.Add(2*time.Second))
	defer iterator.Close()

	totalIdsFound := 0
	for ; iterator.Valid(); iterator.Next() {
		totalIdsFound++
	}

	require.Equal(t, 2, totalIdsFound)
}
