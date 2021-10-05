package types_test

import (
	"testing"

	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

//
// NETWORK ACTOR
//
func TestNewNetworkActor_SetRole(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r")
	require.NoError(t, err)

	actor := types.NewDefaultActor(addr)
	require.False(t, actor.HasRole(types.RoleValidator))

	actor.SetRole(types.RoleValidator)

	require.True(t, actor.HasRole(types.RoleValidator))
}

func TestNewNetworkActor_RemoveRole(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r")
	require.NoError(t, err)

	actor := types.NewDefaultActor(addr)
	actor.SetRole(types.RoleValidator)
	actor.SetRole(types.RoleSudo)
	require.True(t, actor.HasRole(types.RoleValidator))
	require.True(t, actor.HasRole(types.RoleSudo))

	actor.RemoveRole(types.RoleSudo)
	require.True(t, actor.HasRole(types.RoleValidator))
	require.False(t, actor.HasRole(types.RoleSudo))
}

func TestNewNetworkActor_Status(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r")
	require.NoError(t, err)

	actor := types.NewDefaultActor(addr)
	require.Equal(t, types.Active, actor.Status)

	// Active Actor
	actor = types.NewNetworkActor(
		addr,
		types.Roles{},
		types.Active,
		[]types.VoteOption{},
		types.NewPermissions(nil, nil),
		1,
	)

	require.True(t, actor.IsActive())
	actor.Deactivate()
	require.False(t, actor.IsActive())
	require.True(t, actor.IsInactive())
}

func TestNewDefaultActor_CanVote(t *testing.T) {
	actor := types.NewNetworkActor(
		sdk.AccAddress{0x0},
		types.Roles{},
		types.Active,
		[]types.VoteOption{types.OptionYes, types.OptionAbstain},
		nil,
		123,
	)

	require.True(t, actor.CanVote(types.OptionYes))
	require.True(t, actor.CanVote(types.OptionAbstain))
	require.False(t, actor.CanVote(types.OptionNo))
	require.False(t, actor.CanVote(types.OptionNoWithVeto))
}

func TestGetVetoActorsFromList(t *testing.T) {
	actors := []types.NetworkActor{
		types.NewNetworkActor(sdk.AccAddress{0x0}, types.Roles{}, types.Active, []types.VoteOption{types.OptionAbstain}, nil, 123),
		types.NewNetworkActor(sdk.AccAddress{0x1}, types.Roles{}, types.Active, []types.VoteOption{types.OptionAbstain, types.OptionNoWithVeto}, nil, 123),
		types.NewNetworkActor(sdk.AccAddress{0x2}, types.Roles{}, types.Active, []types.VoteOption{types.OptionAbstain}, nil, 123),
		types.NewNetworkActor(sdk.AccAddress{0x3}, types.Roles{}, types.Active, []types.VoteOption{types.OptionAbstain}, nil, 123),
		types.NewNetworkActor(sdk.AccAddress{0x4}, types.Roles{}, types.Active, []types.VoteOption{types.OptionAbstain}, nil, 123),
		types.NewNetworkActor(sdk.AccAddress{0x5}, types.Roles{}, types.Active, []types.VoteOption{types.OptionAbstain}, nil, 123),
		types.NewNetworkActor(sdk.AccAddress{0x6}, types.Roles{}, types.Active, []types.VoteOption{types.OptionAbstain}, nil, 123),
		types.NewNetworkActor(sdk.AccAddress{0x7}, types.Roles{}, types.Active, []types.VoteOption{types.OptionAbstain, types.OptionNoWithVeto}, nil, 123),
		types.NewNetworkActor(sdk.AccAddress{0x8}, types.Roles{}, types.Active, []types.VoteOption{types.OptionAbstain}, nil, 123),
		types.NewNetworkActor(sdk.AccAddress{0x9}, types.Roles{}, types.Active, []types.VoteOption{types.OptionAbstain}, nil, 123),
	}

	actorsWithVeto := types.GetActorsWithVoteWithVeto(actors)
	require.Equal(t, actorsWithVeto[0], actors[1])
	require.Equal(t, actorsWithVeto[1], actors[7])
	require.Len(t, actorsWithVeto, 2)
}
