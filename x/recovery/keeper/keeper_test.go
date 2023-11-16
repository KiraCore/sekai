package keeper_test

import (
	"testing"

	simapp "github.com/KiraCore/sekai/app"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *simapp.SekaiApp
}

func (suite *KeeperTestSuite) SetupTest() {
	checkTx := false
	app := simapp.Setup(checkTx)

	suite.ctx = app.BaseApp.NewContext(checkTx, tmproto.Header{Height: 1})
	suite.app = app
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
