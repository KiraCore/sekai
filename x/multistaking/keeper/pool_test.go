package keeper_test

import (
	"github.com/KiraCore/sekai/x/multistaking/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestLastPoolIdGetSet() {
	// get default last pool id
	lastPoolId := suite.app.MultiStakingKeeper.GetLastPoolId(suite.ctx)
	suite.Require().Equal(lastPoolId, uint64(0))

	// set last pool id to new value
	newPoolId := uint64(2)
	suite.app.MultiStakingKeeper.SetLastPoolId(suite.ctx, newPoolId)

	// check last pool id update
	lastPoolId = suite.app.MultiStakingKeeper.GetLastPoolId(suite.ctx)
	suite.Require().Equal(lastPoolId, newPoolId)
}

func (suite *KeeperTestSuite) TestStakingPoolGetSet() {
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	// get pool by validator
	_, found := suite.app.MultiStakingKeeper.GetStakingPoolByValidator(suite.ctx, sdk.ValAddress(addr).String())
	suite.Require().False(found)

	// get all pools when not available
	allPools := suite.app.MultiStakingKeeper.GetAllStakingPools(suite.ctx)
	suite.Require().Len(allPools, 0)

	pools := []types.StakingPool{
		{
			Id:         1,
			Validator:  sdk.ValAddress(addr).String(),
			Enabled:    true,
			Commission: sdk.ZeroDec(),
			Slashed:    sdk.ZeroDec(),
		},
		{
			Id:         2,
			Validator:  sdk.ValAddress(addr2).String(),
			Enabled:    true,
			Commission: sdk.ZeroDec(),
			Slashed:    sdk.ZeroDec(),
		},
	}

	for _, pool := range pools {
		suite.app.MultiStakingKeeper.SetStakingPool(suite.ctx, pool)
	}

	// check pools existance by validator
	for _, pool := range pools {
		p, found := suite.app.MultiStakingKeeper.GetStakingPoolByValidator(suite.ctx, pool.Validator)
		suite.Require().True(found)
		suite.Require().Equal(pool, p)
	}

	// check pools for whole export
	allPools = suite.app.MultiStakingKeeper.GetAllStakingPools(suite.ctx)
	suite.Require().Len(allPools, 2)
}
