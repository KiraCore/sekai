package keeper_test

import (
	"testing"

	simapp "github.com/KiraCore/sekai/app"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *simapp.SekaiApp
}

func (suite *KeeperTestSuite) SetupTest() {
	checkTx := false
	suite.app = simapp.Setup(checkTx)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, tmproto.Header{Height: 1})
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
