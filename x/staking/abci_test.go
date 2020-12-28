package staking_test

import (
	"testing"

	"github.com/KiraCore/sekai/x/staking"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/KiraCore/sekai/simapp"
	types2 "github.com/KiraCore/sekai/x/staking/types"
)

func TestItUpdatesTheValidatorSetBasedOnPendingValidators(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, types.TokensFromConsensusPower(10))
	addr1 := addrs[0]
	valAddr1 := types.ValAddress(addr1)

	pubKey, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em")
	require.NoError(t, err)

	validator1, err := types2.NewValidator(
		"validator 1",
		"some-web.com",
		"A Social",
		"My Identity",
		types.NewDec(1234),
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
		prepareFunc func(app *simapp.SimApp, ctx types.Context, validator types2.Validator)
	}{
		{
			name: "remove because it is paused",
			prepareFunc: func(app *simapp.SimApp, ctx types.Context, validator types2.Validator) {
				err := app.CustomStakingKeeper.Pause(ctx, validator.ValKey)
				require.NoError(t, err)
			},
		},
		{
			name: "remove because it is inactive",
			prepareFunc: func(app *simapp.SimApp, ctx types.Context, validator types2.Validator) {
				err := app.CustomStakingKeeper.Inactivate(ctx, validator.ValKey)
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, types.TokensFromConsensusPower(10))
			addr1 := addrs[0]
			valAddr1 := types.ValAddress(addr1)

			pubKey, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em")
			require.NoError(t, err)

			validator1, err := types2.NewValidator(
				"validator 1",
				"some-web.com",
				"A Social",
				"My Identity",
				types.NewDec(1234),
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
