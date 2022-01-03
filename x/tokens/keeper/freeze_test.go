package keeper_test

import (
	"github.com/KiraCore/sekai/x/tokens/types"
)

func (suite *KeeperTestSuite) TestTokenBlacklist() {
	ctx := suite.ctx

	// check genesis
	blackWhites := suite.app.TokensKeeper.GetTokenBlackWhites(ctx)
	suite.Require().Equal(blackWhites.Blacklisted, []string{"frozen"})

	// try adding two tokens
	suite.app.TokensKeeper.AddTokensToBlacklist(ctx, []string{"frozen1", "frozen2"})
	blackWhites = suite.app.TokensKeeper.GetTokenBlackWhites(ctx)
	suite.Require().Equal(blackWhites.Blacklisted, []string{"frozen", "frozen1", "frozen2"})

	// try adding one token
	suite.app.TokensKeeper.AddTokensToBlacklist(ctx, []string{"frozen3"})
	blackWhites = suite.app.TokensKeeper.GetTokenBlackWhites(ctx)
	suite.Require().Equal(blackWhites.Blacklisted, []string{"frozen", "frozen1", "frozen2", "frozen3"})

	// try adding no changes
	suite.app.TokensKeeper.AddTokensToBlacklist(ctx, []string{})
	blackWhites = suite.app.TokensKeeper.GetTokenBlackWhites(ctx)
	suite.Require().Equal(blackWhites.Blacklisted, []string{"frozen", "frozen1", "frozen2", "frozen3"})

	// try adding empty denom
	suite.app.TokensKeeper.AddTokensToBlacklist(ctx, []string{""})
	blackWhites = suite.app.TokensKeeper.GetTokenBlackWhites(ctx)
	suite.Require().Equal(blackWhites.Blacklisted, []string{"frozen", "frozen1", "frozen2", "frozen3", ""})

	// try removing blacklisted
	suite.app.TokensKeeper.RemoveTokensFromBlacklist(ctx, []string{"frozen2"})
	blackWhites = suite.app.TokensKeeper.GetTokenBlackWhites(ctx)
	suite.Require().Equal(blackWhites.Blacklisted, []string{"frozen", "frozen1", "", "frozen3"})

	// try removing not blacklisted
	suite.app.TokensKeeper.RemoveTokensFromBlacklist(ctx, []string{"frozen4"})
	blackWhites = suite.app.TokensKeeper.GetTokenBlackWhites(ctx)
	suite.Require().Equal(blackWhites.Blacklisted, []string{"frozen", "frozen1", "", "frozen3"})
}

func (suite *KeeperTestSuite) TestTokenWhitelist() {
	ctx := suite.ctx

	// check genesis
	blackWhites := suite.app.TokensKeeper.GetTokenBlackWhites(ctx)
	suite.Require().Equal(blackWhites.Whitelisted, []string{"ukex"})

	// try adding two tokens
	suite.app.TokensKeeper.AddTokensToWhitelist(ctx, []string{"white1", "white2"})
	blackWhites = suite.app.TokensKeeper.GetTokenBlackWhites(ctx)
	suite.Require().Equal(blackWhites.Whitelisted, []string{"ukex", "white1", "white2"})

	// try adding one token
	suite.app.TokensKeeper.AddTokensToWhitelist(ctx, []string{"white3"})
	blackWhites = suite.app.TokensKeeper.GetTokenBlackWhites(ctx)
	suite.Require().Equal(blackWhites.Whitelisted, []string{"ukex", "white1", "white2", "white3"})

	// try adding no changes
	suite.app.TokensKeeper.AddTokensToWhitelist(ctx, []string{})
	blackWhites = suite.app.TokensKeeper.GetTokenBlackWhites(ctx)
	suite.Require().Equal(blackWhites.Whitelisted, []string{"ukex", "white1", "white2", "white3"})

	// try adding empty denom
	suite.app.TokensKeeper.AddTokensToWhitelist(ctx, []string{""})
	blackWhites = suite.app.TokensKeeper.GetTokenBlackWhites(ctx)
	suite.Require().Equal(blackWhites.Whitelisted, []string{"ukex", "white1", "white2", "white3", ""})

	// try removing whitelisted
	suite.app.TokensKeeper.RemoveTokensFromWhitelist(ctx, []string{"white2"})
	blackWhites = suite.app.TokensKeeper.GetTokenBlackWhites(ctx)
	suite.Require().Equal(blackWhites.Whitelisted, []string{"ukex", "white1", "", "white3"})

	// try removing not whitelisted
	suite.app.TokensKeeper.RemoveTokensFromWhitelist(ctx, []string{"white4"})
	blackWhites = suite.app.TokensKeeper.GetTokenBlackWhites(ctx)
	suite.Require().Equal(blackWhites.Whitelisted, []string{"ukex", "white1", "", "white3"})
}

func (suite *KeeperTestSuite) TestTokenBlackWhiteSetGet() {
	ctx := suite.ctx

	blackWhites := types.TokensWhiteBlack{
		Whitelisted: []string{"newwhite"},
		Blacklisted: []string{"newblack"},
	}
	suite.app.TokensKeeper.SetTokenBlackWhites(ctx, &blackWhites)
	bw := suite.app.TokensKeeper.GetTokenBlackWhites(ctx)
	suite.Require().NotNil(bw)
	suite.Require().Equal(blackWhites, *bw)
}
