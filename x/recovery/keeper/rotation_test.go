package keeper_test

import (
	"github.com/KiraCore/sekai/x/recovery/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestRotationHistorySetGetDelete() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr3 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	records := []types.Rotation{
		{
			Address: addr1.String(),
			Rotated: addr2.String(),
		},
		{
			Address: addr2.String(),
			Rotated: addr3.String(),
		},
	}

	for _, record := range records {
		suite.app.RecoveryKeeper.SetRotationHistory(suite.ctx, record)
	}

	for _, record := range records {
		r := suite.app.RecoveryKeeper.GetRotationHistory(suite.ctx, record.Address)
		suite.Require().Equal(r, record)
	}

	allRecords := suite.app.RecoveryKeeper.GetAllRotationHistory(suite.ctx)
	suite.Require().Len(allRecords, 2)

	suite.app.RecoveryKeeper.DeleteRotationHistory(suite.ctx, records[0])

	allRecords = suite.app.RecoveryKeeper.GetAllRotationHistory(suite.ctx)
	suite.Require().Len(allRecords, 1)
}
