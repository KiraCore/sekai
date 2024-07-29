package keeper_test

import (
	"github.com/KiraCore/sekai/x/layer2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestDappSetGetDelete() {
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

	for _, dapp := range dapps {
		c := suite.app.Layer2Keeper.GetDapp(suite.ctx, dapp.Name)
		suite.Require().Equal(c, dapp)
	}

	allCollectives := suite.app.Layer2Keeper.GetAllDapps(suite.ctx)
	suite.Require().Len(allCollectives, 2)

	suite.app.Layer2Keeper.DeleteDapp(suite.ctx, dapps[0].Name)

	allCollectives = suite.app.Layer2Keeper.GetAllDapps(suite.ctx)
	suite.Require().Len(allCollectives, 1)

	dapp := suite.app.Layer2Keeper.GetDapp(suite.ctx, dapps[0].Name)
	suite.Require().Equal(dapp.Name, "")
}

func (suite *KeeperTestSuite) TestUserDappBondSetGetDelete() {
	userBonds := []types.UserDappBond{
		{
			User:     "kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r",
			DappName: "dapp1",
			Bond:     sdk.NewInt64Coin("ukex", 1000_000),
		},
		{
			User:     "kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r",
			DappName: "dapp2",
			Bond:     sdk.NewInt64Coin("ukex", 1000_000),
		},
	}

	for _, userBond := range userBonds {
		suite.app.Layer2Keeper.SetUserDappBond(suite.ctx, userBond)
	}

	for _, userBond := range userBonds {
		c := suite.app.Layer2Keeper.GetUserDappBond(suite.ctx, userBond.DappName, userBond.User)
		suite.Require().Equal(c, userBond)
	}

	allBonds := suite.app.Layer2Keeper.GetUserDappBonds(suite.ctx, userBonds[0].DappName)
	suite.Require().Len(allBonds, 1)

	allBonds = suite.app.Layer2Keeper.GetAllUserDappBonds(suite.ctx)
	suite.Require().Len(allBonds, 2)

	suite.app.Layer2Keeper.DeleteUserDappBond(suite.ctx, userBonds[0].DappName, userBonds[0].User)

	allBonds = suite.app.Layer2Keeper.GetUserDappBonds(suite.ctx, userBonds[0].DappName)
	suite.Require().Len(allBonds, 0)

	userBond := suite.app.Layer2Keeper.GetUserDappBond(suite.ctx, userBonds[0].DappName, userBonds[0].User)
	suite.Require().Equal(userBond.DappName, "")

	allBonds = suite.app.Layer2Keeper.GetAllUserDappBonds(suite.ctx)
	suite.Require().Len(allBonds, 1)
}

// TODO: add test for ExecuteDappRemove
