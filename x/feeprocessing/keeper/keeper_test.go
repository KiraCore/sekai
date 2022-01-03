package keeper_test

import (
	"bytes"
	"testing"

	simapp "github.com/KiraCore/sekai/app"
	kiratypes "github.com/KiraCore/sekai/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestNewKeeper_SenderCoinsHistory(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	addr := addrs[0]

	savedFees := app.FeeProcessingKeeper.GetSenderCoinsHistory(ctx, addr)
	require.True(t, savedFees.IsEqual(sdk.Coins{}))

	fees := sdk.Coins{sdk.NewInt64Coin("ukex", 100)}
	app.FeeProcessingKeeper.SetSenderCoinsHistory(ctx, addr, fees)

	savedFees = app.FeeProcessingKeeper.GetSenderCoinsHistory(ctx, addr)
	require.True(t, savedFees.IsEqual(fees))
}

func TestNewKeeper_Executions(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	addr := addrs[0]

	// initial executions listing is empty
	executions := app.FeeProcessingKeeper.GetExecutionsStatus(ctx)
	require.True(t, len(executions) == 0)

	msg1 := tokenstypes.NewMsgUpsertTokenRate(addr, "ukex", sdk.NewDec(1), true)
	app.FeeProcessingKeeper.AddExecutionStart(ctx, msg1)
	executions = app.FeeProcessingKeeper.GetExecutionsStatus(ctx)
	require.True(t, len(executions) == 1)

	msg2 := tokenstypes.NewMsgUpsertTokenAlias(addr, "KEX", "Kira", "", 10, []string{"ukex"})
	app.FeeProcessingKeeper.AddExecutionStart(ctx, msg2)
	executions = app.FeeProcessingKeeper.GetExecutionsStatus(ctx)
	require.True(t, len(executions) == 2)

	msg3 := tokenstypes.NewMsgUpsertTokenRate(addr, "ukex", sdk.NewDec(1), true)
	app.FeeProcessingKeeper.AddExecutionStart(ctx, msg3)
	executions = app.FeeProcessingKeeper.GetExecutionsStatus(ctx)
	require.True(t, len(executions) == 3)

	app.FeeProcessingKeeper.SetExecutionStatusSuccess(ctx, msg1)
	app.FeeProcessingKeeper.SetExecutionStatusSuccess(ctx, msg3)
	executions = app.FeeProcessingKeeper.GetExecutionsStatus(ctx)
	successFlaggedCount := 0
	for _, exec := range executions {
		if bytes.Equal(exec.FeePayer, msg1.Proposer) && exec.MsgType == msg1.Type() && exec.Success == true {
			successFlaggedCount += 1
		}
	}
	require.True(t, successFlaggedCount == 2)

	successFlaggedCount = 0
	for _, exec := range executions {
		if bytes.Equal(exec.FeePayer, msg1.Proposer) && exec.MsgType == msg2.Type() && exec.Success == true {
			successFlaggedCount += 1
		}
	}
	require.True(t, successFlaggedCount == 0)
}

func TestNewKeeper_SendCoinsFromAccountToModule(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	addr := addrs[0]

	initialBalance := app.BankKeeper.GetBalance(ctx, addr, "ukex")
	coins := sdk.Coins{sdk.NewInt64Coin("ukex", 10000)}
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, coins)

	fees := sdk.Coins{sdk.NewInt64Coin("ukex", 100)}
	app.FeeProcessingKeeper.SendCoinsFromAccountToModule(ctx, addr, authtypes.FeeCollectorName, fees)

	balance := app.BankKeeper.GetBalance(ctx, addr, "ukex")
	require.Equal(t, balance.Amount.Int64(), initialBalance.Amount.Int64()+int64(10000-100))

	feeCollectorAcc := app.AccountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
	balance = app.BankKeeper.GetBalance(ctx, feeCollectorAcc.GetAddress(), "ukex")
	require.True(t, balance.Amount.Int64() == 100)

	savedFees := app.FeeProcessingKeeper.GetSenderCoinsHistory(ctx, addr)
	require.True(t, savedFees.IsEqual(fees))
}

func TestNewKeeper_SendCoinsFromModuleToAccount(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	addr := addrs[0]
	initialBalance := app.BankKeeper.GetBalance(ctx, addr, "ukex")

	coins := sdk.Coins{sdk.NewInt64Coin("ukex", 10000)}
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, coins)

	fees := sdk.Coins{sdk.NewInt64Coin("ukex", 100)}
	returnFees := sdk.Coins{sdk.NewInt64Coin("ukex", 10)}
	app.FeeProcessingKeeper.SendCoinsFromAccountToModule(ctx, addr, authtypes.FeeCollectorName, fees)
	app.FeeProcessingKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, addr, returnFees)

	balance := app.BankKeeper.GetBalance(ctx, addr, "ukex")
	require.True(t, balance.Amount.Int64() == initialBalance.Amount.Int64()+(10000-100+10))

	feeCollectorAcc := app.AccountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
	balance = app.BankKeeper.GetBalance(ctx, feeCollectorAcc.GetAddress(), "ukex")
	require.True(t, balance.Amount.Int64() == 100-10)

	savedFees := app.FeeProcessingKeeper.GetSenderCoinsHistory(ctx, addr)
	require.True(t, savedFees.IsEqual(fees.Sub(returnFees)))
}

func TestNewKeeper_ProcessExecutionFeeReturn(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 3, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	addr := addrs[0]
	addr2 := addrs[1]
	addr3 := addrs[2]

	initialBalance := app.BankKeeper.GetBalance(ctx, addr2, "ukex")

	coins := sdk.Coins{sdk.NewInt64Coin("ukex", 10000)}
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, coins)
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr2, coins)
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr3, coins)

	app.CustomGovKeeper.SetExecutionFee(ctx, &govtypes.ExecutionFee{
		Name:              kiratypes.MsgTypeUpsertTokenRate,
		TransactionType:   kiratypes.MsgTypeUpsertTokenRate,
		ExecutionFee:      1000,
		FailureFee:        100,
		Timeout:           0,
		DefaultParameters: 0,
	})

	// check failure fee
	fees := sdk.Coins{sdk.NewInt64Coin("ukex", 1000)}
	app.FeeProcessingKeeper.SendCoinsFromAccountToModule(ctx, addr, authtypes.FeeCollectorName, fees)
	msg := tokenstypes.NewMsgUpsertTokenRate(addr, "ukex", sdk.NewDec(1), true)
	app.FeeProcessingKeeper.AddExecutionStart(ctx, msg)
	app.FeeProcessingKeeper.ProcessExecutionFeeReturn(ctx)

	feeCollectorAcc := app.AccountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
	balance := app.BankKeeper.GetBalance(ctx, feeCollectorAcc.GetAddress(), "ukex")
	require.True(t, balance.Amount.Int64() == 100)

	executions := app.FeeProcessingKeeper.GetExecutionsStatus(ctx)
	require.True(t, len(executions) == 0)

	// check success fee
	app.FeeProcessingKeeper.SendCoinsFromAccountToModule(ctx, addr, authtypes.FeeCollectorName, fees)
	app.FeeProcessingKeeper.AddExecutionStart(ctx, msg)
	app.FeeProcessingKeeper.SetExecutionStatusSuccess(ctx, msg)
	app.FeeProcessingKeeper.ProcessExecutionFeeReturn(ctx)

	balance = app.BankKeeper.GetBalance(ctx, feeCollectorAcc.GetAddress(), "ukex")
	require.True(t, balance.Amount.Int64() == 100+1000)

	executions = app.FeeProcessingKeeper.GetExecutionsStatus(ctx)
	require.True(t, len(executions) == 0)

	// check success return when two message types are same but addresses are different
	app.FeeProcessingKeeper.SendCoinsFromAccountToModule(ctx, addr2, authtypes.FeeCollectorName, fees)
	app.FeeProcessingKeeper.SendCoinsFromAccountToModule(ctx, addr3, authtypes.FeeCollectorName, fees)
	msg2 := tokenstypes.NewMsgUpsertTokenRate(addr2, "ukex", sdk.NewDec(1), true)
	msg3 := tokenstypes.NewMsgUpsertTokenRate(addr3, "ukex", sdk.NewDec(1), true)
	app.FeeProcessingKeeper.AddExecutionStart(ctx, msg3)
	app.FeeProcessingKeeper.AddExecutionStart(ctx, msg2)
	app.FeeProcessingKeeper.SetExecutionStatusSuccess(ctx, msg2)
	app.FeeProcessingKeeper.ProcessExecutionFeeReturn(ctx)

	balance = app.BankKeeper.GetBalance(ctx, addr2, "ukex")
	require.Equal(t, balance.Amount.Int64(), initialBalance.Amount.Int64()+(10000-1000)) // success fee
	balance = app.BankKeeper.GetBalance(ctx, addr3, "ukex")
	require.Equal(t, balance.Amount.Int64(), initialBalance.Amount.Int64()+(10000-100)) // failure fee
}
