package keeper_test

import (
	"github.com/KiraCore/sekai/x/layer2/keeper"
	"github.com/KiraCore/sekai/x/layer2/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (suite *KeeperTestSuite) TestCreateDappProposal() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	valAddr := sdk.ValAddress(addr1)
	pubkeys := simtestutil.CreateTestPubKeys(1)
	pubKey := pubkeys[0]
	val, err := stakingtypes.NewValidator(valAddr, pubKey)
	suite.Require().NoError(err)

	val.Status = stakingtypes.Active
	suite.app.CustomStakingKeeper.AddValidator(suite.ctx, val)

	coins := sdk.Coins{sdk.NewInt64Coin("ukex", 10000000000)}
	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr1, coins)
	suite.Require().NoError(err)

	// create dapp
	msgServer := keeper.NewMsgServerImpl(suite.app.Layer2Keeper)
	msg := &types.MsgCreateDappProposal{
		Sender: addr1.String(),
		Dapp: types.Dapp{
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
			TotalBond:     sdk.Coin{},
			CreationTime:  0,
		},
		Bond: sdk.NewInt64Coin("ukex", 10000000000),
	}

	_, err = msgServer.CreateDappProposal(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)

	// check dapp correctly crteated
	dapp := suite.app.Layer2Keeper.GetDapp(suite.ctx, msg.Dapp.Name)
	suite.Require().NotEqual(dapp.Name, "")
}

func (suite *KeeperTestSuite) TestBondDappProposal() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	valAddr := sdk.ValAddress(addr1)
	pubkeys := simtestutil.CreateTestPubKeys(1)
	pubKey := pubkeys[0]
	val, err := stakingtypes.NewValidator(valAddr, pubKey)
	suite.Require().NoError(err)

	val.Status = stakingtypes.Active
	suite.app.CustomStakingKeeper.AddValidator(suite.ctx, val)

	coins := sdk.Coins{sdk.NewInt64Coin("ukex", 20000000000)}
	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr1, coins)
	suite.Require().NoError(err)

	// create dapp
	msgServer := keeper.NewMsgServerImpl(suite.app.Layer2Keeper)
	createMsg := &types.MsgCreateDappProposal{
		Sender: addr1.String(),
		Dapp: types.Dapp{
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
			TotalBond:     sdk.Coin{},
			CreationTime:  0,
		},
		Bond: sdk.NewInt64Coin("ukex", 10000000000),
	}

	_, err = msgServer.CreateDappProposal(sdk.WrapSDKContext(suite.ctx), createMsg)
	suite.Require().NoError(err)

	// bond dapp
	bondMsg := &types.MsgBondDappProposal{
		Sender:   addr1.String(),
		DappName: createMsg.Dapp.Name,
		Bond:     sdk.NewInt64Coin("ukex", 10000000000),
	}

	_, err = msgServer.BondDappProposal(sdk.WrapSDKContext(suite.ctx), bondMsg)
	suite.Require().NoError(err)

	userBond := suite.app.Layer2Keeper.GetUserDappBond(suite.ctx, bondMsg.DappName, addr1.String())
	suite.Require().Equal(userBond.Bond.String(), "20000000000ukex")
}

func (suite *KeeperTestSuite) TestReclaimDappBondProposal() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	valAddr := sdk.ValAddress(addr1)
	pubkeys := simtestutil.CreateTestPubKeys(1)
	pubKey := pubkeys[0]
	val, err := stakingtypes.NewValidator(valAddr, pubKey)
	suite.Require().NoError(err)

	val.Status = stakingtypes.Active
	suite.app.CustomStakingKeeper.AddValidator(suite.ctx, val)

	coins := sdk.Coins{sdk.NewInt64Coin("ukex", 20000000000)}
	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr1, coins)
	suite.Require().NoError(err)

	// create dapp
	msgServer := keeper.NewMsgServerImpl(suite.app.Layer2Keeper)
	createMsg := &types.MsgCreateDappProposal{
		Sender: addr1.String(),
		Dapp: types.Dapp{
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
			TotalBond:     sdk.Coin{},
			CreationTime:  0,
		},
		Bond: sdk.NewInt64Coin("ukex", 10000000000),
	}

	_, err = msgServer.CreateDappProposal(sdk.WrapSDKContext(suite.ctx), createMsg)
	suite.Require().NoError(err)

	// bond dapp
	reclaimMsg := &types.MsgReclaimDappBondProposal{
		Sender:   addr1.String(),
		DappName: createMsg.Dapp.Name,
		Bond:     sdk.NewInt64Coin("ukex", 10000000),
	}

	_, err = msgServer.ReclaimDappBondProposal(sdk.WrapSDKContext(suite.ctx), reclaimMsg)
	suite.Require().NoError(err)

	userBond := suite.app.Layer2Keeper.GetUserDappBond(suite.ctx, reclaimMsg.DappName, addr1.String())
	suite.Require().Equal(userBond.Bond.String(), "9990000000ukex")
}

// TODO: add test for JoinDappVerifierWithBond
// TODO: add test for ExitDapp
// TODO: add test for PauseDappTx
// TODO: add test for UnPauseDappTx
// TODO: add test for ReactivateDappTx
// TODO: add test for ExecuteDappTx
// TODO: add test for TransitionDappTx
// TODO: add test for DenounceLeaderTx
// TODO: add test for ApproveDappTransitionTx
// TODO: add test for RejectDappTransitionTx
// TODO: add test for TransferDappTx
// TODO: add test for RedeemDappPoolTx
// TODO: add test for SwapDappPoolTx
// TODO: add test for ConvertDappPoolTx
// TODO: add test for MintCreateFtTx
// TODO: add test for MintCreateNftTx
// TODO: add test for MintIssueTx
// TODO: add test for MintBurnTx
