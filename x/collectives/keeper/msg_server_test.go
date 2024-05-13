package keeper_test

import (
	"github.com/KiraCore/sekai/x/collectives/keeper"
	"github.com/KiraCore/sekai/x/collectives/types"
	multistakingkeeper "github.com/KiraCore/sekai/x/multistaking/keeper"
	multistakingtypes "github.com/KiraCore/sekai/x/multistaking/types"
	spendingtypes "github.com/KiraCore/sekai/x/spending/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (suite *KeeperTestSuite) TestCreateCollective() {
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
		sdk.NewDecWithPrec(30, 2), 86400, 3000,
	)

	_, err = msgServer.CreateCollective(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)

	// check collective correctly crteated
	collective := suite.app.CollectivesKeeper.GetCollective(suite.ctx, msg.Name)
	suite.Require().NotEqual(collective.Name, "")
}

func (suite *KeeperTestSuite) TestContributeCollective() {
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
		sdk.NewCoins(sdk.NewInt64Coin("v1/ukex", 100_000)),
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
		sdk.NewDecWithPrec(30, 2), 86400, 3000,
	)

	_, err = msgServer.CreateCollective(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)

	// check collective correctly crteated
	collective := suite.app.CollectivesKeeper.GetCollective(suite.ctx, msg.Name)
	suite.Require().NotEqual(collective.Name, "")

	// contribute collective
	contributeMsg := types.NewMsgBondCollective(
		addr1, "collective1", sdk.NewCoins(sdk.NewInt64Coin("v1/ukex", 900_000)),
	)
	_, err = msgServer.ContributeCollective(sdk.WrapSDKContext(suite.ctx), contributeMsg)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) TestDonateCollective() {
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
		sdk.NewDecWithPrec(30, 2), 86400, 3000,
	)

	_, err = msgServer.CreateCollective(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)

	// check collective correctly crteated
	collective := suite.app.CollectivesKeeper.GetCollective(suite.ctx, msg.Name)
	suite.Require().NotEqual(collective.Name, "")

	// donate collective
	donateMsg := types.NewMsgDonateCollective(
		addr1, "collective1", 1000, sdk.NewDecWithPrec(1, 1), false,
	)
	_, err = msgServer.DonateCollective(sdk.WrapSDKContext(suite.ctx), donateMsg)
	suite.Require().NoError(err)

	// donate collective again with lower period
	donateMsg = types.NewMsgDonateCollective(
		addr1, "collective1", 100, sdk.NewDecWithPrec(1, 1), false,
	)
	_, err = msgServer.DonateCollective(sdk.WrapSDKContext(suite.ctx), donateMsg)
	suite.Require().Error(err)

	// donate collective again with higher donation and lock
	donateMsg = types.NewMsgDonateCollective(
		addr1, "collective1", 1000, sdk.NewDecWithPrec(2, 1), true,
	)
	_, err = msgServer.DonateCollective(sdk.WrapSDKContext(suite.ctx), donateMsg)
	suite.Require().NoError(err)

	// try to change donation period shorter
	donateMsg = types.NewMsgDonateCollective(
		addr1, "collective1", 100, sdk.NewDecWithPrec(1, 1), false,
	)
	_, err = msgServer.DonateCollective(sdk.WrapSDKContext(suite.ctx), donateMsg)
	suite.Require().Error(err)
}

func (suite *KeeperTestSuite) TestWithdrawCollective() {
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
		sdk.NewDecWithPrec(30, 2), 86400, 3000,
	)

	_, err = msgServer.CreateCollective(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)

	// check collective correctly crteated
	collective := suite.app.CollectivesKeeper.GetCollective(suite.ctx, msg.Name)
	suite.Require().NotEqual(collective.Name, "")

	// withdraw collective
	withdrawMsg := types.NewMsgWithdrawCollective(
		addr1, "collective1",
	)
	_, err = msgServer.WithdrawCollective(sdk.WrapSDKContext(suite.ctx), withdrawMsg)
	suite.Require().NoError(err)
}
