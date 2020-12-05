package cli_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/KiraCore/sekai/x/gov/client/cli"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"

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

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
