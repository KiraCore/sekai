package keeper_test

import (
	"github.com/KiraCore/sekai/x/layer2/keeper"
	"github.com/KiraCore/sekai/x/layer2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestTokenInfoSetGetDelete() {
	infos := []types.TokenInfo{
		{
			TokenType:   "adr20",
			Denom:       "ku/bridgebtc",
			Name:        "Bridge BTC",
			Symbol:      "BTC",
			Icon:        "",
			Description: "",
			Website:     "",
			Social:      "",
			Decimals:    8,
			Cap:         sdk.NewInt(1000_000_000),
			Supply:      sdk.ZeroInt(),
			Holders:     0,
			Fee:         sdk.ZeroInt(),
			Owner:       "kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d",
		},
		{
			TokenType:   "adr43",
			Denom:       "ku/punk",
			Name:        "Bridge PUNK",
			Symbol:      "PUNK",
			Icon:        "",
			Description: "",
			Website:     "",
			Social:      "",
			Decimals:    0,
			Cap:         sdk.NewInt(1000_000_000),
			Supply:      sdk.ZeroInt(),
			Holders:     0,
			Fee:         sdk.ZeroInt(),
			Owner:       "kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d",
			Hash:        "",
			Metadata:    "",
		},
	}

	for _, info := range infos {
		suite.app.Layer2Keeper.SetTokenInfo(suite.ctx, info)
	}

	for _, info := range infos {
		c := suite.app.Layer2Keeper.GetTokenInfo(suite.ctx, info.Denom)
		suite.Require().Equal(c, info)
	}

	allInfos := suite.app.Layer2Keeper.GetTokenInfos(suite.ctx)
	suite.Require().Len(allInfos, 2)

	suite.app.Layer2Keeper.DeleteTokenInfo(suite.ctx, infos[0].Denom)

	allInfos = suite.app.Layer2Keeper.GetTokenInfos(suite.ctx)
	suite.Require().Len(allInfos, 1)
}

func (suite *KeeperTestSuite) TestPow10() {
	suite.Require().Equal(keeper.Pow10(0), sdk.NewInt(1))
	suite.Require().Equal(keeper.Pow10(1), sdk.NewInt(10))
	suite.Require().Equal(keeper.Pow10(2), sdk.NewInt(100))
	suite.Require().Equal(keeper.Pow10(3), sdk.NewInt(1000))
}
