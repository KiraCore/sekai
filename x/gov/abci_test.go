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

func TestEndBlocker(t *testing.T) {
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
			name: "proposal passes",
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

					// Only 4 first users vote yes. We reach quorum but not half of the votes are yes.
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			addrs := tt.prepareScenario(app, ctx)

			ctx = ctx.WithBlockTime(time.Now().Add(time.Second * 10)) // We make that proposal passes.

			gov.EndBlocker(ctx, app.CustomGovKeeper)

			tt.validateScenario(t, app, ctx, addrs)
		})
	}
}
