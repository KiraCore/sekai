package keeper_test

import (
	"time"

	"github.com/KiraCore/sekai/x/collectives/keeper"
	"github.com/KiraCore/sekai/x/collectives/types"
	multistakingkeeper "github.com/KiraCore/sekai/x/multistaking/keeper"
	multistakingtypes "github.com/KiraCore/sekai/x/multistaking/types"
	spendingtypes "github.com/KiraCore/sekai/x/spending/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (suite *KeeperTestSuite) TestDistributeCollectiveRewards() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	valAddr := sdk.ValAddress(addr1)
	pubkeys := simtestutil.CreateTestPubKeys(1)
	pubKey := pubkeys[0]
	val, err := stakingtypes.NewValidator(valAddr, pubKey)
	suite.Require().NoError(err)

	properties := suite.app.CustomGovKeeper.GetNetworkProperties(suite.ctx)
	properties.MinCollectiveBond = 1
	suite.app.CustomGovKeeper.SetNetworkProperties(suite.ctx, properties)

	val.Status = stakingtypes.Active
	suite.app.CustomStakingKeeper.AddValidator(suite.ctx, val)

	stakingPool := multistakingtypes.StakingPool{
		Id:        1,
		Validator: valAddr.String(),
		Enabled:   true,
	}
	spendingPool := spendingtypes.SpendingPool{
		Name:     "spendingpool1",
		Balances: []sdk.Coin{},
	}

	suite.app.MultiStakingKeeper.SetStakingPool(suite.ctx, stakingPool)
	suite.app.SpendingKeeper.SetSpendingPool(suite.ctx, spendingPool)

	coins := sdk.Coins{sdk.NewInt64Coin("ukex", 1000000)}
	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr1, coins)
	suite.Require().NoError(err)
	msMsgServer := multistakingkeeper.NewMsgServerImpl(suite.app.MultiStakingKeeper, suite.app.BankKeeper, suite.app.CustomGovKeeper, suite.app.CustomStakingKeeper)
	_, err = msMsgServer.Delegate(sdk.WrapSDKContext(suite.ctx), &multistakingtypes.MsgDelegate{
		DelegatorAddress: addr1.String(),
		ValidatorAddress: valAddr.String(),
		Amounts:          coins,
	})
	suite.Require().NoError(err)

	// create collective
	msgServer := keeper.NewMsgServerImpl(suite.app.CollectivesKeeper)
	msg := types.NewMsgCreateCollective(
		addr1, "collective1", "collective1-desc",
		sdk.NewCoins(sdk.NewInt64Coin("v1/ukex", 1000_000)),
		types.DepositWhitelist{Any: true},
		types.OwnersWhitelist{
			Roles:    []uint64{1},
			Accounts: []string{addr1.String()},
		},
		[]types.WeightedSpendingPool{
			{
				Name:   "spendingpool1",
				Weight: sdk.NewDec(1),
			},
		},
		0, 86400, 1000000,
		sdk.NewDecWithPrec(30, 2),
		86400, 3000,
	)

	_, err = msgServer.CreateCollective(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)

	// check collective correctly created
	collective := suite.app.CollectivesKeeper.GetCollective(suite.ctx, msg.Name)
	suite.Require().NotEqual(collective.Name, "")

	// allocate reward to multistaking module
	allocation := sdk.Coins{sdk.NewInt64Coin("ukex", 1000000)}
	pool, found := suite.app.MultiStakingKeeper.GetStakingPoolByValidator(suite.ctx, valAddr.String())
	suite.Require().True(found)

	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, allocation)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, minttypes.ModuleName, authtypes.FeeCollectorName, allocation)
	suite.Require().NoError(err)
	suite.app.MultiStakingKeeper.IncreasePoolRewards(suite.ctx, pool, allocation)

	// distribute collective rewards
	err = suite.app.CollectivesKeeper.DistributeCollectiveRewards(suite.ctx, collective)
	suite.Require().NoError(err)

	// check spending pool balance increased
	updatedSpendingPool := suite.app.SpendingKeeper.GetSpendingPool(suite.ctx, spendingPool.Name)
	suite.Require().NotNil(updatedSpendingPool)
	suite.Require().NotEqual(sdk.Coins(updatedSpendingPool.Balances).String(), "")
}

func (suite *KeeperTestSuite) TestEndBlocker() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	valAddr := sdk.ValAddress(addr1)
	pubkeys := simtestutil.CreateTestPubKeys(1)
	pubKey := pubkeys[0]
	val, err := stakingtypes.NewValidator(valAddr, pubKey)
	suite.Require().NoError(err)

	properties := suite.app.CustomGovKeeper.GetNetworkProperties(suite.ctx)
	properties.MinCollectiveBond = 1
	suite.app.CustomGovKeeper.SetNetworkProperties(suite.ctx, properties)

	val.Status = stakingtypes.Active
	suite.app.CustomStakingKeeper.AddValidator(suite.ctx, val)

	stakingPool := multistakingtypes.StakingPool{
		Id:        1,
		Validator: valAddr.String(),
		Enabled:   true,
	}
	spendingPool := spendingtypes.SpendingPool{
		Name:     "spendingpool1",
		Balances: []sdk.Coin{},
	}

	suite.app.MultiStakingKeeper.SetStakingPool(suite.ctx, stakingPool)
	suite.app.SpendingKeeper.SetSpendingPool(suite.ctx, spendingPool)

	coins := sdk.Coins{sdk.NewInt64Coin("ukex", 1000000)}
	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr1, coins)
	suite.Require().NoError(err)
	msMsgServer := multistakingkeeper.NewMsgServerImpl(suite.app.MultiStakingKeeper, suite.app.BankKeeper, suite.app.CustomGovKeeper, suite.app.CustomStakingKeeper)
	_, err = msMsgServer.Delegate(sdk.WrapSDKContext(suite.ctx), &multistakingtypes.MsgDelegate{
		DelegatorAddress: addr1.String(),
		ValidatorAddress: valAddr.String(),
		Amounts:          coins,
	})
	suite.Require().NoError(err)

	// create collective
	msgServer := keeper.NewMsgServerImpl(suite.app.CollectivesKeeper)
	msg := types.NewMsgCreateCollective(
		addr1, "collective1", "collective1-desc",
		sdk.NewCoins(sdk.NewInt64Coin("v1/ukex", 1000_000)),
		types.DepositWhitelist{Any: true},
		types.OwnersWhitelist{
			Roles:    []uint64{1},
			Accounts: []string{addr1.String()},
		},
		[]types.WeightedSpendingPool{
			{
				Name:   "spendingpool1",
				Weight: sdk.NewDec(1),
			},
		},
		0, 86400, 1000000,
		sdk.NewDecWithPrec(30, 2),
		86400, 3000,
	)

	_, err = msgServer.CreateCollective(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)

	// check collective correctly crteated
	collective := suite.app.CollectivesKeeper.GetCollective(suite.ctx, msg.Name)
	suite.Require().NotEqual(collective.Name, "")

	// allocate reward to multistaking module
	allocation := sdk.Coins{sdk.NewInt64Coin("ukex", 1000000)}
	pool, found := suite.app.MultiStakingKeeper.GetStakingPoolByValidator(suite.ctx, valAddr.String())
	suite.Require().True(found)

	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, allocation)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, minttypes.ModuleName, authtypes.FeeCollectorName, allocation)
	suite.Require().NoError(err)
	suite.app.MultiStakingKeeper.IncreasePoolRewards(suite.ctx, pool, allocation)

	// run endblocker
	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Second * 10000000))
	suite.app.CollectivesKeeper.EndBlocker(suite.ctx)
	suite.Require().NoError(err)

	// check spending pool balance increased
	updatedSpendingPool := suite.app.SpendingKeeper.GetSpendingPool(suite.ctx, spendingPool.Name)
	suite.Require().NotNil(updatedSpendingPool)
	suite.Require().NotEqual(sdk.Coins(updatedSpendingPool.Balances).String(), "")

	// withdraw collective
	withdrawMsg := types.NewMsgWithdrawCollective(
		addr1, "collective1",
	)
	_, err = msgServer.WithdrawCollective(sdk.WrapSDKContext(suite.ctx), withdrawMsg)
	suite.Require().NoError(err)

	// run endblocker again
	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Second * 10000000))
	suite.app.CollectivesKeeper.EndBlocker(suite.ctx)
	suite.Require().NoError(err)

	// check collective correctly removed
	collective = suite.app.CollectivesKeeper.GetCollective(suite.ctx, msg.Name)
	suite.Require().Equal(collective.Name, "")
}
