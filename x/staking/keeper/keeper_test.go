package keeper_test

import (
	"os"
	"testing"

	simapp "github.com/KiraCore/sekai/app"
	appparams "github.com/KiraCore/sekai/app/params"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/staking/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	appparams.SetConfig()
	os.Exit(m.Run())
}

func TestKeeper_AddValidator(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	validators := createValidators(t, app, ctx, 2)
	validator := validators[0]

	app.CustomStakingKeeper.AddValidator(ctx, validator)

	// Get By Validator Address.
	getValidator, err := app.CustomStakingKeeper.GetValidator(ctx, validator.ValKey)
	require.NoError(t, err)
	require.True(t, validator.Equal(getValidator))

	// Non existing validator Addr.
	_, err = app.CustomStakingKeeper.GetValidator(ctx, sdk.ValAddress("non existing"))
	require.EqualError(t, err, "validator not found")

	// Get by AccAddress.
	getValidator, err = app.CustomStakingKeeper.GetValidatorByAccAddress(ctx, sdk.AccAddress(validator.ValKey))
	require.NoError(t, err)
	require.True(t, validator.Equal(getValidator))

	// Non existing AccAddress.
	_, err = app.CustomStakingKeeper.GetValidatorByAccAddress(ctx, sdk.AccAddress("non existing"))
	require.EqualError(t, err, "validator not found")

	// Get by ConsAddress
	getByConsAddrValidator, err := app.CustomStakingKeeper.GetValidatorByConsAddr(ctx, validator.GetConsAddr())
	require.NoError(t, err)
	require.True(t, validator.Equal(getByConsAddrValidator))

	// Non existing moniker
	_, err = app.CustomStakingKeeper.GetValidatorByMoniker(ctx, "UnexistingMoniker")
	require.EqualError(t, err, "validator with moniker UnexistingMoniker not found")
}

func TestKeeper_GetMonikerByAddress(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	validators := createValidators(t, app, ctx, 2)
	validator := validators[0]

	app.CustomStakingKeeper.AddValidator(ctx, validator)
	err := app.CustomGovKeeper.RegisterIdentityRecords(ctx, sdk.AccAddress(validator.ValKey), []govtypes.IdentityInfoEntry{
		{
			Key:  "moniker",
			Info: "node0",
		},
	})
	require.NoError(t, err)

	// get moniker by address
	moniker := app.CustomStakingKeeper.GetMonikerByAddress(ctx, sdk.AccAddress(validator.ValKey))
	require.Equal(t, moniker, "node0")

	moniker = app.CustomStakingKeeper.GetMonikerByAddress(ctx, sdk.AccAddress("non existing"))
	require.Equal(t, moniker, "")
}

func TestKeeper_GetValidatorSet(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	validators := createValidators(t, app, ctx, 2)
	validator1 := validators[0]
	validator2 := validators[1]

	app.CustomStakingKeeper.AddValidator(ctx, validator1)
	app.CustomStakingKeeper.AddValidator(ctx, validator2)

	validatorSet := app.CustomStakingKeeper.GetValidatorSet(ctx)
	require.Equal(t, 3, len(validatorSet))

	// Iterate validators
	passed := 0
	app.CustomStakingKeeper.IterateValidators(ctx, func(index int64, validator *types.Validator) (stop bool) {
		passed++
		return false
	})
	require.Equal(t, 3, passed)

	// Iterate validators with stop
	passed = 0
	app.CustomStakingKeeper.IterateValidators(ctx, func(index int64, validator *types.Validator) (stop bool) {
		passed++
		return true
	})
	require.Equal(t, 1, passed)
}

func TestKeeper_GetPendingValidators(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	validators := createValidators(t, app, ctx, 2)
	validator1 := validators[0]
	validator2 := validators[1]

	app.CustomStakingKeeper.AddPendingValidator(ctx, validator1)
	app.CustomStakingKeeper.AddPendingValidator(ctx, validator2)

	validatorSet := app.CustomStakingKeeper.GetPendingValidatorSet(ctx)
	require.Equal(t, 2, len(validatorSet))

	app.CustomStakingKeeper.RemovePendingValidator(ctx, validator1)
	validatorSet = app.CustomStakingKeeper.GetPendingValidatorSet(ctx)
	require.Equal(t, 1, len(validatorSet))
}

func createValidators(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context, accNum int) (validators []types.Validator) {
	addrs := simapp.AddTestAddrsIncremental(app, ctx, accNum, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))

	for _, addr := range addrs {
		valAddr := sdk.ValAddress(addr)
		pubkeys := simapp.CreateTestPubKeys(1)
		pubKey := pubkeys[0]

		validator, err := types.NewValidator(
			valAddr,
			pubKey,
		)
		require.NoError(t, err)
		validators = append(validators, validator)
	}

	return
}
