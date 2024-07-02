package keeper_test

import (
	"github.com/KiraCore/sekai/x/layer2/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (suite *KeeperTestSuite) TestLpTokenPrice() {
	dapp := types.Dapp{
		Name:        "dapp1",
		Denom:       "dapp1",
		Description: "dapp1 description",
		Status:      types.Active,
		Website:     "",
		Logo:        "",
		Social:      "",
		Docs:        "",
		Controllers: types.Controllers{
			Whitelist: types.AccountRange{
				Roles:     []uint64{1},
				Addresses: []string{"kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r"},
			},
		},
		Bin: []types.BinaryInfo{
			{
				Name:      "dapp1",
				Hash:      "dapp1",
				Source:    "dapp1",
				Reference: "dapp1",
				Type:      "dapp1",
			},
		},
		Pool: types.LpPoolConfig{
			Ratio:   sdk.OneDec(),
			Deposit: "",
			Drip:    86400,
		},
		Issuance: types.IssuanceConfig{
			Deposit:  "",
			Premint:  sdk.OneInt(),
			Postmint: sdk.OneInt(),
			Time:     1680141605,
		},
		VoteQuorum:    sdk.NewDecWithPrec(30, 2),
		VotePeriod:    86400,
		VoteEnactment: 3000,
		UpdateTimeMax: 60,
		ExecutorsMin:  1,
		ExecutorsMax:  2,
		VerifiersMin:  1,
		TotalBond:     sdk.NewInt64Coin("ukex", 10000),
		CreationTime:  0,
		PoolFee:       sdk.NewDec(1),
	}

	suite.app.Layer2Keeper.SetDapp(suite.ctx, dapp)

	lpToken := dapp.LpToken()
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.Coins{sdk.NewInt64Coin(lpToken, 10000)})
	suite.Require().NoError(err)

	price := suite.app.Layer2Keeper.LpTokenPrice(suite.ctx, dapp)
	suite.Require().Equal(price, sdk.OneDec())
}

func (suite *KeeperTestSuite) TestOnCollectFee() {
	coins := sdk.Coins{sdk.NewInt64Coin("ukex", 10000)}
	err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins)
	suite.Require().NoError(err)

	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, moduleAddr)
	suite.Require().False(balances.Empty())

	err = suite.app.Layer2Keeper.OnCollectFee(suite.ctx, coins)
	suite.Require().NoError(err)

	balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, moduleAddr)
	suite.Require().True(balances.Empty())
}

func (suite *KeeperTestSuite) TestRedeemDappPoolTx() {
	dapp := types.Dapp{
		Name:        "dapp1",
		Denom:       "dapp1",
		Description: "dapp1 description",
		Status:      types.Active,
		Website:     "",
		Logo:        "",
		Social:      "",
		Docs:        "",
		Controllers: types.Controllers{
			Whitelist: types.AccountRange{
				Roles:     []uint64{1},
				Addresses: []string{"kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r"},
			},
		},
		Bin: []types.BinaryInfo{
			{
				Name:      "dapp1",
				Hash:      "dapp1",
				Source:    "dapp1",
				Reference: "dapp1",
				Type:      "dapp1",
			},
		},
		Pool: types.LpPoolConfig{
			Ratio:   sdk.OneDec(),
			Deposit: "",
			Drip:    86400,
		},
		Issuance: types.IssuanceConfig{
			Deposit:  "",
			Premint:  sdk.OneInt(),
			Postmint: sdk.OneInt(),
			Time:     1680141605,
		},
		VoteQuorum:    sdk.NewDecWithPrec(30, 2),
		VotePeriod:    86400,
		VoteEnactment: 3000,
		UpdateTimeMax: 60,
		ExecutorsMin:  1,
		ExecutorsMax:  2,
		VerifiersMin:  1,
		TotalBond:     sdk.NewInt64Coin("ukex", 10000),
		CreationTime:  0,
		PoolFee:       sdk.NewDec(1),
	}

	suite.app.Layer2Keeper.SetDapp(suite.ctx, dapp)

	lpToken := dapp.LpToken()
	lpCoins := sdk.Coins{sdk.NewInt64Coin(lpToken, 10000)}
	err := suite.app.TokensKeeper.MintCoins(suite.ctx, minttypes.ModuleName, lpCoins)
	suite.Require().NoError(err)

	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr, lpCoins)
	suite.Require().NoError(err)

	err = suite.app.TokensKeeper.MintCoins(suite.ctx, types.ModuleName, sdk.Coins{dapp.TotalBond})
	suite.Require().NoError(err)

	price := suite.app.Layer2Keeper.LpTokenPrice(suite.ctx, dapp)
	suite.Require().Equal(price, sdk.OneDec())

	out, err := suite.app.Layer2Keeper.RedeemDappPoolTx(suite.ctx, addr, dapp, sdk.NewDecWithPrec(1, 1), sdk.NewInt64Coin(lpToken, 1000))
	suite.Require().NoError(err)
	suite.Require().Equal(out.String(), "819ukex")

	out, err = suite.app.Layer2Keeper.SwapDappPoolTx(suite.ctx, addr, dapp, sdk.NewDecWithPrec(1, 1), sdk.NewInt64Coin("ukex", 800))
	suite.Require().NoError(err)
	suite.Require().Equal(out.String(), "667lp/dapp1")
}

func (suite *KeeperTestSuite) TestConvertDappPoolTx() {
	dapps := []types.Dapp{
		{
			Name:        "dapp1",
			Denom:       "dapp1",
			Description: "dapp1 description",
			Status:      types.Active,
			Website:     "",
			Logo:        "",
			Social:      "",
			Docs:        "",
			Controllers: types.Controllers{
				Whitelist: types.AccountRange{
					Roles:     []uint64{1},
					Addresses: []string{"kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r"},
				},
			},
			Bin: []types.BinaryInfo{
				{
					Name:      "dapp1",
					Hash:      "dapp1",
					Source:    "dapp1",
					Reference: "dapp1",
					Type:      "dapp1",
				},
			},
			Pool: types.LpPoolConfig{
				Ratio:   sdk.OneDec(),
				Deposit: "",
				Drip:    86400,
			},
			Issuance: types.IssuanceConfig{
				Deposit:  "",
				Premint:  sdk.OneInt(),
				Postmint: sdk.OneInt(),
				Time:     1680141605,
			},
			VoteQuorum:    sdk.NewDecWithPrec(30, 2),
			VotePeriod:    86400,
			VoteEnactment: 3000,
			UpdateTimeMax: 60,
			ExecutorsMin:  1,
			ExecutorsMax:  2,
			VerifiersMin:  1,
			TotalBond:     sdk.NewInt64Coin("ukex", 10000),
			CreationTime:  0,
			PoolFee:       sdk.NewDec(1),
		},
		{
			Name:        "dapp2",
			Denom:       "dapp2",
			Description: "dapp2 description",
			Status:      types.Active,
			Website:     "",
			Logo:        "",
			Social:      "",
			Docs:        "",
			Controllers: types.Controllers{
				Whitelist: types.AccountRange{
					Roles:     []uint64{1},
					Addresses: []string{"kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r"},
				},
			},
			Bin: []types.BinaryInfo{
				{
					Name:      "dapp2",
					Hash:      "dapp2",
					Source:    "dapp2",
					Reference: "dapp2",
					Type:      "dapp2",
				},
			},
			Pool: types.LpPoolConfig{
				Ratio:   sdk.OneDec(),
				Deposit: "",
				Drip:    86400,
			},
			Issuance: types.IssuanceConfig{
				Deposit:  "",
				Premint:  sdk.OneInt(),
				Postmint: sdk.OneInt(),
				Time:     1680141605,
			},
			VoteQuorum:    sdk.NewDecWithPrec(30, 2),
			VotePeriod:    86400,
			VoteEnactment: 3000,
			UpdateTimeMax: 60,
			ExecutorsMin:  1,
			ExecutorsMax:  2,
			VerifiersMin:  1,
			TotalBond:     sdk.NewInt64Coin("ukex", 10000),
			CreationTime:  0,
			PoolFee:       sdk.NewDec(1),
		},
	}

	for _, dapp := range dapps {
		suite.app.Layer2Keeper.SetDapp(suite.ctx, dapp)
	}
	dapp1 := dapps[0]
	dapp2 := dapps[1]

	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	lpToken1 := dapp1.LpToken()
	lpCoins1 := sdk.Coins{sdk.NewInt64Coin(lpToken1, 10000)}
	err := suite.app.TokensKeeper.MintCoins(suite.ctx, minttypes.ModuleName, lpCoins1)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr, lpCoins1)
	suite.Require().NoError(err)

	lpToken2 := dapp2.LpToken()
	lpCoins2 := sdk.Coins{sdk.NewInt64Coin(lpToken2, 10000)}
	err = suite.app.TokensKeeper.MintCoins(suite.ctx, types.ModuleName, lpCoins2)
	suite.Require().NoError(err)

	err = suite.app.TokensKeeper.MintCoins(suite.ctx, types.ModuleName, sdk.Coins{dapp1.TotalBond})
	suite.Require().NoError(err)

	price := suite.app.Layer2Keeper.LpTokenPrice(suite.ctx, dapp1)
	suite.Require().Equal(price, sdk.OneDec())

	out, err := suite.app.Layer2Keeper.ConvertDappPoolTx(suite.ctx, addr, dapp1, dapp2, sdk.NewInt64Coin(lpToken1, 1000))
	suite.Require().NoError(err)
	suite.Require().Equal(out.String(), "218lp/dapp2")
}
