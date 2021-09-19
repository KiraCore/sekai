package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/x/upgrade/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestKeeperPlanGetSet(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	plan, err := app.UpgradeKeeper.GetNextPlan(ctx)
	require.NoError(t, err)
	require.Nil(t, plan)

	newPlan := types.Plan{
		UpgradeTime:          1,
		RollbackChecksum:     "checksum",
		MaxEnrolmentDuration: 2,
		Name:                 "plan",
	}

	app.UpgradeKeeper.SaveNextPlan(ctx, newPlan)

	plan, err = app.UpgradeKeeper.GetNextPlan(ctx)
	require.NoError(t, err)
	require.Equal(t, plan, &newPlan)
}

func TestPlanExecutionWithHandler(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})
	acc1 := sdk.AccAddress("test________________")

	upgradeTime := time.Now()

	t.Log("Verify that a panic happens at the upgrade time/height")
	newCtx := ctx.WithBlockHeight(10).WithBlockTime(upgradeTime.Add(time.Second))

	t.Log("Verify that the upgrade can be successfully applied with a handler")
	app.UpgradeKeeper.SetUpgradeHandler("test", func(ctx sdk.Context, plan types.Plan) {
		err := app.BankKeeper.SetBalance(ctx, acc1, sdk.NewInt64Coin("ukex", 10000))
		if err != nil {
			panic(err)
		}
	})
	require.NotPanics(t, func() {
		app.UpgradeKeeper.ApplyUpgradePlan(newCtx, types.Plan{
			UpgradeTime:          upgradeTime.Unix(),
			RollbackChecksum:     "",
			MaxEnrolmentDuration: 0,
			Name:                 "test",
		})
	})

	plan, err := app.UpgradeKeeper.GetNextPlan(ctx)
	require.Nil(t, plan)
	require.NoError(t, err)

	coin := app.BankKeeper.GetBalance(ctx, acc1, "ukex")
	require.Equal(t, coin, sdk.NewInt64Coin("ukex", 10000))
}

func TestPlanExecutionWithoutHandler(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	upgradeTime := time.Now()
	newCtx := ctx.WithBlockHeight(10).WithBlockTime(upgradeTime.Add(time.Second))

	require.Panics(t, func() {
		app.UpgradeKeeper.ApplyUpgradePlan(newCtx, types.Plan{
			UpgradeTime:          upgradeTime.Unix(),
			RollbackChecksum:     "",
			MaxEnrolmentDuration: 0,
			Name:                 "test",
		})
	})
}

func TestNoPlanExecutionBeforeTime(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	upgradeTime := time.Now()
	newCtx := ctx.WithBlockHeight(9).WithBlockTime(upgradeTime)

	app.UpgradeKeeper.SetUpgradeHandler("test", func(ctx sdk.Context, plan types.Plan) {
	})

	plan := types.Plan{
		UpgradeTime:          upgradeTime.Unix(),
		RollbackChecksum:     "",
		MaxEnrolmentDuration: 0,
		Name:                 "test",
	}
	require.NotPanics(t, func() {
		app.UpgradeKeeper.ApplyUpgradePlan(newCtx, plan)
	})

	newCurrentPlan, err := app.UpgradeKeeper.GetCurrentPlan(ctx)
	require.NoError(t, err)
	require.NotNil(t, newCurrentPlan)
	require.Equal(t, *newCurrentPlan, plan)
}
