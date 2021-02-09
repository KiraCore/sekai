package keeper_test

import (
	"fmt"
	"testing"
	"time"

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

	valKeys := app.CustomStakingKeeper.GetRemovingValidatorSet(ctx)
	require.Len(t, valKeys, 1)
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

func TestRemoveFromRemovingValidatorQueue(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	validators := createValidators(t, app, ctx, 1)
	validator1 := validators[0]

	app.CustomStakingKeeper.AddValidator(ctx, validator1)
	err := app.CustomStakingKeeper.Inactivate(ctx, validator1.ValKey)
	require.NoError(t, err)

	valKeys := app.CustomStakingKeeper.GetRemovingValidatorSet(ctx)
	require.Len(t, valKeys, 1)

	app.CustomStakingKeeper.RemoveRemovingValidator(ctx, validator1)

	valKeys = app.CustomStakingKeeper.GetRemovingValidatorSet(ctx)
	require.Len(t, valKeys, 0)
}

func TestReactivatingValidator(t *testing.T) {
	tests := []struct {
		name                string
		prepareDeactivation func(app *simapp.SimApp, ctx sdk.Context, validator types.Validator)
		prepare             func(app *simapp.SimApp, ctx sdk.Context, validator types.Validator)
	}{
		{
			name: "reactivating from pause",
			prepareDeactivation: func(app *simapp.SimApp, ctx sdk.Context, validator types.Validator) {
				err := app.CustomStakingKeeper.Pause(ctx, validator.ValKey)
				require.NoError(t, err)
			},
			prepare: func(app *simapp.SimApp, ctx sdk.Context, validator types.Validator) {
				err := app.CustomStakingKeeper.Unpause(ctx, validator.ValKey)
				require.NoError(t, err)
			},
		},
		{
			name: "reactivating from inactivate",
			prepareDeactivation: func(app *simapp.SimApp, ctx sdk.Context, validator types.Validator) {
				err := app.CustomStakingKeeper.Inactivate(ctx, validator.ValKey)
				require.NoError(t, err)
			},
			prepare: func(app *simapp.SimApp, ctx sdk.Context, validator types.Validator) {
				err := app.CustomStakingKeeper.Activate(ctx, validator.ValKey)
				require.NoError(t, err)
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
			app.CustomStakingKeeper.AddValidator(ctx, validator1)

			tt.prepareDeactivation(app, ctx, validator1)

			validator, err := app.CustomStakingKeeper.GetValidator(ctx, validator1.ValKey)
			require.NoError(t, err)
			require.False(t, validator.IsActive())

			tt.prepare(app, ctx, validator1)

			validator, err = app.CustomStakingKeeper.GetValidator(ctx, validator1.ValKey)
			require.NoError(t, err)
			require.True(t, validator.IsActive())

			// And it is included in the set of ReactivatingValidators
			reactivatingVals := app.CustomStakingKeeper.GetReactivatingValidatorSet(ctx)
			require.Len(t, reactivatingVals, 1)
		})
	}
}

func TestValidatorJail_Errors(t *testing.T) {
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			validators := createValidators(t, app, ctx, 1)
			validator1 := validators[0]

			tt.prepareScenario(app, ctx, validators[0])

			err := app.CustomStakingKeeper.Jail(ctx, validator1.GetValKey())
			require.EqualError(t, err, tt.expectedError.Error())
		})
	}
}

func TestJailValidator(t *testing.T) {
	blockTime := time.Now()

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{
		Time: blockTime,
	})

	validators := createValidators(t, app, ctx, 1)
	validator1 := validators[0]

	app.CustomStakingKeeper.AddValidator(ctx, validator1)
	err := app.CustomStakingKeeper.Jail(ctx, validator1.ValKey)
	require.NoError(t, err)

	inactiveValidator, err := app.CustomStakingKeeper.GetValidator(ctx, validator1.ValKey)
	require.NoError(t, err)
	require.True(t, inactiveValidator.IsJailed())

	valKeys := app.CustomStakingKeeper.GetRemovingValidatorSet(ctx)
	require.Len(t, valKeys, 1)

	// It saved the jailing info.
	valInfo, found := app.CustomStakingKeeper.GetValidatorJailInfo(ctx, validator1.ValKey)
	require.True(t, found)
	require.True(t, valInfo.Time.Equal(blockTime))
}

func TestValidatorUnjail_Errors(t *testing.T) {
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			validators := createValidators(t, app, ctx, 1)
			validator1 := validators[0]

			tt.prepareScenario(app, ctx, validators[0])

			err := app.CustomStakingKeeper.Unjail(ctx, validator1.GetValKey())
			require.EqualError(t, err, tt.expectedError.Error())
		})
	}
}

func TestUnjailValidator(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	validators := createValidators(t, app, ctx, 1)
	validator1 := validators[0]

	app.CustomStakingKeeper.AddValidator(ctx, validator1)
	err := app.CustomStakingKeeper.Jail(ctx, validator1.ValKey)
	require.NoError(t, err)

	inactiveValidator, err := app.CustomStakingKeeper.GetValidator(ctx, validator1.ValKey)
	require.NoError(t, err)
	require.True(t, inactiveValidator.IsJailed())
	_, found := app.CustomStakingKeeper.GetValidatorJailInfo(ctx, validator1.ValKey)
	require.True(t, found)

	err = app.CustomStakingKeeper.Unjail(ctx, validator1.ValKey)
	require.NoError(t, err)

	inactiveValidator, err = app.CustomStakingKeeper.GetValidator(ctx, validator1.ValKey)
	require.NoError(t, err)
	require.True(t, inactiveValidator.IsActive())
	_, found = app.CustomStakingKeeper.GetValidatorJailInfo(ctx, validator1.ValKey)
	require.False(t, found)
}
