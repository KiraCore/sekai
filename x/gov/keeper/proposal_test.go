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

	proposal := types.NewProposalAssignPermission(1, addr, types.PermClaimValidator, ctx.BlockTime(), ctx.BlockTime().Add(10*time.Minute), ctx.BlockTime().Add(20*time.Minute))
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
			endTime,
		)

		err := app.CustomGovKeeper.SaveProposal(ctx, proposal)
		require.NoError(t, err)
		app.CustomGovKeeper.AddToActiveProposals(ctx, proposal)
	}

	// We only get until endtime of the second proposal.
	iterator := app.CustomGovKeeper.GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx, baseEndTime.Add(2*time.Second))
	defer iterator.Close()
	requireIteratorCount(t, iterator, 2)

	// We remove one ActiveProposal, the first
	proposal, found := app.CustomGovKeeper.GetProposal(ctx, 1)
	require.True(t, found)
	app.CustomGovKeeper.RemoveActiveProposal(ctx, proposal)

	// We then only get 1 proposal.
	iterator = app.CustomGovKeeper.GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx, baseEndTime.Add(2*time.Second))
	defer iterator.Close()
	requireIteratorCount(t, iterator, 1)
}

func TestKeeper_AddProposalToEnactmentQueue(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, types2.TokensFromConsensusPower(10))
	addr := addrs[0]

	baseEndTime := time.Now()
	for _, i := range []uint64{1, 2, 3} {
		enactmentEndTime := baseEndTime.Add(time.Duration(i) * time.Second)
		proposal := types.NewProposalAssignPermission(
			i,
			addr,
			types.PermSetPermissions,
			baseEndTime,
			baseEndTime,
			enactmentEndTime,
		)

		err := app.CustomGovKeeper.SaveProposal(ctx, proposal)
		require.NoError(t, err)

		app.CustomGovKeeper.AddToEnactmentProposals(ctx, proposal)
	}

	// We only get until endtime of the second proposal.
	iterator := app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, baseEndTime.Add(2*time.Second))
	defer iterator.Close()
	requireIteratorCount(t, iterator, 2)

	// We remove one Proposal from the Enactment list, the first
	proposal, found := app.CustomGovKeeper.GetProposal(ctx, 1)
	require.True(t, found)
	app.CustomGovKeeper.RemoveEnactmentProposal(ctx, proposal)

	// We then only get 1 proposal.
	iterator = app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, baseEndTime.Add(2*time.Second))
	defer iterator.Close()
	requireIteratorCount(t, iterator, 1)
}

func TestKeeper_GetProposalVotesIterator(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, types2.TokensFromConsensusPower(10))
	addr1 := addrs[0]
	addr2 := addrs[1]

	proposal1 := types.NewProposalAssignPermission(1, addr1, types.PermSetPermissions, time.Now(), time.Now().Add(1*time.Second), time.Now().Add(10*time.Second))
	proposal2 := types.NewProposalAssignPermission(2, addr2, types.PermClaimCouncilor, time.Now(), time.Now().Add(1*time.Second), time.Now().Add(10*time.Second))

	err := app.CustomGovKeeper.SaveProposal(ctx, proposal1)
	require.NoError(t, err)
	err = app.CustomGovKeeper.SaveProposal(ctx, proposal2)
	require.NoError(t, err)

	// 1st proposal has 2 votes
	vote1 := types.NewVote(proposal1.ProposalId, addr1, types.OptionYes)
	vote2 := types.NewVote(proposal1.ProposalId, addr2, types.OptionYes)
	app.CustomGovKeeper.SaveVote(ctx, vote1)
	app.CustomGovKeeper.SaveVote(ctx, vote2)

	// 2nd proposal has 1 vote
	v1 := types.NewVote(proposal2.ProposalId, addr1, types.OptionYes)
	app.CustomGovKeeper.SaveVote(ctx, v1)

	// We iterate the 1st proposal
	iterator := app.CustomGovKeeper.GetProposalVotesIterator(ctx, proposal1.ProposalId)
	require.True(t, iterator.Valid())
	totalVotes := 0
	for ; iterator.Valid(); iterator.Next() {
		totalVotes++
	}
	require.Equal(t, 2, totalVotes)

	// We iterate the 2nd proposal
	iterator = app.CustomGovKeeper.GetProposalVotesIterator(ctx, proposal2.ProposalId)
	require.True(t, iterator.Valid())
	totalVotes = 0
	for ; iterator.Valid(); iterator.Next() {
		totalVotes++
	}
	require.Equal(t, 1, totalVotes)
}
