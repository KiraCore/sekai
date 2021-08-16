package keeper_test

import (
	"testing"
	"time"

	stakingtypes "github.com/KiraCore/sekai/x/staking/types"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/x/gov/types"
)

func TestQuerier_PermissionsByAddress(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	addr1 := addrs[0]
	addr2 := addrs[1]

	permissions := types.NewPermissions(
		[]types.PermValue{
			types.PermClaimValidator,
		},
		nil,
	)

	networkActor := types.NewNetworkActor(
		addr1,
		types.Roles{},
		1,
		nil,
		permissions,
		123,
	)

	app.CustomGovKeeper.SaveNetworkActor(ctx, networkActor)

	querier := app.CustomGovKeeper
	resp, err := querier.PermissionsByAddress(sdk.WrapSDKContext(ctx), &types.PermissionsByAddressRequest{ValAddr: addr1})
	require.NoError(t, err)

	require.Equal(t, permissions, resp.Permissions)

	// Get permissions by address that is not saved.
	_, err = querier.PermissionsByAddress(sdk.WrapSDKContext(ctx), &types.PermissionsByAddressRequest{ValAddr: addr2})
	require.EqualError(t, err, stakingtypes.ErrNetworkActorNotFound.Error())
}

func TestQuerier_RolesByAddress(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	addr1 := addrs[0]
	addr2 := addrs[1]

	networkActor := types.NewNetworkActor(
		addr1,
		types.Roles{
			1, 2, 3,
		},
		1,
		nil,
		types.NewPermissions(
			[]types.PermValue{
				types.PermClaimValidator,
			},
			nil,
		),
		123,
	)

	app.CustomGovKeeper.SaveNetworkActor(ctx, networkActor)

	querier := app.CustomGovKeeper

	resp, err := querier.RolesByAddress(sdk.WrapSDKContext(ctx), &types.RolesByAddressRequest{ValAddr: addr1})
	require.NoError(t, err)

	require.Equal(t,
		[]uint64{0x1, 0x2, 0x3},
		resp.Roles,
	)

	// Get roles for actor that does not exist
	_, err = querier.RolesByAddress(sdk.WrapSDKContext(ctx), &types.RolesByAddressRequest{ValAddr: addr2})
	require.EqualError(t, err, "network actor not found")
}

func TestQuerier_Proposal(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))

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

	app.CustomGovKeeper.SaveVote(ctx, types.Vote{
		ProposalId: proposalID,
		Voter:      addrs[0],
		Option:     types.OptionNo,
	})

	querier := app.CustomGovKeeper

	resp, err := querier.Proposal(
		sdk.WrapSDKContext(ctx),
		&types.QueryProposalRequest{ProposalId: proposalID},
	)
	require.NoError(t, err)
	require.Equal(t, proposalID, resp.Proposal.ProposalId)
	require.Len(t, resp.Votes, 1)
	require.Equal(t, proposal.Description, resp.Proposal.Description)
}

func TestQuerier_CouncilorByAddress(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	addr1 := addrs[0]
	addr2 := addrs[1]

	councilor := types.NewCouncilor(
		"TheMoniker",
		addr1,
	)

	app.CustomGovKeeper.SaveCouncilor(ctx, councilor)

	querier := app.CustomGovKeeper

	resp, err := querier.CouncilorByAddress(
		sdk.WrapSDKContext(ctx),
		&types.CouncilorByAddressRequest{ValAddr: addr1},
	)
	require.NoError(t, err)
	require.Equal(t, councilor, resp.Councilor)

	// Councilor by Moniker
	resp, err = querier.CouncilorByMoniker(
		sdk.WrapSDKContext(ctx),
		&types.CouncilorByMonikerRequest{
			Moniker: councilor.Moniker,
		},
	)
	require.NoError(t, err)
	require.Equal(t, councilor, resp.Councilor)

	// Non existing Councilor
	resp, err = querier.CouncilorByAddress(
		sdk.WrapSDKContext(ctx),
		&types.CouncilorByAddressRequest{ValAddr: addr2},
	)
	require.Error(t, types.ErrCouncilorNotFound)
}
