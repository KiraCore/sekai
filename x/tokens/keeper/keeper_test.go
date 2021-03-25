package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/KiraCore/sekai/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)


type KeeperTestSuite struct {
	suite.Suite

	ctx     sdk.Context
	app     *simapp.SimApp
}

func (suite *KeeperTestSuite) SetupTest() {
	checkTx := false
	suite.app = simapp.Setup(checkTx)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, tmproto.Header{Height: 1})
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
