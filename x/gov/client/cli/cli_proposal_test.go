package cli_test

import (
	"fmt"

	"github.com/KiraCore/sekai/x/gov/client/cli"
	"github.com/KiraCore/sekai/x/gov/types"
	stakingcli "github.com/KiraCore/sekai/x/staking/client/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s IntegrationTestSuite) TestCreateProposalAssignPermission() {
	// Query permissions for role Validator
	val := s.network.Validators[0]

	// We create some random address where we will give perms.
	addr, err := sdk.AccAddressFromBech32("kira1alzyfq40zjsveet87jlg8jxetwqmr0a2x50lqq")
	s.Require().NoError(err)

	cmd := cli.GetTxProposalAssignPermission()
	clientCtx := val.ClientCtx
	_, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("%d", types.PermClaimValidator),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=%s", stakingcli.FlagAddr, addr.String()),
		fmt.Sprintf("--%s=%s", cli.FlagTitle, "title"),
		fmt.Sprintf("--%s=%s", cli.FlagDescription, "some desc"),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)

	cmd = cli.GetTxVoteProposal()
	_, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("%d", 1), // Proposal ID
		fmt.Sprintf("%d", types.OptionYes),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
}

func (s IntegrationTestSuite) TestCreateProposalUpsertDataRegistry() {
	// Query permissions for role Validator
	val := s.network.Validators[0]

	cmd := cli.GetTxProposalUpsertDataRegistry()
	clientCtx := val.ClientCtx
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("%s", "theKey"),
		fmt.Sprintf("%s", "theHash"),
		fmt.Sprintf("%s", "theReference"),
		fmt.Sprintf("%s", "theEncoding"),
		fmt.Sprintf("%d", 12345),
		fmt.Sprintf("--%s=%s", cli.FlagTitle, "title"),
		fmt.Sprintf("--%s=%s", cli.FlagDescription, "some desc"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
	fmt.Printf("%s", out.String())

	// Vote Proposal
	cmd = cli.GetTxVoteProposal()
	out, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("%d", 1), // Proposal ID
		fmt.Sprintf("%d", types.OptionYes),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
	fmt.Printf("%s", out.String())
}

func (s IntegrationTestSuite) TestCreateProposalSetNetworkProperty() {
	// Query permissions for role Validator
	val := s.network.Validators[0]

	cmd := cli.GetTxProposalSetNetworkProperty()
	clientCtx := val.ClientCtx
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("%s", "MIN_TX_FEE"),
		fmt.Sprintf("%d", 12345),
		fmt.Sprintf("--%s=%s", cli.FlagTitle, "title"),
		fmt.Sprintf("--%s=%s", cli.FlagDescription, "some desc"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
	fmt.Printf("%s", out.String())

	// Vote Proposal
	cmd = cli.GetTxVoteProposal()
	out, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("%d", 1), // Proposal ID
		fmt.Sprintf("%d", types.OptionYes),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
	fmt.Printf("%s", out.String())
}

func (s IntegrationTestSuite) TestCreateProposalCreateRole() {
	// Query permissions for role Validator
	val := s.network.Validators[0]

	cmd := cli.GetTxProposalCreateRole()
	clientCtx := val.ClientCtx
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("%d", 12345),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=%s", cli.FlagTitle, "title"),
		fmt.Sprintf("--%s=%s", cli.FlagDescription, "some desc"),
		fmt.Sprintf("--%s=%s", cli.FlagWhitelistPerms, "1,2,3"),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
	fmt.Printf("%s", out.String())

	// Vote Proposal
	cmd = cli.GetTxVoteProposal()
	out, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("%d", 1), // Proposal ID
		fmt.Sprintf("%d", types.OptionYes),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
	fmt.Printf("%s", out.String())
}
