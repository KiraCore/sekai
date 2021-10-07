package cli_test

import (
	"fmt"
	"testing"

	"github.com/KiraCore/sekai/app"
	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/testutil/network"
	"github.com/KiraCore/sekai/x/gov/client/cli"
	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
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
	encCfg := app.MakeEncodingConfig()
	cfg.Codec = encCfg.Marshaler
	cfg.TxConfig = encCfg.TxConfig

	cfg.NumValidators = 1

	// customize proposal end time and enactment time
	govGen := types.DefaultGenesis()
	// govGen.NetworkProperties.ProposalEndTime = 1
	// govGen.NetworkProperties.ProposalEnactmentTime = 2
	govGenRaw := encCfg.Marshaler.MustMarshalJSON(govGen)

	genesis := app.ModuleBasics.DefaultGenesis(encCfg.Marshaler)
	genesis[types.ModuleName] = govGenRaw
	cfg.GenesisState = genesis

	cfg.AppConstructor = func(val network.Validator) servertypes.Application {
		return app.NewInitApp(
			val.Ctx.Logger, dbm.NewMemDB(), nil, true, make(map[int64]bool), val.Ctx.Config.RootDir, 0,
			simapp.MakeEncodingConfig(),
			simapp.EmptyAppOptions{},
			baseapp.SetPruning(storetypes.NewPruningOptionsFromString(val.AppConfig.Pruning)),
			baseapp.SetMinGasPrices(val.AppConfig.MinGasPrices),
		)
	}

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s IntegrationTestSuite) TestRolePermissions_QueryCommand_DefaultRolePerms() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := cli.GetCmdQueryRolePermissions()
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		"2", // RoleValidator
	})
	s.Require().NoError(err)

	var perms types.Permissions
	val.ClientCtx.JSONCodec.MustUnmarshalJSON(out.Bytes(), &perms)
	s.Require().True(perms.IsWhitelisted(types.PermClaimValidator))
}

func (s IntegrationTestSuite) TestClaimCouncilor_HappyPath() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	s.SetCouncilor(val.Address)

	err := s.network.WaitForNextBlock()
	s.Require().NoError(err)

	// Query command
	// Mandatory flags
	cmd := cli.GetCmdQueryCouncilRegistry()

	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{""})
	s.Require().Error(err)

	// From address
	out.Reset()

	cmd = cli.GetCmdQueryCouncilRegistry()

	out, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", cli.FlagAddress, val.Address.String()),
	})
	s.Require().NoError(err)

	var councilorByAddress types.Councilor
	err = val.ClientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), &councilorByAddress)
	s.Require().NoError(err)
	s.Require().Equal(val.Moniker, councilorByAddress.Moniker)
	s.Require().Equal(val.Address, councilorByAddress.Address)

	// From Moniker
	cmd = cli.GetCmdQueryCouncilRegistry()
	out, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", cli.FlagMoniker, val.Moniker),
	})
	s.Require().NoError(err)

	var councilorByMoniker types.Councilor
	err = val.ClientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), &councilorByMoniker)
	s.Require().NoError(err)
	s.Require().Equal(val.Moniker, councilorByMoniker.Moniker)
	s.Require().Equal(val.Address, councilorByMoniker.Address)
}

func (s IntegrationTestSuite) TestProposalAndVoteSetPoorNetworkMessages_HappyPath() {
	val := s.network.Validators[0]

	// create proposal for setting poor network msgs
	result := s.SetPoorNetworkMessages("AAA,BBB")
	s.Require().Contains(result.RawLog, "SetPoorNetworkMessages")

	// query for proposals
	s.QueryProposals()

	// set permission to vote on proposal
	s.WhitelistPermission(val.Address, "19") // 19 is permission for vote on poor network message set proposal

	// vote on the proposal
	s.VoteWithValidator0(1, types.OptionYes)

	// check votes
	s.QueryProposalVotes(1)

	// check proposal status until gov process it
	s.network.WaitForNextBlock()

	// query poor network messages
	s.QueryPoorNetworkMessages()
}

func (s IntegrationTestSuite) TestProposalAndVotePoorNetworkMaxBankSend_HappyPath() {
	val := s.network.Validators[0]

	// set min validators to 2
	s.SetNetworkProperties(1, 10000, 2)

	// try setting network property by governance to allow more amount sending
	s.SetNetworkPropertyProposal("POOR_NETWORK_MAX_BANK_SEND", 100000000)

	// vote on the proposal
	s.VoteWithValidator0(1, types.OptionYes)

	// check votes
	s.QueryProposalVotes(1)

	// check proposal status until gov process it
	s.network.WaitForNextBlock()

	// try sending after modification of poor network bank send param
	s.SendValue(val.ClientCtx, val.Address, val.Address, sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100000000)))
}

func (s IntegrationTestSuite) TestPoorNetworkRestrictions_HappyPath() {
	val := s.network.Validators[0]

	// whitelist permission for modifying network properties
	s.WhitelistPermission(val.Address, "7")

	// test poor network messages after modifying min_validators section
	s.SetNetworkProperties(1, 10000, 2)

	// set permission for upsert token rate
	s.WhitelistPermission(val.Address, "8")

	// try running upser token rate which is not allowed on poor network
	result := s.UpsertRate("mykex", "1.5", true)
	s.Require().Contains(result.RawLog, "invalid transaction type on poor network")

	// try sending more than allowed amount via bank send
	s.SendValue(val.ClientCtx, val.Address, val.Address, sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100000000)))
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
