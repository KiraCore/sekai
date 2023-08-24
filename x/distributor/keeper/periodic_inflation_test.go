package keeper_test

import (
	"time"

	"github.com/KiraCore/sekai/x/distributor/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestPeriodicSnapshot() {
	ctx := suite.ctx
	now := time.Now().UTC()
	snapshot := suite.app.DistrKeeper.GetPeriodicSnapshot(ctx)
	suite.Require().Equal(snapshot.SnapshotTime, int64(0))

	newSnapshot := types.SupplySnapshot{
		SnapshotTime:   now.Unix(),
		SnapshotAmount: sdk.NewInt(1000000),
	}
	suite.app.DistrKeeper.SetPeriodicSnapshot(ctx, newSnapshot)
	snapshot = suite.app.DistrKeeper.GetPeriodicSnapshot(ctx)
	suite.Require().Equal(snapshot, newSnapshot)
}
