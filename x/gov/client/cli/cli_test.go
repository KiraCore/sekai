package cli_test

import (
	"context"
	"fmt"
	"testing"

	cli3 "github.com/cosmos/cosmos-sdk/x/bank/client/cli"

	"github.com/cosmos/cosmos-sdk/client/flags"

	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	types3 "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/gov/client/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/testutil"

	"github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/testutil/network"
	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store/types"
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

	cmd := cli.GetTxClaimGovernanceCmd()
	_, out := testutil.ApplyMockIO(cmd)
	clientCtx := val.ClientCtx.WithOutput(out).WithOutputFormat("json")

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	cmd.SetArgs(
		[]string{
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=%s", flags.FlagFees, types3.NewCoins(types3.NewCoin(s.cfg.BondDenom, types3.NewInt(10))).String()),
			fmt.Sprintf("--%s=%s", cli.FlagAddress, val.Address.String()),
			fmt.Sprintf("--%s=%s", cli.FlagMoniker, val.Moniker),
		},
	)

	err := cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	fmt.Printf("%s\n", out.String())

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	// Query command
	// Mandatory flags
	out.Reset()

	cmd = cli.GetCmdQueryCouncilRegistry()
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

func (s IntegrationTestSuite) sendValue(cCtx client.Context, from types3.AccAddress, to types3.AccAddress, coin types3.Coin) {
	cmd := cli3.NewSendTxCmd()
	_, out := testutil.ApplyMockIO(cmd)
	cCtx = cCtx.WithOutput(out).WithOutputFormat("json")

	cmd.SetArgs(
		[]string{
			from.String(),
			to.String(),
			coin.String(),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=%s", flags.FlagFees, types3.NewCoins(types3.NewCoin(s.cfg.BondDenom, types3.NewInt(10))).String()),
		},
	)

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &cCtx)

	err := cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
