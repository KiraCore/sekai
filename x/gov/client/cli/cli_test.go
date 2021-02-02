package cli_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/KiraCore/sekai/x/gov/client/cli"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"

	"github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/testutil/network"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"
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

func (s IntegrationTestSuite) TestRolePermissions_QueryCommand_DefaultRolePerms() {
	val := s.network.Validators[0]

	cmd := cli.GetCmdQueryRolePermissions()
	_, out := testutil.ApplyMockIO(cmd)
	clientCtx := val.ClientCtx.WithOutput(out).WithOutputFormat("json")

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	cmd.SetArgs([]string{
		"2", // RoleValidator
	})

	err := cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	var perms customgovtypes.Permissions
	val.ClientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &perms)

	s.Require().True(perms.IsWhitelisted(customgovtypes.PermClaimValidator))
}

func (s IntegrationTestSuite) TestClaimCouncilor_HappyPath() {
	val := s.network.Validators[0]

	s.SetCouncilor(val.Address)

	err := s.network.WaitForNextBlock()
	s.Require().NoError(err)

	// Query command
	// Mandatory flags
	cmd := cli.GetCmdQueryCouncilRegistry()

	_, out := testutil.ApplyMockIO(cmd)
	clientCtx := val.ClientCtx.WithOutput(out).WithOutputFormat("json")

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	cmd.SetArgs([]string{
		"",
	})

	err = cmd.ExecuteContext(ctx)
	s.Require().Error(err)

	// From address
	out.Reset()

	cmd = cli.GetCmdQueryCouncilRegistry()
	cmd.SetArgs([]string{
		fmt.Sprintf("--%s=%s", cli.FlagAddress, val.Address.String()),
	})

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	var councilorByAddress customgovtypes.Councilor
	err = val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &councilorByAddress)
	s.Require().NoError(err)
	s.Require().Equal(val.Moniker, councilorByAddress.Moniker)
	s.Require().Equal(val.Address, councilorByAddress.Address)

	// From Moniker
	out.Reset()

	cmd = cli.GetCmdQueryCouncilRegistry()
	cmd.SetArgs([]string{
		fmt.Sprintf("--%s=%s", cli.FlagMoniker, val.Moniker),
	})

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	var councilorByMoniker customgovtypes.Councilor
	err = val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &councilorByMoniker)
	s.Require().NoError(err)
	s.Require().Equal(val.Moniker, councilorByMoniker.Moniker)
	s.Require().Equal(val.Address, councilorByMoniker.Address)
}

func (s IntegrationTestSuite) TestProposalAndVoteSetPoorNetworkMessages_HappyPath() {
	val := s.network.Validators[0]
	// # create proposal for setting poor network msgs
	// sekaid tx customgov proposal set-poor-network-msgs AAA,BBB --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=1000ukex --yes
	s.SetPoorNetworkMessages("AAA,BBB")
	// # query for proposals
	// sekaid query customgov proposals
	s.QueryProposals()
	// # set permission to vote on proposal
	// sekaid tx customgov permission whitelist-permission --permission=19 --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes
	s.WhitelistPermission(val.Address, "19") // 19 is permission for vote on poor network message set proposal
	// # vote on the proposal
	// sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes
	s.VoteWithValidator0(1, customgovtypes.OptionYes)
	// # check votes
	// sekaid query customgov votes 1
	s.QueryProposalVotes(1)
	// # wait until vote end time finish
	// sekaid query customgov proposals
	// TODO: this takes long time and for now skip waiting
	// # query poor network messages
	// sekaid query customgov poor-network-messages
	s.QueryPoorNetworkMessages()
}

func (s IntegrationTestSuite) TestProposalAndVotePoorNetworkMaxBankSend_HappyPath() {
	// TODO: complete scenarios
	// # try setting network property by governance to allow more amount sending
	// sekaid tx customgov proposal set-network-property POOR_NETWORK_MAX_BANK_SEND 100000000 --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes
	// sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes
	// # try sending after modification of poor network bank send param
	// sekaid tx bank send validator $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) 100000000ukex --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
}

func (s IntegrationTestSuite) TestPoorNetworkRestrictions_HappyPath() {
	// TODO: complete scenarios
	// # whitelist permission for modifying network properties
	// sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=7 --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
	// # test poor network messages after modifying min_validators section
	// sekaid tx customgov set-network-properties --from validator --min_validators="2" --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
	// # set permission for upsert token rate
	// sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermUpsertTokenRate --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
	// # try running upser token rate which is not allowed on poor network
	// sekaid tx tokens upsert-rate --from validator --keyring-backend=test --denom="mykex" --rate="1.5" --fee_payments=true --chain-id=testing --fees=100ukex --home=$HOME/.sekaid  --yes
	// # try sending more than allowed amount via bank send
	// sekaid tx bank send validator $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) 100000000ukex --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
