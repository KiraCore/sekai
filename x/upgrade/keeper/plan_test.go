package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/KiraCore/sekai/simapp"
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
		MinHaltTime:          1,
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

	minHaltTime := time.Now()

	t.Log("Verify that a panic happens at the upgrade time/height")
	newCtx := ctx.WithBlockHeight(10).WithBlockTime(minHaltTime.Add(time.Second))

	t.Log("Verify that the upgrade can be successfully applied with a handler")
	app.UpgradeKeeper.SetUpgradeHandler("test", func(ctx sdk.Context, plan types.Plan) {})
	require.NotPanics(t, func() {
		app.UpgradeKeeper.ApplyUpgradePlan(newCtx, types.Plan{
			Height:               10,
			MinHaltTime:          minHaltTime.Unix(),
			RollbackChecksum:     "",
			MaxEnrolmentDuration: 0,
			Name:                 "test",
		})
	})

	plan, err := app.UpgradeKeeper.GetUpgradePlan(ctx)
	require.Nil(t, plan)
	require.NoError(t, err)
}

func TestPlanExecutionWithoutHandler(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	minHaltTime := time.Now()

	newCtx := ctx.WithBlockHeight(10).WithBlockTime(minHaltTime.Add(time.Second))

	require.Panics(t, func() {
		app.UpgradeKeeper.ApplyUpgradePlan(newCtx, types.Plan{
			Height:               10,
			MinHaltTime:          minHaltTime.Unix(),
			RollbackChecksum:     "",
			MaxEnrolmentDuration: 0,
			Name:                 "test",
		})
	})
}

func TestNoPlanExecutionBeforeTimeOrHeight(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	minHaltTime := time.Now()

	newCtx := ctx.WithBlockHeight(9).WithBlockTime(minHaltTime)

	require.NotPanics(t, func() {
		app.UpgradeKeeper.ApplyUpgradePlan(newCtx, types.Plan{
			Height:               10,
			MinHaltTime:          minHaltTime.Unix(),
			RollbackChecksum:     "",
			MaxEnrolmentDuration: 0,
			Name:                 "test",
		})
	})
}
