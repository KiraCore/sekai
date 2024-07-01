package keeper_test

import (
	"github.com/KiraCore/sekai/x/collectives/types"
	spendingtypes "github.com/KiraCore/sekai/x/spending/types"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (suite *KeeperTestSuite) TestCollectiveSetGetDelete() {
	collectives := []types.Collective{
		{
			Name:             "collective1",
			Description:      "collective1 description",
			Status:           types.CollectiveActive,
			DepositWhitelist: types.DepositWhitelist{Any: true},
			OwnersWhitelist: types.OwnersWhitelist{
				Roles:    []uint64{1},
				Accounts: []string{"kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r"},
			},
			SpendingPools: []types.WeightedSpendingPool{
				{
					Name:   "spendingpool1",
					Weight: sdk.NewDec(1),
				},
			},
			ClaimStart:       0,
			ClaimPeriod:      1000,
			ClaimEnd:         100000,
			VoteQuorum:       sdk.NewDecWithPrec(30, 2),
			VotePeriod:       86400,
			VoteEnactment:    3000,
			Donations:        []sdk.Coin{sdk.NewInt64Coin("ukex", 1000_000)},
			Rewards:          []sdk.Coin{sdk.NewInt64Coin("ukex", 1000_000)},
			LastDistribution: 0,
			Bonds:            []sdk.Coin{sdk.NewInt64Coin("ukex", 1000_000)},
			CreationTime:     0,
		},
		{
			Name:             "collective2",
			Description:      "collective2 description",
			Status:           types.CollectiveInactive,
			DepositWhitelist: types.DepositWhitelist{Any: true},
			OwnersWhitelist: types.OwnersWhitelist{
				Roles:    []uint64{1},
				Accounts: []string{"kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r"},
			},
			SpendingPools: []types.WeightedSpendingPool{
				{
					Name:   "spendingpool1",
					Weight: sdk.NewDec(1),
				},
			},
			ClaimStart:       0,
			ClaimPeriod:      1000,
			ClaimEnd:         100000,
			VoteQuorum:       sdk.NewDecWithPrec(30, 2),
			VotePeriod:       86400,
			VoteEnactment:    3000,
			Donations:        []sdk.Coin{sdk.NewInt64Coin("ukex", 1000_000)},
			Rewards:          []sdk.Coin{sdk.NewInt64Coin("ukex", 1000_000)},
			LastDistribution: 0,
			Bonds:            []sdk.Coin{sdk.NewInt64Coin("ukex", 1000_000)},
			CreationTime:     0,
		},
	}

	for _, collective := range collectives {
		suite.app.CollectivesKeeper.SetCollective(suite.ctx, collective)
	}

	for _, collective := range collectives {
		c := suite.app.CollectivesKeeper.GetCollective(suite.ctx, collective.Name)
		suite.Require().Equal(c, collective)
	}

	allCollectives := suite.app.CollectivesKeeper.GetAllCollectives(suite.ctx)
	suite.Require().Len(allCollectives, 2)

	suite.app.CollectivesKeeper.DeleteCollective(suite.ctx, collectives[0].Name)

	allCollectives = suite.app.CollectivesKeeper.GetAllCollectives(suite.ctx)
	suite.Require().Len(allCollectives, 1)

	collective := suite.app.CollectivesKeeper.GetCollective(suite.ctx, collectives[0].Name)
	suite.Require().Equal(collective.Name, "")
}

func (suite *KeeperTestSuite) TestCollectiveContributerSetGetDelete() {
	contributers := []types.CollectiveContributor{
		{
			Address:      "kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r",
			Name:         "collective1",
			Bonds:        []sdk.Coin{sdk.NewInt64Coin("ukex", 1000_000)},
			Locking:      1000,
			Donation:     sdk.NewDecWithPrec(1, 1), // 10%
			DonationLock: true,
		},
		{
			Address:      "kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r",
			Name:         "collective2",
			Bonds:        []sdk.Coin{sdk.NewInt64Coin("ukex", 1000_000)},
			Locking:      1000,
			Donation:     sdk.NewDecWithPrec(1, 1), // 10%
			DonationLock: false,
		},
	}

	for _, contributer := range contributers {
		suite.app.CollectivesKeeper.SetCollectiveContributer(suite.ctx, contributer)
	}

	for _, contributer := range contributers {
		c := suite.app.CollectivesKeeper.GetCollectiveContributer(suite.ctx, contributer.Name, contributer.Address)
		suite.Require().Equal(c, contributer)
	}

	allContributers := suite.app.CollectivesKeeper.GetCollectiveContributers(suite.ctx, contributers[0].Name)
	suite.Require().Len(allContributers, 1)

	suite.app.CollectivesKeeper.DeleteCollectiveContributer(suite.ctx, contributers[0].Name, contributers[0].Address)

	allContributers = suite.app.CollectivesKeeper.GetCollectiveContributers(suite.ctx, contributers[0].Name)
	suite.Require().Len(allContributers, 0)

	contributer := suite.app.CollectivesKeeper.GetCollectiveContributer(suite.ctx, contributers[0].Name, contributers[0].Address)
	suite.Require().Equal(contributer.Name, "")
}

func (suite *KeeperTestSuite) TestSendDonation() {
	pubKey := newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50")
	addr := sdk.AccAddress(pubKey.Address())
	//  case collective does not exists
	coins := sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000_000))
	err := suite.app.CollectivesKeeper.SendDonation(suite.ctx, "collective1", addr, coins)
	suite.Require().Error(err)

	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, minttypes.ModuleName, types.ModuleName, coins)
	suite.Require().NoError(err)

	// case donations amount lower than withdraw
	collective := types.Collective{
		Name:             "collective1",
		Description:      "collective1 description",
		Status:           types.CollectiveActive,
		DepositWhitelist: types.DepositWhitelist{Any: true},
		OwnersWhitelist: types.OwnersWhitelist{
			Roles:    []uint64{1},
			Accounts: []string{"kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r"},
		},
		SpendingPools: []types.WeightedSpendingPool{
			{
				Name:   "spendingpool1",
				Weight: sdk.NewDec(1),
			},
		},
		ClaimStart:       0,
		ClaimPeriod:      1000,
		ClaimEnd:         100000,
		VoteQuorum:       sdk.NewDecWithPrec(30, 2),
		VotePeriod:       86400,
		VoteEnactment:    3000,
		Donations:        []sdk.Coin{sdk.NewInt64Coin("ukex", 100_000)},
		Rewards:          []sdk.Coin{sdk.NewInt64Coin("ukex", 1000_000)},
		LastDistribution: 0,
		Bonds:            []sdk.Coin{sdk.NewInt64Coin("ukex", 1000_000)},
		CreationTime:     0,
	}
	suite.app.CollectivesKeeper.SetCollective(suite.ctx, collective)
	err = suite.app.CollectivesKeeper.SendDonation(suite.ctx, "collective1", addr, coins)
	suite.Require().Error(err)

	collective.Donations = coins
	suite.app.CollectivesKeeper.SetCollective(suite.ctx, collective)
	err = suite.app.CollectivesKeeper.SendDonation(suite.ctx, "collective1", addr, coins)
	suite.Require().NoError(err)

	collective = suite.app.CollectivesKeeper.GetCollective(suite.ctx, collective.Name)

	// check balance correctly increased on withdrawal address
	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr)
	suite.Require().Equal(balances, coins)
	// check collective donations amount decreased by withdrawal amount
	suite.Require().Equal(collective.Donations, []sdk.Coin(nil))
}

func (suite *KeeperTestSuite) GetBondsValue() {
	// get bonds value for ukex
	coins := sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000_000))
	value := suite.app.CollectivesKeeper.GetBondsValue(suite.ctx, coins)
	suite.Require().Equal(value, sdk.NewDec(1000_000))

	// get bonds value including not registered token
	coins = sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000_000), sdk.NewInt64Coin("zzz", 1000))
	value = suite.app.CollectivesKeeper.GetBondsValue(suite.ctx, coins)
	suite.Require().Equal(value, sdk.NewDec(1000_000))

	// get bonds value with newly registered token
	suite.app.TokensKeeper.UpsertTokenInfo(suite.ctx, tokenstypes.TokenInfo{
		Denom:   "zzz",
		FeeRate: sdk.NewDec(10),
	})
	coins = sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000_000), sdk.NewInt64Coin("zzz", 1000))
	value = suite.app.CollectivesKeeper.GetBondsValue(suite.ctx, coins)
	suite.Require().Equal(value, sdk.NewDec(1010_000))
}

func (suite *KeeperTestSuite) WithdrawCollective() {
	coins := []sdk.Coin{sdk.NewInt64Coin("ukex", 1000_000)}
	collective := types.Collective{
		Name:             "collective1",
		Description:      "collective1 description",
		Status:           types.CollectiveActive,
		DepositWhitelist: types.DepositWhitelist{Any: true},
		OwnersWhitelist: types.OwnersWhitelist{
			Roles:    []uint64{1},
			Accounts: []string{"kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r"},
		},
		SpendingPools: []types.WeightedSpendingPool{
			{
				Name:   "spendingpool1",
				Weight: sdk.NewDec(1),
			},
		},
		ClaimStart:       0,
		ClaimPeriod:      1000,
		ClaimEnd:         100000,
		VoteQuorum:       sdk.NewDecWithPrec(30, 2),
		VotePeriod:       86400,
		VoteEnactment:    3000,
		Donations:        []sdk.Coin{sdk.NewInt64Coin("ukex", 100_000)},
		Rewards:          []sdk.Coin{sdk.NewInt64Coin("ukex", 1000_000)},
		LastDistribution: 0,
		Bonds:            []sdk.Coin{sdk.NewInt64Coin("ukex", 2000_000)},
		CreationTime:     0,
	}

	contributer := types.CollectiveContributor{
		Address:      "kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r",
		Name:         "collective1",
		Bonds:        coins,
		Locking:      1000,
		Donation:     sdk.NewDecWithPrec(5, 1), // 50%
		DonationLock: true,
	}

	suite.app.CollectivesKeeper.SetCollective(suite.ctx, collective)
	suite.app.CollectivesKeeper.SetCollectiveContributer(suite.ctx, contributer)

	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, collective.GetCollectiveAddress(), coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, collective.GetCollectiveDonationAddress(), coins)
	suite.Require().NoError(err)

	err = suite.app.CollectivesKeeper.WithdrawCollective(suite.ctx, collective, contributer)
	suite.Require().NoError(err)

	// check collective bonds decreased
	updatedCollective := suite.app.CollectivesKeeper.GetCollective(suite.ctx, collective.Name)
	suite.Require().Equal(updatedCollective.Bonds, []sdk.Coin{sdk.NewInt64Coin("ukex", 100_000)})

	// check collective contributor removed
	updatedContributer := suite.app.CollectivesKeeper.GetCollectiveContributer(suite.ctx, contributer.Name, contributer.Address)
	suite.Require().Equal(updatedContributer.Name, "")

	// check partial sent from collective address
	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, collective.GetCollectiveAddress())
	suite.Require().Equal(balances, []sdk.Coin{sdk.NewInt64Coin("ukex", 50_000)})
	// check partial sent from collective donation address
	balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, collective.GetCollectiveDonationAddress())
	suite.Require().Equal(balances, []sdk.Coin{sdk.NewInt64Coin("ukex", 50_000)})
}

func (suite *KeeperTestSuite) ExecuteCollectiveRemove() {
	coins := []sdk.Coin{sdk.NewInt64Coin("ukex", 1000_000)}
	collective := types.Collective{
		Name:             "collective1",
		Description:      "collective1 description",
		Status:           types.CollectiveActive,
		DepositWhitelist: types.DepositWhitelist{Any: true},
		OwnersWhitelist: types.OwnersWhitelist{
			Roles:    []uint64{1},
			Accounts: []string{"kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r"},
		},
		SpendingPools: []types.WeightedSpendingPool{
			{
				Name:   "spendingpool1",
				Weight: sdk.NewDec(1),
			},
		},
		ClaimStart:       0,
		ClaimPeriod:      1000,
		ClaimEnd:         100000,
		VoteQuorum:       sdk.NewDecWithPrec(30, 2),
		VotePeriod:       86400,
		VoteEnactment:    3000,
		Donations:        []sdk.Coin{sdk.NewInt64Coin("ukex", 100_000)},
		Rewards:          []sdk.Coin{sdk.NewInt64Coin("ukex", 1000_000)},
		LastDistribution: 0,
		Bonds:            []sdk.Coin{sdk.NewInt64Coin("ukex", 2000_000)},
		CreationTime:     0,
	}

	contributer := types.CollectiveContributor{
		Address:      "kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r",
		Name:         "collective1",
		Bonds:        coins,
		Locking:      1000,
		Donation:     sdk.NewDecWithPrec(5, 1), // 50%
		DonationLock: true,
	}

	pool := spendingtypes.SpendingPool{
		Name:     "spendingpool1",
		Balances: []sdk.Coin{},
	}

	suite.app.CollectivesKeeper.SetCollective(suite.ctx, collective)
	suite.app.CollectivesKeeper.SetCollectiveContributer(suite.ctx, contributer)
	suite.app.SpendingKeeper.SetSpendingPool(suite.ctx, pool)

	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, collective.GetCollectiveAddress(), coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, collective.GetCollectiveDonationAddress(), coins)
	suite.Require().NoError(err)

	err = suite.app.CollectivesKeeper.ExecuteCollectiveRemove(suite.ctx, collective)
	suite.Require().NoError(err)

	// check collective removed
	updatedCollective := suite.app.CollectivesKeeper.GetCollective(suite.ctx, collective.Name)
	suite.Require().Equal(updatedCollective.Name, "")

	// check collective contributor removed
	updatedContributer := suite.app.CollectivesKeeper.GetCollectiveContributer(suite.ctx, contributer.Name, contributer.Address)
	suite.Require().Equal(updatedContributer.Name, "")

	// check rewards are distributed after removal
	spendingPool := suite.app.SpendingKeeper.GetSpendingPool(suite.ctx, "spendingpool1")
	suite.Require().NotNil(spendingPool)
	suite.Require().False(sdk.Coins(spendingPool.Balances).Empty())
}
