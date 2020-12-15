package slashing_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/KiraCore/sekai/x/slashing"
	"github.com/KiraCore/sekai/x/slashing/keeper"
	"github.com/KiraCore/sekai/x/slashing/testslashing"
	"github.com/KiraCore/sekai/x/slashing/types"
	"github.com/KiraCore/sekai/x/staking"
	"github.com/KiraCore/sekai/x/staking/teststaking"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestCannotActivateUnlessInactivated(t *testing.T) {
	// initial setup
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	pks := simapp.CreateTestPubKeys(1)
	simapp.AddTestAddrsFromPubKeys(app, ctx, pks, sdk.TokensFromConsensusPower(200))

	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)
	slh := slashing.NewHandler(app.SlashingKeeper)
	addr, val := sdk.ValAddress(pks[0].Address()), pks[0]

	amt := tstaking.CreateValidatorWithValPower(addr, val, 100, true)
	staking.EndBlocker(ctx, app.StakingKeeper)
	require.Equal(
		t, app.BankKeeper.GetAllBalances(ctx, sdk.AccAddress(addr)),
		sdk.Coins{sdk.NewCoin(app.StakingKeeper.GetParams(ctx).BondDenom, InitTokens.Sub(amt))},
	)
	require.Equal(t, amt, app.StakingKeeper.Validator(ctx, addr).GetBondedTokens())

	// assert non-activated validator can't be activated
	res, err := slh(ctx, types.NewMsgActivate(addr))
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, errors.Is(types.ErrValidatorNotJailed, err))
}

func TestInvalidMsg(t *testing.T) {
	k := keeper.Keeper{}
	h := slashing.NewHandler(k)

	res, err := h(sdk.NewContext(nil, tmproto.Header{}, false, nil), testdata.NewTestMsg())
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "unrecognized slashing message type"))
}

// Test a validator through uptime, downtime, revocation,
// unrevocation, starting height reset, and revocation again
func TestHandleAbsentValidator(t *testing.T) {
	// initial setup
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Unix(0, 0)})
	pks := simapp.CreateTestPubKeys(1)
	simapp.AddTestAddrsFromPubKeys(app, ctx, pks, sdk.TokensFromConsensusPower(200))
	app.SlashingKeeper.SetParams(ctx, testslashing.TestParams())

	power := int64(100)
	addr, val := sdk.ValAddress(pks[0].Address()), pks[0]
	slh := slashing.NewHandler(app.SlashingKeeper)
	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)

	amt := tstaking.CreateValidatorWithValPower(addr, val, power, true)
	staking.EndBlocker(ctx, app.StakingKeeper)

	require.Equal(
		t, app.BankKeeper.GetAllBalances(ctx, sdk.AccAddress(addr)),
		sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.GetParams(ctx).BondDenom, InitTokens.Sub(amt))),
	)
	require.Equal(t, amt, app.StakingKeeper.Validator(ctx, addr).GetBondedTokens())

	// will exist since the validator has been bonded
	info, found := app.SlashingKeeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(val.Address()))
	require.True(t, found)
	require.Equal(t, int64(0), info.StartHeight)
	require.Equal(t, int64(0), info.IndexOffset)
	require.Equal(t, int64(0), info.MissedBlocksCounter)
	require.Equal(t, time.Unix(0, 0).UTC(), info.JailedUntil)
	height := int64(0)

	// 1000 first blocks OK
	for ; height < app.SlashingKeeper.SignedBlocksWindow(ctx); height++ {
		ctx = ctx.WithBlockHeight(height)
		app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), power, true)
	}
	info, found = app.SlashingKeeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(val.Address()))
	require.True(t, found)
	require.Equal(t, int64(0), info.StartHeight)
	require.Equal(t, int64(0), info.MissedBlocksCounter)

	// 500 blocks missed
	for ; height < app.SlashingKeeper.SignedBlocksWindow(ctx)+(app.SlashingKeeper.SignedBlocksWindow(ctx)-app.SlashingKeeper.MinSignedPerWindow(ctx)); height++ {
		ctx = ctx.WithBlockHeight(height)
		app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), power, false)
	}
	info, found = app.SlashingKeeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(val.Address()))
	require.True(t, found)
	require.Equal(t, int64(0), info.StartHeight)
	require.Equal(t, app.SlashingKeeper.SignedBlocksWindow(ctx)-app.SlashingKeeper.MinSignedPerWindow(ctx), info.MissedBlocksCounter)

	// validator should be bonded still
	validator, _ := app.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val))
	require.Equal(t, stakingtypes.Bonded, validator.GetStatus())

	bondPool := app.StakingKeeper.GetBondedPool(ctx)
	require.True(sdk.IntEq(t, amt, app.BankKeeper.GetBalance(ctx, bondPool.GetAddress(), app.StakingKeeper.BondDenom(ctx)).Amount))

	// 501st block missed
	ctx = ctx.WithBlockHeight(height)
	app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), power, false)
	info, found = app.SlashingKeeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(val.Address()))
	require.True(t, found)
	require.Equal(t, int64(0), info.StartHeight)
	// counter now reset to zero
	require.Equal(t, int64(0), info.MissedBlocksCounter)

	// end block
	staking.EndBlocker(ctx, app.StakingKeeper)

	// validator should have been jailed
	validator, _ = app.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val))
	require.Equal(t, stakingtypes.Unbonding, validator.GetStatus())

	slashAmt := amt.ToDec().Mul(app.SlashingKeeper.SlashFractionDowntime(ctx)).RoundInt64()

	// validator should have been slashed
	require.Equal(t, amt.Int64()-slashAmt, validator.GetTokens().Int64())

	// 502nd block *also* missed (since the LastCommit would have still included the just-unbonded validator)
	height++
	ctx = ctx.WithBlockHeight(height)
	app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), power, false)
	info, found = app.SlashingKeeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(val.Address()))
	require.True(t, found)
	require.Equal(t, int64(0), info.StartHeight)
	require.Equal(t, int64(1), info.MissedBlocksCounter)

	// end block
	staking.EndBlocker(ctx, app.StakingKeeper)

	// validator should not have been slashed any more, since it was already jailed
	validator, _ = app.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val))
	require.Equal(t, amt.Int64()-slashAmt, validator.GetTokens().Int64())

	// unrevocation should fail prior to jail expiration
	res, err := slh(ctx, types.NewMsgActivate(addr))
	require.Error(t, err)
	require.Nil(t, res)

	// unrevocation should succeed after jail expiration
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: time.Unix(1, 0).Add(app.SlashingKeeper.DowntimeJailDuration(ctx))})
	res, err = slh(ctx, types.NewMsgActivate(addr))
	require.NoError(t, err)
	require.NotNil(t, res)

	// end block
	staking.EndBlocker(ctx, app.StakingKeeper)

	// validator should be rebonded now
	validator, _ = app.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val))
	require.Equal(t, stakingtypes.Bonded, validator.GetStatus())

	// validator should have been slashed
	require.Equal(t, amt.Int64()-slashAmt, app.BankKeeper.GetBalance(ctx, bondPool.GetAddress(), app.StakingKeeper.BondDenom(ctx)).Amount.Int64())

	// Validator start height should not have been changed
	info, found = app.SlashingKeeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(val.Address()))
	require.True(t, found)
	require.Equal(t, int64(0), info.StartHeight)
	// we've missed 2 blocks more than the maximum, so the counter was reset to 0 at 1 block more and is now 1
	require.Equal(t, int64(1), info.MissedBlocksCounter)

	// validator should not be immediately jailed again
	height++
	ctx = ctx.WithBlockHeight(height)
	app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), power, false)
	validator, _ = app.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val))
	require.Equal(t, stakingtypes.Bonded, validator.GetStatus())

	// 500 signed blocks
	nextHeight := height + app.SlashingKeeper.MinSignedPerWindow(ctx) + 1
	for ; height < nextHeight; height++ {
		ctx = ctx.WithBlockHeight(height)
		app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), power, false)
	}

	// end block
	staking.EndBlocker(ctx, app.StakingKeeper)

	// validator should be jailed again after 500 unsigned blocks
	nextHeight = height + app.SlashingKeeper.MinSignedPerWindow(ctx) + 1
	for ; height <= nextHeight; height++ {
		ctx = ctx.WithBlockHeight(height)
		app.SlashingKeeper.HandleValidatorSignature(ctx, val.Address(), power, false)
	}

	// end block
	staking.EndBlocker(ctx, app.StakingKeeper)

	validator, _ = app.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val))
	require.Equal(t, stakingtypes.Unbonding, validator.GetStatus())
}
