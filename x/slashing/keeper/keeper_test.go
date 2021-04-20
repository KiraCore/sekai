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
	ctx = ctx.WithBlockHeight(1)

	// Validator created
	tstaking.CreateValidator(addr, val, true)

	staking.EndBlocker(ctx, app.CustomStakingKeeper)

	// Now a validator, for two blocks
	app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 100, true)
	ctx = ctx.WithBlockHeight(2)
	app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 100, false)

	info, found := app.CustomSlashingKeeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(val.Address()))
	require.True(t, found)
	require.Equal(t, int64(1), info.StartHeight)
	require.Equal(t, int64(0), info.Mischance)
	require.Equal(t, int64(1), info.MissedBlocksCounter)
	require.Equal(t, int64(1), info.ProducedBlocksCounter)
	require.Equal(t, time.Unix(0, 0).UTC(), info.InactiveUntil)

	// validator should be active still, should not have been inactivated
	validator, _ := app.CustomStakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val))
	require.Equal(t, stakingtypes.Active, validator.GetStatus())
}

// Test missed blocks, produced blocks, mischance counter changes
func TestMissedBlockAndRankStreakCounter(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addrDels := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(200))
	valAddrs := simapp.ConvertAddrsToValAddrs(addrDels)
	pks := simapp.CreateTestPubKeys(1)
	addr, val := valAddrs[0], pks[0]
	valAddr := sdk.ValAddress(addr)
	tstaking := teststaking.NewHelper(t, ctx, app.CustomStakingKeeper, app.CustomGovKeeper)
	ctx = ctx.WithBlockHeight(1)

	// Validator created
	tstaking.CreateValidator(addr, val, true)

	staking.EndBlocker(ctx, app.CustomStakingKeeper)

	// Now a validator, for two blocks
	app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 100, true)
	ctx = ctx.WithBlockHeight(2)
	app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 100, false)

	v := tstaking.CheckValidator(valAddr, stakingtypes.Active, false)
	require.Equal(t, v.Rank, int64(0))
	require.Equal(t, v.Streak, int64(0))

	info, found := app.CustomSlashingKeeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(val.Address()))
	require.True(t, found)
	require.Equal(t, int64(0), info.Mischance)
	require.Equal(t, int64(1), info.MissedBlocksCounter)
	require.Equal(t, int64(1), info.ProducedBlocksCounter)

	height := ctx.BlockHeight() + 1
	for i := int64(0); i < 10; i++ {
		ctx = ctx.WithBlockHeight(height + i)
		app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 100, false)
	}
	ctx = ctx.WithBlockHeight(height + 10)
	app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 100, true)
	info, found = app.CustomSlashingKeeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(val.Address()))
	require.True(t, found)
	require.Equal(t, int64(0), info.Mischance)
	require.Equal(t, int64(11), info.MissedBlocksCounter)
	require.Equal(t, int64(2), info.ProducedBlocksCounter)

	v = tstaking.CheckValidator(valAddr, stakingtypes.Active, false)
	require.Equal(t, v.Rank, int64(1))
	require.Equal(t, v.Streak, int64(1))

	// sign 100 blocks successfully
	height = ctx.BlockHeight() + 1
	for i := int64(0); i < 100; i++ {
		ctx = ctx.WithBlockHeight(height + i)
		app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 100, true)
	}

	v = tstaking.CheckValidator(valAddr, stakingtypes.Active, false)
	require.Equal(t, v.Rank, int64(101))
	require.Equal(t, v.Streak, int64(101))

	// miss one block
	ctx = ctx.WithBlockHeight(height + 100)
	app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 100, false)
	v = tstaking.CheckValidator(valAddr, stakingtypes.Active, false)
	require.Equal(t, v.Rank, int64(91))
	require.Equal(t, v.Streak, int64(0))

	app.CustomStakingKeeper.Inactivate(ctx, valAddr)
	v = tstaking.CheckValidator(valAddr, stakingtypes.Inactive, true)
	require.Equal(t, v.Rank, int64(45))
	require.Equal(t, v.Streak, int64(0))

	app.CustomStakingKeeper.Activate(ctx, valAddr)
	// miss 5 blocks
	height = ctx.BlockHeight() + 1
	for i := int64(0); i < 5; i++ {
		ctx = ctx.WithBlockHeight(height + i)
		app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 100, false)
	}
	v = tstaking.CheckValidator(valAddr, stakingtypes.Active, false)
	require.Equal(t, v.Rank, int64(0))
	require.Equal(t, v.Streak, int64(0))
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
	for ; height < 1000; height++ {
		ctx = ctx.WithBlockHeight(height)
		app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), power, true)
	}

	properties := app.CustomGovKeeper.GetNetworkProperties(ctx)
	// miss 11 blocks for mischance confidence
	for ; height < 1000+int64(properties.MischanceConfidence)+1; height++ {
		ctx = ctx.WithBlockHeight(height)
		app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), power, false)
	}

	// info correctness after the overflow of mischance confidence
	info, found := app.CustomSlashingKeeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(val.Address()))
	require.True(t, found)
	require.Equal(t, int64(1), info.Mischance)
	require.Equal(t, int64(999), info.LastPresentBlock)

	// miss 110 blocks after mischance confidence happen
	for ; height < 1000+int64(properties.MaxMischance+properties.MischanceConfidence)+1; height++ {
		ctx = ctx.WithBlockHeight(height)
		app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), power, false)
	}

	// end block
	staking.EndBlocker(ctx, app.CustomStakingKeeper)

	// validator should have been inactivated
	validator, _ := app.CustomStakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val))
	require.Equal(t, stakingtypes.Inactive, validator.GetStatus())

	// another block missed
	ctx = ctx.WithBlockHeight(height)
	app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), power, false)

	// validator should be in inactive status yet
	validator, _ = app.CustomStakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val))
	require.Equal(t, stakingtypes.Inactive, validator.GetStatus())
}

// Test a validator dipping in and out of the validator set
// Ensure that counters for mischance, last block produced are correct and uptime counters are reset correctly
func TestValidatorDippingInAndOut(t *testing.T) {
	// initial setup
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	app.CustomSlashingKeeper.SetParams(ctx, testslashing.TestParams())

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
		app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), power, true)
	}

	// add one more validator into the set
	tstaking.CreateValidator(sdk.ValAddress(pks[1].Address()), pks[1], true)
	validatorUpdates := staking.EndBlocker(ctx, app.CustomStakingKeeper)
	require.Equal(t, 1, len(validatorUpdates))
	tstaking.CheckValidator(valAddr, stakingtypes.Active, false)
	tstaking.CheckValidator(sdk.ValAddress(pks[1].Address()), stakingtypes.Active, false)

	// 600 more blocks happened
	height = 700
	ctx = ctx.WithBlockHeight(height)

	validatorUpdates = staking.EndBlocker(ctx, app.CustomStakingKeeper)
	require.Equal(t, 0, len(validatorUpdates))
	tstaking.CheckValidator(valAddr, stakingtypes.Active, false)

	// shouldn't be inactive/kicked yet
	tstaking.CheckValidator(valAddr, stakingtypes.Active, false)

	// validator misses 500 more blocks, 501 total
	latest := height
	for ; height < latest+501; height++ {
		ctx = ctx.WithBlockHeight(height)
		app.CustomSlashingKeeper.HandleValidatorSignature(ctx, addr, 1, false)
	}

	// should now be inactive & kicked
	staking.EndBlocker(ctx, app.CustomStakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Inactive, true)

	// check all the signing information
	signInfo, found := app.CustomSlashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
	require.True(t, found)
	require.Equal(t, int64(390), signInfo.Mischance)
	require.Equal(t, int64(99), signInfo.LastPresentBlock)
	require.Equal(t, int64(501), signInfo.MissedBlocksCounter)
	require.Equal(t, int64(100), signInfo.ProducedBlocksCounter)

	// some blocks pass
	height = int64(5000)
	ctx = ctx.WithBlockHeight(height)

	// Try pausing on inactive node here, should fail
	err := app.CustomStakingKeeper.Pause(ctx, valAddr)
	require.Error(t, err)

	// validator rejoins and starts signing again
	app.CustomStakingKeeper.Activate(ctx, valAddr)
	app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 1, true)
	height++

	// validator should not be kicked since we reset counter/array when it was jailed
	staking.EndBlocker(ctx, app.CustomStakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Active, false)

	// Try pausing on active node here, should success
	err = app.CustomStakingKeeper.Pause(ctx, valAddr)
	require.NoError(t, err)
	staking.EndBlocker(ctx, app.CustomStakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Paused, false)

	// validator misses 501 blocks
	latest = height
	for ; height < latest+501; height++ {
		ctx = ctx.WithBlockHeight(height)
		app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 1, false)
	}

	// validator should not be in inactive status since node is paused
	staking.EndBlocker(ctx, app.CustomStakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Paused, false)

	// Try activating paused node: should unpause but it's activating - should fail
	err = app.CustomStakingKeeper.Activate(ctx, valAddr)
	require.Error(t, err)
	staking.EndBlocker(ctx, app.CustomStakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Paused, false)

	// Unpause node and it should be active
	err = app.CustomStakingKeeper.Unpause(ctx, valAddr)
	require.NoError(t, err)
	staking.EndBlocker(ctx, app.CustomStakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Active, false)

	// After reentering after unpause, check if signature info is recovered correctly
	signInfo, found = app.CustomSlashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
	require.True(t, found)
	require.Equal(t, int64(491), signInfo.Mischance)
	require.Equal(t, int64(5000), signInfo.LastPresentBlock)
	require.Equal(t, int64(1002), signInfo.MissedBlocksCounter)
	require.Equal(t, int64(101), signInfo.ProducedBlocksCounter)

	// Miss another 501 blocks
	latest = height
	for ; height < latest+501; height++ {
		ctx = ctx.WithBlockHeight(height)
		app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 1, false)
	}

	// validator should be in inactive status
	staking.EndBlocker(ctx, app.CustomStakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Inactive, true)
}

func TestValidatorLifecycle(t *testing.T) {
	// initial setup
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	app.CustomSlashingKeeper.SetParams(ctx, testslashing.TestParams())
	properties := app.CustomGovKeeper.GetNetworkProperties(ctx)

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
		app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), power, true)
	}

	// check info
	signInfo, found := app.CustomSlashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
	require.True(t, found)
	require.Equal(t, consAddr.String(), signInfo.Address)
	require.Equal(t, int64(0), signInfo.StartHeight)
	require.Equal(t, time.Unix(0, 0).UTC(), signInfo.InactiveUntil.UTC())
	require.Equal(t, false, signInfo.Tombstoned)
	require.Equal(t, int64(0), signInfo.Mischance)
	require.Equal(t, int64(99), signInfo.LastPresentBlock)
	require.Equal(t, int64(0), signInfo.MissedBlocksCounter)
	require.Equal(t, int64(100), signInfo.ProducedBlocksCounter)

	// add one more validator into the set
	tstaking.CreateValidator(sdk.ValAddress(pks[1].Address()), pks[1], true)
	validatorUpdates := staking.EndBlocker(ctx, app.CustomStakingKeeper)
	require.Equal(t, 1, len(validatorUpdates))
	tstaking.CheckValidator(valAddr, stakingtypes.Active, false)
	tstaking.CheckValidator(sdk.ValAddress(pks[1].Address()), stakingtypes.Active, false)

	// 600 more blocks happened
	height = 100

	// validator misses 11 blocks
	latest := height
	for ; height < latest+int64(properties.MischanceConfidence)+1; height++ {
		ctx = ctx.WithBlockHeight(height)
		app.CustomSlashingKeeper.HandleValidatorSignature(ctx, addr, 1, false)
	}

	// should now have 1 mischance
	signInfo, found = app.CustomSlashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
	require.True(t, found)
	require.Equal(t, consAddr.String(), signInfo.Address)
	require.Equal(t, int64(0), signInfo.StartHeight)
	require.Equal(t, time.Unix(0, 0).UTC(), signInfo.InactiveUntil.UTC())
	require.Equal(t, false, signInfo.Tombstoned)
	require.Equal(t, int64(1), signInfo.Mischance)
	require.Equal(t, int64(99), signInfo.LastPresentBlock)
	require.Equal(t, int64(11), signInfo.MissedBlocksCounter)
	require.Equal(t, int64(100), signInfo.ProducedBlocksCounter)

	// validator misses 100 blocks
	latest = height
	for ; height < latest+int64(properties.MaxMischance); height++ {
		ctx = ctx.WithBlockHeight(height)
		app.CustomSlashingKeeper.HandleValidatorSignature(ctx, addr, 1, false)
	}

	// should now be inactive & kicked
	staking.EndBlocker(ctx, app.CustomStakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Inactive, true)

	// some blocks pass
	height = int64(5000)
	ctx = ctx.WithBlockHeight(height)

	// Try pausing on inactive node here, should fail
	err := app.CustomStakingKeeper.Pause(ctx, valAddr)
	require.Error(t, err)

	// validator rejoins and starts signing again
	app.CustomStakingKeeper.Activate(ctx, valAddr)
	app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 1, true)
	height++

	// validator should not be kicked since we reset counter/array when it was jailed
	staking.EndBlocker(ctx, app.CustomStakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Active, false)

	// Try pausing on active node here, should success
	err = app.CustomStakingKeeper.Pause(ctx, valAddr)
	require.NoError(t, err)
	staking.EndBlocker(ctx, app.CustomStakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Paused, false)

	// validator misses 501 blocks
	latest = height
	for ; height < latest+501; height++ {
		ctx = ctx.WithBlockHeight(height)
		app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 1, false)
	}

	// validator should not be in inactive status since node is paused
	staking.EndBlocker(ctx, app.CustomStakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Paused, false)

	// Try activating paused node: should unpause but it's activating - should fail
	err = app.CustomStakingKeeper.Activate(ctx, valAddr)
	require.Error(t, err)
	staking.EndBlocker(ctx, app.CustomStakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Paused, false)

	// Unpause node and it should be active
	err = app.CustomStakingKeeper.Unpause(ctx, valAddr)
	require.NoError(t, err)
	staking.EndBlocker(ctx, app.CustomStakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Active, false)

	// After reentering after unpause, check if signature info is recovered correctly
	signInfo, found = app.CustomSlashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
	require.True(t, found)
	require.Equal(t, consAddr.String(), signInfo.Address)
	require.Equal(t, int64(0), signInfo.StartHeight)
	require.Equal(t, ctx.BlockTime().Add(app.CustomSlashingKeeper.DowntimeInactiveDuration(ctx)).String(), signInfo.InactiveUntil.String())
	require.Equal(t, false, signInfo.Tombstoned)
	require.Equal(t, int64(491), signInfo.Mischance)
	require.Equal(t, int64(5000), signInfo.LastPresentBlock)
	require.Equal(t, int64(622), signInfo.MissedBlocksCounter)
	require.Equal(t, int64(101), signInfo.ProducedBlocksCounter)

	// Miss another 120 blocks
	latest = height
	for ; height < latest+120; height++ {
		ctx = ctx.WithBlockHeight(height)
		app.CustomSlashingKeeper.HandleValidatorSignature(ctx, val.Address(), 1, false)
	}

	// validator should be in inactive status
	staking.EndBlocker(ctx, app.CustomStakingKeeper)
	tstaking.CheckValidator(valAddr, stakingtypes.Inactive, true)
}
