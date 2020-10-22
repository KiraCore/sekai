package gov_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/types/errors"

	types2 "github.com/cosmos/cosmos-sdk/x/gov/types"

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
	os.Exit(m.Run())
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

			actor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
			require.True(t, found)

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

			actor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
			require.True(t, found)

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

			actor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
			require.True(t, found)

			if tt.checkWhitelist {
				require.True(t, actor.Permissions.IsWhitelisted(types.PermClaimValidator))
			} else {
				require.True(t, actor.Permissions.IsBlacklisted(types.PermClaimValidator))
			}
		})
	}
}

func TestNewHandler_SetNetworkProperties(t *testing.T) {
	changeFeeAddr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	sudoAddr, err := sdk.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	tests := []struct {
		name       string
		msg        sdk.Msg
		desiredErr string
	}{
		{
			name: "Success run with ChangeTxFee permission",
			msg: &types.MsgSetNetworkProperties{
				NetworkProperties: &types.NetworkProperties{
					MinTxFee: 100,
					MaxTxFee: 1000,
				},
				Proposer: changeFeeAddr,
			},
			desiredErr: "",
		},
		{
			name: "Failure run without ChangeTxFee permission",
			msg: &types.MsgSetNetworkProperties{
				NetworkProperties: &types.NetworkProperties{
					MinTxFee: 100,
					MaxTxFee: 1000,
				},
				Proposer: sudoAddr,
			},
			desiredErr: "not enough permissions",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})
			// First we set Role Sudo to proposer Actor
			proposerActor := types.NewDefaultActor(sudoAddr)
			proposerActor.SetRole(types.RoleSudo)
			require.NoError(t, err)
			app.CustomGovKeeper.SaveNetworkActor(ctx, proposerActor)

			handler := gov.NewHandler(app.CustomGovKeeper)

			// set change fee permission to addr
			_, err = handler(ctx, &types.MsgWhitelistPermissions{
				Proposer:   sudoAddr,
				Address:    changeFeeAddr,
				Permission: uint32(types.PermChangeTxFee),
			})
			require.NoError(t, err)

			_, err = handler(ctx, tt.msg)
			if tt.desiredErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.desiredErr)
			}
		})
	}
}

func TestNewHandler_SetExecutionFee(t *testing.T) {
	execFeeSetAddr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	sudoAddr, err := sdk.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	tests := []struct {
		name       string
		msg        types.MsgSetExecutionFee
		desiredErr string
	}{
		{
			name: "Success run with ChangeTxFee permission",
			msg: types.MsgSetExecutionFee{
				Name:              "network-properties",
				TransactionType:   "network-properties",
				ExecutionFee:      10000,
				FailureFee:        1000,
				Timeout:           1,
				DefaultParameters: 2,
				Proposer:          execFeeSetAddr,
			},
			desiredErr: "",
		},
		{
			name: "Success run without ChangeTxFee permission",
			msg: types.MsgSetExecutionFee{
				Name:              "network-properties-2",
				TransactionType:   "network-properties-2",
				ExecutionFee:      10000,
				FailureFee:        1000,
				Timeout:           1,
				DefaultParameters: 2,
				Proposer:          sudoAddr,
			},
			desiredErr: "PermChangeTxFee: not enough permissions",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})
			// First we set Role Sudo to proposer Actor
			proposerActor := types.NewDefaultActor(sudoAddr)
			proposerActor.SetRole(types.RoleSudo)
			require.NoError(t, err)
			app.CustomGovKeeper.SaveNetworkActor(ctx, proposerActor)

			handler := gov.NewHandler(app.CustomGovKeeper)

			// set change fee permission to addr
			_, err = handler(ctx, &types.MsgWhitelistPermissions{
				Proposer:   sudoAddr,
				Address:    execFeeSetAddr,
				Permission: uint32(types.PermChangeTxFee),
			})
			require.NoError(t, err)

			_, err = handler(ctx, &tt.msg)
			if tt.desiredErr == "" {
				require.NoError(t, err)
				execFee := app.CustomGovKeeper.GetExecutionFee(ctx, tt.msg.TransactionType)
				require.Equal(t, tt.msg.Name, execFee.Name)
				require.Equal(t, tt.msg.TransactionType, execFee.TransactionType)
				require.Equal(t, tt.msg.ExecutionFee, execFee.ExecutionFee)
				require.Equal(t, tt.msg.FailureFee, execFee.FailureFee)
				require.Equal(t, tt.msg.Timeout, execFee.Timeout)
				require.Equal(t, tt.msg.DefaultParameters, execFee.DefaultParameters)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.desiredErr)
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

			councilor, found := app.CustomGovKeeper.GetCouncilor(ctx, addr)
			require.True(t, found)

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

				perms, found := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
				require.True(t, found)

				err2 = perms.AddToBlacklist(types.PermSetPermissions)
				require.NoError(t, err2)

				app.CustomGovKeeper.SavePermissionsForRole(ctx, types.RoleValidator, &perms)
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

	perms, found := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
	require.True(t, found)
	require.False(t, perms.IsWhitelisted(types.PermSetPermissions))

	msg := types.NewMsgWhitelistRolePermission(
		addr,
		uint32(types.RoleValidator),
		uint32(types.PermSetPermissions),
	)

	handler := gov.NewHandler(app.CustomGovKeeper)
	_, err = handler(ctx, msg)
	require.NoError(t, err)

	perms, found = app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
	require.True(t, found)
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

				_, found := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
				require.True(t, found)

				err2 = app.CustomGovKeeper.WhitelistRolePermission(ctx, types.RoleValidator, types.PermSetPermissions)
				require.NoError(t, err2)
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

	perms, found := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
	require.True(t, found)
	require.False(t, perms.IsBlacklisted(types.PermSetPermissions))

	msg := types.NewMsgBlacklistRolePermission(
		addr,
		uint32(types.RoleValidator),
		uint32(types.PermSetPermissions),
	)

	handler := gov.NewHandler(app.CustomGovKeeper)
	_, err = handler(ctx, msg)
	require.NoError(t, err)

	perms, found = app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
	require.True(t, found)
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

	perms, found := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
	require.True(t, found)
	require.True(t, perms.IsWhitelisted(types.PermClaimValidator))

	msg := types.NewMsgRemoveWhitelistRolePermission(
		addr,
		uint32(types.RoleValidator),
		uint32(types.PermClaimValidator),
	)

	handler := gov.NewHandler(app.CustomGovKeeper)
	_, err = handler(ctx, msg)
	require.NoError(t, err)

	perms, found = app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
	require.True(t, found)
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

	_, found := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
	require.True(t, found)

	// Set some blacklist value
	err = app.CustomGovKeeper.BlacklistRolePermission(ctx, types.RoleValidator, types.PermClaimCouncilor)
	require.NoError(t, err)

	// Check if it is blacklisted.
	perms, found := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
	require.True(t, found)
	require.True(t, perms.IsBlacklisted(types.PermClaimCouncilor))

	msg := types.NewMsgRemoveBlacklistRolePermission(
		addr,
		uint32(types.RoleValidator),
		uint32(types.PermClaimCouncilor),
	)

	handler := gov.NewHandler(app.CustomGovKeeper)
	_, err = handler(ctx, msg)
	require.NoError(t, err)

	perms, found = app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleValidator)
	require.True(t, found)
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
				app.CustomGovKeeper.CreateRole(ctx, types.Role(1234))
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

	_, found := app.CustomGovKeeper.GetPermissionsForRole(ctx, 1234)
	require.False(t, found)

	handler := gov.NewHandler(app.CustomGovKeeper)
	_, err = handler(ctx, types.NewMsgCreateRole(
		addr,
		1234,
	))
	require.NoError(t, err)

	_, found = app.CustomGovKeeper.GetPermissionsForRole(ctx, 1234)
	require.True(t, found)
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

				app.CustomGovKeeper.CreateRole(ctx, types.Role(3))
				err2 = app.CustomGovKeeper.WhitelistRolePermission(ctx, types.Role(3), types.PermClaimValidator)
				require.NoError(t, err2)

				networkActor := types.NewDefaultActor(addr)
				app.CustomGovKeeper.AssignRoleToActor(ctx, networkActor, types.Role(3))
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
	app.CustomGovKeeper.CreateRole(ctx, types.Role(3))
	err = app.CustomGovKeeper.WhitelistRolePermission(ctx, types.Role(3), types.PermSetPermissions)
	require.NoError(t, err)

	msg := types.NewMsgAssignRole(proposerAddr, addr, 3)

	handler := gov.NewHandler(app.CustomGovKeeper)
	_, err = handler(ctx, msg)
	require.NoError(t, err)

	actor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
	require.True(t, found)

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

				app.CustomGovKeeper.CreateRole(ctx, types.Role(3))
				err2 = app.CustomGovKeeper.WhitelistRolePermission(ctx, types.Role(3), types.PermClaimValidator)
				require.NoError(t, err2)
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
	app.CustomGovKeeper.SavePermissionsForRole(ctx, types.Role(3), types.NewPermissions([]types.PermValue{types.PermSetPermissions}, nil))
	actor := types.NewDefaultActor(addr)
	actor.SetRole(3)
	app.CustomGovKeeper.SaveNetworkActor(ctx, actor)

	actor, found := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
	require.True(t, found)
	require.True(t, actor.HasRole(3))

	msg := types.NewMsgRemoveRole(proposerAddr, addr, 3)
	handler := gov.NewHandler(app.CustomGovKeeper)
	_, err = handler(ctx, msg)
	require.NoError(t, err)

	actor, found = app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
	require.True(t, found)

	require.False(t, actor.HasRole(3))
}

func TestHandler_CreateProposalAssignPermission_Errors(t *testing.T) {
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
			"Proposer does not have Perm",
			types.NewMsgProposalAssignPermission(
				proposerAddr, addr, types.PermClaimValidator,
			),
			func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {},
			errors.Wrap(types.ErrNotEnoughPermissions, "PermCreateSetPermissionsProposal"),
		},
		{
			"address already has that permission",
			types.NewMsgProposalAssignPermission(
				proposerAddr, addr, types.PermClaimValidator,
			),
			func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {
				proposerActor := types.NewDefaultActor(proposerAddr)
				err2 := app.CustomGovKeeper.AddWhitelistPermission(ctx, proposerActor, types.PermCreateSetPermissionsProposal)
				require.NoError(t, err2)

				actor := types.NewDefaultActor(addr)
				err2 = app.CustomGovKeeper.AddWhitelistPermission(ctx, actor, types.PermClaimValidator)
				require.NoError(t, err2)
			},
			fmt.Errorf("permission already whitelisted: error adding to whitelist"),
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

func TestHandler_ProposalAssignPermission(t *testing.T) {
	proposerAddr, err := sdk.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{
		Time: time.Now(),
	})

	// Set proposer Permissions
	proposerActor := types.NewDefaultActor(proposerAddr)
	err2 := app.CustomGovKeeper.AddWhitelistPermission(ctx, proposerActor, types.PermCreateSetPermissionsProposal)
	require.NoError(t, err2)

	handler := gov.NewHandler(app.CustomGovKeeper)
	res, err := handler(
		ctx,
		types.NewMsgProposalAssignPermission(proposerAddr, addr, types.PermValue(1)),
	)
	require.NoError(t, err)
	require.Equal(t, types2.GetProposalIDBytes(1), res.Data)

	savedProposal, found := app.CustomGovKeeper.GetProposal(ctx, 1)
	require.True(t, found)
	require.Equal(t, types.NewProposalAssignPermission(1, addr, types.PermValue(1), ctx.BlockTime(), ctx.BlockTime().Add(time.Minute*10)), savedProposal)

	// Next proposal ID is increased.
	id, err := app.CustomGovKeeper.GetNextProposalID(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(2), id)

	// Is not on finished active proposals.
	iterator := app.CustomGovKeeper.GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx, ctx.BlockTime())
	require.False(t, iterator.Valid())

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Minute * 10))
	iterator = app.CustomGovKeeper.GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx, ctx.BlockTime())
	require.True(t, iterator.Valid())
}

func TestHandler_VoteProposal_Errors(t *testing.T) {
	voterAddr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	tests := []struct {
		name         string
		msg          *types.MsgVoteProposal
		preparePerms func(t *testing.T, app *simapp.SimApp, ctx sdk.Context)
		expectedErr  error
	}{
		{
			"Voter does not have permission to vote this proposal",
			types.NewMsgVoteProposal(
				1, voterAddr, types.OptionAbstain,
			),
			func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {},
			types.ErrUserIsNotCouncilor,
		},
		{
			"Proposal does not exist",
			types.NewMsgVoteProposal(
				1, voterAddr, types.OptionAbstain,
			),
			func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {
				councilor := types.NewCouncilor(
					"test",
					"website",
					"social",
					"identity",
					voterAddr,
				)

				app.CustomGovKeeper.SaveCouncilor(ctx, councilor)

				actor := types.NewNetworkActor(
					voterAddr,
					types.Roles{},
					types.Active,
					[]uint32{},
					types.NewPermissions(nil, nil),
					1,
				)
				app.CustomGovKeeper.SaveNetworkActor(ctx, actor)
			},
			types.ErrProposalDoesNotExist,
		},
		{
			"Voter is not active",
			types.NewMsgVoteProposal(
				1, voterAddr, types.OptionAbstain,
			),
			func(t *testing.T, app *simapp.SimApp, ctx sdk.Context) {
				councilor := types.NewCouncilor(
					"test",
					"website",
					"social",
					"identity",
					voterAddr,
				)
				app.CustomGovKeeper.SaveCouncilor(ctx, councilor)

				actor := types.NewDefaultActor(voterAddr)
				app.CustomGovKeeper.SaveNetworkActor(ctx, actor)
			},
			types.ErrActorIsNotActive,
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

func TestHandler_VoteProposal(t *testing.T) {
	voterAddr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	// Put voter as councilor
	councilor := types.NewCouncilor(
		"test",
		"website",
		"social",
		"identity",
		voterAddr,
	)
	app.CustomGovKeeper.SaveCouncilor(ctx, councilor)

	// Create Voter as active actor.
	actor := types.NewNetworkActor(
		voterAddr,
		types.Roles{},
		types.Active,
		[]uint32{},
		types.NewPermissions(nil, nil),
		1,
	)
	app.CustomGovKeeper.SaveNetworkActor(ctx, actor)

	// Create proposal
	proposal := types.NewProposalAssignPermission(1, voterAddr, types.PermClaimCouncilor, ctx.BlockTime(), ctx.BlockTime().Add(time.Second*1))
	err = app.CustomGovKeeper.SaveProposal(ctx, proposal)
	require.NoError(t, err)

	msg := types.NewMsgVoteProposal(proposal.ProposalId, voterAddr, types.OptionAbstain)
	handler := gov.NewHandler(app.CustomGovKeeper)
	_, err = handler(ctx, msg)
	require.NoError(t, err)

	vote, found := app.CustomGovKeeper.GetVote(ctx, proposal.ProposalId, voterAddr)
	require.True(t, found)
	require.Equal(t, types.NewVote(proposal.ProposalId, voterAddr, types.OptionAbstain), vote)
}

func setPermissionToAddr(t *testing.T, app *simapp.SimApp, ctx sdk.Context, addr sdk.AccAddress, perm types.PermValue) error {
	proposerActor := types.NewDefaultActor(addr)
	err := app.CustomGovKeeper.AddWhitelistPermission(ctx, proposerActor, perm)
	require.NoError(t, err)

	return nil
}
