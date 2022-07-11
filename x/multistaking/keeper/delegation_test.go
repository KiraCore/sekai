package keeper_test

import (
	"github.com/KiraCore/sekai/x/multistaking/keeper"
	"github.com/KiraCore/sekai/x/multistaking/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func (suite *KeeperTestSuite) TestLastUndelegationIdGetSet() {
	// get default last delegation id
	lastDelegationId := suite.app.MultiStakingKeeper.GetLastUndelegationId(suite.ctx)
	suite.Require().Equal(lastDelegationId, uint64(0))

	// set last delegation id to new value
	newDelegationId := uint64(2)
	suite.app.MultiStakingKeeper.SetLastUndelegationId(suite.ctx, newDelegationId)

	// check last delegation id update
	lastDelegationId = suite.app.MultiStakingKeeper.GetLastUndelegationId(suite.ctx)
	suite.Require().Equal(lastDelegationId, newDelegationId)
}

func (suite *KeeperTestSuite) TestUndelegationGetSet() {
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	// get undelegation by id
	_, found := suite.app.MultiStakingKeeper.GetUndelegationById(suite.ctx, 1)
	suite.Require().False(found)

	// check whole undelegations
	allUndelegations := suite.app.MultiStakingKeeper.GetAllUndelegations(suite.ctx)
	suite.Require().Len(allUndelegations, 0)

	undelegations := []types.Undelegation{
		{
			Id:      1,
			Address: addr.String(),
			Expiry:  uint64(suite.ctx.BlockTime().Unix() + 1000),
			Amount:  sdk.Coins{sdk.NewInt64Coin("ukex", 10000)},
		},
		{
			Id:      2,
			Address: addr2.String(),
			Expiry:  uint64(suite.ctx.BlockTime().Unix() + 1000),
			Amount:  sdk.Coins{sdk.NewInt64Coin("ukex", 10000)},
		},
	}

	for _, undelegation := range undelegations {
		suite.app.MultiStakingKeeper.SetUndelegation(suite.ctx, undelegation)
	}

	// check undelegation by id
	for _, undelegation := range undelegations {
		p, found := suite.app.MultiStakingKeeper.GetUndelegationById(suite.ctx, undelegation.Id)
		suite.Require().True(found)
		suite.Require().Equal(undelegation, p)
	}

	// check undelegations for whole export
	allUndelegations = suite.app.MultiStakingKeeper.GetAllUndelegations(suite.ctx)
	suite.Require().Len(allUndelegations, 2)
}

func (suite *KeeperTestSuite) TestPoolDelegatorGetSet() {
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	// get pool delegators at the beginning
	delegators := suite.app.MultiStakingKeeper.GetPoolDelegators(suite.ctx, 1)
	suite.Require().Len(delegators, 0)

	suite.app.MultiStakingKeeper.SetPoolDelegator(suite.ctx, 1, addr)
	suite.app.MultiStakingKeeper.SetPoolDelegator(suite.ctx, 1, addr2)

	// get pool delegators after setting up pool delegators
	delegators = suite.app.MultiStakingKeeper.GetPoolDelegators(suite.ctx, 1)
	suite.Require().Len(delegators, 2)

	suite.app.MultiStakingKeeper.RemovePoolDelegator(suite.ctx, 1, addr2)

	// get pool delegators after removing a pool delegator
	delegators = suite.app.MultiStakingKeeper.GetPoolDelegators(suite.ctx, 1)
	suite.Require().Len(delegators, 1)
}

func (suite *KeeperTestSuite) TestDelegatorRewardsGetSet() {
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	// get pool delegators at the beginning
	rewards := suite.app.MultiStakingKeeper.GetDelegatorRewards(suite.ctx, addr)
	suite.Require().Equal(rewards, sdk.Coins{})

	allocation1 := sdk.Coins{sdk.NewInt64Coin("ukex", 1000000)}
	suite.app.MultiStakingKeeper.IncreaseDelegatorRewards(suite.ctx, addr, allocation1)

	// get pool delegators after setting up pool delegators
	rewards = suite.app.MultiStakingKeeper.GetDelegatorRewards(suite.ctx, addr)
	suite.Require().Equal(rewards, allocation1)

	allocation2 := sdk.Coins{sdk.NewInt64Coin("mkex", 1000000)}
	suite.app.MultiStakingKeeper.IncreaseDelegatorRewards(suite.ctx, addr, allocation2)

	rewards = suite.app.MultiStakingKeeper.GetDelegatorRewards(suite.ctx, addr)
	suite.Require().Equal(rewards, allocation1.Add(allocation2...))

	allRewards := suite.app.MultiStakingKeeper.GetAllDelegatorRewards(suite.ctx)
	suite.Require().Len(allRewards, 1)
	suite.Require().Equal(allRewards[0], types.Rewards{
		Delegator: addr.String(),
		Rewards:   rewards,
	})

	suite.app.MultiStakingKeeper.RemoveDelegatorRewards(suite.ctx, addr)

	// get pool delegators after removing
	rewards = suite.app.MultiStakingKeeper.GetDelegatorRewards(suite.ctx, addr)
	suite.Require().Equal(rewards, sdk.Coins{})
}

func (suite *KeeperTestSuite) TestIncreasePoolRewards() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	valAddr := sdk.ValAddress(addr1)
	coins := sdk.Coins{sdk.NewInt64Coin("ukex", 10000000)}
	suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr1, coins)
	suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr2, coins)

	pubkeys := simapp.CreateTestPubKeys(1)
	pubKey := pubkeys[0]

	val, err := stakingtypes.NewValidator(valAddr, pubKey)
	suite.Require().NoError(err)

	val.Status = stakingtypes.Active
	suite.app.CustomStakingKeeper.AddValidator(suite.ctx, val)

	pool := types.StakingPool{
		Id:        1,
		Validator: valAddr.String(),
		Enabled:   true,
	}
	suite.app.MultiStakingKeeper.SetStakingPool(suite.ctx, pool)
	msgServer := keeper.NewMsgServerImpl(suite.app.MultiStakingKeeper, suite.app.BankKeeper, suite.app.CustomGovKeeper, suite.app.CustomStakingKeeper)
	_, err = msgServer.Delegate(sdk.WrapSDKContext(suite.ctx), &types.MsgDelegate{
		DelegatorAddress: addr1.String(),
		ValidatorAddress: valAddr.String(),
		Amounts:          sdk.Coins{sdk.NewInt64Coin("ukex", 1000000)},
	})
	suite.Require().NoError(err)
	_, err = msgServer.Delegate(sdk.WrapSDKContext(suite.ctx), &types.MsgDelegate{
		DelegatorAddress: addr2.String(),
		ValidatorAddress: valAddr.String(),
		Amounts:          sdk.Coins{sdk.NewInt64Coin("ukex", 1000000)},
	})
	suite.Require().NoError(err)

	allocation := sdk.Coins{sdk.NewInt64Coin("ukex", 1000000)}
	pool, found := suite.app.MultiStakingKeeper.GetStakingPoolByValidator(suite.ctx, valAddr.String())
	suite.Require().True(found)
	suite.app.MultiStakingKeeper.IncreasePoolRewards(suite.ctx, pool, allocation)

	rewards := suite.app.MultiStakingKeeper.GetDelegatorRewards(suite.ctx, addr1)
	suite.Require().Equal(rewards, sdk.Coins{sdk.NewInt64Coin("ukex", 250000)})
	rewards = suite.app.MultiStakingKeeper.GetDelegatorRewards(suite.ctx, addr2)
	suite.Require().Equal(rewards, sdk.Coins{sdk.NewInt64Coin("ukex", 250000)})
}

func (suite *KeeperTestSuite) TestDelegate() {
	testCases := map[string]struct {
		delegateToken   string
		valStatus       stakingtypes.ValidatorStatus
		poolCreate      bool
		maxDelegators   uint64
		preDelegations  int
		mintCoins       sdk.Int
		delegationCoins sdk.Int
		expectErr       bool
	}{
		"inactive validator delegate": {
			"ukex",
			stakingtypes.Paused,
			true,
			3,
			0,
			sdk.NewInt(1000000),
			sdk.NewInt(1000000),
			true,
		},
		"not existing pool delegate": {
			"ukex",
			stakingtypes.Active,
			false,
			3,
			0,
			sdk.NewInt(1000000),
			sdk.NewInt(1000000),
			true,
		},
		"max delegators exceed": {
			"ukex",
			stakingtypes.Active,
			true,
			3,
			3,
			sdk.NewInt(1000000),
			sdk.NewInt(1000000),
			true,
		},
		"not enough amounts on delegator": {
			"ukex",
			stakingtypes.Active,
			true,
			3,
			0,
			sdk.NewInt(1),
			sdk.NewInt(1000000),
			true,
		},
		"not registered token delegate": {
			"ukexxx",
			stakingtypes.Active,
			true,
			3,
			0,
			sdk.NewInt(1000000),
			sdk.NewInt(1000000),
			true,
		},
		"successful delegate": {
			"ukex",
			stakingtypes.Active,
			true,
			3,
			0,
			sdk.NewInt(1000000),
			sdk.NewInt(1000000),
			false,
		},
		"successful delegate with pushout": {
			"ukex",
			stakingtypes.Active,
			true,
			3,
			3,
			sdk.NewInt(100000000),
			sdk.NewInt(100000000),
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		suite.Run(name, func() {
			suite.SetupTest()
			addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
			valAddr := sdk.ValAddress(addr1)
			pubkeys := simapp.CreateTestPubKeys(1)
			pubKey := pubkeys[0]

			val, err := stakingtypes.NewValidator(valAddr, pubKey)
			suite.Require().NoError(err)

			val.Status = tc.valStatus
			suite.app.CustomStakingKeeper.AddValidator(suite.ctx, val)

			if tc.poolCreate {
				pool := types.StakingPool{
					Id:        1,
					Validator: valAddr.String(),
					Enabled:   true,
				}
				suite.app.MultiStakingKeeper.SetStakingPool(suite.ctx, pool)
			}

			params := suite.app.CustomGovKeeper.GetNetworkProperties(suite.ctx)
			params.MaxDelegators = tc.maxDelegators
			suite.app.CustomGovKeeper.SetNetworkProperties(suite.ctx, params)

			coins := sdk.Coins{sdk.NewCoin("ukex", tc.mintCoins)}
			suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
			suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr1, coins)

			for i := 0; i < tc.preDelegations; i++ {
				coins := sdk.Coins{sdk.NewInt64Coin("ukex", 1000000)}
				raddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
				suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
				suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, raddr, coins)
				err = suite.app.MultiStakingKeeper.Delegate(suite.ctx, &types.MsgDelegate{
					DelegatorAddress: raddr.String(),
					ValidatorAddress: valAddr.String(),
					Amounts:          sdk.Coins{sdk.NewInt64Coin("ukex", 1000000)},
				})
				suite.Require().NoError(err)
			}

			err = suite.app.MultiStakingKeeper.Delegate(suite.ctx, &types.MsgDelegate{
				DelegatorAddress: addr1.String(),
				ValidatorAddress: valAddr.String(),
				Amounts:          sdk.Coins{sdk.NewCoin(tc.delegateToken, tc.delegationCoins)},
			})
			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)

				// check share denoms are correctly minted
				balance := suite.app.BankKeeper.GetBalance(suite.ctx, addr1, "v1/ukex")
				suite.Require().True(balance.Amount.IsPositive())

				// check delegator is set as delegator
				isDelegator := suite.app.MultiStakingKeeper.IsPoolDelegator(suite.ctx, 1, addr1)
				suite.Require().True(isDelegator)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestRegisterDelegator() {
	// delegate
	suite.SetupTest()
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	valAddr := sdk.ValAddress(addr1)
	pubkeys := simapp.CreateTestPubKeys(1)
	pubKey := pubkeys[0]

	val, err := stakingtypes.NewValidator(valAddr, pubKey)
	suite.Require().NoError(err)

	val.Status = stakingtypes.Active
	suite.app.CustomStakingKeeper.AddValidator(suite.ctx, val)

	pool := types.StakingPool{
		Id:        1,
		Validator: valAddr.String(),
		Enabled:   true,
	}
	suite.app.MultiStakingKeeper.SetStakingPool(suite.ctx, pool)

	coins := sdk.Coins{sdk.NewInt64Coin("ukex", 1000000)}
	suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr1, coins)

	err = suite.app.MultiStakingKeeper.Delegate(suite.ctx, &types.MsgDelegate{
		DelegatorAddress: addr1.String(),
		ValidatorAddress: valAddr.String(),
		Amounts:          sdk.Coins{sdk.NewInt64Coin("ukex", 1000000)},
	})
	suite.Require().NoError(err)

	// check share denoms are correctly minted
	balance := suite.app.BankKeeper.GetBalance(suite.ctx, addr1, "v1/ukex")
	suite.Require().True(balance.Amount.IsPositive())

	// send minted token to new address
	err = suite.app.BankKeeper.SendCoins(suite.ctx, addr1, addr2, sdk.Coins{balance})
	suite.Require().NoError(err)

	// register delegator
	suite.app.MultiStakingKeeper.RegisterDelegator(suite.ctx, addr2)

	// check if registered as delegator
	isDelegator := suite.app.MultiStakingKeeper.IsPoolDelegator(suite.ctx, 1, addr2)
	suite.Require().True(isDelegator)
}
