package keeper_test

import (
	"github.com/KiraCore/sekai/x/layer2/types"
)

func (suite *KeeperTestSuite) TestDappSessionSetGetDelete() {
	registrars := []types.ExecutionRegistrar{
		{
			DappName:    "dapp1",
			PrevSession: nil,
			CurrSession: nil,
			NextSession: nil,
		},
		{
			DappName:    "dapp2",
			PrevSession: nil,
			CurrSession: nil,
			NextSession: &types.DappSession{
				Leader:     "kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r",
				Start:      1680141605,
				StatusHash: "0c",
				Status:     types.SessionUnscheduled,
				Gateway:    "dapp2.com",
			},
		},
		{
			DappName:    "dapp3",
			PrevSession: nil,
			CurrSession: &types.DappSession{
				Leader:     "kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r",
				Start:      1680141605,
				StatusHash: "0c",
				Status:     types.SessionOngoing,
				Gateway:    "dapp3.com",
			},
			NextSession: &types.DappSession{
				Leader:     "kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r",
				Start:      1680141605,
				StatusHash: "0c",
				Status:     types.SessionUnscheduled,
				Gateway:    "dapp3.com",
			},
		},
		{
			DappName: "dapp4",
			PrevSession: &types.DappSession{
				Leader:     "kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r",
				Start:      1680141605,
				StatusHash: "0c",
				Status:     types.SessionAccepted,
				Gateway:    "dapp4.com",
			},
			CurrSession: &types.DappSession{
				Leader:     "kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r",
				Start:      1680141605,
				StatusHash: "0c",
				Status:     types.SessionOngoing,
				Gateway:    "dapp4.com",
			},
			NextSession: &types.DappSession{
				Leader:     "kira1q24436yrnettd6v4eu6r4t9gycnnddack4jr5r",
				Start:      1680141605,
				StatusHash: "0c",
				Status:     types.SessionUnscheduled,
				Gateway:    "dapp4.com",
			},
		},
	}

	for _, registrar := range registrars {
		suite.app.Layer2Keeper.SetDappSession(suite.ctx, registrar)
	}

	for _, registrar := range registrars {
		c := suite.app.Layer2Keeper.GetDappSession(suite.ctx, registrar.DappName)
		suite.Require().Equal(c, registrar)
	}

	allRegistrars := suite.app.Layer2Keeper.GetAllDappSessions(suite.ctx)
	suite.Require().Len(allRegistrars, 4)

	suite.app.Layer2Keeper.DeleteDappSession(suite.ctx, registrars[0].DappName)

	registrar := suite.app.Layer2Keeper.GetDappSession(suite.ctx, registrars[0].DappName)
	suite.Require().Equal(registrar.DappName, "")

	allRegistrars = suite.app.Layer2Keeper.GetAllDappSessions(suite.ctx)
	suite.Require().Len(allRegistrars, 3)
}

// TODO: add test for ResetNewSession
// TODO: add test for CreateNewSession
