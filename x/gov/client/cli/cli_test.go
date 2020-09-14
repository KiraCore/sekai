package cli_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cli3 "github.com/cosmos/cosmos-sdk/x/bank/client/cli"

	"github.com/cosmos/cosmos-sdk/client/flags"

	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	cli2 "github.com/KiraCore/sekai/x/staking/client/cli"
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

func (s IntegrationTestSuite) TestGetTxSetWhitelistPermissions() {
	val := s.network.Validators[0]
	cmd := cli.GetTxSetWhitelistPermissions()

	_, out := testutil.ApplyMockIO(cmd)

	clientCtx := val.ClientCtx.WithOutput(out).WithOutputFormat("json")

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	// We create some random address where we will give perms.
	addr, err := types3.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	s.Require().NoError(err)

	cmd.SetArgs(
		[]string{
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			fmt.Sprintf("--%s=%s", cli2.FlagAddr, addr.String()),
			fmt.Sprintf("--%s=%s", cli.FlagPermission, "1"),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=%s", flags.FlagFees, types3.NewCoins(types3.NewCoin(s.cfg.BondDenom, types3.NewInt(10))).String()),
		},
	)

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	// We check if the user has the permissions
	cmd = cli.GetCmdQueryPermissions()
	out.Reset()

	cmd.SetArgs(
		[]string{
			addr.String(),
		},
	)

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	var perms customgovtypes.Permissions
	clientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &perms)

	// Validator 1 has permission to Add Permissions.
	s.Require().False(perms.IsWhitelisted(customgovtypes.PermSetPermissions))
	s.Require().True(perms.IsWhitelisted(customgovtypes.PermClaimValidator))
}

func (s IntegrationTestSuite) TestGetTxSetBlacklistPermissions() {
	val := s.network.Validators[0]
	cmd := cli.GetTxSetBlacklistPermissions()

	_, out := testutil.ApplyMockIO(cmd)

	clientCtx := val.ClientCtx.WithOutput(out).WithOutputFormat("json")

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	// We create some random address where we will give perms.
	addr, err := types3.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	s.Require().NoError(err)

	cmd.SetArgs(
		[]string{
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			fmt.Sprintf("--%s=%s", cli2.FlagAddr, addr.String()),
			fmt.Sprintf("--%s=%s", cli.FlagPermission, "1"),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=%s", flags.FlagFees, types3.NewCoins(types3.NewCoin(s.cfg.BondDenom, types3.NewInt(10))).String()),
		},
	)

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)
	s.T().Logf("error %s", out.String())

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	// We check if the user has the permissions
	cmd = cli.GetCmdQueryPermissions()
	out.Reset()

	cmd.SetArgs(
		[]string{
			addr.String(),
		},
	)

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	var perms customgovtypes.Permissions
	clientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &perms)

	// Validator 1 has permission to Add Permissions.
	s.Require().False(perms.IsBlacklisted(customgovtypes.PermSetPermissions))
	s.Require().True(perms.IsBlacklisted(customgovtypes.PermClaimValidator))
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

func (s IntegrationTestSuite) TestGetTxSetWhitelistPermissions_WithUserThatDoesNotHaveSetPermissions() {
	val := s.network.Validators[0]

	// We create some random address where we will give perms.
	newAccount, _, err := val.ClientCtx.Keyring.NewMnemonic("test", keyring.English, "", hd.Secp256k1)
	s.Require().NoError(err)
	s.sendValue(val.ClientCtx, val.Address, newAccount.GetAddress(), types3.NewCoin(s.cfg.BondDenom, types3.NewInt(100)))

	// Now we try to set permissions with a user that does not have.
	cmd := cli.GetTxSetWhitelistPermissions()
	_, out := testutil.ApplyMockIO(cmd)
	clientCtx := val.ClientCtx.WithOutput(out).WithOutputFormat("json")

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	cmd.SetArgs(
		[]string{
			fmt.Sprintf("--%s=%s", flags.FlagFrom, newAccount.GetAddress().String()),
			fmt.Sprintf("--%s=%s", cli2.FlagAddr, val.Address.String()),
			fmt.Sprintf("--%s=%s", cli.FlagPermission, "1"),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=%s", flags.FlagFees, types3.NewCoins(types3.NewCoin(s.cfg.BondDenom, types3.NewInt(10))).String()),
		},
	)

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	strings.Contains(out.String(), "SetPermissions: not enough permissions")
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

func (s IntegrationTestSuite) TestWhitelistRolePermission() {
	// Query permissions for role Validator
	val := s.network.Validators[0]

	cmd := cli.GetCmdQueryRolePermissions()
	_, out := testutil.ApplyMockIO(cmd)
	clientCtx := val.ClientCtx.WithOutput(out).WithOutputFormat("json")

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	cmd.SetArgs([]string{
		"0", // RoleInTest
	})

	err := cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	var perms customgovtypes.Permissions
	val.ClientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &perms)

	s.Require().False(perms.IsWhitelisted(customgovtypes.PermSetPermissions))

	// Send Tx To Whitelist permission
	out.Reset()

	cmd = cli.GetTxWhitelistRolePermission()
	cmd.SetArgs([]string{
		"0", // Role created in test
		"2", // PermSetPermission
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, types3.NewCoins(types3.NewCoin(s.cfg.BondDenom, types3.NewInt(10))).String()),
	})

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	// Query again to check if it has the new permission
	out.Reset()

	cmd = cli.GetCmdQueryRolePermissions()

	cmd.SetArgs([]string{
		"0", // RoleCreatedInTest
	})

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	var newPerms customgovtypes.Permissions
	val.ClientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &newPerms)

	s.Require().True(newPerms.IsWhitelisted(customgovtypes.PermSetPermissions))
}

func (s IntegrationTestSuite) TestBlacklistRolePermission() {
	// Query permissions for role Validator
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
	s.Require().False(perms.IsBlacklisted(customgovtypes.PermClaimCouncilor))

	// Send Tx To Blacklist permission
	out.Reset()

	cmd = cli.GetTxBlacklistRolePermission()
	cmd.SetArgs([]string{
		"2", // RoleValidator
		"3", // PermClaimCouncilor
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, types3.NewCoins(types3.NewCoin(s.cfg.BondDenom, types3.NewInt(10))).String()),
	})

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	// Query again to check if it has the new permission
	out.Reset()

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	cmd = cli.GetCmdQueryRolePermissions()

	cmd.SetArgs([]string{
		"2", // RoleValidator
	})

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	var newPerms customgovtypes.Permissions
	val.ClientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &newPerms)

	s.Require().True(newPerms.IsWhitelisted(customgovtypes.PermClaimValidator))
	s.Require().True(newPerms.IsBlacklisted(customgovtypes.PermClaimCouncilor))
}

func (s IntegrationTestSuite) TestRemoveWhitelistRolePermission() {
	// Query permissions for role Validator
	val := s.network.Validators[0]

	cmd := cli.GetCmdQueryRolePermissions()
	_, out := testutil.ApplyMockIO(cmd)
	clientCtx := val.ClientCtx.WithOutput(out).WithOutputFormat("json")

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	cmd.SetArgs([]string{
		"0", // RoleInTest
	})

	err := cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	var perms customgovtypes.Permissions
	val.ClientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &perms)

	s.Require().True(perms.IsWhitelisted(customgovtypes.PermClaimValidator))

	// Send Tx To Blacklist permission
	out.Reset()

	cmd = cli.GetTxRemoveWhitelistRolePermission()
	cmd.SetArgs([]string{
		"0", // RoleValidator
		"1", // PermClaimValidator
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, types3.NewCoins(types3.NewCoin(s.cfg.BondDenom, types3.NewInt(10))).String()),
	})

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	// Query again to check if it has the new permission
	out.Reset()

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	cmd = cli.GetCmdQueryRolePermissions()

	cmd.SetArgs([]string{
		"0", // RoleInTest
	})

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	var newPerms customgovtypes.Permissions
	val.ClientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &newPerms)

	s.Require().False(newPerms.IsWhitelisted(customgovtypes.PermClaimValidator))
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
