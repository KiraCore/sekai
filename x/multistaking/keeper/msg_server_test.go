package keeper_test

import (
	"time"

	multistakingkeeper "github.com/KiraCore/sekai/x/multistaking/keeper"
	multistakingtypes "github.com/KiraCore/sekai/x/multistaking/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (suite *KeeperTestSuite) TestClaimUndelegation() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	valAddr := sdk.ValAddress(addr1)
	pubkeys := simtestutil.CreateTestPubKeys(1)
	pubKey := pubkeys[0]
	val, err := stakingtypes.NewValidator(valAddr, pubKey)
	suite.Require().NoError(err)

	properties := suite.app.CustomGovKeeper.GetNetworkProperties(suite.ctx)
	properties.MinCollectiveBond = 1
	suite.app.CustomGovKeeper.SetNetworkProperties(suite.ctx, properties)

	val.Status = stakingtypes.Active
	suite.app.CustomStakingKeeper.AddValidator(suite.ctx, val)

	stakingPool := multistakingtypes.StakingPool{
		Id:        1,
		Validator: valAddr.String(),
		Enabled:   true,
	}

	suite.app.MultiStakingKeeper.SetStakingPool(suite.ctx, stakingPool)

	coins := sdk.Coins{sdk.NewInt64Coin("ukex", 1000000)}
	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr1, coins)
	suite.Require().NoError(err)
	msgServer := multistakingkeeper.NewMsgServerImpl(suite.app.MultiStakingKeeper, suite.app.BankKeeper, suite.app.CustomGovKeeper, suite.app.CustomStakingKeeper)
	_, err = msgServer.Delegate(sdk.WrapSDKContext(suite.ctx), &multistakingtypes.MsgDelegate{
		DelegatorAddress: addr1.String(),
		ValidatorAddress: valAddr.String(),
		Amounts:          coins,
	})
	suite.Require().NoError(err)

	_, err = msgServer.Undelegate(sdk.WrapSDKContext(suite.ctx), &multistakingtypes.MsgUndelegate{
		DelegatorAddress: addr1.String(),
		ValidatorAddress: valAddr.String(),
		Amounts:          coins,
	})
	suite.Require().NoError(err)

	undelegations := suite.app.MultiStakingKeeper.GetAllUndelegations(suite.ctx)
	suite.Require().Len(undelegations, 1)
	suite.Require().Equal(undelegations[0], multistakingtypes.Undelegation{
		Id:         1,
		Address:    addr1.String(),
		ValAddress: valAddr.String(),
		Expiry:     uint64(suite.ctx.BlockTime().Unix()) + properties.UnstakingPeriod,
		Amount:     coins,
	})

	// claim before expiry
	_, err = msgServer.ClaimUndelegation(sdk.WrapSDKContext(suite.ctx), &multistakingtypes.MsgClaimUndelegation{
		Sender:         addr1.String(),
		UndelegationId: 1,
	})
	suite.Require().Error(err)

	// claim correctly
	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Second * time.Duration(properties.UnstakingPeriod)))
	_, err = msgServer.ClaimUndelegation(sdk.WrapSDKContext(suite.ctx), &multistakingtypes.MsgClaimUndelegation{
		Sender:         addr1.String(),
		UndelegationId: 1,
	})
	suite.Require().NoError(err)

	// Try claim again
	_, err = msgServer.ClaimUndelegation(sdk.WrapSDKContext(suite.ctx), &multistakingtypes.MsgClaimUndelegation{
		Sender:         addr1.String(),
		UndelegationId: 1,
	})
	suite.Require().Error(err)
}
