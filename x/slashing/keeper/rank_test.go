package keeper_test

import (
	"testing"
	"time"

	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/x/slashing/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestResetWholeValidatorRank(t *testing.T) {
	tests := []struct {
		name            string
		expectedError   error
		prepareScenario func(app *simapp.SekaiApp, ctx sdk.Context, validator stakingtypes.Validator)
	}{
		{
			name:          "check validator status reset",
			expectedError: nil,
			prepareScenario: func(app *simapp.SekaiApp, ctx sdk.Context, validator stakingtypes.Validator) {
				app.CustomStakingKeeper.Inactivate(ctx, validator.ValKey)
			},
		},
		{
			name:          "check validator rank, streak reset",
			expectedError: nil,
			prepareScenario: func(app *simapp.SekaiApp, ctx sdk.Context, validator stakingtypes.Validator) {
				validator.Rank = 10
				validator.Streak = 10
				app.CustomStakingKeeper.AddValidator(ctx, validator)
			},
		},
		{
			name:          "check validator mischance, produced blocks counter, missed blocks counter, last present block reset",
			expectedError: nil,
			prepareScenario: func(app *simapp.SekaiApp, ctx sdk.Context, validator stakingtypes.Validator) {
				info, found := app.CustomSlashingKeeper.GetValidatorSigningInfo(ctx, validator.GetConsAddr())
				if !found {
					panic("validator signing info not found")
				}
				info.StartHeight = 100
				info.InactiveUntil = time.Unix(0, 0)
				info.MischanceConfidence = 0
				info.Mischance = 0
				info.MissedBlocksCounter = 0
				info.ProducedBlocksCounter = 0
				info.LastPresentBlock = 100
				app.CustomSlashingKeeper.SetValidatorSigningInfo(ctx, validator.GetConsAddr(), info)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			validators := createValidators(t, app, ctx, 1)
			app.CustomStakingKeeper.AddValidator(ctx, validators[0])
			app.CustomStakingKeeper.AfterValidatorJoined(ctx, validators[0].GetConsAddr(), validators[0].ValKey)

			tt.prepareScenario(app, ctx, validators[0])

			err := app.CustomSlashingKeeper.ResetWholeValidatorRank(ctx)
			require.NoError(t, err)

			validators = []stakingtypes.Validator{}
			app.CustomStakingKeeper.IterateValidators(ctx, func(index int64, validator *stakingtypes.Validator) (stop bool) {
				validators = append(validators, *validator)
				return false
			})
			require.Equal(t, 1, len(validators))
			require.Equal(t, stakingtypes.Active, validators[0].Status)
			require.Equal(t, int64(0), validators[0].Rank)
			require.Equal(t, int64(0), validators[0].Streak)

			infos := []types.ValidatorSigningInfo{}
			app.CustomSlashingKeeper.IterateValidatorSigningInfos(ctx, func(address sdk.ConsAddress, info types.ValidatorSigningInfo) (stop bool) {
				infos = append(infos, info)
				return false
			})
			require.Len(t, infos, 1)
			require.Equal(t, ctx.BlockHeight(), infos[0].StartHeight)
			require.Equal(t, time.Unix(0, 0).UTC(), infos[0].InactiveUntil.UTC())
			require.Equal(t, int64(0), infos[0].MischanceConfidence)
			require.Equal(t, int64(0), infos[0].Mischance)
			require.Equal(t, int64(0), infos[0].MissedBlocksCounter)
			require.Equal(t, int64(0), infos[0].ProducedBlocksCounter)
		})
	}
}

func createValidators(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context, accNum int) (validators []stakingtypes.Validator) {
	addrs := simapp.AddTestAddrsIncremental(app, ctx, accNum, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))

	for _, addr := range addrs {
		valAddr := sdk.ValAddress(addr)
		pubKey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em")
		require.NoError(t, err)

		validator, err := stakingtypes.NewValidator(
			"validator 1",
			sdk.NewDec(1234),
			valAddr,
			pubKey,
		)
		require.NoError(t, err)
		validators = append(validators, validator)
	}

	return
}
