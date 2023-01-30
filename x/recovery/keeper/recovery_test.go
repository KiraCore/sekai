package keeper_test

import (
	"github.com/KiraCore/sekai/x/recovery/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func (suite *KeeperTestSuite) TestRecoveryRecordSetGetDelete() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	records := []types.RecoveryRecord{
		{
			Address:   addr1.String(),
			Challenge: "12345678",
			Nonce:     "001100110022",
		},
		{
			Address:   addr2.String(),
			Challenge: "87654321",
			Nonce:     "001100110022",
		},
	}

	for _, record := range records {
		suite.app.RecoveryKeeper.SetRecoveryRecord(suite.ctx, record)
	}

	for _, record := range records {
		r, err := suite.app.RecoveryKeeper.GetRecoveryRecord(suite.ctx, record.Address)
		suite.Require().NoError(err)
		suite.Require().Equal(r, record)
	}

	allRecords := suite.app.RecoveryKeeper.GetAllRecoveryRecords(suite.ctx)
	suite.Require().Len(allRecords, 2)

	addr := suite.app.RecoveryKeeper.GetRecoveryAddressFromChallenge(suite.ctx, records[0].Challenge)
	suite.Require().Equal(addr, records[0].Address)

	suite.app.RecoveryKeeper.DeleteRecoveryRecord(suite.ctx, records[0])

	allRecords = suite.app.RecoveryKeeper.GetAllRecoveryRecords(suite.ctx)
	suite.Require().Len(allRecords, 1)

	_, err := suite.app.RecoveryKeeper.GetRecoveryRecord(suite.ctx, records[0].Address)
	suite.Require().Error(err)
}
