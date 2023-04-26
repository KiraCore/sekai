package keeper_test

import (
	"github.com/KiraCore/sekai/x/layer2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestBridgeRegistrarHelperSetGet() {
	helper := types.BridgeRegistrarHelper{
		NextUser:  1,
		NextXam:   1,
		NextToken: 1,
	}
	suite.app.Layer2Keeper.SetBridgeRegistrarHelper(suite.ctx, helper)

	saved := suite.app.Layer2Keeper.GetBridgeRegistrarHelper(suite.ctx)
	suite.Require().Equal(saved, helper)
}

func (suite *KeeperTestSuite) TestBridgeAccountSetGet() {
	infos := []types.BridgeAccount{
		{
			Index:    1,
			Address:  "kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d",
			DappName: "",
			Balances: []types.BridgeBalance(nil),
		},
		{
			Index:    2,
			Address:  "kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f",
			DappName: "dapp1",
			Balances: []types.BridgeBalance{
				{
					BridgeTokenIndex: 1,
					Amount:           sdk.OneInt(),
				},
			},
		},
	}

	for _, info := range infos {
		suite.app.Layer2Keeper.SetBridgeAccount(suite.ctx, info)
	}

	for _, info := range infos {
		c := suite.app.Layer2Keeper.GetBridgeAccount(suite.ctx, info.Index)
		suite.Require().Equal(c, info)
	}

	accounts := suite.app.Layer2Keeper.GetBridgeAccounts(suite.ctx)
	suite.Require().Len(accounts, 2)
}

func (suite *KeeperTestSuite) TestBridgeTokenSetGet() {
	infos := []types.BridgeToken{
		{
			Index: 1,
			Denom: "btc",
		},
		{
			Index: 2,
			Denom: "eth",
		},
	}

	for _, info := range infos {
		suite.app.Layer2Keeper.SetBridgeToken(suite.ctx, info)
	}

	for _, info := range infos {
		c := suite.app.Layer2Keeper.GetBridgeToken(suite.ctx, info.Index)
		suite.Require().Equal(c, info)
	}

	tokens := suite.app.Layer2Keeper.GetBridgeTokens(suite.ctx)
	suite.Require().Len(tokens, 2)
}

func (suite *KeeperTestSuite) TestXAMSetGet() {
	infos := []types.XAM{
		{
			Req: types.XAMRequest{
				Amounts:         []types.BridgeBalance(nil),
				SourceDapp:      1,
				SourceAccount:   1,
				DestDapp:        2,
				DestBeneficiary: 1,
				Xam:             "",
			},
			Res: types.XAMResponse{
				Xid: 1,
				Irc: 0,
				Src: 0,
				Drc: 0,
				Irm: 0,
				Srm: 0,
				Drm: 0,
			},
		},
		{
			Req: types.XAMRequest{
				Amounts: []types.BridgeBalance{
					{
						BridgeTokenIndex: 1,
						Amount:           sdk.OneInt(),
					},
				},
				SourceDapp:      1,
				SourceAccount:   1,
				DestDapp:        2,
				DestBeneficiary: 1,
				Xam:             "",
			},
			Res: types.XAMResponse{
				Xid: 2,
				Irc: 0,
				Src: 0,
				Drc: 0,
				Irm: 0,
				Srm: 0,
				Drm: 0,
			},
		},
	}

	for _, info := range infos {
		suite.app.Layer2Keeper.SetXAM(suite.ctx, info)
	}

	for _, info := range infos {
		c := suite.app.Layer2Keeper.GetXAM(suite.ctx, info.Res.Xid)
		suite.Require().Equal(c, info)
	}

	xams := suite.app.Layer2Keeper.GetXAMs(suite.ctx)
	suite.Require().Len(xams, 2)
}
