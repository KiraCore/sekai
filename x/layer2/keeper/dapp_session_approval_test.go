package keeper_test

import (
	"github.com/KiraCore/sekai/x/layer2/types"
)

func (suite *KeeperTestSuite) TestDappSessionApprovalSetGetDelete() {
	approvals := []types.DappSessionApproval{
		{
			Approver:   "kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r",
			DappName:   "dapp1",
			IsApproved: true,
		},
		{
			Approver:   "kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r",
			DappName:   "dapp2",
			IsApproved: false,
		},
	}

	for _, approval := range approvals {
		suite.app.Layer2Keeper.SetDappSessionApproval(suite.ctx, approval)
	}

	for _, approval := range approvals {
		c := suite.app.Layer2Keeper.GetDappSessionApproval(suite.ctx, approval.DappName, approval.Approver)
		suite.Require().Equal(c, approval)
	}

	allApprovals := suite.app.Layer2Keeper.GetDappSessionApprovals(suite.ctx, approvals[0].DappName)
	suite.Require().Len(allApprovals, 1)

	allApprovals = suite.app.Layer2Keeper.GetAllDappSessionApprovals(suite.ctx)
	suite.Require().Len(allApprovals, 2)

	suite.app.Layer2Keeper.DeleteDappSessionApproval(suite.ctx, approvals[0].DappName, approvals[0].Approver)

	allApprovals = suite.app.Layer2Keeper.GetDappSessionApprovals(suite.ctx, approvals[0].DappName)
	suite.Require().Len(allApprovals, 0)

	approval := suite.app.Layer2Keeper.GetDappSessionApproval(suite.ctx, approvals[0].DappName, approvals[0].Approver)
	suite.Require().Equal(approval.DappName, "")

	allApprovals = suite.app.Layer2Keeper.GetAllDappSessionApprovals(suite.ctx)
	suite.Require().Len(allApprovals, 1)
}
