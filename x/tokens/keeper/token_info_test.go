package keeper_test

import (
	"github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestTokenInfos() {
	suite.SetupTest()
	ctx := suite.ctx

	// check initial token rate before registration
	rate := suite.app.TokensKeeper.GetTokenInfo(ctx, "stake")
	suite.Require().Nil(rate)
	rates := suite.app.TokensKeeper.GetAllTokenInfos(ctx)
	suite.Require().Len(rates, 4)
	rateMap := suite.app.TokensKeeper.GetTokenInfosByDenom(ctx, []string{"stake"})
	suite.Require().Equal(len(rateMap), 1)
	suite.Require().Nil(rateMap["stake"].Data)

	// upsert token rate and check
	newRate := types.TokenInfo{
		Denom:      "stake",
		FeeRate:    sdk.NewDec(2),
		FeeEnabled: true,
	}
	suite.app.TokensKeeper.UpsertTokenInfo(ctx, newRate)
	rate = suite.app.TokensKeeper.GetTokenInfo(ctx, "stake")
	suite.Require().NotNil(rate)
	rates = suite.app.TokensKeeper.GetAllTokenInfos(ctx)
	suite.Require().Len(rates, 5)
	rateMap = suite.app.TokensKeeper.GetTokenInfosByDenom(ctx, []string{"stake"})
	suite.Require().Equal(len(rateMap), 1)
	suite.Require().NotNil(rateMap["stake"])

	// delete token rate and check
	suite.app.TokensKeeper.DeleteTokenInfo(ctx, "stake")
	rate = suite.app.TokensKeeper.GetTokenInfo(ctx, "stake")
	suite.Require().Nil(rate)
	rates = suite.app.TokensKeeper.GetAllTokenInfos(ctx)
	suite.Require().Len(rates, 4)
	rateMap = suite.app.TokensKeeper.GetTokenInfosByDenom(ctx, []string{"stake"})
	suite.Require().Equal(len(rateMap), 1)
	suite.Require().Nil(rateMap["stake"].Data)
}
