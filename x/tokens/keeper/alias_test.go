package keeper_test

import (
	"github.com/KiraCore/sekai/x/tokens/types"
)

func (suite *KeeperTestSuite) TestTokenAlias() {
	suite.SetupTest()
	ctx := suite.ctx

	// check initial token alias before registration
	alias := suite.app.TokensKeeper.GetTokenAlias(ctx, "stake")
	suite.Require().Nil(alias)
	aliases := suite.app.TokensKeeper.ListTokenAlias(ctx)
	suite.Require().Len(aliases, 1)
	aliasMap := suite.app.TokensKeeper.GetTokenAliasesByDenom(ctx, []string{"stake"})
	suite.Require().Equal(len(aliasMap), 0)
	suite.Require().Nil(aliasMap["stake"])

	// upsert token alias and check
	newAlias := types.TokenAlias{
		Symbol:   "stake",
		Name:     "test token",
		Icon:     "icon_link",
		Decimals: 6,
		Denoms:   []string{"stake"},
	}
	suite.app.TokensKeeper.UpsertTokenAlias(ctx, newAlias)
	alias = suite.app.TokensKeeper.GetTokenAlias(ctx, "stake")
	suite.Require().NotNil(alias)
	aliases = suite.app.TokensKeeper.ListTokenAlias(ctx)
	suite.Require().Len(aliases, 2)
	aliasMap = suite.app.TokensKeeper.GetTokenAliasesByDenom(ctx, []string{"stake"})
	suite.Require().Equal(len(aliasMap), 1)
	suite.Require().NotNil(aliasMap["stake"])

	// delete token alias and check
	suite.app.TokensKeeper.DeleteTokenAlias(ctx, "stake")
	alias = suite.app.TokensKeeper.GetTokenAlias(ctx, "stake")
	suite.Require().Nil(alias)
	aliases = suite.app.TokensKeeper.ListTokenAlias(ctx)
	suite.Require().Len(aliases, 1)
	aliasMap = suite.app.TokensKeeper.GetTokenAliasesByDenom(ctx, []string{"stake"})
	suite.Require().Equal(len(aliasMap), 0)
	suite.Require().Nil(aliasMap["stake"])
}
