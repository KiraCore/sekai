package types_test

import (
	"testing"

	customgovtypes "github.com/KiraCore/sekai/x/gov/types"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

//
// NETWORK ACTOR
//
func TestNewNetworkActor_SetRole(t *testing.T) {
	addr, err := types.AccAddressFromBech32("kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r")
	require.NoError(t, err)

	actor := customgovtypes.NewDefaultActor(addr)
	require.False(t, actor.HasRole(customgovtypes.RoleValidator))

	actor.SetRole(customgovtypes.RoleValidator)

	require.True(t, actor.HasRole(customgovtypes.RoleValidator))
}

func TestNewNetworkActor_RemoveRole(t *testing.T) {
	addr, err := types.AccAddressFromBech32("kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r")
	require.NoError(t, err)

	actor := customgovtypes.NewDefaultActor(addr)
	actor.SetRole(customgovtypes.RoleValidator)
	actor.SetRole(customgovtypes.RoleSudo)
	require.True(t, actor.HasRole(customgovtypes.RoleValidator))
	require.True(t, actor.HasRole(customgovtypes.RoleSudo))

	actor.RemoveRole(customgovtypes.RoleSudo)
	require.True(t, actor.HasRole(customgovtypes.RoleValidator))
	require.False(t, actor.HasRole(customgovtypes.RoleSudo))
}

func TestNewNetworkActor_Status(t *testing.T) {
	addr, err := types.AccAddressFromBech32("kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r")
	require.NoError(t, err)

	actor := customgovtypes.NewDefaultActor(addr)
	require.Equal(t, customgovtypes.Undefined, actor.Status)

	// Active Actor
	actor = customgovtypes.NewNetworkActor(
		addr,
		customgovtypes.Roles{},
		customgovtypes.Active,
		[]customgovtypes.VoteOption{},
		customgovtypes.NewPermissions(nil, nil),
		1,
	)

	require.True(t, actor.IsActive())
}

func TestNewDefaultActor_CanVote(t *testing.T) {
	actor := customgovtypes.NewNetworkActor(
		types.AccAddress{0x0},
		customgovtypes.Roles{},
		customgovtypes.Active,
		[]customgovtypes.VoteOption{customgovtypes.OptionYes, customgovtypes.OptionAbstain},
		nil,
		123,
	)

	require.True(t, actor.CanVote(customgovtypes.OptionYes))
	require.True(t, actor.CanVote(customgovtypes.OptionAbstain))
	require.False(t, actor.CanVote(customgovtypes.OptionNo))
	require.False(t, actor.CanVote(customgovtypes.OptionNoWithVeto))
}

func TestGetVetoActorsFromList(t *testing.T) {
	actors := []customgovtypes.NetworkActor{
		customgovtypes.NewNetworkActor(types.AccAddress{0x0}, customgovtypes.Roles{}, customgovtypes.Active, []customgovtypes.VoteOption{customgovtypes.OptionAbstain}, nil, 123),
		customgovtypes.NewNetworkActor(types.AccAddress{0x1}, customgovtypes.Roles{}, customgovtypes.Active, []customgovtypes.VoteOption{customgovtypes.OptionAbstain, customgovtypes.OptionNoWithVeto}, nil, 123),
		customgovtypes.NewNetworkActor(types.AccAddress{0x2}, customgovtypes.Roles{}, customgovtypes.Active, []customgovtypes.VoteOption{customgovtypes.OptionAbstain}, nil, 123),
		customgovtypes.NewNetworkActor(types.AccAddress{0x3}, customgovtypes.Roles{}, customgovtypes.Active, []customgovtypes.VoteOption{customgovtypes.OptionAbstain}, nil, 123),
		customgovtypes.NewNetworkActor(types.AccAddress{0x4}, customgovtypes.Roles{}, customgovtypes.Active, []customgovtypes.VoteOption{customgovtypes.OptionAbstain}, nil, 123),
		customgovtypes.NewNetworkActor(types.AccAddress{0x5}, customgovtypes.Roles{}, customgovtypes.Active, []customgovtypes.VoteOption{customgovtypes.OptionAbstain}, nil, 123),
		customgovtypes.NewNetworkActor(types.AccAddress{0x6}, customgovtypes.Roles{}, customgovtypes.Active, []customgovtypes.VoteOption{customgovtypes.OptionAbstain}, nil, 123),
		customgovtypes.NewNetworkActor(types.AccAddress{0x7}, customgovtypes.Roles{}, customgovtypes.Active, []customgovtypes.VoteOption{customgovtypes.OptionAbstain, customgovtypes.OptionNoWithVeto}, nil, 123),
		customgovtypes.NewNetworkActor(types.AccAddress{0x8}, customgovtypes.Roles{}, customgovtypes.Active, []customgovtypes.VoteOption{customgovtypes.OptionAbstain}, nil, 123),
		customgovtypes.NewNetworkActor(types.AccAddress{0x9}, customgovtypes.Roles{}, customgovtypes.Active, []customgovtypes.VoteOption{customgovtypes.OptionAbstain}, nil, 123),
	}

	actorsWithVeto := customgovtypes.GetActorsWithVoteWithVeto(actors)
	require.Equal(t, actorsWithVeto[0], actors[1])
	require.Equal(t, actorsWithVeto[1], actors[7])
	require.Len(t, actorsWithVeto, 2)
}
