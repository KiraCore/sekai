package keeper_test

import (
	"time"

	"github.com/KiraCore/sekai/x/spending/keeper"
	"github.com/KiraCore/sekai/x/spending/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (suite *KeeperTestSuite) TestEndBlocker() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	coins := sdk.Coins{sdk.NewInt64Coin("ukex", 1000000)}
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr1, coins)
	suite.Require().NoError(err)

	// create collective
	msgServer := keeper.NewMsgServerImpl(suite.app.SpendingKeeper, suite.app.CustomGovKeeper, suite.app.BankKeeper)
	msg := types.NewMsgCreateSpendingPool(
		"spendingpool1",
		0, 0,
		sdk.NewDecCoins(sdk.NewDecCoin("ukex", sdk.NewInt(1))),
		sdk.NewDecWithPrec(30, 2),
		86400, 3000,
		types.PermInfo{
			OwnerRoles:    []uint64{1},
			OwnerAccounts: []string{addr1.String()},
		},
		types.WeightedPermInfo{
			Roles: []types.WeightedRole{
				{
					Role:   1,
					Weight: sdk.NewDec(1),
				},
			},
			Accounts: []types.WeightedAccount{
				{
					Account: addr1.String(),
					Weight:  sdk.NewDec(2),
				},
			},
		},
		addr1,
		true,
		43200,
	)

	_, err = msgServer.CreateSpendingPool(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)
	pool := suite.app.SpendingKeeper.GetSpendingPool(suite.ctx, msg.Name)
	suite.Require().NotNil(pool)

	registerClaimInfo := types.NewMsgRegisterSpendingPoolBeneficiary(
		"spendingpool1",
		addr1,
	)

	_, err = msgServer.RegisterSpendingPoolBeneficiary(sdk.WrapSDKContext(suite.ctx), registerClaimInfo)
	suite.Require().NoError(err)

	// allocate reward to spending pool
	allocation := sdk.Coins{sdk.NewInt64Coin("ukex", 1000000)}
	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, allocation)
	suite.Require().NoError(err)
	suite.app.SpendingKeeper.DepositSpendingPoolFromModule(suite.ctx, minttypes.ModuleName, msg.Name, allocation)

	// run endblocker
	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Second * 100000))
	suite.app.SpendingKeeper.EndBlocker(suite.ctx)
	suite.Require().NoError(err)

	// check rate change and last dynamic rate calc time
	pool = suite.app.SpendingKeeper.GetSpendingPool(suite.ctx, msg.Name)
	suite.Require().NotNil(pool)
	suite.Require().Equal(sdk.DecCoins(pool.Rates).String(), "11.574074074074074074ukex")
	suite.Require().Equal(pool.LastDynamicRateCalcTime, uint64(suite.ctx.BlockTime().Unix()))
}
