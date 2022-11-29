package keeper_test

func (suite *KeeperTestSuite) TestProcessUBIRecordDynamic() {
	ctx := suite.ctx

	// check genesis ubi records
	records := suite.app.UbiKeeper.GetUBIRecords(ctx)
	suite.Require().Len(records, 1)
	suite.Require().True(records[0].Dynamic)

	// process it and check
	suite.app.UbiKeeper.ProcessUBIRecord(ctx, records[0])
	pool1 := suite.app.SpendingKeeper.GetSpendingPool(ctx, records[0].Pool)

	// process again and check
	suite.app.UbiKeeper.ProcessUBIRecord(ctx, records[0])
	pool2 := suite.app.SpendingKeeper.GetSpendingPool(ctx, records[0].Pool)
	suite.Require().Equal(pool1.Balances, pool2.Balances)

	// update to not dynamic and check
	records[0].Dynamic = false
	suite.app.UbiKeeper.ProcessUBIRecord(ctx, records[0])
	pool3 := suite.app.SpendingKeeper.GetSpendingPool(ctx, records[0].Pool)
	suite.Require().NotEqual(pool1.Balances, pool3.Balances)
}
