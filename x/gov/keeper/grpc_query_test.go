package keeper_test

import (
	"testing"
	"time"

	stakingtypes "github.com/KiraCore/sekai/x/staking/types"

	"github.com/stretchr/testify/require"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
		[]uint64{},
		1,
		nil,
		permissions,
		123,
	)

	app.CustomGovKeeper.SaveNetworkActor(ctx, networkActor)

	querier := app.CustomGovKeeper
	resp, err := querier.PermissionsByAddress(sdk.WrapSDKContext(ctx), &types.PermissionsByAddressRequest{Addr: addr1.String()})
	require.NoError(t, err)

	require.Equal(t, permissions, resp.Permissions)

	// Get permissions by address that is not saved.
	_, err = querier.PermissionsByAddress(sdk.WrapSDKContext(ctx), &types.PermissionsByAddressRequest{Addr: addr2.String()})
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
		[]uint64{
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

	resp, err := querier.RolesByAddress(sdk.WrapSDKContext(ctx), &types.RolesByAddressRequest{Addr: addr1.String()})
	require.NoError(t, err)

	require.Equal(t,
		[]uint64{0x1, 0x2, 0x3},
		resp.RoleIds,
	)

	// Get roles for actor that does not exist
	_, err = querier.RolesByAddress(sdk.WrapSDKContext(ctx), &types.RolesByAddressRequest{Addr: addr2.String()})
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
		types.NewWhitelistAccountPermissionProposal(
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
		addr1,
		types.CouncilorActive,
	)

	app.CustomGovKeeper.SaveCouncilor(ctx, councilor)

	querier := app.CustomGovKeeper

	resp, err := querier.CouncilorByAddress(
		sdk.WrapSDKContext(ctx),
		&types.CouncilorByAddressRequest{Addr: addr1.String()},
	)
	require.NoError(t, err)
	require.Equal(t, councilor, resp.Councilor)

	// Non existing Councilor
	resp, err = querier.CouncilorByAddress(
		sdk.WrapSDKContext(ctx),
		&types.CouncilorByAddressRequest{Addr: addr2.String()},
	)
	require.Error(t, types.ErrCouncilorNotFound)
}

func TestQuerier_CouncilorQueries(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	addr1 := addrs[0]
	addr2 := addrs[1]

	councilor := types.NewCouncilor(
		addr1,
		types.CouncilorActive,
	)

	app.CustomGovKeeper.SaveCouncilor(ctx, councilor)

	networkActor := types.NewNetworkActor(
		addr2,
		[]uint64{types.RoleSudo},
		1,
		nil,
		types.NewPermissions(
			[]types.PermValue{
				types.PermClaimValidator,
			},
			[]types.PermValue{
				types.PermClaimCouncilor,
			},
		),
		123,
	)
	app.CustomGovKeeper.SaveNetworkActor(ctx, networkActor)
	for _, perm := range networkActor.Permissions.Whitelist {
		app.CustomGovKeeper.SetWhitelistAddressPermKey(ctx, networkActor, types.PermValue(perm))
	}
	for _, role := range networkActor.Roles {
		app.CustomGovKeeper.AssignRoleToActor(ctx, networkActor, role)
	}

	querier := app.CustomGovKeeper

	// specific councilor
	resp, err := querier.Councilors(
		sdk.WrapSDKContext(ctx),
		&types.QueryCouncilors{Address: addr1.String()},
	)
	require.NoError(t, err)
	require.Len(t, resp.Councilors, 1)
	require.Equal(t, councilor, resp.Councilors[0])

	// all councilors
	resp, err = querier.Councilors(
		sdk.WrapSDKContext(ctx),
		&types.QueryCouncilors{},
	)
	require.NoError(t, err)
	require.Len(t, resp.Councilors, 1)
	require.Equal(t, councilor, resp.Councilors[0])

	// non-councilors
	nresp, err := querier.NonCouncilors(
		sdk.WrapSDKContext(ctx),
		&types.QueryNonCouncilors{},
	)
	require.NoError(t, err)
	require.Len(t, nresp.NonCouncilors, 1)
	require.Equal(t, networkActor, nresp.NonCouncilors[0])

	wresp, err := querier.AddressesByWhitelistedPermission(
		sdk.WrapSDKContext(ctx),
		&types.QueryAddressesByWhitelistedPermission{Permission: uint32(types.PermClaimValidator)},
	)
	require.NoError(t, err)
	require.Len(t, wresp.Addresses, 1)
	require.Equal(t, addr2.String(), wresp.Addresses[0])

	bresp, err := querier.AddressesByBlacklistedPermission(
		sdk.WrapSDKContext(ctx),
		&types.QueryAddressesByBlacklistedPermission{Permission: uint32(types.PermClaimCouncilor)},
	)
	require.NoError(t, err)
	require.Len(t, bresp.Addresses, 1)
	require.Equal(t, addr2.String(), bresp.Addresses[0])

	rresp, err := querier.AddressesByWhitelistedRole(
		sdk.WrapSDKContext(ctx),
		&types.QueryAddressesByWhitelistedRole{Role: uint32(types.RoleSudo)},
	)
	require.NoError(t, err)
	require.Len(t, rresp.Addresses, 1)
	require.Equal(t, addr2.String(), rresp.Addresses[0])
}
