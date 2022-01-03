package keeper_test

import (
	"testing"
	"time"

	simapp "github.com/KiraCore/sekai/app"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
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
			UpgradeTime:               upgradeTime.Unix(),
			RollbackChecksum:          "",
			MaxEnrolmentDuration:      0,
			Name:                      "test",
			InstateUpgrade:            true,
			RebootRequired:            true,
			ProcessedNoVoteValidators: true,
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
			UpgradeTime:               upgradeTime.Unix(),
			RollbackChecksum:          "",
			MaxEnrolmentDuration:      0,
			Name:                      "test",
			InstateUpgrade:            true,
			RebootRequired:            true,
			SkipHandler:               false,
			ProcessedNoVoteValidators: true,
		})
	})

	require.NotPanics(t, func() {
		app.UpgradeKeeper.ApplyUpgradePlan(newCtx, types.Plan{
			UpgradeTime:               upgradeTime.Unix(),
			RollbackChecksum:          "",
			MaxEnrolmentDuration:      0,
			Name:                      "test",
			InstateUpgrade:            true,
			RebootRequired:            true,
			SkipHandler:               true,
			ProcessedNoVoteValidators: true,
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
		UpgradeTime:               upgradeTime.Unix(),
		RollbackChecksum:          "",
		MaxEnrolmentDuration:      0,
		Name:                      "test",
		InstateUpgrade:            true,
		RebootRequired:            true,
		ProcessedNoVoteValidators: true,
	}
	require.NotPanics(t, func() {
		app.UpgradeKeeper.ApplyUpgradePlan(newCtx, plan)
	})

	newCurrentPlan, err := app.UpgradeKeeper.GetCurrentPlan(ctx)
	require.NoError(t, err)
	require.NotNil(t, newCurrentPlan)
	require.Equal(t, *newCurrentPlan, plan)
}

func TestNoPlanExecutionBeforeNotVotedValidatorsProcess(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})
	acc1 := sdk.AccAddress("test1_______________")
	acc2 := sdk.AccAddress("test2_______________")
	acc3 := sdk.AccAddress("test3_______________")

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

	propoaslID, err := app.CustomGovKeeper.CreateAndSaveProposalWithContent(ctx, "test proposal", "test proposal description", &types.ProposalSoftwareUpgrade{})
	require.NoError(t, err)

	pubkeys := simapp.CreateTestPubKeys(3)

	validator1, err := stakingtypes.NewValidator(sdk.ValAddress(acc1), pubkeys[0])
	require.NoError(t, err)
	validator1.Status = stakingtypes.Inactive
	app.CustomStakingKeeper.AddValidator(ctx, validator1)

	validator2, err := stakingtypes.NewValidator(sdk.ValAddress(acc2), pubkeys[1])
	require.NoError(t, err)
	validator2.Status = stakingtypes.Active
	app.CustomStakingKeeper.AddValidator(ctx, validator2)

	app.CustomGovKeeper.SaveVote(ctx, govtypes.Vote{
		ProposalId: propoaslID,
		Voter:      acc1,
		Option:     govtypes.OptionYes,
	})

	app.CustomGovKeeper.SaveVote(ctx, govtypes.Vote{
		ProposalId: propoaslID,
		Voter:      acc2,
		Option:     govtypes.OptionNo,
	})

	app.CustomGovKeeper.SaveNetworkActor(ctx, govtypes.NetworkActor{
		Address:     acc1,
		Status:      govtypes.Active,
		Permissions: &govtypes.Permissions{Whitelist: []uint32{uint32(govtypes.PermVoteSoftwareUpgradeProposal)}},
	})

	app.CustomGovKeeper.SaveNetworkActor(ctx, govtypes.NetworkActor{
		Address:     acc3,
		Status:      govtypes.Active,
		Permissions: &govtypes.Permissions{Whitelist: []uint32{uint32(govtypes.PermVoteSoftwareUpgradeProposal)}},
	})

	require.NotPanics(t, func() {
		plan := types.Plan{
			UpgradeTime:               upgradeTime.Unix(),
			RollbackChecksum:          "",
			MaxEnrolmentDuration:      0,
			Name:                      "test",
			InstateUpgrade:            true,
			RebootRequired:            true,
			ProposalID:                propoaslID,
			ProcessedNoVoteValidators: false,
		}
		app.UpgradeKeeper.SaveNextPlan(ctx, plan)
		app.UpgradeKeeper.ApplyUpgradePlan(newCtx, plan)
	})

	plan, err := app.UpgradeKeeper.GetNextPlan(newCtx)
	require.NotNil(t, plan)
	require.NoError(t, err)
	require.True(t, plan.ProcessedNoVoteValidators)

	coin := app.BankKeeper.GetBalance(newCtx, acc1, "ukex")
	require.Equal(t, coin, sdk.NewInt64Coin("ukex", 0))

	validator2, err = app.CustomStakingKeeper.GetValidator(newCtx, sdk.ValAddress(acc2))
	require.NoError(t, err)
	require.Equal(t, validator2.Status, stakingtypes.Paused)

	require.NotPanics(t, func() {
		app.UpgradeKeeper.ApplyUpgradePlan(newCtx, *plan)
	})

	plan, err = app.UpgradeKeeper.GetNextPlan(ctx)
	require.Nil(t, plan)
	require.NoError(t, err)

	coin = app.BankKeeper.GetBalance(ctx, acc1, "ukex")
	require.Equal(t, coin, sdk.NewInt64Coin("ukex", 10000))
}
