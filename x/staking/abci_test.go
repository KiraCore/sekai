package staking_test

import (
	"testing"

	"github.com/KiraCore/sekai/x/staking"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	simapp "github.com/KiraCore/sekai/app"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
)

func TestItUpdatesTheValidatorSetBasedOnPendingValidators(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	addr1 := addrs[0]
	valAddr1 := sdk.ValAddress(addr1)

	pubkeys := simapp.CreateTestPubKeys(1)
	pubKey := pubkeys[0]

	validator1, err := stakingtypes.NewValidator(
		valAddr1,
		pubKey,
	)
	require.NoError(t, err)
	app.CustomStakingKeeper.AddPendingValidator(ctx, validator1)

	// First checkings
	validatorSet := app.CustomStakingKeeper.GetValidatorSet(ctx)
	require.Len(t, validatorSet, 0)
	validatorSet = app.CustomStakingKeeper.GetPendingValidatorSet(ctx)
	require.Len(t, validatorSet, 1)

	updates := staking.EndBlocker(ctx, app.CustomStakingKeeper)
	require.Len(t, updates, 1)

	validatorSet = app.CustomStakingKeeper.GetValidatorSet(ctx)
	require.Len(t, validatorSet, 1)
	validatorSet = app.CustomStakingKeeper.GetPendingValidatorSet(ctx)
	require.Len(t, validatorSet, 0)
}

func TestItDoesNotReturnUpdatesIfThereIsNoPending(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	// First checkings
	validatorSet := app.CustomStakingKeeper.GetValidatorSet(ctx)
	require.Len(t, validatorSet, 0)
	validatorSet = app.CustomStakingKeeper.GetPendingValidatorSet(ctx)
	require.Len(t, validatorSet, 0)

	updates := staking.EndBlocker(ctx, app.CustomStakingKeeper)
	require.Len(t, updates, 0)

	validatorSet = app.CustomStakingKeeper.GetValidatorSet(ctx)
	require.Len(t, validatorSet, 0)
	validatorSet = app.CustomStakingKeeper.GetPendingValidatorSet(ctx)
	require.Len(t, validatorSet, 0)
}

func TestItRemovesFromTheValidatorSetWhenInRemovingQueue(t *testing.T) {
	tests := []struct {
		name        string
		prepareFunc func(app *simapp.SekaiApp, ctx sdk.Context, validator stakingtypes.Validator)
	}{
		{
			name: "remove because it is paused",
			prepareFunc: func(app *simapp.SekaiApp, ctx sdk.Context, validator stakingtypes.Validator) {
				err := app.CustomStakingKeeper.Pause(ctx, validator.ValKey)
				require.NoError(t, err)
			},
		},
		{
			name: "remove because it is inactive",
			prepareFunc: func(app *simapp.SekaiApp, ctx sdk.Context, validator stakingtypes.Validator) {
				err := app.CustomStakingKeeper.Inactivate(ctx, validator.ValKey)
				require.NoError(t, err)
			},
		},
		{
			name: "remove because it is jailed",
			prepareFunc: func(app *simapp.SekaiApp, ctx sdk.Context, validator stakingtypes.Validator) {
				err := app.CustomStakingKeeper.Jail(ctx, validator.ValKey)
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
			addr1 := addrs[0]
			valAddr1 := sdk.ValAddress(addr1)

			pubkeys := simapp.CreateTestPubKeys(1)
			pubKey := pubkeys[0]

			validator1, err := stakingtypes.NewValidator(
				valAddr1,
				pubKey,
			)
			require.NoError(t, err)
			app.CustomStakingKeeper.AddValidator(ctx, validator1)

			tt.prepareFunc(app, ctx, validator1)

			updates := staking.EndBlocker(ctx, app.CustomStakingKeeper)
			require.Len(t, updates, 1)

			set := app.CustomStakingKeeper.GetRemovingValidatorSet(ctx)
			require.Len(t, set, 0)
		})
	}
}

func TestItIncludesItBackToValidatorSetOnceReactivatingIt(t *testing.T) {
	tests := []struct {
		name                string
		prepareDeactivation func(app *simapp.SekaiApp, ctx sdk.Context, validator stakingtypes.Validator)
		prepareFunc         func(app *simapp.SekaiApp, ctx sdk.Context, validator stakingtypes.Validator)
	}{
		{
			name: "reactivating from paused",
			prepareDeactivation: func(app *simapp.SekaiApp, ctx sdk.Context, validator stakingtypes.Validator) {
				err := app.CustomStakingKeeper.Pause(ctx, validator.ValKey)
				require.NoError(t, err)

				// We end the block so the validator is paused
				staking.EndBlocker(ctx, app.CustomStakingKeeper)
			},
			prepareFunc: func(app *simapp.SekaiApp, ctx sdk.Context, validator stakingtypes.Validator) {
				err := app.CustomStakingKeeper.Unpause(ctx, validator.ValKey)
				require.NoError(t, err)
			},
		},
		{
			name: "reactivating from inactive",
			prepareDeactivation: func(app *simapp.SekaiApp, ctx sdk.Context, validator stakingtypes.Validator) {
				err := app.CustomStakingKeeper.Inactivate(ctx, validator.ValKey)
				require.NoError(t, err)

				// We end the block so the validator is paused
				staking.EndBlocker(ctx, app.CustomStakingKeeper)
			},
			prepareFunc: func(app *simapp.SekaiApp, ctx sdk.Context, validator stakingtypes.Validator) {
				err := app.CustomStakingKeeper.Activate(ctx, validator.ValKey)
				require.NoError(t, err)
			},
		},
		{
			name: "reactivating from jailed",
			prepareDeactivation: func(app *simapp.SekaiApp, ctx sdk.Context, validator stakingtypes.Validator) {
				err := app.CustomStakingKeeper.Jail(ctx, validator.ValKey)
				require.NoError(t, err)

				// We end the block so the validator is paused
				staking.EndBlocker(ctx, app.CustomStakingKeeper)
			},
			prepareFunc: func(app *simapp.SekaiApp, ctx sdk.Context, validator stakingtypes.Validator) {
				err := app.CustomStakingKeeper.Unjail(ctx, validator.ValKey)
				require.NoError(t, err)
				err = app.CustomStakingKeeper.Activate(ctx, validator.ValKey)
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
			addr1 := addrs[0]
			valAddr1 := sdk.ValAddress(addr1)

			pubkeys := simapp.CreateTestPubKeys(1)
			pubKey := pubkeys[0]

			validator1, err := stakingtypes.NewValidator(
				valAddr1,
				pubKey,
			)
			require.NoError(t, err)
			app.CustomStakingKeeper.AddValidator(ctx, validator1)

			tt.prepareDeactivation(app, ctx, validator1)

			set := app.CustomStakingKeeper.GetReactivatingValidatorSet(ctx)
			require.Len(t, set, 0)

			tt.prepareFunc(app, ctx, validator1)
			set = app.CustomStakingKeeper.GetReactivatingValidatorSet(ctx)
			require.Len(t, set, 1)

			updatedSet := staking.EndBlocker(ctx, app.CustomStakingKeeper)
			require.Len(t, updatedSet, 1)
		})
	}
}
