package keeper_test

import (
	"fmt"
	"testing"

	"github.com/KiraCore/sekai/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/simapp"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestPauseValidator_Errors(t *testing.T) {
	tests := []struct {
		name            string
		expectedError   error
		prepareScenario func(app *simapp.SimApp, ctx sdk.Context, validator types.Validator)
	}{
		{
			name:          "validator does not exist",
			expectedError: fmt.Errorf("validator not found"),
			prepareScenario: func(app *simapp.SimApp, ctx sdk.Context, validator types.Validator) {
			},
		},
		{
			name:          "validator is inactivated",
			expectedError: types.ErrValidatorInactive,
			prepareScenario: func(app *simapp.SimApp, ctx sdk.Context, validator types.Validator) {
				validator.Status = types.Inactive
				app.CustomStakingKeeper.AddValidator(ctx, validator)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			validators := createValidators(t, app, ctx, 1)
			validator1 := validators[0]

			tt.prepareScenario(app, ctx, validators[0])

			err := app.CustomStakingKeeper.Pause(ctx, validator1.GetValKey())
			require.EqualError(t, err, tt.expectedError.Error())
		})
	}
}

func TestPauseValidator(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	validators := createValidators(t, app, ctx, 1)
	validator1 := validators[0]

	app.CustomStakingKeeper.AddValidator(ctx, validator1)

	savedValidator, err := app.CustomStakingKeeper.GetValidator(ctx, validator1.ValKey)
	require.NoError(t, err)
	require.False(t, savedValidator.IsPaused())

	err = app.CustomStakingKeeper.Pause(ctx, savedValidator.ValKey)
	require.NoError(t, err)
	pausedValidator, err := app.CustomStakingKeeper.GetValidator(ctx, validator1.ValKey)
	require.NoError(t, err)
	require.True(t, pausedValidator.IsPaused())

	valKeys := app.CustomStakingKeeper.GetRemovingValidatorSet(ctx)
	require.Len(t, valKeys, 1)
}

func TestUnpauseValidator_Errors(t *testing.T) {
	tests := []struct {
		name            string
		expectedError   error
		prepareScenario func(app *simapp.SimApp, ctx sdk.Context, validator types.Validator)
	}{
		{
			name:          "validator does not exist",
			expectedError: fmt.Errorf("validator not found"),
			prepareScenario: func(app *simapp.SimApp, ctx sdk.Context, validator types.Validator) {
			},
		},
		{
			name:          "validator is inactivated",
			expectedError: types.ErrValidatorInactive,
			prepareScenario: func(app *simapp.SimApp, ctx sdk.Context, validator types.Validator) {
				validator.Status = types.Inactive
				app.CustomStakingKeeper.AddValidator(ctx, validator)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			validators := createValidators(t, app, ctx, 1)
			validator1 := validators[0]

			tt.prepareScenario(app, ctx, validators[0])

			err := app.CustomStakingKeeper.Unpause(ctx, validator1.GetValKey())
			require.EqualError(t, err, tt.expectedError.Error())
		})
	}
}

func TestUnpauseValidator(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	validators := createValidators(t, app, ctx, 2)
	validator1 := validators[0]

	app.CustomStakingKeeper.AddValidator(ctx, validator1)
	err := app.CustomStakingKeeper.Pause(ctx, validator1.ValKey)
	require.NoError(t, err)

	err = app.CustomStakingKeeper.Unpause(ctx, validator1.ValKey)
	require.NoError(t, err)
	unpausedValidator, err := app.CustomStakingKeeper.GetValidator(ctx, validator1.ValKey)
	require.NoError(t, err)
	require.False(t, unpausedValidator.IsPaused())
}

func TestValidatorInactivate_Errors(t *testing.T) {
	tests := []struct {
		name            string
		expectedError   error
		prepareScenario func(app *simapp.SimApp, ctx sdk.Context, validator types.Validator)
	}{
		{
			name:          "validator does not exist",
			expectedError: fmt.Errorf("validator not found"),
			prepareScenario: func(app *simapp.SimApp, ctx sdk.Context, validator types.Validator) {
			},
		},
		{
			name:          "validator is paused",
			expectedError: types.ErrValidatorPaused,
			prepareScenario: func(app *simapp.SimApp, ctx sdk.Context, validator types.Validator) {
				validator.Status = types.Paused
				app.CustomStakingKeeper.AddValidator(ctx, validator)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			validators := createValidators(t, app, ctx, 1)
			validator1 := validators[0]

			tt.prepareScenario(app, ctx, validators[0])

			err := app.CustomStakingKeeper.Inactivate(ctx, validator1.GetValKey())
			require.EqualError(t, err, tt.expectedError.Error())
		})
	}
}

func TestInactiveValidator(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	validators := createValidators(t, app, ctx, 1)
	validator1 := validators[0]

	app.CustomStakingKeeper.AddValidator(ctx, validator1)
	err := app.CustomStakingKeeper.Inactivate(ctx, validator1.ValKey)
	require.NoError(t, err)

	inactiveValidator, err := app.CustomStakingKeeper.GetValidator(ctx, validator1.ValKey)
	require.NoError(t, err)
	require.True(t, inactiveValidator.IsInactivated())
}

func TestValidatorActivate_Errors(t *testing.T) {
	tests := []struct {
		name            string
		expectedError   error
		prepareScenario func(app *simapp.SimApp, ctx sdk.Context, validator types.Validator)
	}{
		{
			name:          "validator does not exist",
			expectedError: fmt.Errorf("validator not found"),
			prepareScenario: func(app *simapp.SimApp, ctx sdk.Context, validator types.Validator) {
			},
		},
		{
			name:          "validator is paused",
			expectedError: types.ErrValidatorPaused,
			prepareScenario: func(app *simapp.SimApp, ctx sdk.Context, validator types.Validator) {
				validator.Status = types.Paused
				app.CustomStakingKeeper.AddValidator(ctx, validator)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			validators := createValidators(t, app, ctx, 1)
			validator1 := validators[0]

			tt.prepareScenario(app, ctx, validators[0])

			err := app.CustomStakingKeeper.Activate(ctx, validator1.GetValKey())
			require.EqualError(t, err, tt.expectedError.Error())
		})
	}
}

func TestActivateValidator(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	validators := createValidators(t, app, ctx, 1)
	validator1 := validators[0]

	app.CustomStakingKeeper.AddValidator(ctx, validator1)
	err := app.CustomStakingKeeper.Inactivate(ctx, validator1.ValKey)
	require.NoError(t, err)

	inactiveValidator, err := app.CustomStakingKeeper.GetValidator(ctx, validator1.ValKey)
	require.NoError(t, err)
	require.True(t, inactiveValidator.IsInactivated())

	err = app.CustomStakingKeeper.Activate(ctx, validator1.ValKey)
	require.NoError(t, err)

	inactiveValidator, err = app.CustomStakingKeeper.GetValidator(ctx, validator1.ValKey)
	require.NoError(t, err)
	require.True(t, inactiveValidator.IsActive())
}
