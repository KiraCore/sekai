package keeper_test

import (
	"time"

	"github.com/KiraCore/sekai/x/distributor/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (suite *KeeperTestSuite) TestYearStartSnapshot() {
	ctx := suite.ctx
	now := time.Now().UTC()
	snapshot := suite.app.DistrKeeper.GetYearStartSnapshot(ctx)
	suite.Require().Equal(snapshot.SnapshotTime, int64(0))

	newSnapshot := types.YearStartSnapshot{
		SnapshotTime:   now.Unix(),
		SnapshotAmount: sdk.NewInt(1000000),
	}
	suite.app.DistrKeeper.SetYearStartSnapshot(ctx, newSnapshot)
	snapshot = suite.app.DistrKeeper.GetYearStartSnapshot(ctx)
	suite.Require().Equal(snapshot, newSnapshot)

	ctx = ctx.WithBlockTime(now)
	inflationPossible := suite.app.DistrKeeper.InflationPossible(ctx)
	suite.Require().True(inflationPossible)

	err := suite.app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.Coins{sdk.NewInt64Coin("ukex", 2000000)})
	suite.Require().NoError(err)
	inflationPossible = suite.app.DistrKeeper.InflationPossible(ctx)
	suite.Require().False(inflationPossible)
}

// func (k Keeper) InflationPossible(ctx sdk.Context) bool {
// 	snapshot := k.GetYearStartSnapshot(ctx)
// 	if snapshot.SnapshotAmount.IsNil() || snapshot.SnapshotAmount.IsZero() {
// 		return true
// 	}
// 	yearlyInflation := k.gk.GetNetworkProperties(ctx).MaxAnnualInflation
// 	currSupply := k.bk.GetSupply(ctx, k.BondDenom(ctx))

// 	month := int64(86400 * 30)
// 	currTimeGone := ctx.BlockTime().Unix() - snapshot.SnapshotTime
// 	monthIndex := (currTimeGone + month - 1) / month
// 	currInflation := currSupply.Amount.ToDec().Quo(sdk.Dec(snapshot.SnapshotAmount.ToDec())).Sub(sdk.OneDec())
// 	if currInflation.GTE(yearlyInflation.Mul(sdk.NewDec(monthIndex)).Quo(sdk.NewDec(12))) {
// 		return false
// 	}
// 	return true
// }
