package keeper_test

import (
	"testing"
	"time"

	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestKeeperPlanGetSet(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	now := time.Now()
	ctx = ctx.WithBlockTime(now)

	plan, err := app.UpgradeKeeper.GetNextPlan(ctx)
	require.NoError(t, err)
	require.Nil(t, plan)

	newPlan := types.Plan{
		UpgradeTime:          now.Add(time.Second).Unix(),
		RollbackChecksum:     "checksum",
		MaxEnrolmentDuration: 2,
		Name:                 "plan",
		InstateUpgrade:       true,
		RebootRequired:       true,
	}

	err = app.UpgradeKeeper.SaveNextPlan(ctx, newPlan)
	require.NoError(t, err)

	plan, err = app.UpgradeKeeper.GetNextPlan(ctx)
	require.NoError(t, err)
	require.Equal(t, plan, &newPlan)

	newPlan.UpgradeTime = 0
	err = app.UpgradeKeeper.SaveNextPlan(ctx, newPlan)
	require.Error(t, err)
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
		coins := sdk.Coins{sdk.NewInt64Coin("ukex", 10000)}
		err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
		if err != nil {
			panic(err)
		}
		err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, acc1, coins)
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
			InstateUpgrade:       true,
			RebootRequired:       true,
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
			InstateUpgrade:       true,
			RebootRequired:       true,
			SkipHandler:          false,
		})
	})

	require.NotPanics(t, func() {
		app.UpgradeKeeper.ApplyUpgradePlan(newCtx, types.Plan{
			UpgradeTime:          upgradeTime.Unix(),
			RollbackChecksum:     "",
			MaxEnrolmentDuration: 0,
			Name:                 "test",
			InstateUpgrade:       true,
			RebootRequired:       true,
			SkipHandler:          true,
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
		InstateUpgrade:       true,
		RebootRequired:       true,
	}
	require.NotPanics(t, func() {
		app.UpgradeKeeper.ApplyUpgradePlan(newCtx, plan)
	})

	newCurrentPlan, err := app.UpgradeKeeper.GetCurrentPlan(ctx)
	require.NoError(t, err)
	require.NotNil(t, newCurrentPlan)
	require.Equal(t, *newCurrentPlan, plan)
}
