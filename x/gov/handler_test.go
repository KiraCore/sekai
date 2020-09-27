package gov_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/x/gov"
	"github.com/KiraCore/sekai/x/gov/types"
)

func TestMain(m *testing.M) {
	app.SetConfig()
	m.Run()
}

// When a network actor has not been saved before, it creates one with default params
// and sets the permissions.
func TestHandler_MsgWhitelistPermissions_ActorDoesNotExist(t *testing.T) {
	proposerAddr, err := sdk.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	tests := []struct {
		name             string
		msg              sdk.Msg
		checkWhitelisted bool
	}{
		{
			"Msg Whitelist Permissions",
			&types.MsgWhitelistPermissions{
				Proposer:   proposerAddr,
				Address:    addr,
				Permission: uint32(types.PermClaimValidator),
			},
			true,
		},
		{
			"Msg Blacklist Permissions",
			&types.MsgBlacklistPermissions{
				Proposer:   proposerAddr,
				Address:    addr,
				Permission: uint32(types.PermClaimValidator),
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			err := setPermissionToAddr(t, app, ctx, proposerAddr, types.PermSetPermissions)
			require.NoError(t, err)

			handler := gov.NewHandler(app.CustomGovKeeper)
			_, err = handler(ctx, tt.msg)
			require.NoError(t, err)

			actor, err := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
			require.NoError(t, err)

			if tt.checkWhitelisted {
				require.True(t, actor.Permissions.IsWhitelisted(types.PermClaimValidator))
			} else {
				require.True(t, actor.Permissions.IsBlacklisted(types.PermClaimValidator))
			}
		})
	}
}

// When a network actor has already permissions it just appends the permission.
func TestNewHandler_SetPermissions_ActorWithPerms(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	proposerAddr, err := sdk.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	tests := []struct {
		name             string
		msg              sdk.Msg
		checkWhitelisted bool
	}{
		{
			name: "actor with Whitelist Permissions",
			msg: &types.MsgWhitelistPermissions{
				Proposer:   proposerAddr,
				Address:    addr,
				Permission: uint32(types.PermClaimValidator),
			},
			checkWhitelisted: true,
		},
		{
			name: "actor with Blacklist Permissions",
			msg: &types.MsgBlacklistPermissions{
				Proposer:   proposerAddr,
				Address:    addr,
				Permission: uint32(types.PermClaimValidator),
			},
			checkWhitelisted: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			err := setPermissionToAddr(t, app, ctx, proposerAddr, types.PermSetPermissions)
			require.NoError(t, err)

			// Add some perms before to the actor.
			actor := types.NewDefaultActor(addr)
			if tt.checkWhitelisted {
				err = actor.Permissions.AddToWhitelist(types.PermSetPermissions)
			} else {
				err = actor.Permissions.AddToBlacklist(types.PermSetPermissions)
			}
			require.NoError(t, err)

			app.CustomGovKeeper.SaveNetworkActor(ctx, actor)

			// Call the handler to add some permissions.
			handler := gov.NewHandler(app.CustomGovKeeper)
			_, err = handler(ctx, tt.msg)
			require.NoError(t, err)

			actor, err = app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
			require.NoError(t, err)

			if tt.checkWhitelisted {
				require.True(t, actor.Permissions.IsWhitelisted(types.PermClaimValidator))
				require.True(t, actor.Permissions.IsWhitelisted(types.PermSetPermissions)) // This permission was already set before callid add permission.
			} else {
				require.True(t, actor.Permissions.IsBlacklisted(types.PermClaimValidator))
				require.True(t, actor.Permissions.IsBlacklisted(types.PermSetPermissions)) // This permission was already set before callid add permission.
			}
		})
	}
}

func TestNewHandler_SetPermissionsWithoutSetPermissions(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	proposerAddr, err := sdk.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	tests := []struct {
		name string
		msg  sdk.Msg
	}{
		{
			name: "MsgWhitelist",
			msg: &types.MsgWhitelistPermissions{
				Proposer:   proposerAddr,
				Address:    addr,
				Permission: uint32(types.PermClaimValidator),
			},
		},
		{
			name: "MsgBlacklist",
			msg: &types.MsgBlacklistPermissions{
				Proposer:   proposerAddr,
				Address:    addr,
				Permission: uint32(types.PermClaimValidator),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			handler := gov.NewHandler(app.CustomGovKeeper)
			_, err = handler(ctx, tt.msg)
			require.EqualError(t, err, "PermSetPermissions: not enough permissions")
		})
	}
}

func TestNewHandler_SetPermissions_ProposerHasRoleSudo(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	proposerAddr, err := sdk.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	tests := []struct {
		name           string
		msg            sdk.Msg
		checkWhitelist bool
	}{
		{
			name: "MsgWhitelist",
			msg: &types.MsgWhitelistPermissions{
				Proposer:   proposerAddr,
				Address:    addr,
				Permission: uint32(types.PermClaimValidator),
			},
			checkWhitelist: true,
		},
		{
			name: "MsgBlacklist",
			msg: &types.MsgBlacklistPermissions{
				Proposer:   proposerAddr,
				Address:    addr,
				Permission: uint32(types.PermClaimValidator),
			},
			checkWhitelist: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			// First we set Role Sudo to proposer Actor
			proposerActor := types.NewDefaultActor(proposerAddr)
			proposerActor.SetRole(types.RoleSudo)
			require.NoError(t, err)
			app.CustomGovKeeper.SaveNetworkActor(ctx, proposerActor)

			handler := gov.NewHandler(app.CustomGovKeeper)
			_, err = handler(ctx, tt.msg)
			require.NoError(t, err)

			actor, err := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
			require.NoError(t, err)

			if tt.checkWhitelist {
				require.True(t, actor.Permissions.IsWhitelisted(types.PermClaimValidator))
			} else {
				require.True(t, actor.Permissions.IsBlacklisted(types.PermClaimValidator))
			}
		})
	}
}

func TestHandler_ClaimCouncilor_Fails(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	tests := []struct {
		name        string
		msg         sdk.Msg
		expectedErr error
	}{
		{
			name: "not enough permissions",
			msg: &types.MsgClaimCouncilor{
				Moniker:  "",
				Website:  "",
				Social:   "",
				Identity: "",
				Address:  addr,
			},
			expectedErr: fmt.Errorf("PermClaimCouncilor: not enough permissions"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			handler := gov.NewHandler(app.CustomGovKeeper)
			_, err := handler(ctx, tt.msg)
			require.EqualError(t, err, tt.expectedErr.Error())
		})
	}
}

func TestHandler_ClaimCouncilor_HappyPath(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	tests := []struct {
		name string
		msg  *types.MsgClaimCouncilor
	}{
		{
			name: "all correct",
			msg: &types.MsgClaimCouncilor{
				Moniker:  "TheMoniker",
				Website:  "TheWebsite",
				Social:   "The Social",
				Identity: "The Identity",
				Address:  addr,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			err = setPermissionToAddr(t, app, ctx, addr, types.PermClaimCouncilor)
			require.NoError(t, err)

			handler := gov.NewHandler(app.CustomGovKeeper)
			_, err := handler(ctx, tt.msg)
			require.NoError(t, err)

			expectedCouncilor := types.NewCouncilor(
				tt.msg.Moniker,
				tt.msg.Website,
				tt.msg.Social,
				tt.msg.Identity,
				tt.msg.Address,
			)

			councilor, err := app.CustomGovKeeper.GetCouncilor(ctx, addr)
			require.NoError(t, err)

			require.Equal(t, expectedCouncilor, councilor)
		})
	}
}

func TestHandler_WhitelistRolePermissions_Errors(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	tests := []struct {
		name         string
		msg          *types.MsgWhitelistRolePermission
		preparePerms func(t *testing.T, app *simapp.SimApp, ctx sdk.Context)
		expectedErr  error
	}{
		{
			name: "address without SetPermissions perm",
			msg: types.NewMsgWhitelistRolePermission(
				addr,
				uint32(types.RoleValidator),
				uint32(types.PermSetPermissions),
			),
			preparePerms: func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {
				return
			},
			expectedErr: fmt.Errorf("PermSetPermissions: not enough permissions"),
		},
		{
			name: "role does not exist",
			msg: types.NewMsgWhitelistRolePermission(
				addr,
				10000,
				1,
			),
			preparePerms: func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {
				err2 := setPermissionToAddr(t, app, ctx, addr, types.PermSetPermissions)
				require.NoError(t, err2)
			},
			expectedErr: fmt.Errorf("role does not exist"),
		},
		{
			name: "permission is blacklisted",
			msg: types.NewMsgWhitelistRolePermission(
				addr,
				uint32(types.RoleValidator),
				uint32(types.PermSetPermissions),
			),
			preparePerms: func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {
				err2 := setPermissionToAddr(t, app, ctx, addr, types.PermSetPermissions)
				require.NoError(t, err2)

				perms := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
				require.NotNil(t, perms)

				err2 = perms.AddToBlacklist(types.PermSetPermissions)
				require.NoError(t, err2)

				app.CustomGovKeeper.SetPermissionsForRole(ctx, types.RoleValidator, perms)
			},
			expectedErr: fmt.Errorf("permission is already blacklisted: error adding to whitelist"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			tt.preparePerms(t, app, ctx)

			handler := gov.NewHandler(app.CustomGovKeeper)
			_, err := handler(ctx, tt.msg)
			require.EqualError(t, err, tt.expectedErr.Error())
		})
	}
}

func TestHandler_WhitelistRolePermissions(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	err = setPermissionToAddr(t, app, ctx, addr, types.PermSetPermissions)
	require.NoError(t, err)

	perms := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
	require.NotNil(t, perms)
	require.False(t, perms.IsWhitelisted(types.PermSetPermissions))

	msg := types.NewMsgWhitelistRolePermission(
		addr,
		uint32(types.RoleValidator),
		uint32(types.PermSetPermissions),
	)

	handler := gov.NewHandler(app.CustomGovKeeper)
	_, err = handler(ctx, msg)
	require.NoError(t, err)

	perms = app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
	require.NotNil(t, perms)
	require.True(t, perms.IsWhitelisted(types.PermSetPermissions))
}

func TestHandler_BlacklistRolePermissions_Errors(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	tests := []struct {
		name         string
		msg          *types.MsgBlacklistRolePermission
		preparePerms func(t *testing.T, app *simapp.SimApp, ctx sdk.Context)
		expectedErr  error
	}{
		{
			name: "address without SetPermissions perm",
			msg: types.NewMsgBlacklistRolePermission(
				addr,
				uint32(types.RoleValidator),
				uint32(types.PermSetPermissions),
			),
			preparePerms: func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {},
			expectedErr:  fmt.Errorf("PermSetPermissions: not enough permissions"),
		},
		{
			name: "role does not exist",
			msg: types.NewMsgBlacklistRolePermission(
				addr,
				10000,
				1,
			),
			preparePerms: func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {
				err2 := setPermissionToAddr(t, app, ctx, addr, types.PermSetPermissions)
				require.NoError(t, err2)
			},
			expectedErr: fmt.Errorf("role does not exist"),
		},
		{
			name: "permission is whitelisted",
			msg: types.NewMsgBlacklistRolePermission(
				addr,
				uint32(types.RoleValidator),
				uint32(types.PermSetPermissions),
			),
			preparePerms: func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {
				err2 := setPermissionToAddr(t, app, ctx, addr, types.PermSetPermissions)
				require.NoError(t, err2)

				perms := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
				require.NotNil(t, perms)

				err2 = perms.AddToWhitelist(types.PermSetPermissions)
				require.NoError(t, err2)

				app.CustomGovKeeper.SetPermissionsForRole(ctx, types.RoleValidator, perms)
			},
			expectedErr: fmt.Errorf("permission is already whitelisted: error adding to blacklist"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			tt.preparePerms(t, app, ctx)

			handler := gov.NewHandler(app.CustomGovKeeper)
			_, err := handler(ctx, tt.msg)
			require.EqualError(t, err, tt.expectedErr.Error())
		})
	}
}

func TestHandler_BlacklistRolePermissions(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	err = setPermissionToAddr(t, app, ctx, addr, types.PermSetPermissions)
	require.NoError(t, err)

	perms := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
	require.NotNil(t, perms)
	require.False(t, perms.IsBlacklisted(types.PermSetPermissions))

	msg := types.NewMsgBlacklistRolePermission(
		addr,
		uint32(types.RoleValidator),
		uint32(types.PermSetPermissions),
	)

	handler := gov.NewHandler(app.CustomGovKeeper)
	_, err = handler(ctx, msg)
	require.NoError(t, err)

	perms = app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
	require.NotNil(t, perms)
	require.True(t, perms.IsBlacklisted(types.PermSetPermissions))
}

func TestHandler_RemoveWhitelistRolePermissions_Errors(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	tests := []struct {
		name         string
		msg          *types.MsgRemoveWhitelistRolePermission
		preparePerms func(t *testing.T, app *simapp.SimApp, ctx sdk.Context)
		expectedErr  error
	}{
		{
			name: "address without SetPermissions perm",
			msg: types.NewMsgRemoveWhitelistRolePermission(
				addr,
				uint32(types.RoleValidator),
				uint32(types.PermSetPermissions),
			),
			preparePerms: func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {},
			expectedErr:  fmt.Errorf("PermSetPermissions: not enough permissions"),
		},
		{
			name: "role does not exist",
			msg: types.NewMsgRemoveWhitelistRolePermission(
				addr,
				10000,
				1,
			),
			preparePerms: func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {
				err2 := setPermissionToAddr(t, app, ctx, addr, types.PermSetPermissions)
				require.NoError(t, err2)
			},
			expectedErr: fmt.Errorf("role does not exist"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			tt.preparePerms(t, app, ctx)

			handler := gov.NewHandler(app.CustomGovKeeper)
			_, err := handler(ctx, tt.msg)
			require.EqualError(t, err, tt.expectedErr.Error())
		})
	}
}

func TestHandler_RemoveWhitelistRolePermissions(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	err = setPermissionToAddr(t, app, ctx, addr, types.PermSetPermissions)
	require.NoError(t, err)

	perms := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
	require.NotNil(t, perms)
	require.True(t, perms.IsWhitelisted(types.PermClaimValidator))

	msg := types.NewMsgRemoveWhitelistRolePermission(
		addr,
		uint32(types.RoleValidator),
		uint32(types.PermClaimValidator),
	)

	handler := gov.NewHandler(app.CustomGovKeeper)
	_, err = handler(ctx, msg)
	require.NoError(t, err)

	perms = app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
	require.NotNil(t, perms)
	require.False(t, perms.IsWhitelisted(types.PermClaimValidator))
}

func TestHandler_RemoveBlacklistRolePermissions_Errors(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	tests := []struct {
		name         string
		msg          *types.MsgRemoveBlacklistRolePermission
		preparePerms func(t *testing.T, app *simapp.SimApp, ctx sdk.Context)
		expectedErr  error
	}{
		{
			name: "address without SetPermissions perm",
			msg: types.NewMsgRemoveBlacklistRolePermission(
				addr,
				uint32(types.RoleValidator),
				uint32(types.PermSetPermissions),
			),
			preparePerms: func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {},
			expectedErr:  fmt.Errorf("PermSetPermissions: not enough permissions"),
		},
		{
			name: "role does not exist",
			msg: types.NewMsgRemoveBlacklistRolePermission(
				addr,
				10000,
				1,
			),
			preparePerms: func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {
				err2 := setPermissionToAddr(t, app, ctx, addr, types.PermSetPermissions)
				require.NoError(t, err2)
			},
			expectedErr: fmt.Errorf("role does not exist"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			tt.preparePerms(t, app, ctx)

			handler := gov.NewHandler(app.CustomGovKeeper)
			_, err := handler(ctx, tt.msg)
			require.EqualError(t, err, tt.expectedErr.Error())
		})
	}
}

func TestHandler_RemoveBlacklistRolePermissions(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	err = setPermissionToAddr(t, app, ctx, addr, types.PermSetPermissions)
	require.NoError(t, err)

	perms := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
	require.NotNil(t, perms)

	// Set some blacklist value
	err = perms.AddToBlacklist(types.PermClaimCouncilor)
	require.NoError(t, err)
	app.CustomGovKeeper.SetPermissionsForRole(ctx, types.RoleValidator, perms)

	// Check if it is blacklisted.
	perms = app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
	require.NotNil(t, perms)
	require.True(t, perms.IsBlacklisted(types.PermClaimCouncilor))

	msg := types.NewMsgRemoveBlacklistRolePermission(
		addr,
		uint32(types.RoleValidator),
		uint32(types.PermClaimCouncilor),
	)

	handler := gov.NewHandler(app.CustomGovKeeper)
	_, err = handler(ctx, msg)
	require.NoError(t, err)

	perms = app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
	require.NotNil(t, perms)
	require.False(t, perms.IsBlacklisted(types.PermClaimCouncilor))
}

func TestHandler_CreateRole_Errors(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	tests := []struct {
		name         string
		msg          *types.MsgCreateRole
		preparePerms func(t *testing.T, app *simapp.SimApp, ctx sdk.Context)
		expectedErr  error
	}{
		{
			"fails when no perms",
			types.NewMsgCreateRole(
				addr,
				10,
			),
			func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {},
			fmt.Errorf("PermSetPermissions: not enough permissions"),
		},
		{
			"fails when role already exists",
			types.NewMsgCreateRole(
				addr,
				1234,
			),
			func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {
				err2 := setPermissionToAddr(t, app, ctx, addr, types.PermSetPermissions)
				require.NoError(t, err2)
				app.CustomGovKeeper.SetPermissionsForRole(ctx, types.Role(1234), types.NewPermissions(nil, nil))
			},
			fmt.Errorf("role already exist"),
		},
	}

	for _, tt := range tests {
		app := simapp.Setup(false)
		ctx := app.NewContext(false, tmproto.Header{})

		tt.preparePerms(t, app, ctx)

		handler := gov.NewHandler(app.CustomGovKeeper)
		_, err := handler(ctx, tt.msg)
		require.EqualError(t, err, tt.expectedErr.Error())
	}
}

func TestHandler_CreateRole(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	err = setPermissionToAddr(t, app, ctx, addr, types.PermSetPermissions)
	require.NoError(t, err)

	role := app.CustomGovKeeper.GetPermissionsForRole(ctx, 1234)
	require.Nil(t, role)

	handler := gov.NewHandler(app.CustomGovKeeper)
	_, err = handler(ctx, types.NewMsgCreateRole(
		addr,
		1234,
	))
	require.NoError(t, err)

	role = app.CustomGovKeeper.GetPermissionsForRole(ctx, 1234)
	require.NotNil(t, role)
}

func TestHandler_AssignRole_Errors(t *testing.T) {
	proposerAddr, err := sdk.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	tests := []struct {
		name         string
		msg          *types.MsgAssignRole
		preparePerms func(t *testing.T, app *simapp.SimApp, ctx sdk.Context)
		expectedErr  error
	}{
		{
			"fails when no perms",
			types.NewMsgAssignRole(
				proposerAddr, addr, 3,
			),
			func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {},
			fmt.Errorf("PermSetPermissions: not enough permissions"),
		},
		{
			"fails when role does not exist",
			types.NewMsgAssignRole(
				proposerAddr, addr, 3,
			),
			func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {
				err2 := setPermissionToAddr(t, app, ctx, proposerAddr, types.PermSetPermissions)
				require.NoError(t, err2)
			},
			types.ErrRoleDoesNotExist,
		},
		{
			"role already assigned",
			types.NewMsgAssignRole(
				proposerAddr, addr, 3,
			),
			func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {
				err2 := setPermissionToAddr(t, app, ctx, proposerAddr, types.PermSetPermissions)
				require.NoError(t, err2)

				app.CustomGovKeeper.SetPermissionsForRole(ctx, types.Role(3), types.NewPermissions([]types.PermValue{
					types.PermClaimValidator,
				}, nil))

				networkActor := types.NewDefaultActor(addr)
				networkActor.SetRole(3)
				app.CustomGovKeeper.SaveNetworkActor(ctx, networkActor)
			},
			types.ErrRoleAlreadyAssigned,
		},
	}

	for _, tt := range tests {
		app := simapp.Setup(false)
		ctx := app.NewContext(false, tmproto.Header{})

		tt.preparePerms(t, app, ctx)

		handler := gov.NewHandler(app.CustomGovKeeper)
		_, err := handler(ctx, tt.msg)
		require.EqualError(t, err, tt.expectedErr.Error())
	}
}

func TestHandler_AssignRole(t *testing.T) {
	proposerAddr, err := sdk.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	// Set permissions to proposer.
	err = setPermissionToAddr(t, app, ctx, proposerAddr, types.PermSetPermissions)
	require.NoError(t, err)

	// Create role
	app.CustomGovKeeper.SetPermissionsForRole(ctx, types.Role(3), types.NewPermissions([]types.PermValue{types.PermSetPermissions}, nil))

	msg := types.NewMsgAssignRole(proposerAddr, addr, 3)

	handler := gov.NewHandler(app.CustomGovKeeper)
	_, err = handler(ctx, msg)
	require.NoError(t, err)

	actor, err := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
	require.NoError(t, err)

	require.True(t, actor.HasRole(types.Role(3)))
}

func TestHandler_RemoveRole_Errors(t *testing.T) {
	proposerAddr, err := sdk.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	tests := []struct {
		name         string
		msg          *types.MsgRemoveRole
		preparePerms func(t *testing.T, app *simapp.SimApp, ctx sdk.Context)
		expectedErr  error
	}{
		{
			"fails when no perms",
			types.NewMsgRemoveRole(
				proposerAddr, addr, 3,
			),
			func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {},
			fmt.Errorf("PermSetPermissions: not enough permissions"),
		},
		{
			"fails when role does not exist",
			types.NewMsgRemoveRole(
				proposerAddr, addr, 3,
			),
			func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {
				err2 := setPermissionToAddr(t, app, ctx, proposerAddr, types.PermSetPermissions)
				require.NoError(t, err2)
			},
			types.ErrRoleDoesNotExist,
		},
		{
			"role not assigned",
			types.NewMsgRemoveRole(
				proposerAddr, addr, 3,
			),
			func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {
				err2 := setPermissionToAddr(t, app, ctx, proposerAddr, types.PermSetPermissions)
				require.NoError(t, err2)

				app.CustomGovKeeper.SetPermissionsForRole(ctx, types.Role(3), types.NewPermissions([]types.PermValue{
					types.PermClaimValidator,
				}, nil))

				networkActor := types.NewDefaultActor(addr)
				app.CustomGovKeeper.SaveNetworkActor(ctx, networkActor)
			},
			types.ErrRoleNotAssigned,
		},
	}

	for _, tt := range tests {
		app := simapp.Setup(false)
		ctx := app.NewContext(false, tmproto.Header{})

		tt.preparePerms(t, app, ctx)

		handler := gov.NewHandler(app.CustomGovKeeper)
		_, err := handler(ctx, tt.msg)
		require.EqualError(t, err, tt.expectedErr.Error())
	}
}

func TestHandler_RemoveRoles(t *testing.T) {
	proposerAddr, err := sdk.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	// Set permissions to proposer.
	err = setPermissionToAddr(t, app, ctx, proposerAddr, types.PermSetPermissions)
	require.NoError(t, err)

	// Set new role and set permission to actor.
	app.CustomGovKeeper.SetPermissionsForRole(ctx, types.Role(3), types.NewPermissions([]types.PermValue{types.PermSetPermissions}, nil))
	actor := types.NewDefaultActor(addr)
	actor.SetRole(3)
	app.CustomGovKeeper.SaveNetworkActor(ctx, actor)

	actor, err = app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
	require.NoError(t, err)
	require.True(t, actor.HasRole(3))

	msg := types.NewMsgRemoveRole(proposerAddr, addr, 3)
	handler := gov.NewHandler(app.CustomGovKeeper)
	_, err = handler(ctx, msg)
	require.NoError(t, err)

	actor, err = app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
	require.NoError(t, err)

	require.False(t, actor.HasRole(3))
}

func TestHandler_ProposalAssignPermission_Errors(t *testing.T) {
	proposerAddr, err := sdk.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	tests := []struct {
		name         string
		msg          *types.MsgProposalAssignPermission
		preparePerms func(t *testing.T, app *simapp.SimApp, ctx sdk.Context)
		expectedErr  error
	}{
		{
			"NetworkActor does not exist",
			types.NewMsgProposalAssignPermission(
				proposerAddr, addr, 3,
			),
			func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {},
			fmt.Errorf("network actor not found"),
		},
		//{
		//	"fails when role does not exist",
		//	types.NewMsgProposalAssignPermission(
		//		proposerAddr, addr, 3,
		//	),
		//	func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {
		//		err2 := setPermissionToAddr(t, app, ctx, proposerAddr, types.PermSetPermissions)
		//		require.NoError(t, err2)
		//	},
		//	types.ErrRoleDoesNotExist,
		//},
		//{
		//	"role not assigned",
		//	types.NewMsgProposalAssignPermission(
		//		proposerAddr, addr, 3,
		//	),
		//	func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {
		//		err2 := setPermissionToAddr(t, app, ctx, proposerAddr, types.PermSetPermissions)
		//		require.NoError(t, err2)
		//
		//		app.CustomGovKeeper.SetPermissionsForRole(ctx, types.Role(3), types.NewPermissions([]types.PermValue{
		//			types.PermClaimValidator,
		//		}, nil))
		//
		//		networkActor := types.NewDefaultActor(addr)
		//		app.CustomGovKeeper.SaveNetworkActor(ctx, networkActor)
		//	},
		//	types.ErrRoleNotAssigned,
		//},
	}

	for _, tt := range tests {
		app := simapp.Setup(false)
		ctx := app.NewContext(false, tmproto.Header{})

		tt.preparePerms(t, app, ctx)

		handler := gov.NewHandler(app.CustomGovKeeper)
		_, err := handler(ctx, tt.msg)
		require.EqualError(t, err, tt.expectedErr.Error())
	}
}

func setPermissionToAddr(t *testing.T, app *simapp.SimApp, ctx sdk.Context, addr sdk.AccAddress, perm types.PermValue) error {
	proposerActor := types.NewDefaultActor(addr)
	err := proposerActor.Permissions.AddToWhitelist(perm)
	require.NoError(t, err)

	app.CustomGovKeeper.SaveNetworkActor(ctx, proposerActor)

	return nil
}
