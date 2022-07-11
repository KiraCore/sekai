package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func (suite *KeeperTestSuite) TestAllocateTokens() {
	suite.SetupTest()
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	consAddr := sdk.ConsAddress(addr1)
	coins := sdk.Coins{sdk.NewInt64Coin("ukex", 10000000)}
	suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr1, coins)
	suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr2, coins)

	oldTreasury := suite.app.DistrKeeper.GetFeesTreasury(suite.ctx)
	suite.app.DistrKeeper.AllocateTokens(suite.ctx, 10, 10, consAddr, []abci.VoteInfo{})
	newTreasury := suite.app.DistrKeeper.GetFeesTreasury(suite.ctx)
	suite.Require().True(oldTreasury.DenomsSubsetOf(newTreasury))

	// TODO: add case for validator exit case
	// TODO: add case for staking pool exist case
	// TODO: check tokens are distributed correctly
}
