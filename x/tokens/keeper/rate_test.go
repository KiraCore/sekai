package keeper_test

import (
	"github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestTokenRates() {
	suite.SetupTest()
	ctx := suite.ctx

	// check initial token rate before registration
	rate := suite.app.TokensKeeper.GetTokenRate(ctx, "stake")
	suite.Require().Nil(rate)
	rates := suite.app.TokensKeeper.ListTokenRate(ctx)
	suite.Require().Len(rates, 4)
	rateMap := suite.app.TokensKeeper.GetTokenRatesByDenom(ctx, []string{"stake"})
	suite.Require().Equal(len(rateMap), 0)
	suite.Require().Nil(rateMap["stake"])

	// upsert token rate and check
	newRate := types.TokenRate{
		Denom:       "stake",
		Rate:        sdk.NewDec(2),
		FeePayments: true,
	}
	suite.app.TokensKeeper.UpsertTokenRate(ctx, newRate)
	rate = suite.app.TokensKeeper.GetTokenRate(ctx, "stake")
	suite.Require().NotNil(rate)
	rates = suite.app.TokensKeeper.ListTokenRate(ctx)
	suite.Require().Len(rates, 5)
	rateMap = suite.app.TokensKeeper.GetTokenRatesByDenom(ctx, []string{"stake"})
	suite.Require().Equal(len(rateMap), 1)
	suite.Require().NotNil(rateMap["stake"])

	// delete token rate and check
	suite.app.TokensKeeper.DeleteTokenRate(ctx, "stake")
	rate = suite.app.TokensKeeper.GetTokenRate(ctx, "stake")
	suite.Require().Nil(rate)
	rates = suite.app.TokensKeeper.ListTokenRate(ctx)
	suite.Require().Len(rates, 4)
	rateMap = suite.app.TokensKeeper.GetTokenRatesByDenom(ctx, []string{"stake"})
	suite.Require().Equal(len(rateMap), 0)
	suite.Require().Nil(rateMap["stake"])
}
