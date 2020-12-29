package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/x/slashing/testslashing"
	"github.com/KiraCore/sekai/x/staking"
	"github.com/KiraCore/sekai/x/staking/teststaking"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Test a new validator entering the validator set
// Ensure that SigningInfo.StartHeight is set correctly
// and that they are not immediately inactivated
func TestHandleNewValidator(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addrDels := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(200))
	valAddrs := simapp.ConvertAddrsToValAddrs(addrDels)
	pks := simapp.CreateTestPubKeys(1)
	addr, val := valAddrs[0], pks[0]
	tstaking := teststaking.NewHelper(t, ctx, app.CustomStakingKeeper, app.CustomGovKeeper)
	ctx = ctx.WithBlockHeight(app.SlashingKeeper.SignedBlocksWindow(ctx) + 1)

	// Validator created
	tstaking.CreateValidator(addr, val, true)

	staking.EndBlocker(ctx, app.CustomStakingKeeper)

	// Now a validator, for two blocks
	app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 100, true)
	ctx = ctx.WithBlockHeight(app.SlashingKeeper.SignedBlocksWindow(ctx) + 2)
	app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 100, false)

	info, found := app.SlashingKeeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(val.Address()))
	require.True(t, found)
	require.Equal(t, app.SlashingKeeper.SignedBlocksWindow(ctx)+1, info.StartHeight)
	require.Equal(t, int64(2), info.IndexOffset)
	require.Equal(t, int64(1), info.MissedBlocksCounter)
	require.Equal(t, time.Unix(0, 0).UTC(), info.JailedUntil)

	// validator should be bonded still, should not have been jailed or slashed
	validator, _ := app.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val))
	require.Equal(t, stakingtypes.Active, validator.GetStatus())
	bondPool := app.StakingKeeper.GetBondedPool(ctx)
	expTokens := sdk.TokensFromConsensusPower(100)
	require.True(t, expTokens.Equal(app.BankKeeper.GetBalance(ctx, bondPool.GetAddress(), app.StakingKeeper.BondDenom(ctx)).Amount))
}

// Test an inactivated validator being "down" twice
// Ensure that they're only inactivated once
func TestHandleAlreadyInactive(t *testing.T) {
	// initial setup
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addrDels := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(200))
	valAddrs := simapp.ConvertAddrsToValAddrs(addrDels)
	pks := simapp.CreateTestPubKeys(1)
	addr, val := valAddrs[0], pks[0]
	power := int64(100)
	tstaking := teststaking.NewHelper(t, ctx, app.CustomStakingKeeper, app.CustomGovKeeper)

	tstaking.CreateValidator(addr, val, true)

	staking.EndBlocker(ctx, app.CustomStakingKeeper)

	// 1000 first blocks OK
	height := int64(0)
	for ; height < app.SlashingKeeper.SignedBlocksWindow(ctx); height++ {
		ctx = ctx.WithBlockHeight(height)
		app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), power, true)
	}

	// 501 blocks missed
	for ; height < app.SlashingKeeper.SignedBlocksWindow(ctx)+(app.SlashingKeeper.SignedBlocksWindow(ctx)-app.SlashingKeeper.MinSignedPerWindow(ctx))+1; height++ {
		ctx = ctx.WithBlockHeight(height)
		app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), power, false)
	}

	// end block
	staking.EndBlocker(ctx, app.CustomStakingKeeper)

	// validator should have been inactivated
	validator, _ := app.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val))
	require.Equal(t, stakingtypes.Inactive, validator.GetStatus())

	// another block missed
	ctx = ctx.WithBlockHeight(height)
	app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), power, false)

	// validator should be in inactive status yet
	validator, _ = app.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val))
	require.Equal(t, stakingtypes.Inactive, validator.GetStatus())
}

// Test a validator dipping in and out of the validator set
// Ensure that missed blocks are tracked correctly and that
// the start height of the signing info is reset correctly
func TestValidatorDippingInAndOut(t *testing.T) {
	// initial setup
	// TestParams set the SignedBlocksWindow to 1000 and MaxMissedBlocksPerWindow to 500
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	app.CustomSlashingKeeper.SetParams(ctx, testslashing.TestParams())

	params := app.StakingKeeper.GetParams(ctx)
	params.MaxValidators = 1
	app.StakingKeeper.SetParams(ctx, params)
	power := int64(100)

	pks := simapp.CreateTestPubKeys(3)
	simapp.AddTestAddrsFromPubKeys(app, ctx, pks, sdk.TokensFromConsensusPower(200))

	addr, val := pks[0].Address(), pks[0]
	consAddr := sdk.ConsAddress(addr)
	tstaking := teststaking.NewHelper(t, ctx, app.CustomStakingKeeper, app.CustomGovKeeper)
	valAddr := sdk.ValAddress(addr)

	tstaking.CreateValidator(valAddr, val, true)
	staking.EndBlocker(ctx, app.CustomStakingKeeper)

	// 100 first blocks OK
	height := int64(0)
	for ; height < int64(100); height++ {
		ctx = ctx.WithBlockHeight(height)
		app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), power, true)
	}

	// kick first validator out of validator set
	tstaking.CreateValidator(sdk.ValAddress(pks[1].Address()), pks[1], true)
	validatorUpdates := staking.EndBlocker(ctx, app.CustomStakingKeeper)
	require.Equal(t, 2, len(validatorUpdates))
	tstaking.CheckValidator(valAddr, stakingtypes.Inactive, false)

	// 600 more blocks happened
	height = 700
	ctx = ctx.WithBlockHeight(height)

	validatorUpdates = staking.EndBlocker(ctx, app.CustomStakingKeeper)
	require.Equal(t, 2, len(validatorUpdates))
	tstaking.CheckValidator(valAddr, stakingtypes.Active, false)

	// shouldn't be inactive/kicked yet
	tstaking.CheckValidator(valAddr, stakingtypes.Active, false)

	// validator misses 500 more blocks, 501 total
	latest := height
	for ; height < latest+500; height++ {
		ctx = ctx.WithBlockHeight(height)
		app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 1, false)
	}

	// should now be inactive & kicked
	staking.EndBlocker(ctx, app.CustomStakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Inactive, true)

	// check all the signing information
	signInfo, found := app.SlashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
	require.True(t, found)
	require.Equal(t, int64(0), signInfo.MissedBlocksCounter)
	require.Equal(t, int64(0), signInfo.IndexOffset)
	// array should be cleared
	for offset := int64(0); offset < app.SlashingKeeper.SignedBlocksWindow(ctx); offset++ {
		missed := app.SlashingKeeper.GetValidatorMissedBlockBitArray(ctx, consAddr, offset)
		require.False(t, missed)
	}

	// some blocks pass
	height = int64(5000)
	ctx = ctx.WithBlockHeight(height)

	// validator rejoins and starts signing again
	app.StakingKeeper.Unjail(ctx, consAddr)
	app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 1, true)
	height++

	// validator should not be kicked since we reset counter/array when it was jailed
	staking.EndBlocker(ctx, app.CustomStakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Active, false)

	// validator misses 501 blocks
	latest = height
	for ; height < latest+501; height++ {
		ctx = ctx.WithBlockHeight(height)
		app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 1, false)
	}

	// validator should now be jailed & kicked
	staking.EndBlocker(ctx, app.CustomStakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Inactive, true)
}
