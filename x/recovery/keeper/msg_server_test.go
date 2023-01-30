package keeper_test

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/KiraCore/sekai/app"
	collectivestypes "github.com/KiraCore/sekai/x/collectives/types"
	custodytypes "github.com/KiraCore/sekai/x/custody/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/recovery/keeper"
	"github.com/KiraCore/sekai/x/recovery/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func init() {
	app.SetConfig()
}

func (suite *KeeperTestSuite) TestRegisterRecoverySecret() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	// create recovery record
	msgServer := keeper.NewMsgServerImpl(suite.app.RecoveryKeeper)
	msg := types.NewMsgRegisterRecoverySecret(
		addr1.String(), "123456", "111111", "",
	)

	_, err := msgServer.RegisterRecoverySecret(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)

	// check recovery correctly crteated
	record, err := suite.app.RecoveryKeeper.GetRecoveryRecord(suite.ctx, addr1.String())
	suite.Require().NoError(err)
	suite.Require().Equal(record, types.RecoveryRecord{
		Address:   addr1.String(),
		Challenge: "123456",
		Nonce:     "111111",
	})

	// try another execution without proof
	_, err = msgServer.RegisterRecoverySecret(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().Error(err)
}

func (suite *KeeperTestSuite) TestRotateRecoveryAddress() {
	pubkey1 := secp256k1.GenPrivKey().PubKey()
	pubkey2 := secp256k1.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pubkey1.Address())
	addr2 := sdk.AccAddress(pubkey2.Address())

	acc1 := authtypes.NewBaseAccount(addr1, pubkey1, 0, 0)
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc1)

	// bank module setup
	coins := sdk.Coins{sdk.NewInt64Coin("ukex", 1000_000)}
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr1, coins)
	suite.Require().NoError(err)

	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, keeper.RecoveryFee)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr1, keeper.RecoveryFee)
	suite.Require().NoError(err)

	// collectives module setup
	contributer := collectivestypes.CollectiveContributor{
		Address:      addr1.String(),
		Name:         "collective1",
		Bonds:        []sdk.Coin{sdk.NewInt64Coin("ukex", 1000_000)},
		Locking:      1000,
		Donation:     sdk.NewDecWithPrec(1, 1), // 10%
		DonationLock: true,
	}
	suite.app.CollectivesKeeper.SetCollectiveContributer(suite.ctx, contributer)

	// gov module setup
	councilor := govtypes.NewCouncilor(
		addr1,
		govtypes.CouncilorActive,
	)
	suite.app.CustomGovKeeper.SaveCouncilor(suite.ctx, councilor)

	// multistaking module
	suite.app.MultiStakingKeeper.SetDelegatorRewards(suite.ctx, addr1, coins)

	// staking module
	valAddr := sdk.ValAddress(addr1)
	pubkeys := simapp.CreateTestPubKeys(1)
	pubKey := pubkeys[0]
	val, err := stakingtypes.NewValidator(valAddr, pubKey)
	suite.Require().NoError(err)
	val.Status = stakingtypes.Active
	suite.app.CustomStakingKeeper.AddValidator(suite.ctx, val)

	// custody settings
	settings := custodytypes.CustodySettings{
		CustodyEnabled: true,
		CustodyMode:    1,
		UsePassword:    true,
		UseWhiteList:   true,
		UseLimits:      true,
		Key:            "key",
	}
	suite.app.CustodyKeeper.SetCustodyRecord(suite.ctx, custodytypes.CustodyRecord{
		Address:         addr1,
		CustodySettings: &settings,
	})

	// recovery record set
	privKey, err := hex.DecodeString("10a0fbe01030000122300000000000")
	suite.Require().NoError(err)
	proof := sha256.Sum256(privKey)
	challenge := sha256.Sum256(proof[:])

	suite.app.RecoveryKeeper.SetRecoveryRecord(suite.ctx, types.RecoveryRecord{
		Address:   addr1.String(),
		Challenge: hex.EncodeToString(challenge[:]),
		Nonce:     "111111",
	})

	// invalid proof
	msg := types.NewMsgRotateRecoveryAddress(
		addr1.String(), addr1.String(), addr2.String(), "",
	)

	msgServer := keeper.NewMsgServerImpl(suite.app.RecoveryKeeper)
	cacheCtx, _ := suite.ctx.CacheContext()
	_, err = msgServer.RotateRecoveryAddress(sdk.WrapSDKContext(cacheCtx), msg)
	suite.Require().Error(err)

	// valid proof
	msg.Proof = hex.EncodeToString(proof[:])
	_, err = msgServer.RotateRecoveryAddress(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)

	// check bank module transfer
	addr1Coins := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1)
	suite.Require().Equal(addr1Coins, sdk.Coins{})
	addr2Coins := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr2)
	suite.Require().Equal(addr2Coins, coins)

	// check collective contributer recovery
	addr1Contributer := suite.app.CollectivesKeeper.GetCollectiveContributer(suite.ctx, "collective1", addr1.String())
	suite.Require().Equal(addr1Contributer, collectivestypes.CollectiveContributor{})
	addr2Contributer := suite.app.CollectivesKeeper.GetCollectiveContributer(suite.ctx, "collective1", addr2.String())
	suite.Require().Equal(addr2Contributer.Address, addr2.String())
	suite.Require().Equal(addr2Contributer.Bonds, contributer.Bonds)

	// check gov: councilor recovery
	_, found := suite.app.CustomGovKeeper.GetCouncilor(suite.ctx, addr1)
	suite.Require().False(found)
	addr2Councilor, found := suite.app.CustomGovKeeper.GetCouncilor(suite.ctx, addr2)
	suite.Require().True(found)
	suite.Require().Equal(addr2Councilor.Status, councilor.Status)

	// TODO: check gov: identity records recovery
	// TODO: check gov: identity records verification requests recovery
	// TODO: check gov: network actor recovery
	// TODO: check gov: votes recovery

	// check multistaking delegator rewards recovery
	addr1Rewards := suite.app.MultiStakingKeeper.GetDelegatorRewards(suite.ctx, addr1)
	suite.Require().Equal(addr1Rewards, sdk.Coins{})
	addr2Rewards := suite.app.MultiStakingKeeper.GetDelegatorRewards(suite.ctx, addr2)
	suite.Require().Equal(addr2Rewards, coins)

	// TODO: check multistaking is delegator recovery
	// TODO: check multistaking validator address recovery
	// TODO: check spending pool claim info recovery

	// check validator address recovery
	_, err = suite.app.CustomStakingKeeper.GetValidator(suite.ctx, sdk.ValAddress(addr1))
	suite.Require().Error(err)
	val2, err := suite.app.CustomStakingKeeper.GetValidator(suite.ctx, sdk.ValAddress(addr2))
	suite.Require().NoError(err)
	suite.Require().Equal(val2.Status, val.Status)

	// check custody settings recovery
	addr1Custody := suite.app.CustodyKeeper.GetCustodyInfoByAddress(suite.ctx, addr1)
	suite.Require().Nil(addr1Custody)
	addr2Custody := suite.app.CustodyKeeper.GetCustodyInfoByAddress(suite.ctx, addr2)
	suite.Require().NotNil(addr2Custody)
	suite.Require().Equal(*addr2Custody, settings)
}
