package cli_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"

	"github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/testutil/network"
	"github.com/KiraCore/sekai/x/ixp/client/cli"
	customtypes "github.com/KiraCore/sekai/x/ixp/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	app.SetConfig()
	s.T().Log("setting up integration test suite")

	cfg := network.DefaultConfig()
	encodingConfig := simapp.MakeEncodingConfig()
	cfg.Codec = encodingConfig.Marshaler
	cfg.TxConfig = encodingConfig.TxConfig

	cfg.NumValidators = 1

	cfg.AppConstructor = func(val network.Validator) servertypes.Application {
		return app.NewInitApp(
			val.Ctx.Logger, dbm.NewMemDB(), nil, true, make(map[int64]bool), val.Ctx.Config.RootDir, 0,
			app.MakeEncodingConfig(),
			baseapp.SetPruning(types.NewPruningOptionsFromString(val.AppConfig.Pruning)),
			baseapp.SetMinGasPrices(val.AppConfig.MinGasPrices),
		)
	}

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestCreateOrderBook_AndQueriers() {
	s.T().SkipNow()
	val := s.network.Validators[0]

	cmd := cli.CreateOrderBook()
	_, out := testutil.ApplyMockIO(cmd)
	clientCtx := val.ClientCtx.WithOutput(out)

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	cmd.SetArgs(
		[]string{
			"base",
			"quote",
			"mnemonic",
			// fmt.Sprintf("--%s=%s", cli.FlagXXX, "XXXX"),
		},
	)

	err := cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	height, err := s.network.LatestHeight()
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(height + 2)
	s.Require().NoError(err)

	query := cli.GetOrderBooksCmd()
	query.SetArgs(
		[]string{
			"Quote",
			"quote",
			// fmt.Sprintf("--%s=%s", cli.FlagXXX, "XXXX"),
		},
	)

	out.Reset()

	clientCtx = clientCtx.WithOutputFormat("json")
	err = query.ExecuteContext(ctx)
	s.Require().NoError(err)

	var orderbookResp customtypes.GetOrderBooksResponse
	clientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &orderbookResp)

	orderBooks := orderbookResp.Orderbooks
	s.Require().Len(orderBooks, 1)
	orderbook := orderBooks[0]

	s.Require().Equal("index", orderbook.Index)
	s.Require().Equal("base", orderbook.Base)
	s.Require().Equal("quote", orderbook.Quote)
	s.Require().Equal("mnemonic", orderbook.Mnemonic)
	s.Require().Equal("curator", orderbook.Curator)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
