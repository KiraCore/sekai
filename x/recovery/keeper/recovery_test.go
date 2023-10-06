package keeper_test

import (
	"github.com/KiraCore/sekai/x/recovery/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (suite *KeeperTestSuite) TestRecoveryTokenSetGetDelete() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	records := []types.RecoveryToken{
		{
			Address:          addr1.String(),
			Token:            "rr_moniker1",
			RrSupply:         sdk.NewInt(10000),
			UnderlyingTokens: sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000000)),
		},
		{
			Address:          addr2.String(),
			Token:            "rr_moniker2",
			RrSupply:         sdk.NewInt(20000),
			UnderlyingTokens: sdk.NewCoins(sdk.NewInt64Coin("ukex", 2000000)),
		},
	}

	for _, record := range records {
		suite.app.RecoveryKeeper.SetRecoveryToken(suite.ctx, record)
	}

	for _, record := range records {
		r, err := suite.app.RecoveryKeeper.GetRecoveryToken(suite.ctx, record.Address)
		suite.Require().NoError(err)
		suite.Require().Equal(r, record)
	}

	allRecords := suite.app.RecoveryKeeper.GetAllRecoveryTokens(suite.ctx)
	suite.Require().Len(allRecords, 2)

	recoveryToken, err := suite.app.RecoveryKeeper.GetRecoveryTokenByDenom(suite.ctx, records[0].Token)
	suite.Require().NoError(err)
	suite.Require().Equal(recoveryToken, records[0])

	err = suite.app.RecoveryKeeper.IncreaseRecoveryTokenUnderlying(suite.ctx, addr1, sdk.NewCoins(sdk.NewInt64Coin("reward", 10000)))
	suite.Require().NoError(err)
	recoveryToken, err = suite.app.RecoveryKeeper.GetRecoveryToken(suite.ctx, records[0].Address)
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.Coins(recoveryToken.UnderlyingTokens).String(), "10000reward,1000000ukex")

	suite.app.RecoveryKeeper.DeleteRecoveryToken(suite.ctx, records[0])

	allRecords = suite.app.RecoveryKeeper.GetAllRecoveryTokens(suite.ctx)
	suite.Require().Len(allRecords, 1)

	_, err = suite.app.RecoveryKeeper.GetRecoveryToken(suite.ctx, records[0].Address)
	suite.Require().Error(err)

	_, err = suite.app.RecoveryKeeper.GetRecoveryTokenByDenom(suite.ctx, records[0].Token)
	suite.Require().Error(err)
}
