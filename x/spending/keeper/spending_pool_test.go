package keeper_test

import (
	"github.com/KiraCore/sekai/x/spending/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestSpendingPoolSetGet() {
	// check genesis has 1 pool
	allSpendingPools := suite.app.SpendingKeeper.GetAllSpendingPools(suite.ctx)
	suite.Require().Len(allSpendingPools, 1)

	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	pools := []types.SpendingPool{
		{
			Name:          "spendingpool1",
			ClaimStart:    0,
			ClaimEnd:      0,
			Rates:         sdk.NewDecCoins(sdk.NewDecCoin("ukex", sdk.NewInt(1))),
			VoteQuorum:    sdk.NewDecWithPrec(30, 2),
			VotePeriod:    3000,
			VoteEnactment: 1000,
			Owners: &types.PermInfo{
				OwnerRoles:    []uint64{1},
				OwnerAccounts: []string{addr1.String()},
			},
			Beneficiaries: &types.WeightedPermInfo{
				Roles: []types.WeightedRole{
					{
						Role:   1,
						Weight: sdk.NewDec(1),
					},
				},
				Accounts: []types.WeightedAccount{
					{
						Account: addr1.String(),
						Weight:  sdk.NewDec(2),
					},
				},
			},
			Balances:                []sdk.Coin{sdk.NewInt64Coin("ukex", 1000_000)},
			DynamicRate:             false,
			DynamicRatePeriod:       0,
			LastDynamicRateCalcTime: 0,
		},
		{
			Name:          "spendingpool2",
			ClaimStart:    0,
			ClaimEnd:      0,
			Rates:         sdk.NewDecCoins(sdk.NewDecCoin("ukex", sdk.NewInt(1))),
			VoteQuorum:    sdk.NewDecWithPrec(30, 2),
			VotePeriod:    3000,
			VoteEnactment: 1000,
			Owners: &types.PermInfo{
				OwnerRoles:    []uint64{1},
				OwnerAccounts: []string{addr1.String()},
			},
			Beneficiaries: &types.WeightedPermInfo{
				Roles: []types.WeightedRole{
					{
						Role:   1,
						Weight: sdk.NewDec(1),
					},
				},
				Accounts: []types.WeightedAccount{
					{
						Account: addr1.String(),
						Weight:  sdk.NewDec(2),
					},
				},
			},
			Balances:                []sdk.Coin{sdk.NewInt64Coin("ukex", 1000_000)},
			DynamicRate:             true,
			DynamicRatePeriod:       43200,
			LastDynamicRateCalcTime: 0,
		},
	}

	for _, pool := range pools {
		suite.app.SpendingKeeper.SetSpendingPool(suite.ctx, pool)
	}

	for _, pool := range pools {
		p := suite.app.SpendingKeeper.GetSpendingPool(suite.ctx, pool.Name)
		suite.Require().NotNil(p)
		suite.Require().Equal(*p, pool)
	}

	// check 2 pools added
	allSpendingPools = suite.app.SpendingKeeper.GetAllSpendingPools(suite.ctx)
	suite.Require().Len(allSpendingPools, 1+2)
}

func (suite *KeeperTestSuite) TestClaimInfoSetGet() {
	// check genesis has empty claim infos
	allClaimInfos := suite.app.SpendingKeeper.GetAllClaimInfos(suite.ctx)
	suite.Require().Len(allClaimInfos, 0)

	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	infos := []types.ClaimInfo{
		{
			Account:   addr1.String(),
			PoolName:  "spendingpool1",
			LastClaim: 0,
		},
		{
			Account:   addr2.String(),
			PoolName:  "spendingpool1",
			LastClaim: 0,
		},
	}

	for _, info := range infos {
		suite.app.SpendingKeeper.SetClaimInfo(suite.ctx, info)
	}

	p := suite.app.SpendingKeeper.GetClaimInfo(suite.ctx, "spendingpool1", addr1)
	suite.Require().NotNil(p)
	suite.Require().Equal(*p, infos[0])
	p = suite.app.SpendingKeeper.GetClaimInfo(suite.ctx, "spendingpool1", addr2)
	suite.Require().NotNil(p)
	suite.Require().Equal(*p, infos[1])

	// check 2 infos added
	allClaimInfos = suite.app.SpendingKeeper.GetAllClaimInfos(suite.ctx)
	suite.Require().Len(allClaimInfos, 2)

	// check 2 pool claim infos added
	poolClaimInfos := suite.app.SpendingKeeper.GetPoolClaimInfos(suite.ctx, "spendingpool1")
	suite.Require().Len(poolClaimInfos, 2)
}

// TODO: CreateSpendingPool
// TODO: ClaimSpendingPool
// TODO: DepositSpendingPoolFromModule
// TODO: DepositSpendingPoolFromAccount
