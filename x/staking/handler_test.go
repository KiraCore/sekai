package staking_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/KiraCore/sekai/app"
	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/x/gov"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/staking"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestMain(m *testing.M) {
	app.SetConfig()
	os.Exit(m.Run())
}

func TestNewHandler_MsgClaimValidator_HappyPath(t *testing.T) {
	valAddr1, err := types.ValAddressFromBech32("kiravaloper15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq38f2fp")
	require.NoError(t, err)

	pubKey, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	// First we give user permissions
	networkActor := govtypes.NewNetworkActor(
		types.AccAddress(valAddr1),
		nil,
		1,
		nil,
		govtypes.NewPermissions([]govtypes.PermValue{
			govtypes.PermClaimValidator,
		}, nil),
		1,
	)
	app.CustomGovKeeper.SaveNetworkActor(ctx, networkActor)

	handler := staking.NewHandler(app.CustomStakingKeeper, app.CustomGovKeeper)

	theMsg, err := stakingtypes.NewMsgClaimValidator(
		"aMoniker",
		types.NewDec(1234),
		valAddr1,
		pubKey,
	)
	require.NoError(t, err)

	validatorSet := app.CustomStakingKeeper.GetPendingValidatorSet(ctx)
	require.Len(t, validatorSet, 0)

	_, err = handler(ctx, theMsg)
	require.NoError(t, err)

	validatorSet = app.CustomStakingKeeper.GetPendingValidatorSet(ctx)
	require.Len(t, validatorSet, 1)
}

func TestNewHandler_MsgClaimValidator_Errors(t *testing.T) {
	valAddr1, err := types.ValAddressFromBech32("kiravaloper15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq38f2fp")
	require.NoError(t, err)

	pubKey, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em")
	require.NoError(t, err)

	tests := []struct {
		name          string
		prepareFunc   func(ctx types.Context, app *simapp.SekaiApp)
		expectedError error
	}{
		{
			"user does not have permissions",
			func(ctx types.Context, app *simapp.SekaiApp) {},
			errors.Wrap(govtypes.ErrNotEnoughPermissions, "PermClaimValidator"),
		},
		{
			"validator already exist",
			func(ctx types.Context, app *simapp.SekaiApp) {
				// First we give user permissions
				networkActor := govtypes.NewNetworkActor(
					types.AccAddress(valAddr1),
					nil,
					1,
					nil,
					govtypes.NewPermissions([]govtypes.PermValue{
						govtypes.PermClaimValidator,
					}, nil),
					1,
				)
				app.CustomGovKeeper.SaveNetworkActor(ctx, networkActor)

				validator, err := stakingtypes.NewValidator(
					"aMoniker",
					types.NewDec(1234),
					valAddr1,
					pubKey,
				)
				require.NoError(t, err)
				app.CustomStakingKeeper.AddValidator(ctx, validator)
			},
			stakingtypes.ErrValidatorAlreadyClaimed,
		},
		{
			"validator moniker exists",
			func(ctx types.Context, app *simapp.SekaiApp) {
				pubKey, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, "kiravalconspub1zcjduepqylc5k8r40azmw0xa7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsqgrkp48")
				require.NoError(t, err)

				valAddr2, err := types.ValAddressFromBech32("kiravaloper15ky9du8a2wlstz6fpx3p4mqpryrm5cgqeyf3v0")
				require.NoError(t, err)

				networkActor := govtypes.NewNetworkActor(
					types.AccAddress(valAddr1),
					nil,
					1,
					nil,
					govtypes.NewPermissions([]govtypes.PermValue{
						govtypes.PermClaimValidator,
					}, nil),
					1,
				)
				app.CustomGovKeeper.SaveNetworkActor(ctx, networkActor)

				validator, err := stakingtypes.NewValidator(
					"aMoniker", // Other validator with repeated moniker.
					types.NewDec(1234),
					valAddr2,
					pubKey,
				)
				require.NoError(t, err)
				app.CustomStakingKeeper.AddValidator(ctx, validator)
			},
			stakingtypes.ErrValidatorMonikerExists,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			tt.prepareFunc(ctx, app)

			handler := staking.NewHandler(app.CustomStakingKeeper, app.CustomGovKeeper)

			theMsg, err := stakingtypes.NewMsgClaimValidator(
				"aMoniker",
				types.NewDec(1234),
				valAddr1,
				pubKey,
			)
			require.NoError(t, err)

			_, err = handler(ctx, theMsg)
			require.EqualError(t, err, tt.expectedError.Error())
		})
	}
}

func TestNewHandler_SetPermissions_ActorWithRole(t *testing.T) {
	valAddr1, err := types.ValAddressFromBech32("kiravaloper15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq38f2fp")
	require.NoError(t, err)

	pubKey, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	// Save network actor With Role Validator
	networkActor := govtypes.NewDefaultActor(types.AccAddress(valAddr1))
	networkActor.SetRole(govtypes.RoleValidator)
	app.CustomGovKeeper.SaveNetworkActor(ctx, networkActor)

	handler := staking.NewHandler(app.CustomStakingKeeper, app.CustomGovKeeper)

	theMsg, err := stakingtypes.NewMsgClaimValidator(
		"aMoniker",
		types.NewDec(1234),
		valAddr1,
		pubKey,
	)
	require.NoError(t, err)

	_, err = handler(ctx, theMsg)
	require.NoError(t, err)

	validatorSet := app.CustomStakingKeeper.GetPendingValidatorSet(ctx)
	require.Len(t, validatorSet, 1)
}

func TestHandler_ProposalUnjailValidator_Errors(t *testing.T) {
	proposerAddr, err := types.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)
	valAddr := types.ValAddress(proposerAddr)

	pubKey, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em")
	require.NoError(t, err)

	tests := []struct {
		name        string
		expectedErr error
		prepareFunc func(ctx types.Context, app *simapp.SekaiApp)
	}{
		{
			name:        "not enough permissions to create validator",
			expectedErr: errors.Wrap(govtypes.ErrNotEnoughPermissions, govtypes.PermCreateUnjailValidatorProposal.String()),
			prepareFunc: func(ctx types.Context, app *simapp.SekaiApp) {},
		},
		{
			name:        "validator does not exist",
			expectedErr: fmt.Errorf("validator not found"),
			prepareFunc: func(ctx types.Context, app *simapp.SekaiApp) {
				proposerActor := govtypes.NewDefaultActor(proposerAddr)
				err2 := app.CustomGovKeeper.AddWhitelistPermission(ctx, proposerActor, govtypes.PermCreateUnjailValidatorProposal)
				require.NoError(t, err2)
			},
		},
		{
			name:        "validator is not jailed",
			expectedErr: fmt.Errorf("validator is not jailed"),
			prepareFunc: func(ctx types.Context, app *simapp.SekaiApp) {
				proposerActor := govtypes.NewDefaultActor(proposerAddr)
				err2 := app.CustomGovKeeper.AddWhitelistPermission(ctx, proposerActor, govtypes.PermCreateUnjailValidatorProposal)
				require.NoError(t, err2)

				val, err := stakingtypes.NewValidator("Moniker", types.NewDec(123), valAddr, pubKey)
				require.NoError(t, err)

				app.CustomStakingKeeper.AddValidator(ctx, val)
			},
		},
		{
			name:        "it passed the time when validator cannot be unjailed",
			expectedErr: fmt.Errorf("time to unjail passed"),
			prepareFunc: func(ctx types.Context, app *simapp.SekaiApp) {
				networkProperties := app.CustomGovKeeper.GetNetworkProperties(ctx)
				networkProperties.JailMaxTime = 300 // 300 seconds = 5 min
				app.CustomGovKeeper.SetNetworkProperties(ctx, networkProperties)

				proposerActor := govtypes.NewDefaultActor(proposerAddr)
				err2 := app.CustomGovKeeper.AddWhitelistPermission(ctx, proposerActor, govtypes.PermCreateUnjailValidatorProposal)
				require.NoError(t, err2)

				val, err := stakingtypes.NewValidator("Moniker", types.NewDec(123), valAddr, pubKey)
				require.NoError(t, err)

				app.CustomStakingKeeper.AddValidator(ctx, val)

				// Jail Validator
				err = app.CustomStakingKeeper.Jail(ctx, val.ValKey)
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{
				Time: time.Now(),
			})

			tt.prepareFunc(ctx, app)

			// After 10 minutes
			ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Minute * 10))

			handler := gov.NewHandler(app.CustomGovKeeper)
			proposal := stakingtypes.NewUnjailValidatorProposal(
				proposerAddr,
				"thehash",
				"theReference",
			)
			msg, err := govtypes.NewMsgSubmitProposal(proposerAddr, "title", "some desc", proposal)
			require.NoError(t, err)
			_, err = handler(
				ctx,
				msg,
			)
			require.EqualError(t, err, tt.expectedErr.Error())
		})
	}
}

func TestHandler_ProposalUnjailValidator(t *testing.T) {
	proposerAddr, err := types.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	valAddr := types.ValAddress(proposerAddr)

	pubKey, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{
		Time: time.Now(),
	})

	// Set proposer Permissions
	proposerActor := govtypes.NewDefaultActor(proposerAddr)
	err2 := app.CustomGovKeeper.AddWhitelistPermission(ctx, proposerActor, govtypes.PermCreateUnjailValidatorProposal)
	require.NoError(t, err2)

	properties := app.CustomGovKeeper.GetNetworkProperties(ctx)
	properties.ProposalEndTime = 10
	app.CustomGovKeeper.SetNetworkProperties(ctx, properties)

	val, err := stakingtypes.NewValidator("Moniker", types.NewDec(123), valAddr, pubKey)
	require.NoError(t, err)
	app.CustomStakingKeeper.AddValidator(ctx, val)

	err = app.CustomStakingKeeper.Jail(ctx, val.ValKey)
	require.NoError(t, err)

	handler := gov.NewHandler(app.CustomGovKeeper)
	proposal := stakingtypes.NewUnjailValidatorProposal(
		proposerAddr,
		"thehash",
		"theReference",
	)
	msg, err := govtypes.NewMsgSubmitProposal(proposerAddr, "title", "some desc", proposal)
	require.NoError(t, err)
	_, err = handler(
		ctx,
		msg,
	)
	require.NoError(t, err)

	savedProposal, found := app.CustomGovKeeper.GetProposal(ctx, 1)
	require.True(t, found)

	expectedSavedProposal, err := govtypes.NewProposal(
		1,
		"title",
		"some desc",
		stakingtypes.NewUnjailValidatorProposal(
			proposerAddr,
			"thehash",
			"theReference",
		),
		ctx.BlockTime(),
		ctx.BlockTime().Add(time.Second*time.Duration(properties.ProposalEndTime)),
		ctx.BlockTime().Add(time.Second*time.Duration(properties.ProposalEndTime)+
			time.Second*time.Duration(properties.ProposalEnactmentTime),
		),
		ctx.BlockHeight()+2,
		ctx.BlockHeight()+3,
	)
	require.NoError(t, err)
	require.Equal(t, expectedSavedProposal, savedProposal)

	// Next proposal ID is increased.
	id := app.CustomGovKeeper.GetNextProposalID(ctx)
	require.Equal(t, uint64(2), id)

	// Is not on finished active proposals.
	iterator := app.CustomGovKeeper.GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx, ctx.BlockTime())
	require.False(t, iterator.Valid())

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Minute * 10))
	iterator = app.CustomGovKeeper.GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx, ctx.BlockTime())
	require.True(t, iterator.Valid())
}

// TODO: should add more tests for various types of cases by network properties status
