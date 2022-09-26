package keeper_test

import (
	"time"

	"github.com/KiraCore/sekai/x/basket/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestMintActions() {
	suite.SetupTest()

	now := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(now.Add(time.Hour))

	baskets := []types.Basket{
		{
			Id:           1,
			LimitsPeriod: 1800,
		},
		{
			Id:           2,
			LimitsPeriod: 7200,
		},
	}

	for _, basket := range baskets {
		suite.app.BasketKeeper.SetBasket(suite.ctx, basket)
	}

	actions := []types.AmountAtTime{
		{
			BasketId: 1,
			Time:     uint64(now.Unix()),
			Amount:   sdk.NewInt(1000),
		},
		{
			BasketId: 1,
			Time:     uint64(now.Add(time.Second).Unix()),
			Amount:   sdk.NewInt(10000),
		},
		{
			BasketId: 2,
			Time:     uint64(now.Unix()),
			Amount:   sdk.NewInt(1000),
		},
		{
			BasketId: 2,
			Time:     uint64(now.Add(time.Hour).Unix()),
			Amount:   sdk.NewInt(10000),
		},
	}

	for _, action := range actions {
		suite.app.BasketKeeper.SetMintAmount(suite.ctx, time.Unix(int64(action.Time), 0), action.BasketId, action.Amount)
	}

	for _, action := range actions {
		amount := suite.app.BasketKeeper.GetMintAmount(suite.ctx, action.BasketId, time.Unix(int64(action.Time), 0))
		suite.Require().Equal(amount, action.Amount)
	}

	allAmounts := suite.app.BasketKeeper.GetAllMintAmounts(suite.ctx)
	suite.Require().Len(allAmounts, len(actions))

	historicalAmount := suite.app.BasketKeeper.GetLimitsPeriodMintAmount(suite.ctx, 1, 1800)
	suite.Require().Equal(historicalAmount, sdk.NewInt(0))
	historicalAmount = suite.app.BasketKeeper.GetLimitsPeriodMintAmount(suite.ctx, 2, 1800)
	suite.Require().Equal(historicalAmount, sdk.NewInt(10000))
	historicalAmount = suite.app.BasketKeeper.GetLimitsPeriodMintAmount(suite.ctx, 2, 7200)
	suite.Require().Equal(historicalAmount, sdk.NewInt(11000))

	suite.app.BasketKeeper.RegisterMintAction(suite.ctx, 1, sdk.NewInt(1000))
	historicalAmount = suite.app.BasketKeeper.GetLimitsPeriodMintAmount(suite.ctx, 1, 1800)
	suite.Require().Equal(historicalAmount, sdk.NewInt(1000))

	suite.app.BasketKeeper.ClearOldMintAmounts(suite.ctx, 1, 1800)
	historicalAmount = suite.app.BasketKeeper.GetLimitsPeriodMintAmount(suite.ctx, 1, 1800000)
	suite.Require().Equal(historicalAmount, sdk.NewInt(1000))
}

func (suite *KeeperTestSuite) TestBurnActions() {
	suite.SetupTest()

	now := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(now.Add(time.Hour))

	baskets := []types.Basket{
		{
			Id:           1,
			LimitsPeriod: 1800,
		},
		{
			Id:           2,
			LimitsPeriod: 7200,
		},
	}

	for _, basket := range baskets {
		suite.app.BasketKeeper.SetBasket(suite.ctx, basket)
	}

	actions := []types.AmountAtTime{
		{
			BasketId: 1,
			Time:     uint64(now.Unix()),
			Amount:   sdk.NewInt(1000),
		},
		{
			BasketId: 1,
			Time:     uint64(now.Add(time.Second).Unix()),
			Amount:   sdk.NewInt(10000),
		},
		{
			BasketId: 2,
			Time:     uint64(now.Unix()),
			Amount:   sdk.NewInt(1000),
		},
		{
			BasketId: 2,
			Time:     uint64(now.Add(time.Hour).Unix()),
			Amount:   sdk.NewInt(10000),
		},
	}

	for _, action := range actions {
		suite.app.BasketKeeper.SetBurnAmount(suite.ctx, time.Unix(int64(action.Time), 0), action.BasketId, action.Amount)
	}

	for _, action := range actions {
		amount := suite.app.BasketKeeper.GetBurnAmount(suite.ctx, action.BasketId, time.Unix(int64(action.Time), 0))
		suite.Require().Equal(amount, action.Amount)
	}

	allAmounts := suite.app.BasketKeeper.GetAllBurnAmounts(suite.ctx)
	suite.Require().Len(allAmounts, len(actions))

	historicalAmount := suite.app.BasketKeeper.GetLimitsPeriodBurnAmount(suite.ctx, 1, 1800)
	suite.Require().Equal(historicalAmount, sdk.NewInt(0))
	historicalAmount = suite.app.BasketKeeper.GetLimitsPeriodBurnAmount(suite.ctx, 2, 1800)
	suite.Require().Equal(historicalAmount, sdk.NewInt(10000))
	historicalAmount = suite.app.BasketKeeper.GetLimitsPeriodBurnAmount(suite.ctx, 2, 7200)
	suite.Require().Equal(historicalAmount, sdk.NewInt(11000))

	suite.app.BasketKeeper.RegisterBurnAction(suite.ctx, 1, sdk.NewInt(1000))
	historicalAmount = suite.app.BasketKeeper.GetLimitsPeriodBurnAmount(suite.ctx, 1, 1800)
	suite.Require().Equal(historicalAmount, sdk.NewInt(1000))

	suite.app.BasketKeeper.ClearOldBurnAmounts(suite.ctx, 1, 1800)
	historicalAmount = suite.app.BasketKeeper.GetLimitsPeriodBurnAmount(suite.ctx, 1, 1800000)
	suite.Require().Equal(historicalAmount, sdk.NewInt(1000))
}

func (suite *KeeperTestSuite) TestSwapActions() {
	suite.SetupTest()

	now := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(now.Add(time.Hour))

	baskets := []types.Basket{
		{
			Id:           1,
			LimitsPeriod: 1800,
		},
		{
			Id:           2,
			LimitsPeriod: 7200,
		},
	}

	for _, basket := range baskets {
		suite.app.BasketKeeper.SetBasket(suite.ctx, basket)
	}

	actions := []types.AmountAtTime{
		{
			BasketId: 1,
			Time:     uint64(now.Unix()),
			Amount:   sdk.NewInt(1000),
		},
		{
			BasketId: 1,
			Time:     uint64(now.Add(time.Second).Unix()),
			Amount:   sdk.NewInt(10000),
		},
		{
			BasketId: 2,
			Time:     uint64(now.Unix()),
			Amount:   sdk.NewInt(1000),
		},
		{
			BasketId: 2,
			Time:     uint64(now.Add(time.Hour).Unix()),
			Amount:   sdk.NewInt(10000),
		},
	}

	for _, action := range actions {
		suite.app.BasketKeeper.SetSwapAmount(suite.ctx, time.Unix(int64(action.Time), 0), action.BasketId, action.Amount)
	}

	for _, action := range actions {
		amount := suite.app.BasketKeeper.GetSwapAmount(suite.ctx, action.BasketId, time.Unix(int64(action.Time), 0))
		suite.Require().Equal(amount, action.Amount)
	}

	allAmounts := suite.app.BasketKeeper.GetAllSwapAmounts(suite.ctx)
	suite.Require().Len(allAmounts, len(actions))

	historicalAmount := suite.app.BasketKeeper.GetLimitsPeriodSwapAmount(suite.ctx, 1, 1800)
	suite.Require().Equal(historicalAmount, sdk.NewInt(0))
	historicalAmount = suite.app.BasketKeeper.GetLimitsPeriodSwapAmount(suite.ctx, 2, 1800)
	suite.Require().Equal(historicalAmount, sdk.NewInt(10000))
	historicalAmount = suite.app.BasketKeeper.GetLimitsPeriodSwapAmount(suite.ctx, 2, 7200)
	suite.Require().Equal(historicalAmount, sdk.NewInt(11000))

	suite.app.BasketKeeper.RegisterSwapAction(suite.ctx, 1, sdk.NewInt(1000))
	historicalAmount = suite.app.BasketKeeper.GetLimitsPeriodSwapAmount(suite.ctx, 1, 1800)
	suite.Require().Equal(historicalAmount, sdk.NewInt(1000))

	suite.app.BasketKeeper.ClearOldSwapAmounts(suite.ctx, 1, 1800)
	historicalAmount = suite.app.BasketKeeper.GetLimitsPeriodSwapAmount(suite.ctx, 1, 1800000)
	suite.Require().Equal(historicalAmount, sdk.NewInt(1000))
}
