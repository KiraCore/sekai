package keeper_test

import (
	"github.com/KiraCore/sekai/x/layer2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestDappOperatorSetGetDelete() {
	operators := []types.DappOperator{
		{
			DappName:       "dapp1",
			Operator:       "kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r",
			Executor:       true,
			Verifier:       true,
			Interx:         "kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r",
			Status:         types.OperatorActive,
			Rank:           1,
			Streak:         1,
			BondedLpAmount: sdk.ZeroInt(),
		},
		{
			DappName:       "dapp1",
			Operator:       "kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f",
			Executor:       true,
			Verifier:       false,
			Interx:         "kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f",
			Status:         types.OperatorActive,
			Rank:           1,
			Streak:         1,
			BondedLpAmount: sdk.OneInt(),
		},
		{
			DappName:       "dapp2",
			Operator:       "kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f",
			Executor:       false,
			Verifier:       true,
			Interx:         "kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f",
			Status:         types.OperatorActive,
			Rank:           1,
			Streak:         1,
			BondedLpAmount: sdk.OneInt(),
		},
	}

	for _, operator := range operators {
		suite.app.Layer2Keeper.SetDappOperator(suite.ctx, operator)
	}

	for _, operator := range operators {
		c := suite.app.Layer2Keeper.GetDappOperator(suite.ctx, operator.DappName, operator.Operator)
		suite.Require().Equal(c, operator)
	}

	allOperators := suite.app.Layer2Keeper.GetDappOperators(suite.ctx, operators[0].DappName)
	suite.Require().Len(allOperators, 2)

	allOperators = suite.app.Layer2Keeper.GetAllDappOperators(suite.ctx)
	suite.Require().Len(allOperators, 3)

	allVerifiers := suite.app.Layer2Keeper.GetDappVerifiers(suite.ctx, operators[0].DappName)
	suite.Require().Len(allVerifiers, 1)

	allExecutors := suite.app.Layer2Keeper.GetDappExecutors(suite.ctx, operators[0].DappName)
	suite.Require().Len(allExecutors, 2)

	suite.app.Layer2Keeper.DeleteDappOperator(suite.ctx, operators[0].DappName, operators[0].Operator)

	allOperators = suite.app.Layer2Keeper.GetDappOperators(suite.ctx, operators[0].DappName)
	suite.Require().Len(allOperators, 1)

	operator := suite.app.Layer2Keeper.GetDappOperator(suite.ctx, operators[0].DappName, operators[0].Operator)
	suite.Require().Equal(operator.DappName, "")

	allOperators = suite.app.Layer2Keeper.GetAllDappOperators(suite.ctx)
	suite.Require().Len(allOperators, 2)
}

// TODO: add test for ExecuteJoinDappProposal
