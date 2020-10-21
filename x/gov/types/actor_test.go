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
		[]uint32{},
		customgovtypes.NewPermissions(nil, nil),
		1,
	)

	require.True(t, actor.IsActive())
}
