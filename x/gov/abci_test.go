package gov_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/x/gov"
	"github.com/KiraCore/sekai/x/gov/types"
)

func TestEndBlocker_ActiveProposal(t *testing.T) {
	tests := []struct {
		name             string
		prepareScenario  func(app *simapp.SimApp, ctx sdk.Context) []sdk.AccAddress
		validateScenario func(t *testing.T, app *simapp.SimApp, ctx sdk.Context, addrs []sdk.AccAddress)
	}{
		{
			name: "proposal passes: quorum not reached",
			prepareScenario: func(app *simapp.SimApp, ctx sdk.Context) []sdk.AccAddress {
				addrs := simapp.AddTestAddrsIncremental(app, ctx, 10, sdk.NewInt(100))

				proposalID := uint64(1234)
				proposal := types.NewProposalAssignPermission(
					proposalID,
					addrs[0],
					types.PermSetPermissions,
					time.Now(),
					time.Now(),
					time.Now(),
				)

				err := app.CustomGovKeeper.SaveProposal(ctx, proposal)
				require.NoError(t, err)
				app.CustomGovKeeper.AddToActiveProposals(ctx, proposal)

				// We set permissions to Vote The proposal to all the actors. 10 in total.
				for i, addr := range addrs {
					actor := types.NewDefaultActor(addr)
					err := app.CustomGovKeeper.AddWhitelistPermission(ctx, actor, types.PermVoteSetPermissionProposal)
					require.NoError(t, err)

					// Only 3 first users vote yes. We dont reach Quorum.
					if i < 3 {
						vote := types.NewVote(proposalID, addr, types.OptionYes)
						app.CustomGovKeeper.SaveVote(ctx, vote)
					}
				}

				return addrs
			},
			validateScenario: func(t *testing.T, app *simapp.SimApp, ctx sdk.Context, addrs []sdk.AccAddress) {
				actor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addrs[0])
				require.True(t, found)
				require.False(t, actor.Permissions.IsWhitelisted(types.PermSetPermissions))
			},
		},
		{
			name: "proposal passes and joins Enactment place",
			prepareScenario: func(app *simapp.SimApp, ctx sdk.Context) []sdk.AccAddress {
				addrs := simapp.AddTestAddrsIncremental(app, ctx, 10, sdk.NewInt(100))

				proposalID := uint64(1234)
				proposal := types.NewProposalAssignPermission(
					proposalID,
					addrs[0],
					types.PermSetPermissions,
					time.Now(),
					time.Now().Add(10*time.Second),
					time.Now().Add(20*time.Second),
				)

				err := app.CustomGovKeeper.SaveProposal(ctx, proposal)
				require.NoError(t, err)
				app.CustomGovKeeper.AddToActiveProposals(ctx, proposal)

				// We set permissions to Vote The proposal to all the actors. 10 in total.
				for i, addr := range addrs {
					actor := types.NewDefaultActor(addr)
					err := app.CustomGovKeeper.AddWhitelistPermission(ctx, actor, types.PermVoteSetPermissionProposal)
					require.NoError(t, err)

					// Only 4 first users vote yes. We reach quorum but not half of the votes are yes.
					if i < 4 {
						vote := types.NewVote(proposalID, addr, types.OptionYes)
						app.CustomGovKeeper.SaveVote(ctx, vote)
					}
				}

				iterator := app.CustomGovKeeper.GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx, time.Now().Add(10*time.Second))
				requireIteratorCount(t, iterator, 1)

				iterator = app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 0)

				return addrs
			},
			validateScenario: func(t *testing.T, app *simapp.SimApp, ctx sdk.Context, addrs []sdk.AccAddress) {
				actor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addrs[0])
				require.True(t, found)
				require.False(t, actor.Permissions.IsWhitelisted(types.PermSetPermissions))

				// We check that is not in the ActiveProposals
				iterator := app.CustomGovKeeper.GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx, time.Now().Add(15*time.Second))
				requireIteratorCount(t, iterator, 0)

				// And it is in the EnactmentProposals
				iterator = app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 1)

				proposal, found := app.CustomGovKeeper.GetProposal(ctx, 1234)
				require.True(t, found)
				require.Equal(t, types.Passed, proposal.Result)
			},
		},
		{
			name: "Passed proposal in enactment is applied and removed from enactment list",
			prepareScenario: func(app *simapp.SimApp, ctx sdk.Context) []sdk.AccAddress {
				addrs := simapp.AddTestAddrsIncremental(app, ctx, 10, sdk.NewInt(100))

				actor := types.NewDefaultActor(addrs[0])
				app.CustomGovKeeper.SaveNetworkActor(ctx, actor)

				proposalID := uint64(1234)
				proposal := types.NewProposalAssignPermission(
					proposalID,
					addrs[0],
					types.PermSetPermissions,
					time.Now(),
					time.Now().Add(10*time.Second),
					time.Now().Add(20*time.Second),
				)

				proposal.Result = types.Passed
				err := app.CustomGovKeeper.SaveProposal(ctx, proposal)
				require.NoError(t, err)

				app.CustomGovKeeper.AddToEnactmentProposals(ctx, proposal)

				iterator := app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 1)

				return addrs
			},
			validateScenario: func(t *testing.T, app *simapp.SimApp, ctx sdk.Context, addrs []sdk.AccAddress) {
				iterator := app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 0)

				actor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addrs[0])
				require.True(t, found)

				require.True(t, actor.Permissions.IsWhitelisted(types.PermSetPermissions))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			addrs := tt.prepareScenario(app, ctx)

			ctx = ctx.WithBlockTime(time.Now().Add(time.Second * 25))

			gov.EndBlocker(ctx, app.CustomGovKeeper)

			tt.validateScenario(t, app, ctx, addrs)
		})
	}
}

func requireIteratorCount(t *testing.T, iterator sdk.Iterator, expectedCount int) {
	c := 0
	for ; iterator.Valid(); iterator.Next() {
		c++
	}

	require.Equal(t, expectedCount, c)
}
