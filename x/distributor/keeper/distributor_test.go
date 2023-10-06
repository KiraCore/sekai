package keeper_test

import (
	"time"

	"github.com/KiraCore/sekai/x/distributor/types"
	recoverytypes "github.com/KiraCore/sekai/x/recovery/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (suite *KeeperTestSuite) TestAllocateTokensToRecoveryTokenValidator() {
	suite.SetupTest()
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	valAddr := sdk.ValAddress(addr1)
	pubkeys := simtestutil.CreateTestPubKeys(1)
	pubKey := pubkeys[0]
	consAddr := sdk.ConsAddress(pubKey.Address())
	val, err := stakingtypes.NewValidator(valAddr, pubKey)
	suite.Require().NoError(err)

	coins := sdk.Coins{sdk.NewInt64Coin("ukex", 10000000)}

	suite.app.CustomStakingKeeper.AddValidator(suite.ctx, val)
	suite.app.RecoveryKeeper.SetRecoveryToken(suite.ctx, recoverytypes.RecoveryToken{
		Address:          addr1.String(),
		Token:            "rr/validator1",
		RrSupply:         sdk.NewInt(1000_000),
		UnderlyingTokens: coins,
	})

	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr1, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr2, coins)
	suite.Require().NoError(err)

	suite.app.DistrKeeper.AllocateTokens(suite.ctx, 10, 10, consAddr, []abci.VoteInfo{})
	balance := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1)
	suite.Require().Equal(balance, coins)

	// recoveryToken, err := suite.app.RecoveryKeeper.GetRecoveryToken(suite.ctx, addr1.String())
	// suite.Require().NoError(err)
	// suite.Require().NotEqual(coins.String(), sdk.Coins(recoveryToken.UnderlyingTokens).String())
}

func (suite *KeeperTestSuite) TestAllocateTokens() {
	suite.SetupTest()
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	consAddr := sdk.ConsAddress(addr1)
	coins := sdk.Coins{sdk.NewInt64Coin("ukex", 10000000)}
	suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr1, coins)
	suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr2, coins)

	supply := suite.app.BankKeeper.GetSupply(suite.ctx, "ukex")
	suite.Require().Equal(supply.Amount, sdk.NewInt(20000000))

	now := time.Now()
	suite.ctx = suite.ctx.WithBlockTime(now)
	suite.app.DistrKeeper.SetPeriodicSnapshot(suite.ctx, types.SupplySnapshot{
		SnapshotAmount: supply.Amount,
		SnapshotTime:   now.Unix(),
	})
	future := now.Add(time.Hour * 24 * 365)
	suite.ctx = suite.ctx.WithBlockTime(future)
	oldTreasury := suite.app.DistrKeeper.GetFeesTreasury(suite.ctx)
	suite.app.DistrKeeper.AllocateTokens(suite.ctx, 10, 10, consAddr, []abci.VoteInfo{})
	newTreasury := suite.app.DistrKeeper.GetFeesTreasury(suite.ctx)
	suite.Require().True(oldTreasury.DenomsSubsetOf(newTreasury))

	supply = suite.app.BankKeeper.GetSupply(suite.ctx, "ukex")
	suite.Require().Equal(supply.Amount, sdk.NewInt(23597535))

	// TODO: add case for validator exit case
	// TODO: add case for staking pool exist case
	// TODO: check tokens are distributed correctly
}
