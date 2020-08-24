package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
)

func TestPermissions_IsBlacklisted(t *testing.T) {
	perms := Permissions{
		Blacklist: []uint32{
			uint32(PermClaimValidator),
		},
		Whitelist: nil,
	}

	require.True(t, perms.IsBlacklisted(PermClaimValidator))
	require.False(t, perms.IsBlacklisted(PermClaimGovernanceSeat))
}

func TestPermissions_IsWhitelisted(t *testing.T) {
	perms := NewPermissions([]PermValue{PermClaimGovernanceSeat}, nil)

	require.True(t, perms.IsWhitelisted(PermClaimGovernanceSeat))
	require.False(t, perms.IsWhitelisted(PermClaimValidator))
}

func TestPermissions_AddWhitelist(t *testing.T) {
	perms := NewPermissions(nil, nil)

	require.False(t, perms.IsWhitelisted(PermClaimGovernanceSeat))
	perms.AddToWhitelist(PermClaimGovernanceSeat)
	require.True(t, perms.IsWhitelisted(PermClaimGovernanceSeat))
}

func TestPermissions_AddBlacklist(t *testing.T) {
	perms := NewPermissions(nil, nil)

	require.False(t, perms.IsBlacklisted(PermClaimGovernanceSeat))
	perms.AddToBlacklist(PermClaimGovernanceSeat)
	require.True(t, perms.IsBlacklisted(PermClaimGovernanceSeat))
}

func TestNewNetworkActor_HasPermission(t *testing.T) {
	tests := []struct {
		name             string
		perms            Permissions
		permissionToTest PermValue
		expectedHasPerm  bool
	}{
		{
			"permission is whitelisted",
			NewPermissions([]PermValue{PermClaimGovernanceSeat}, nil),
			PermClaimGovernanceSeat,
			true,
		},
		{
			"permission is blacklisted",
			NewPermissions(nil, []PermValue{PermClaimGovernanceSeat}),
			PermClaimGovernanceSeat,
			false,
		},
		{
			"not blacklisted or whitelisted",
			NewPermissions(nil, nil),
			PermClaimGovernanceSeat,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actor := NewNetworkActor(
				types.AccAddress("hola"),
				nil,
				1,
				nil,
				tt.perms,
				123,
			)

			hasPerm := actor.HasPermissionFor(tt.permissionToTest)
			require.Equal(t, tt.expectedHasPerm, hasPerm)
		})
	}
}

func TestNetworkActor_AddPermission(t *testing.T) {
	actor := NewNetworkActor(
		types.AccAddress("hola"),
		nil,
		1,
		nil,
		NewPermissions(nil, nil),
		123,
	)

	require.False(t, actor.HasPermissionFor(PermClaimGovernanceSeat))

	actor.AddPermission(PermClaimGovernanceSeat)
	require.True(t, actor.HasPermissionFor(PermClaimGovernanceSeat))
}
