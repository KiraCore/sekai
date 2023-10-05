package keeper_test

import (
	"github.com/KiraCore/sekai/x/multistaking/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestRRTokenHolderSetGetDelete() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr3 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	holders := []sdk.AccAddress{
		addr1,
		addr2,
		addr3,
	}

	for _, holder := range holders {
		suite.app.RecoveryKeeper.SetRRTokenHolder(suite.ctx, "rr/val1", holder)
	}

	for _, holder := range holders {
		isHolder := suite.app.RecoveryKeeper.IsRRTokenHolder(suite.ctx, "rr/val1", holder)
		suite.Require().True(isHolder)
	}

	allHolders := suite.app.RecoveryKeeper.GetRRTokenHolders(suite.ctx, "rr/val1")
	suite.Require().Len(allHolders, 3)

	suite.app.RecoveryKeeper.RemoveRRTokenHolder(suite.ctx, "rr/val1", holders[0])

	allHolders = suite.app.RecoveryKeeper.GetRRTokenHolders(suite.ctx, "rr/val1")
	suite.Require().Len(allHolders, 2)
}

func (suite *KeeperTestSuite) TestRRTokenHolderRewardsSetGetDelete() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr3 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	rewards := []types.Rewards{
		{
			Delegator: addr1.String(),
			Rewards:   sdk.NewCoins(sdk.NewInt64Coin("token1", 1000)),
		},
		{
			Delegator: addr2.String(),
			Rewards:   sdk.NewCoins(sdk.NewInt64Coin("token2", 1000)),
		},
		{
			Delegator: addr3.String(),
			Rewards:   sdk.NewCoins(sdk.NewInt64Coin("token3", 1000)),
		},
	}

	for _, reward := range rewards {
		holder := sdk.MustAccAddressFromBech32(reward.Delegator)
		suite.app.RecoveryKeeper.SetRRTokenHolderRewards(suite.ctx, holder, reward.Rewards)
	}

	for _, reward := range rewards {
		holder := sdk.MustAccAddressFromBech32(reward.Delegator)
		coins := suite.app.RecoveryKeeper.GetRRTokenHolderRewards(suite.ctx, holder)
		suite.Require().Equal(coins.String(), sdk.Coins(reward.Rewards).String())
	}

	allRewards := suite.app.RecoveryKeeper.GetAllRRHolderRewards(suite.ctx)
	suite.Require().Len(allRewards, 3)

	suite.app.RecoveryKeeper.RemoveRRTokenHolderRewards(suite.ctx, addr1)

	allRewards = suite.app.RecoveryKeeper.GetAllRRHolderRewards(suite.ctx)
	suite.Require().Len(allRewards, 2)

	suite.app.RecoveryKeeper.IncreaseRRTokenHolderRewards(suite.ctx, addr1, sdk.NewCoins(sdk.NewInt64Coin("token2", 1000)))
	allRewards = suite.app.RecoveryKeeper.GetAllRRHolderRewards(suite.ctx)
	suite.Require().Len(allRewards, 3)
}

// TODO: add test for UnregisterNotEnoughAmountHolder
