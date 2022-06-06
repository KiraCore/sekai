package keeper_test

import (
	"github.com/KiraCore/sekai/x/multistaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

// SetPoolDelegator(ctx sdk.Context, poolId uint64, delegator sdk.AccAddress)
// RemovePoolDelegator(ctx sdk.Context, poolId uint64, delegator sdk.AccAddress)
// GetPoolDelegators(ctx sdk.Context, poolId uint64) []sdk.AccAddress {

// IncreasePoolRewards(ctx sdk.Context, pool types.StakingPool, rewards sdk.Coins) {
// IncreaseDelegatorRewards(ctx sdk.Context, delegator sdk.AccAddress, amounts sdk.Coins) {
// GetDelegatorRewards(ctx sdk.Context, delegator sdk.AccAddress) sdk.Coins {
