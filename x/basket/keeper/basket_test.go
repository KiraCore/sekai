package keeper_test

import (
	"time"

	"github.com/KiraCore/sekai/x/basket/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestLastBasketIdGetSet() {
	suite.SetupTest()

	now := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(now.Add(time.Hour))

	suite.app.BasketKeeper.SetLastBasketId(suite.ctx, 1)
	lastBasketId := suite.app.BasketKeeper.GetLastBasketId(suite.ctx)
	suite.Require().Equal(lastBasketId, uint64(1))
}

func (suite *KeeperTestSuite) TestBasketSetGetDelete() {
	baskets := []types.Basket{
		{
			Id:              1,
			Suffix:          "usd",
			Description:     "usd basket",
			LimitsPeriod:    1800,
			Amount:          sdk.NewInt(1000),
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
					Denom:     "usd1",
					Weight:    sdk.NewDec(1),
					Amount:    sdk.NewInt(500),
					Deposits:  true,
					Withdraws: true,
					Swaps:     true,
				},
				{
					Denom:     "usd2",
					Weight:    sdk.NewDec(10),
					Amount:    sdk.NewInt(50),
					Deposits:  true,
					Withdraws: true,
					Swaps:     true,
				},
			},
			Surplus: sdk.NewCoins(sdk.NewInt64Coin("usd2", 1)),
		},
		{
			Id:              2,
			Suffix:          "euro",
			Description:     "euro basket",
			LimitsPeriod:    1800,
			Amount:          sdk.NewInt(1000),
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
					Denom:     "euro1",
					Weight:    sdk.NewDec(1),
					Amount:    sdk.NewInt(500),
					Deposits:  true,
					Withdraws: true,
					Swaps:     true,
				},
				{
					Denom:     "euro2",
					Weight:    sdk.NewDec(10),
					Amount:    sdk.NewInt(50),
					Deposits:  true,
					Withdraws: true,
					Swaps:     true,
				},
			},
			Surplus: sdk.NewCoins(sdk.NewInt64Coin("euro2", 1)),
		},
	}

	for _, basket := range baskets {
		suite.app.BasketKeeper.SetBasket(suite.ctx, basket)
	}

	for _, basket := range baskets {
		b, err := suite.app.BasketKeeper.GetBasketById(suite.ctx, basket.Id)
		suite.Require().NoError(err)
		suite.Require().Equal(b, basket)

		b, err = suite.app.BasketKeeper.GetBasketByDenom(suite.ctx, basket.GetBasketDenom())
		suite.Require().NoError(err)
		suite.Require().Equal(b, basket)
	}

	allBaskets := suite.app.BasketKeeper.GetAllBaskets(suite.ctx)
	suite.Require().Len(allBaskets, 2)

	suite.app.BasketKeeper.DeleteBasket(suite.ctx, baskets[0])

	allBaskets = suite.app.BasketKeeper.GetAllBaskets(suite.ctx)
	suite.Require().Len(allBaskets, 1)

	_, err := suite.app.BasketKeeper.GetBasketById(suite.ctx, baskets[0].Id)
	suite.Require().Error(err)

	_, err = suite.app.BasketKeeper.GetBasketByDenom(suite.ctx, baskets[0].GetBasketDenom())
	suite.Require().Error(err)
}

func (suite *KeeperTestSuite) TestCreateBasket() {
	sampleBasket := types.Basket{
		Id:              1,
		Suffix:          "usd",
		Description:     "usd basket",
		LimitsPeriod:    1800,
		Amount:          sdk.NewInt(1000),
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
				Denom:     "usd1",
				Weight:    sdk.NewDec(1),
				Amount:    sdk.NewInt(500),
				Deposits:  true,
				Withdraws: true,
				Swaps:     true,
			},
			{
				Denom:     "usd2",
				Weight:    sdk.NewDec(10),
				Amount:    sdk.NewInt(50),
				Deposits:  true,
				Withdraws: true,
				Swaps:     true,
			},
		},
		Surplus: sdk.NewCoins(sdk.NewInt64Coin("usd2", 1)),
	}

	err := suite.app.BasketKeeper.CreateBasket(suite.ctx, sampleBasket)
	suite.Require().NoError(err)

	baskets := suite.app.BasketKeeper.GetAllBaskets(suite.ctx)
	suite.Require().Len(baskets, 1)

	basket := baskets[0]
	lastBasketId := suite.app.BasketKeeper.GetLastBasketId(suite.ctx)
	suite.Require().Equal(lastBasketId, basket.Id)
	// surplus zero check
	suite.Require().Nil(basket.Surplus)

	// token amount to be zero
	for _, token := range basket.Tokens {
		suite.Require().Equal(token.Amount, sdk.ZeroInt())
	}

	// total number of tokens zero check
	basket = sampleBasket
	basket.Tokens = []types.BasketToken{}
	err = suite.app.BasketKeeper.CreateBasket(suite.ctx, basket)
	suite.Require().Error(err)

	// token weight zero check
	basket = sampleBasket
	basket.Tokens[0].Weight = sdk.ZeroDec()
	err = suite.app.BasketKeeper.CreateBasket(suite.ctx, basket)
	suite.Require().Error(err)

	// duplicated denom check
	basket = sampleBasket
	basket.Tokens[1].Denom = basket.Tokens[0].Denom
	err = suite.app.BasketKeeper.CreateBasket(suite.ctx, basket)
	suite.Require().Error(err)

	// invalid denom check
	basket = sampleBasket
	basket.Tokens[0].Denom = "1"
	err = suite.app.BasketKeeper.CreateBasket(suite.ctx, basket)
	suite.Require().Error(err)
}

func (suite *KeeperTestSuite) TestEditBasket() {
	sampleBasket := types.Basket{
		Id:              1,
		Suffix:          "usd",
		Description:     "usd basket",
		LimitsPeriod:    1800,
		Amount:          sdk.NewInt(1000),
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
				Denom:     "usd1",
				Weight:    sdk.NewDec(1),
				Amount:    sdk.NewInt(500),
				Deposits:  true,
				Withdraws: true,
				Swaps:     true,
			},
			{
				Denom:     "usd2",
				Weight:    sdk.NewDec(10),
				Amount:    sdk.NewInt(50),
				Deposits:  true,
				Withdraws: true,
				Swaps:     true,
			},
		},
		Surplus: sdk.NewCoins(sdk.NewInt64Coin("usd2", 1)),
	}

	suite.app.BasketKeeper.SetBasket(suite.ctx, sampleBasket)

	// successful update check
	basket := sampleBasket
	basket.Tokens[0].Deposits = false
	err := suite.app.BasketKeeper.EditBasket(suite.ctx, basket)
	suite.Require().NoError(err)
	savedBasket, err := suite.app.BasketKeeper.GetBasketById(suite.ctx, basket.Id)
	suite.Require().NoError(err)
	suite.Require().Equal(savedBasket, basket)

	// surplus amount does not change check
	basket = sampleBasket
	basket.Surplus = sdk.NewCoins(sdk.NewInt64Coin("usd1", 10000000))
	err = suite.app.BasketKeeper.EditBasket(suite.ctx, basket)
	suite.Require().NoError(err)
	savedBasket, err = suite.app.BasketKeeper.GetBasketById(suite.ctx, basket.Id)
	suite.Require().NoError(err)
	suite.Require().Equal(savedBasket.Surplus, sampleBasket.Surplus)

	// token amount derivate from previous state or zero
	basket.Tokens = append(basket.Tokens,
		types.BasketToken{
			Denom:     "usd3",
			Weight:    sdk.NewDec(1),
			Amount:    sdk.NewInt(500),
			Deposits:  true,
			Withdraws: true,
			Swaps:     true,
		})
	err = suite.app.BasketKeeper.EditBasket(suite.ctx, basket)
	suite.Require().NoError(err)
	savedBasket, err = suite.app.BasketKeeper.GetBasketById(suite.ctx, basket.Id)
	suite.Require().NoError(err)
	suite.Require().Len(savedBasket.Tokens, 3)
	suite.Require().Equal(savedBasket.Tokens[0].Amount, sampleBasket.Tokens[0].Amount)
	suite.Require().Equal(savedBasket.Tokens[1].Amount, sampleBasket.Tokens[1].Amount)
	suite.Require().Equal(savedBasket.Tokens[2].Amount, sdk.ZeroInt())

	// basket existance check
	basket = sampleBasket
	basket.Id = 1000
	err = suite.app.BasketKeeper.EditBasket(suite.ctx, basket)
	suite.Require().Error(err)

	// total number of tokens zero check
	basket = sampleBasket
	basket.Tokens = []types.BasketToken{}
	err = suite.app.BasketKeeper.EditBasket(suite.ctx, basket)
	suite.Require().Error(err)

	// token weight zero check
	basket = sampleBasket
	basket.Tokens[0].Weight = sdk.ZeroDec()
	err = suite.app.BasketKeeper.EditBasket(suite.ctx, basket)
	suite.Require().Error(err)

	// duplicated denom check
	basket = sampleBasket
	basket.Tokens[1].Denom = basket.Tokens[0].Denom
	err = suite.app.BasketKeeper.EditBasket(suite.ctx, basket)
	suite.Require().Error(err)

	// invalid denom check
	basket = sampleBasket
	basket.Tokens[0].Denom = "1"
	err = suite.app.BasketKeeper.EditBasket(suite.ctx, basket)
	suite.Require().Error(err)
}
