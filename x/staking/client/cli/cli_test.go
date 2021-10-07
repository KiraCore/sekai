package cli_test

import (
	"fmt"
	"testing"

	customgovcli "github.com/KiraCore/sekai/x/gov/client/cli"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/client/flags"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/stretchr/testify/suite"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/app"
	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/testutil/network"
	"github.com/KiraCore/sekai/x/staking/client/cli"
	customtypes "github.com/KiraCore/sekai/x/staking/types"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
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

func (s *IntegrationTestSuite) TestQueryValidator() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := cli.GetCmdQueryValidator()
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", cli.FlagValAddr, val.ValAddress.String()),
	})
	s.Require().NoError(err)

	var respValidator customtypes.Validator
	clientCtx.JSONCodec.MustUnmarshalJSON(out.Bytes(), &respValidator)

	s.Require().Equal(val.ValAddress, respValidator.ValKey)

	var pubkey cryptotypes.PubKey
	err = s.cfg.Codec.UnpackAny(respValidator.PubKey, &pubkey)
	s.Require().NoError(err)
	s.Require().Equal(val.PubKey, pubkey)

	// Query by Acc Addrs.
	cmd = cli.GetCmdQueryValidator()
	out, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", cli.FlagAddr, val.Address.String()),
	})
	s.Require().NoError(err)

	clientCtx.JSONCodec.MustUnmarshalJSON(out.Bytes(), &respValidator)

	s.Require().Equal(val.ValAddress, respValidator.ValKey)

	err = s.cfg.Codec.UnpackAny(respValidator.PubKey, &pubkey)
	s.Require().NoError(err)
	s.Require().Equal(val.PubKey, pubkey)

	// Query by moniker.
	cmd = cli.GetCmdQueryValidator()
	out, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", cli.FlagMoniker, val.Moniker),
	})
	s.Require().NoError(err)

	clientCtx.JSONCodec.MustUnmarshalJSON(out.Bytes(), &respValidator)

	s.Require().Equal(val.ValAddress, respValidator.ValKey)

	err = s.cfg.Codec.UnpackAny(respValidator.PubKey, &pubkey)
	s.Require().NoError(err)
	s.Require().Equal(val.PubKey, pubkey)
}

func (s *IntegrationTestSuite) TestQueryValidator_Errors() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	nonExistingAddr, err := sdk.ValAddressFromBech32("kiravaloper15ky9du8a2wlstz6fpx3p4mqpjyrm5cgpv3al5n")
	s.Require().NoError(err)

	cmd := cli.GetCmdQueryValidator()
	_, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", cli.FlagValAddr, nonExistingAddr.String()),
	})
	s.Require().EqualError(err, "rpc error: code = InvalidArgument desc = validator not found: key not found: invalid request")

	// Non existing moniker.
	cmd = cli.GetCmdQueryValidator()
	_, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", cli.FlagAddr, sdk.AccAddress(nonExistingAddr).String()),
	})
	s.Require().EqualError(err, "rpc error: code = InvalidArgument desc = validator not found: key not found: invalid request")

	// Non existing moniker.
	cmd = cli.GetCmdQueryValidator()
	_, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", cli.FlagMoniker, "weirdMoniker"),
	})
	s.Require().EqualError(err, "rpc error: code = InvalidArgument desc = validator with moniker weirdMoniker not found: key not found: invalid request")
}

func (s IntegrationTestSuite) TestCreateProposalUnjailValidator() {
	// Query permissions for role Validator
	val := s.network.Validators[0]

	clientCtx := val.ClientCtx.WithOutputFormat("json")
	out, err := clitestutil.ExecTestCLICmd(
		clientCtx,
		cli.GetTxProposalUnjailValidatorCmd(),
		[]string{
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			fmt.Sprintf("--%s=%s", cli.FlagTitle, "title"),
			fmt.Sprintf("--%s=%s", cli.FlagDescription, "some desc"),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
			"theReference",
			"theHash",
		},
	)
	s.Require().NoError(err)
	fmt.Printf("%s", out.String())

	// Vote Proposal
	out, err = clitestutil.ExecTestCLICmd(
		clientCtx,
		customgovcli.GetTxVoteProposal(),
		[]string{
			fmt.Sprintf("%d", 1), // Proposal ID
			fmt.Sprintf("%d", govtypes.OptionYes),
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
		},
	)
	s.Require().NoError(err)
	fmt.Printf("%s", out.String())
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
