package cli_test

import (
	"context"
	"fmt"
	"testing"

	cli3 "github.com/KiraCore/sekai/x/gov/client/cli"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	types3 "github.com/cosmos/cosmos-sdk/types"

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

func (s *IntegrationTestSuite) TestUpsertTokenAliasAndQuery() {
	val := s.network.Validators[0]

	s.WhitelistPermissions(val.Address, customgovtypes.PermUpsertTokenAlias)

	cmd := cli.GetTxUpsertTokenAliasCmd()
	_, out := testutil.ApplyMockIO(cmd)
	clientCtx := val.ClientCtx.WithOutput(out).WithOutputFormat("json")

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	cmd.SetArgs(
		[]string{
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			fmt.Sprintf("--%s=%s", cli.FlagSymbol, "ETH"),
			fmt.Sprintf("--%s=%s", cli.FlagName, "Ethereum"),
			fmt.Sprintf("--%s=%s", cli.FlagIcon, "myiconurl"),
			fmt.Sprintf("--%s=%d", cli.FlagDecimals, 6),
			fmt.Sprintf("--%s=%s", cli.FlagDenoms, "finney"),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=%s", flags.FlagFees, types3.NewCoins(types3.NewCoin(s.cfg.BondDenom, types3.NewInt(100))).String()),
		},
	)

	err := cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	height, err := s.network.LatestHeight()
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(height + 2)
	s.Require().NoError(err)

	out.Reset()

	query := cli.GetCmdQueryTokenAlias()
	query.SetArgs([]string{"ETH"})

	clientCtx = clientCtx.WithOutputFormat("json")
	err = query.ExecuteContext(ctx)
	s.Require().NoError(err)

	var tokenAlias tokenstypes.TokenAlias
	clientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &tokenAlias)

	s.Require().Equal(tokenAlias.Symbol, "ETH")
	s.Require().Equal(tokenAlias.Name, "Ethereum")
	s.Require().Equal(tokenAlias.Icon, "myiconurl")
	s.Require().Equal(tokenAlias.Decimals, uint32(6))
	s.Require().Equal(tokenAlias.Denoms, []string{"finney"})
}

func (s *IntegrationTestSuite) TestUpsertTokenRateAndQuery() {
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
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=%s", flags.FlagFees, types3.NewCoins(types3.NewCoin(s.cfg.BondDenom, types3.NewInt(100))).String()),
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

	var tokenRate tokenstypes.TokenRate
	clientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &tokenRate)

	s.Require().Equal(tokenRate.Denom, "ubtc")
	s.Require().Equal(tokenRate.Rate, types3.NewDec(10))
	s.Require().Equal(tokenRate.FeePayments, true)
}

func (s IntegrationTestSuite) TestCreateProposalUpsertTokenRates() {
	// Query permissions for role Validator
	val := s.network.Validators[0]

	cmd := cli.GetTxProposalUpsertTokenRatesCmd()
	_, out := testutil.ApplyMockIO(cmd)
	clientCtx := val.ClientCtx.WithOutput(out).WithOutputFormat("json")

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	cmd.SetArgs([]string{
		fmt.Sprintf("%s", "theKey"),
		fmt.Sprintf("%s", "theHash"),
		fmt.Sprintf("%s", "theReference"),
		fmt.Sprintf("%s", "theEncoding"),
		fmt.Sprintf("%d", 12345),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, types3.NewCoins(types3.NewCoin(s.cfg.BondDenom, types3.NewInt(100))).String()),
	})

	err := cmd.ExecuteContext(ctx)
	s.Require().NoError(err)
	fmt.Printf("%s", out.String())

	// Vote Proposal
	out.Reset()
	cmd = cli3.GetTxVoteProposal()
	cmd.SetArgs([]string{
		fmt.Sprintf("%d", 1), // Proposal ID
		fmt.Sprintf("%d", customgovtypes.OptionYes),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, types3.NewCoins(types3.NewCoin(s.cfg.BondDenom, types3.NewInt(100))).String()),
	})

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)
	fmt.Printf("%s", out.String())
}

func (s IntegrationTestSuite) TestCreateProposalUpsertDataRegistry() {
	// Query permissions for role Validator
	val := s.network.Validators[0]

	cmd := cli.GetTxProposalUpsertTokenAliasCmd()
	_, out := testutil.ApplyMockIO(cmd)
	clientCtx := val.ClientCtx.WithOutput(out).WithOutputFormat("json")

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	cmd.SetArgs([]string{
		fmt.Sprintf("%s", "theKey"),
		fmt.Sprintf("%s", "theHash"),
		fmt.Sprintf("%s", "theReference"),
		fmt.Sprintf("%s", "theEncoding"),
		fmt.Sprintf("%d", 12345),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, types3.NewCoins(types3.NewCoin(s.cfg.BondDenom, types3.NewInt(100))).String()),
	})

	err := cmd.ExecuteContext(ctx)
	s.Require().NoError(err)
	fmt.Printf("%s", out.String())

	// Vote Proposal
	out.Reset()
	cmd = cli3.GetTxVoteProposal()
	cmd.SetArgs([]string{
		fmt.Sprintf("%d", 1), // Proposal ID
		fmt.Sprintf("%d", customgovtypes.OptionYes),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, types3.NewCoins(types3.NewCoin(s.cfg.BondDenom, types3.NewInt(100))).String()),
	})

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)
	fmt.Printf("%s", out.String())
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
