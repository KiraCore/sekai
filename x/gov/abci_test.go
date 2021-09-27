package gov_test

import (
	"testing"
	"time"

	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/x/gov"
	"github.com/KiraCore/sekai/x/gov/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestEndBlocker_ActiveProposal(t *testing.T) {
	tests := []struct {
		name              string
		prepareScenario   func(app *simapp.SekaiApp, ctx sdk.Context) []sdk.AccAddress
		validateScenario  func(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context, addrs []sdk.AccAddress)
		blockHeightChange int64
	}{
		{
			name: "proposal passes: min block height for proposal voting time not reached",
			prepareScenario: func(app *simapp.SekaiApp, ctx sdk.Context) []sdk.AccAddress {
				addrs := simapp.AddTestAddrsIncremental(app, ctx, 10, sdk.NewInt(100))

				proposalID := uint64(1234)
				proposal, err := types.NewProposal(
					proposalID,
					"title",
					"some desc",
					types.NewAssignPermissionProposal(
						addrs[0],
						types.PermSetPermissions,
					),
					time.Now(),
					time.Now().Add(10*time.Second),
					time.Now().Add(20*time.Second),
					ctx.BlockHeight()+2,
					ctx.BlockHeight()+3,
				)
				require.NoError(t, err)

				app.CustomGovKeeper.SaveProposal(ctx, proposal)
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
			validateScenario: func(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context, addrs []sdk.AccAddress) {
				actor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addrs[0])
				require.True(t, found)
				require.False(t, actor.Permissions.IsWhitelisted(types.PermSetPermissions))

				// We check that is not in the ActiveProposals
				iterator := app.CustomGovKeeper.GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx, time.Now().Add(15*time.Second))
				requireIteratorCount(t, iterator, 1)

				// And it is in the EnactmentProposals
				iterator = app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 0)

				proposal, found := app.CustomGovKeeper.GetProposal(ctx, 1234)
				require.True(t, found)
				require.Equal(t, types.Pending, proposal.Result)
			},
			blockHeightChange: 1,
		},
		{
			name: "proposal passes: quorum not reached",
			prepareScenario: func(app *simapp.SekaiApp, ctx sdk.Context) []sdk.AccAddress {
				addrs := simapp.AddTestAddrsIncremental(app, ctx, 10, sdk.NewInt(100))

				proposalID := uint64(1234)
				proposal, err := types.NewProposal(
					proposalID,
					"title",
					"some desc",
					types.NewAssignPermissionProposal(
						addrs[0],
						types.PermSetPermissions,
					),
					time.Now(),
					time.Now(),
					time.Now(),
					ctx.BlockHeight()+2,
					ctx.BlockHeight()+3,
				)
				require.NoError(t, err)

				app.CustomGovKeeper.SaveProposal(ctx, proposal)
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
			validateScenario: func(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context, addrs []sdk.AccAddress) {
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
				require.Equal(t, types.QuorumNotReached, proposal.Result)
			},
			blockHeightChange: 3,
		},
		{
			name: "proposal passes and joins Enactment place",
			prepareScenario: func(app *simapp.SekaiApp, ctx sdk.Context) []sdk.AccAddress {
				addrs := simapp.AddTestAddrsIncremental(app, ctx, 10, sdk.NewInt(100))

				proposalID := uint64(1234)
				proposal, err := types.NewProposal(
					proposalID,
					"title",
					"some desc",
					types.NewAssignPermissionProposal(
						addrs[0],
						types.PermSetPermissions,
					),
					time.Now(),
					time.Now().Add(10*time.Second),
					time.Now().Add(20*time.Second),
					ctx.BlockHeight()+2,
					ctx.BlockHeight()+3,
				)
				require.NoError(t, err)

				app.CustomGovKeeper.SaveProposal(ctx, proposal)
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
			validateScenario: func(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context, addrs []sdk.AccAddress) {
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
				require.Equal(t, types.Enactment, proposal.Result)
			},
			blockHeightChange: 3,
		},
		{
			name: "Passed proposal in enactment is applied and min block height for enactment not reached",
			prepareScenario: func(app *simapp.SekaiApp, ctx sdk.Context) []sdk.AccAddress {
				addrs := simapp.AddTestAddrsIncremental(app, ctx, 10, sdk.NewInt(100))

				actor := types.NewDefaultActor(addrs[0])
				app.CustomGovKeeper.SaveNetworkActor(ctx, actor)

				proposalID := uint64(1234)
				proposal, err := types.NewProposal(
					proposalID,
					"title",
					"some desc",
					types.NewAssignPermissionProposal(
						addrs[0],
						types.PermSetPermissions,
					),
					time.Now(),
					time.Now().Add(10*time.Second),
					time.Now().Add(20*time.Second),
					ctx.BlockHeight()+2,
					ctx.BlockHeight()+3,
				)
				require.NoError(t, err)

				proposal.Result = types.Enactment
				app.CustomGovKeeper.SaveProposal(ctx, proposal)

				app.CustomGovKeeper.AddToEnactmentProposals(ctx, proposal)

				iterator := app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 1)

				return addrs
			},
			validateScenario: func(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context, addrs []sdk.AccAddress) {
				iterator := app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 1)
			},
			blockHeightChange: 0,
		},
		{
			name: "Passed proposal in enactment is applied and removed from enactment list: Assign permission",
			prepareScenario: func(app *simapp.SekaiApp, ctx sdk.Context) []sdk.AccAddress {
				addrs := simapp.AddTestAddrsIncremental(app, ctx, 10, sdk.NewInt(100))

				actor := types.NewDefaultActor(addrs[0])
				app.CustomGovKeeper.SaveNetworkActor(ctx, actor)

				proposalID := uint64(1234)
				proposal, err := types.NewProposal(
					proposalID,
					"title",
					"some desc",
					types.NewAssignPermissionProposal(
						addrs[0],
						types.PermSetPermissions,
					),
					time.Now(),
					time.Now().Add(10*time.Second),
					time.Now().Add(20*time.Second),
					ctx.BlockHeight()+2,
					ctx.BlockHeight()+3,
				)
				require.NoError(t, err)

				proposal.Result = types.Enactment
				app.CustomGovKeeper.SaveProposal(ctx, proposal)

				app.CustomGovKeeper.AddToEnactmentProposals(ctx, proposal)

				iterator := app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 1)

				return addrs
			},
			validateScenario: func(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context, addrs []sdk.AccAddress) {
				iterator := app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 0)

				actor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addrs[0])
				require.True(t, found)

				require.True(t, actor.Permissions.IsWhitelisted(types.PermSetPermissions))
			},
			blockHeightChange: 3,
		},
		{
			name: "Passed proposal in enactment is applied and removed from enactment list, actor does not exist",
			prepareScenario: func(app *simapp.SekaiApp, ctx sdk.Context) []sdk.AccAddress {
				addrs := simapp.AddTestAddrsIncremental(app, ctx, 10, sdk.NewInt(100))

				proposalID := uint64(1234)
				proposal, err := types.NewProposal(
					proposalID,
					"title",
					"some desc",
					types.NewAssignPermissionProposal(
						addrs[0],
						types.PermSetPermissions,
					),
					time.Now(),
					time.Now().Add(10*time.Second),
					time.Now().Add(20*time.Second),
					ctx.BlockHeight()+2,
					ctx.BlockHeight()+3,
				)
				require.NoError(t, err)

				proposal.Result = types.Enactment
				app.CustomGovKeeper.SaveProposal(ctx, proposal)

				app.CustomGovKeeper.AddToEnactmentProposals(ctx, proposal)

				iterator := app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 1)

				return addrs
			},
			validateScenario: func(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context, addrs []sdk.AccAddress) {
				iterator := app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 0)

				actor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addrs[0])
				require.True(t, found)

				require.True(t, actor.Permissions.IsWhitelisted(types.PermSetPermissions))
			},
			blockHeightChange: 3,
		},
		{
			name: "Passed proposal in enactment is applied and removed from enactment list: Upsert Data Registry",
			prepareScenario: func(app *simapp.SekaiApp, ctx sdk.Context) []sdk.AccAddress {
				addrs := simapp.AddTestAddrsIncremental(app, ctx, 10, sdk.NewInt(100))

				actor := types.NewDefaultActor(addrs[0])
				app.CustomGovKeeper.SaveNetworkActor(ctx, actor)

				proposalID := uint64(1234)
				proposal, err := types.NewProposal(
					proposalID,
					"title",
					"some desc",
					types.NewUpsertDataRegistryProposal(
						"theKey",
						"theHash",
						"theReference",
						"theEncoding",
						1234,
					),
					time.Now(),
					time.Now().Add(10*time.Second),
					time.Now().Add(20*time.Second),
					ctx.BlockHeight()+2,
					ctx.BlockHeight()+3,
				)
				require.NoError(t, err)

				proposal.Result = types.Enactment
				app.CustomGovKeeper.SaveProposal(ctx, proposal)

				app.CustomGovKeeper.AddToEnactmentProposals(ctx, proposal)

				iterator := app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 1)

				return addrs
			},
			validateScenario: func(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context, addrs []sdk.AccAddress) {
				iterator := app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 0)

				entry, found := app.CustomGovKeeper.GetDataRegistryEntry(ctx, "theKey")
				require.True(t, found)

				require.Equal(t, "theHash", entry.Hash)
				require.Equal(t, "theEncoding", entry.Encoding)
				require.Equal(t, "theReference", entry.Reference)
				require.Equal(t, uint64(1234), entry.Size_)
			},
			blockHeightChange: 3,
		},
		{
			name: "Passed proposal in enactment is applied and removed from enactment list: Set Network Property",
			prepareScenario: func(app *simapp.SekaiApp, ctx sdk.Context) []sdk.AccAddress {
				addrs := simapp.AddTestAddrsIncremental(app, ctx, 10, sdk.NewInt(100))

				actor := types.NewDefaultActor(addrs[0])
				app.CustomGovKeeper.SaveNetworkActor(ctx, actor)

				proposalID := uint64(1234)
				proposal, err := types.NewProposal(
					proposalID,
					"title",
					"some desc",
					types.NewSetNetworkPropertyProposal(
						types.MinTxFee,
						types.NetworkPropertyValue{Value: 300},
					),
					time.Now(),
					time.Now().Add(10*time.Second),
					time.Now().Add(20*time.Second),
					ctx.BlockHeight()+2,
					ctx.BlockHeight()+3,
				)
				require.NoError(t, err)

				proposal.Result = types.Enactment
				app.CustomGovKeeper.SaveProposal(ctx, proposal)

				app.CustomGovKeeper.AddToEnactmentProposals(ctx, proposal)

				iterator := app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 1)

				return addrs
			},
			validateScenario: func(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context, addrs []sdk.AccAddress) {
				iterator := app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 0)

				minTxFee, err := app.CustomGovKeeper.GetNetworkProperty(ctx, types.MinTxFee)
				require.NoError(t, err)

				require.Equal(t, uint64(300), minTxFee.Value)
			},
			blockHeightChange: 3,
		},
		{
			name: "Passed proposal in enactment is applied and removed from enactment list: Set Token Alias",
			prepareScenario: func(app *simapp.SekaiApp, ctx sdk.Context) []sdk.AccAddress {
				addrs := simapp.AddTestAddrsIncremental(app, ctx, 10, sdk.NewInt(100))

				actor := types.NewDefaultActor(addrs[0])
				app.CustomGovKeeper.SaveNetworkActor(ctx, actor)

				proposalID := uint64(1234)
				proposal, err := types.NewProposal(
					proposalID,
					"title",
					"some desc",
					tokenstypes.NewUpsertTokenAliasProposal(
						"EUR",
						"Euro",
						"http://www.google.es",
						12,
						[]string{
							"eur",
							"€",
						},
					),
					time.Now(),
					time.Now().Add(10*time.Second),
					time.Now().Add(20*time.Second),
					ctx.BlockHeight()+2,
					ctx.BlockHeight()+3,
				)
				require.NoError(t, err)

				proposal.Result = types.Enactment
				app.CustomGovKeeper.SaveProposal(ctx, proposal)

				app.CustomGovKeeper.AddToEnactmentProposals(ctx, proposal)

				iterator := app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 1)

				return addrs
			},
			validateScenario: func(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context, addrs []sdk.AccAddress) {
				iterator := app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 0)

				token := app.TokensKeeper.GetTokenAlias(ctx, "EUR")
				require.Equal(t, "Euro", token.Name)
			},
			blockHeightChange: 3,
		},
		{
			name: "Passed proposal in enactment is applied and removed from enactment list: Set Token Rates",
			prepareScenario: func(app *simapp.SekaiApp, ctx sdk.Context) []sdk.AccAddress {
				addrs := simapp.AddTestAddrsIncremental(app, ctx, 10, sdk.NewInt(100))

				actor := types.NewDefaultActor(addrs[0])
				app.CustomGovKeeper.SaveNetworkActor(ctx, actor)

				proposalID := uint64(1234)
				proposal, err := types.NewProposal(
					proposalID,
					"title",
					"some desc",
					tokenstypes.NewUpsertTokenRatesProposal(
						"btc",
						sdk.NewDec(1234),
						false,
					),
					time.Now(),
					time.Now().Add(10*time.Second),
					time.Now().Add(20*time.Second),
					ctx.BlockHeight()+2,
					ctx.BlockHeight()+3,
				)
				require.NoError(t, err)

				proposal.Result = types.Enactment
				app.CustomGovKeeper.SaveProposal(ctx, proposal)

				app.CustomGovKeeper.AddToEnactmentProposals(ctx, proposal)

				iterator := app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 1)

				return addrs
			},
			validateScenario: func(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context, addrs []sdk.AccAddress) {
				iterator := app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 0)

				token := app.TokensKeeper.GetTokenRate(ctx, "btc")
				require.Equal(t, sdk.NewDec(1234), token.Rate)
				require.Equal(t, "btc", token.Denom)
				require.Equal(t, false, token.FeePayments)
			},
			blockHeightChange: 3,
		},
		{
			name: "Passed proposal in enactment is applied and removed from enactment list: Unjail Validator",
			prepareScenario: func(app *simapp.SekaiApp, ctx sdk.Context) []sdk.AccAddress {
				addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.NewInt(100))
				valAddr := sdk.ValAddress(addrs[0])
				pubkeys := simapp.CreateTestPubKeys(1)
				pubKey := pubkeys[0]

				actor := types.NewDefaultActor(addrs[0])
				app.CustomGovKeeper.SaveNetworkActor(ctx, actor)

				val, err := stakingtypes.NewValidator(valAddr, pubKey)
				require.NoError(t, err)
				app.CustomStakingKeeper.AddValidator(ctx, val)
				err = app.CustomStakingKeeper.Jail(ctx, val.ValKey)
				require.NoError(t, err)

				proposalID := uint64(1234)
				proposal, err := types.NewProposal(
					proposalID,
					"title",
					"some desc",
					stakingtypes.NewUnjailValidatorProposal(
						addrs[0],
						"theHash",
						"theProposal",
					),
					time.Now(),
					time.Now().Add(10*time.Second),
					time.Now().Add(20*time.Second),
					ctx.BlockHeight()+2,
					ctx.BlockHeight()+3,
				)
				require.NoError(t, err)

				proposal.Result = types.Enactment
				app.CustomGovKeeper.SaveProposal(ctx, proposal)

				app.CustomGovKeeper.AddToEnactmentProposals(ctx, proposal)

				iterator := app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 1)

				return addrs
			},
			validateScenario: func(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context, addrs []sdk.AccAddress) {
				iterator := app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 0)

				validator, err := app.CustomStakingKeeper.GetValidator(ctx, sdk.ValAddress(addrs[0]))
				require.NoError(t, err)

				require.False(t, validator.IsJailed())
			},
			blockHeightChange: 3,
		},
		{
			name: "Passed proposal in enactment is applied and removed from enactment list: Create Role",
			prepareScenario: func(app *simapp.SekaiApp, ctx sdk.Context) []sdk.AccAddress {
				addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.NewInt(100))

				actor := types.NewDefaultActor(addrs[0])
				app.CustomGovKeeper.SaveNetworkActor(ctx, actor)

				proposalID := uint64(1234)
				proposal, err := types.NewProposal(
					proposalID,
					"title",
					"some desc",
					types.NewCreateRoleProposal(
						types.Role(1000),
						[]types.PermValue{
							types.PermClaimValidator,
						},
						[]types.PermValue{
							types.PermChangeTxFee,
						},
					),
					time.Now(),
					time.Now().Add(10*time.Second),
					time.Now().Add(20*time.Second),
					ctx.BlockHeight()+2,
					ctx.BlockHeight()+3,
				)
				require.NoError(t, err)

				proposal.Result = types.Enactment
				app.CustomGovKeeper.SaveProposal(ctx, proposal)

				app.CustomGovKeeper.AddToEnactmentProposals(ctx, proposal)

				iterator := app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 1)

				_, found := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.Role(1000))
				require.False(t, found)

				return addrs
			},
			validateScenario: func(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context, addrs []sdk.AccAddress) {
				iterator := app.CustomGovKeeper.GetEnactmentProposalsWithFinishedEnactmentEndTimeIterator(ctx, time.Now().Add(25*time.Second))
				requireIteratorCount(t, iterator, 0)

				perms, found := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.Role(1000))
				require.True(t, found)
				require.True(t, perms.IsWhitelisted(types.PermClaimValidator))
				require.True(t, perms.IsBlacklisted(types.PermChangeTxFee))
			},
			blockHeightChange: 3,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})
			ctx = ctx.WithBlockTime(time.Now())

			addrs := tt.prepareScenario(app, ctx)

			ctx = ctx.WithBlockTime(time.Now().Add(time.Second * 25))
			ctx = ctx.WithBlockHeight(ctx.BlockHeight() + tt.blockHeightChange)

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
