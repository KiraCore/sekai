package cli_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/KiraCore/sekai/app"
	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/testutil/network"
	"github.com/KiraCore/sekai/x/evidence/client/cli"
	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store/types"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	dbm "github.com/tendermint/tm-db"
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
	encodingConfig := app.MakeEncodingConfig()
	cfg.Codec = encodingConfig.Marshaler
	cfg.TxConfig = encodingConfig.TxConfig

	cfg.NumValidators = 1

	cfg.AppConstructor = func(val network.Validator) servertypes.Application {
		return app.NewInitApp(
			val.Ctx.Logger, dbm.NewMemDB(), nil, true, make(map[int64]bool), val.Ctx.Config.RootDir, 0,
			simapp.MakeEncodingConfig(),
			simapp.EmptyAppOptions{},
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

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestGetQueryCmd() {
	val := s.network.Validators[0]

	testCases := map[string]struct {
		args           []string
		expectedOutput string
		expectErr      bool
	}{
		"non-existent evidence": {
			[]string{"DF0C23E8634E480F84B9D5674A7CDC9816466DEC28A3358F73260F68D28D7660"},
			"evidence DF0C23E8634E480F84B9D5674A7CDC9816466DEC28A3358F73260F68D28D7660 not found",
			true,
		},
		"all evidence (default pagination)": {
			[]string{},
			"evidence: []\npagination:\n  next_key: null\n  total: \"0\"",
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			cmd := cli.GetQueryCmd()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}

			s.Require().Contains(strings.TrimSpace(out.String()), tc.expectedOutput)
		})
	}
}
