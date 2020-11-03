package cli_test

import (
	"context"
	"fmt"
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
	"github.com/KiraCore/sekai/x/tokens/client/cli"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
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

func (s *IntegrationTestSuite) TestUpsertTokenAliasAndQuery() {
	s.T().SkipNow()
	val := s.network.Validators[0]

	cmd := cli.GetTxUpsertTokenAliasCmd()
	_, out := testutil.ApplyMockIO(cmd)
	clientCtx := val.ClientCtx.WithOutput(out)

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	cmd.SetArgs(
		[]string{
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			fmt.Sprintf("--%s=%d", cli.FlagExpiration, 0),
			fmt.Sprintf("--%s=%d", cli.FlagEnactment, 0),
			fmt.Sprintf("--%s=%s", cli.FlagAllowedVoteTypes, "0,1"),
			fmt.Sprintf("--%s=%s", cli.FlagSymbol, "ETH"),
			fmt.Sprintf("--%s=%s", cli.FlagName, "Ethereum"),
			fmt.Sprintf("--%s=%s", cli.FlagIcon, "myiconurl"),
			fmt.Sprintf("--%s=%d", cli.FlagDecimals, 6),
			fmt.Sprintf("--%s=%s", cli.FlagDenoms, "finney"),
		},
	)

	err := cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	height, err := s.network.LatestHeight()
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(height + 2)
	s.Require().NoError(err)

	query := cli.GetCmdQueryTokenAlias()
	query.SetArgs([]string{"ETH"})

	out.Reset()

	clientCtx = clientCtx.WithOutputFormat("json")
	err = query.ExecuteContext(ctx)
	s.Require().NoError(err)

	var tokenAliasResponse tokenstypes.TokenAliasResponse
	clientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &tokenAliasResponse)
	tokenAlias := tokenAliasResponse.Data

	s.Require().Equal(tokenAlias.Expiration, 0)
	s.Require().Equal(tokenAlias.Enactment, 0)
	s.Require().Equal(tokenAlias.AllowedVoteTypes, []tokenstypes.VoteType{tokenstypes.VoteType_yes, tokenstypes.VoteType_no})
	s.Require().Equal(tokenAlias.Symbol, "ETH")
	s.Require().Equal(tokenAlias.Name, "Ethereum")
	s.Require().Equal(tokenAlias.Icon, "myiconurl")
	s.Require().Equal(tokenAlias.Decimals, 6)
	s.Require().Equal(tokenAlias.Denoms, []string{"finney"})
}

func (s *IntegrationTestSuite) TestUpsertTokenRateAndQuery() {
	s.T().SkipNow()
	val := s.network.Validators[0]

	cmd := cli.GetTxUpsertTokenRateCmd()
	_, out := testutil.ApplyMockIO(cmd)
	clientCtx := val.ClientCtx.WithOutput(out)

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	cmd.SetArgs(
		[]string{
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			fmt.Sprintf("--%s=%s", cli.FlagDenom, "ubtc"),
			fmt.Sprintf("--%s=%f", cli.FlagRate, 0.00001),
			fmt.Sprintf("--%s=%s", cli.FlagFeePayments, "true"),
		},
	)

	err := cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	height, err := s.network.LatestHeight()
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(height + 2)
	s.Require().NoError(err)

	query := cli.GetCmdQueryTokenRate()
	query.SetArgs([]string{"ubtc"})

	out.Reset()

	clientCtx = clientCtx.WithOutputFormat("json")
	err = query.ExecuteContext(ctx)
	s.Require().NoError(err)

	var tokenRateResponse tokenstypes.TokenRateResponse
	clientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &tokenRateResponse)
	tokenRate := tokenRateResponse.Data

	s.Require().Equal(tokenRate.Denom, "ubtc")
	s.Require().Equal(tokenRate.Rate, 0.00001)
	s.Require().Equal(tokenRate.FeePayments, true)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}