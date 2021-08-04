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

	plan, err := app.UpgradeKeeper.GetUpgradePlan(ctx)
	require.NoError(t, err)
	require.Nil(t, plan)

	newPlan := types.Plan{
		MinUpgradeTime:       1,
		RollbackChecksum:     "checksum",
		MaxEnrolmentDuration: 2,
		Name:                 "plan",
	}

	app.UpgradeKeeper.SaveUpgradePlan(ctx, newPlan)

	plan, err = app.UpgradeKeeper.GetUpgradePlan(ctx)
	require.NoError(t, err)
	require.Equal(t, plan, &newPlan)
}

func TestPlanExecutionWithHandler(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})
	acc1 := sdk.AccAddress("test________________")

	minUpgradeTime := time.Now()

	t.Log("Verify that a panic happens at the upgrade time/height")
	newCtx := ctx.WithBlockHeight(10).WithBlockTime(minUpgradeTime.Add(time.Second))

	t.Log("Verify that the upgrade can be successfully applied with a handler")
	app.UpgradeKeeper.SetUpgradeHandler("test", func(ctx sdk.Context, plan types.Plan) {
		err := app.BankKeeper.SetBalance(ctx, acc1, sdk.NewInt64Coin("ukex", 10000))
		if err != nil {
			panic(err)
		}
	})
	require.NotPanics(t, func() {
		app.UpgradeKeeper.ApplyUpgradePlan(newCtx, types.Plan{
			Height:               10,
			MinUpgradeTime:       minUpgradeTime.Unix(),
			RollbackChecksum:     "",
			MaxEnrolmentDuration: 0,
			Name:                 "test",
		})
	})

	plan, err := app.UpgradeKeeper.GetUpgradePlan(ctx)
	require.Nil(t, plan)
	require.NoError(t, err)

	coin := app.BankKeeper.GetBalance(ctx, acc1, "ukex")
	require.Equal(t, coin, sdk.NewInt64Coin("ukex", 10000))
}

func TestPlanExecutionWithoutHandler(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	minUpgradeTime := time.Now()

	newCtx := ctx.WithBlockHeight(10).WithBlockTime(minUpgradeTime.Add(time.Second))

	require.Panics(t, func() {
		app.UpgradeKeeper.ApplyUpgradePlan(newCtx, types.Plan{
			Height:               10,
			MinUpgradeTime:       minUpgradeTime.Unix(),
			RollbackChecksum:     "",
			MaxEnrolmentDuration: 0,
			Name:                 "test",
		})
	})
}

func TestNoPlanExecutionBeforeTimeOrHeight(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	minUpgradeTime := time.Now()

	newCtx := ctx.WithBlockHeight(9).WithBlockTime(minUpgradeTime)

	require.NotPanics(t, func() {
		app.UpgradeKeeper.ApplyUpgradePlan(newCtx, types.Plan{
			Height:               10,
			MinUpgradeTime:       minUpgradeTime.Unix(),
			RollbackChecksum:     "",
			MaxEnrolmentDuration: 0,
			Name:                 "test",
		})
	})
}
