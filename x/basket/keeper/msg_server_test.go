package keeper_test

import (
	"time"

	"github.com/KiraCore/sekai/x/basket/keeper"
	"github.com/KiraCore/sekai/x/basket/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestDisableBasketDeposits() {
	suite.SetupTest()
	now := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(now)

	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	suite.app.CustomGovKeeper.SaveNetworkActor(suite.ctx, govtypes.NetworkActor{
		Address: addr2,
		Permissions: &govtypes.Permissions{
			Whitelist: []uint32{uint32(govtypes.PermHandleBasketEmergency)},
		},
	})

	basket := types.Basket{
		Id:              1,
		Suffix:          "usd",
		Description:     "usd basket",
		LimitsPeriod:    3600,
		Amount:          sdk.NewInt(0),
		SwapFee:         sdk.NewDecWithPrec(1, 2), // 1%
		SlipppageFeeMin: sdk.NewDecWithPrec(1, 2), // 1%
		TokensCap:       sdk.NewDecWithPrec(9, 1), // 90%
		MintsMin:        sdk.NewInt(1),
		MintsMax:        sdk.NewInt(100000000),
		MintsDisabled:   false,
		BurnsMin:        sdk.NewInt(1),
		BurnsMax:        sdk.NewInt(100000000),
		BurnsDisabled:   false,
		SwapsMin:        sdk.NewInt(1),
		SwapsMax:        sdk.NewInt(100000000),
		SwapsDisabled:   false,
		Tokens: []types.BasketToken{
			{
				Denom:     "ukex",
				Weight:    sdk.NewDec(1),
				Amount:    sdk.NewInt(0),
				Deposits:  true,
				Withdraws: true,
				Swaps:     true,
			},
			{
				Denom:     "ueth",
				Weight:    sdk.NewDec(10),
				Amount:    sdk.NewInt(0),
				Deposits:  true,
				Withdraws: true,
				Swaps:     true,
			},
		},
		Surplus: sdk.NewCoins(sdk.NewInt64Coin("ueth", 1)),
	}
	suite.app.BasketKeeper.SetBasket(suite.ctx, basket)

	msgServer := keeper.NewMsgServerImpl(suite.app.BasketKeeper, suite.app.CustomGovKeeper)
	_, err := msgServer.DisableBasketDeposits(sdk.WrapSDKContext(suite.ctx), &types.MsgDisableBasketDeposits{
		Sender:   addr1.String(),
		BasketId: 1,
	})
	suite.Require().Error(err)

	_, err = msgServer.DisableBasketDeposits(sdk.WrapSDKContext(suite.ctx), &types.MsgDisableBasketDeposits{
		Sender:   addr2.String(),
		BasketId: 1,
		Disabled: true,
	})
	suite.Require().NoError(err)

	// check disabled
	savedBasket, err := suite.app.BasketKeeper.GetBasketById(suite.ctx, basket.Id)
	suite.Require().NoError(err)
	suite.Require().True(savedBasket.MintsDisabled)
}

func (suite *KeeperTestSuite) TestDisableBasketWithdraws() {
	suite.SetupTest()
	now := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(now)

	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	suite.app.CustomGovKeeper.SaveNetworkActor(suite.ctx, govtypes.NetworkActor{
		Address: addr2,
		Permissions: &govtypes.Permissions{
			Whitelist: []uint32{uint32(govtypes.PermHandleBasketEmergency)},
		},
	})

	basket := types.Basket{
		Id:              1,
		Suffix:          "usd",
		Description:     "usd basket",
		LimitsPeriod:    3600,
		Amount:          sdk.NewInt(0),
		SwapFee:         sdk.NewDecWithPrec(1, 2), // 1%
		SlipppageFeeMin: sdk.NewDecWithPrec(1, 2), // 1%
		TokensCap:       sdk.NewDecWithPrec(9, 1), // 90%
		MintsMin:        sdk.NewInt(1),
		MintsMax:        sdk.NewInt(100000000),
		MintsDisabled:   false,
		BurnsMin:        sdk.NewInt(1),
		BurnsMax:        sdk.NewInt(100000000),
		BurnsDisabled:   false,
		SwapsMin:        sdk.NewInt(1),
		SwapsMax:        sdk.NewInt(100000000),
		SwapsDisabled:   false,
		Tokens: []types.BasketToken{
			{
				Denom:     "ukex",
				Weight:    sdk.NewDec(1),
				Amount:    sdk.NewInt(0),
				Deposits:  true,
				Withdraws: true,
				Swaps:     true,
			},
			{
				Denom:     "ueth",
				Weight:    sdk.NewDec(10),
				Amount:    sdk.NewInt(0),
				Deposits:  true,
				Withdraws: true,
				Swaps:     true,
			},
		},
		Surplus: sdk.NewCoins(sdk.NewInt64Coin("ueth", 1)),
	}
	suite.app.BasketKeeper.SetBasket(suite.ctx, basket)

	msgServer := keeper.NewMsgServerImpl(suite.app.BasketKeeper, suite.app.CustomGovKeeper)
	_, err := msgServer.DisableBasketWithdraws(sdk.WrapSDKContext(suite.ctx), &types.MsgDisableBasketWithdraws{
		Sender:   addr1.String(),
		BasketId: 1,
	})
	suite.Require().Error(err)

	_, err = msgServer.DisableBasketWithdraws(sdk.WrapSDKContext(suite.ctx), &types.MsgDisableBasketWithdraws{
		Sender:   addr2.String(),
		BasketId: 1,
		Disabled: true,
	})
	suite.Require().NoError(err)

	// check disabled
	savedBasket, err := suite.app.BasketKeeper.GetBasketById(suite.ctx, basket.Id)
	suite.Require().NoError(err)
	suite.Require().True(savedBasket.BurnsDisabled)
}

func (suite *KeeperTestSuite) TestDisableBasketSwaps() {
	suite.SetupTest()
	now := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(now)

	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	suite.app.CustomGovKeeper.SaveNetworkActor(suite.ctx, govtypes.NetworkActor{
		Address: addr2,
		Permissions: &govtypes.Permissions{
			Whitelist: []uint32{uint32(govtypes.PermHandleBasketEmergency)},
		},
	})

	basket := types.Basket{
		Id:              1,
		Suffix:          "usd",
		Description:     "usd basket",
		LimitsPeriod:    3600,
		Amount:          sdk.NewInt(0),
		SwapFee:         sdk.NewDecWithPrec(1, 2), // 1%
		SlipppageFeeMin: sdk.NewDecWithPrec(1, 2), // 1%
		TokensCap:       sdk.NewDecWithPrec(9, 1), // 90%
		MintsMin:        sdk.NewInt(1),
		MintsMax:        sdk.NewInt(100000000),
		MintsDisabled:   false,
		BurnsMin:        sdk.NewInt(1),
		BurnsMax:        sdk.NewInt(100000000),
		BurnsDisabled:   false,
		SwapsMin:        sdk.NewInt(1),
		SwapsMax:        sdk.NewInt(100000000),
		SwapsDisabled:   false,
		Tokens: []types.BasketToken{
			{
				Denom:     "ukex",
				Weight:    sdk.NewDec(1),
				Amount:    sdk.NewInt(0),
				Deposits:  true,
				Withdraws: true,
				Swaps:     true,
			},
			{
				Denom:     "ueth",
				Weight:    sdk.NewDec(10),
				Amount:    sdk.NewInt(0),
				Deposits:  true,
				Withdraws: true,
				Swaps:     true,
			},
		},
		Surplus: sdk.NewCoins(sdk.NewInt64Coin("ueth", 1)),
	}
	suite.app.BasketKeeper.SetBasket(suite.ctx, basket)

	msgServer := keeper.NewMsgServerImpl(suite.app.BasketKeeper, suite.app.CustomGovKeeper)
	_, err := msgServer.DisableBasketSwaps(sdk.WrapSDKContext(suite.ctx), &types.MsgDisableBasketSwaps{
		Sender:   addr1.String(),
		BasketId: 1,
	})
	suite.Require().Error(err)

	_, err = msgServer.DisableBasketSwaps(sdk.WrapSDKContext(suite.ctx), &types.MsgDisableBasketSwaps{
		Sender:   addr2.String(),
		BasketId: 1,
		Disabled: true,
	})
	suite.Require().NoError(err)

	// check disabled
	savedBasket, err := suite.app.BasketKeeper.GetBasketById(suite.ctx, basket.Id)
	suite.Require().NoError(err)
	suite.Require().True(savedBasket.SwapsDisabled)
}
