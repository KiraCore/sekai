package cli_test

import (
	"fmt"
	"testing"

	customgovcli "github.com/KiraCore/sekai/x/gov/client/cli"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"

	"github.com/KiraCore/sekai/app"
	simapp "github.com/KiraCore/sekai/app"
	appparams "github.com/KiraCore/sekai/app/params"
	"github.com/KiraCore/sekai/testutil/network"
	"github.com/KiraCore/sekai/x/tokens/client/cli"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	pruningtypes "github.com/cosmos/cosmos-sdk/store/pruning/types"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	appparams.SetConfig()
	s.T().Log("setting up integration test suite")

	cfg := network.DefaultConfig()
	encodingConfig := simapp.MakeEncodingConfig()
	cfg.Codec = encodingConfig.Marshaler
	cfg.TxConfig = encodingConfig.TxConfig

	cfg.NumValidators = 1

	cfg.AppConstructor = func(val network.Validator, chainId string) servertypes.Application {
		return app.NewInitApp(
			val.Ctx.Logger, dbm.NewMemDB(), nil, true, make(map[int64]bool), val.Ctx.Config.RootDir, 0,
			simapp.MakeEncodingConfig(),
			simapp.EmptyAppOptions{},
			baseapp.SetPruning(pruningtypes.NewPruningOptionsFromString(val.AppConfig.Pruning)),
			baseapp.SetMinGasPrices(val.AppConfig.MinGasPrices),
			baseapp.SetChainID(chainId),
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

func (s *IntegrationTestSuite) TestUpsertTokenInfoAndQuery() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	s.WhitelistPermissions(val.Address, govtypes.PermUpsertTokenInfo)

	cmd := cli.GetTxUpsertTokenInfoCmd()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=%s", cli.FlagDescription, "some desc"),
		fmt.Sprintf("--%s=%s", cli.FlagDenom, "ubtc"),
		fmt.Sprintf("--%s=%s", cli.FlagSupply, "0"),
		fmt.Sprintf("--%s=%s", cli.FlagSupplyCap, "0"),
		fmt.Sprintf("--%s=%f", cli.FlagFeeRate, 0.00001),
		fmt.Sprintf("--%s=%d", cli.FlagMintingFee, 2),
		fmt.Sprintf("--%s=%s", cli.FlagFeeEnabled, "true"),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.DefaultDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)

	height, err := s.network.LatestHeight()
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(height + 2)
	s.Require().NoError(err)

	cmd = cli.GetCmdQueryTokenInfo()
	_, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{"ubtc"})
	s.Require().NoError(err)

}

func (s *IntegrationTestSuite) TestGetCmdQueryTokenBlackWhites() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := cli.GetCmdQueryTokenBlackWhites()
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{})
	s.Require().NoError(err)

	var blackWhites tokenstypes.TokenBlackWhitesResponse
	clientCtx.Codec.MustUnmarshalJSON(out.Bytes(), &blackWhites)

	s.Require().Equal(blackWhites.Data.Blacklisted, []string{"frozen"})
	s.Require().Equal(blackWhites.Data.Whitelisted, []string{"ukex"})
}

func (s IntegrationTestSuite) TestCreateProposalUpsertTokenInfo() {
	// Query permissions for role Validator
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := cli.GetTxProposalUpsertTokenInfoCmd()
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", cli.FlagDenom, "ubtc"),
		fmt.Sprintf("--%s=%f", cli.FlagFeeRate, 0.00001),
		fmt.Sprintf("--%s=%s", cli.FlagTitle, "title"),
		fmt.Sprintf("--%s=%s", cli.FlagDescription, "some desc"),
		fmt.Sprintf("--%s=%s", cli.FlagFeeEnabled, "true"),
		fmt.Sprintf("--%s=%s", cli.FlagDecimals, "6"),
		fmt.Sprintf("--%s=%s", cli.FlagStakeCap, "0.3"),
		fmt.Sprintf("--%s=%s", cli.FlagTokenRate, "1"),
		fmt.Sprintf("--%s=%s", cli.FlagSupply, "0"),
		fmt.Sprintf("--%s=%s", cli.FlagSupplyCap, "0"),
		fmt.Sprintf("--%s=%s", cli.FlagMintingFee, "1"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.DefaultDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
	fmt.Printf("%s", out.String())

	// Vote Proposal
	cmd = customgovcli.GetTxVoteProposal()
	out, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("%d", 1), // Proposal ID
		fmt.Sprintf("%d", govtypes.OptionYes),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.DefaultDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
	fmt.Printf("%s", out.String())
}

func (s IntegrationTestSuite) TestTxProposalTokensBlackWhiteChangeCmd() {
	// Query permissions for role Validator
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := cli.GetTxProposalTokensBlackWhiteChangeCmd()
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=true", cli.FlagIsBlacklist),
		fmt.Sprintf("--%s=true", cli.FlagIsAdd),
		fmt.Sprintf("--%s=%s", cli.FlagTitle, "title"),
		fmt.Sprintf("--%s=%s", cli.FlagDescription, "some desc"),
		fmt.Sprintf("--%s=frozen1", cli.FlagTokens),
		fmt.Sprintf("--%s=frozen2", cli.FlagTokens),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.DefaultDenom, sdk.NewInt(100))).String()),
	})

	s.Require().NoError(err)
	fmt.Printf("%s", out.String())

	// Vote Proposal
	out.Reset()
	cmd = customgovcli.GetTxVoteProposal()
	out, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("%d", 1), // Proposal ID
		fmt.Sprintf("%d", govtypes.OptionYes),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.DefaultDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
	fmt.Printf("%s", out.String())
}

func (s *IntegrationTestSuite) TestGetCmdQueryAllTokenInfos() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := cli.GetCmdQueryAllTokenInfos()
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{})
	s.Require().NoError(err)

	var resp tokenstypes.AllTokenInfosResponse
	clientCtx.Codec.MustUnmarshalJSON(out.Bytes(), &resp)

	s.Require().Greater(len(resp.Data), 0)
}

func (s *IntegrationTestSuite) TestGetCmdQueryTokenInfosByDenom() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := cli.GetCmdQueryTokenInfosByDenom()
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{"ukex"})
	s.Require().NoError(err)

	var resp tokenstypes.TokenInfosByDenomResponse
	clientCtx.Codec.MustUnmarshalJSON(out.Bytes(), &resp)

	s.Require().Greater(len(resp.Data), 0)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
