package cli_test

import (
	"context"
	"fmt"

	cli2 "github.com/KiraCore/sekai/x/staking/client/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	types3 "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/gov/client/cli"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/testutil"
)

func (s IntegrationTestSuite) TestCreateProposalAssignPermission() {
	// Query permissions for role Validator
	val := s.network.Validators[0]

	// We create some random address where we will give perms.
	addr, err := types3.AccAddressFromBech32("kira1alzyfq40zjsveet87jlg8jxetwqmr0a2x50lqq")
	s.Require().NoError(err)

	cmd := cli.GetTxProposalAssignPermission()
	_, out := testutil.ApplyMockIO(cmd)
	clientCtx := val.ClientCtx.WithOutput(out).WithOutputFormat("json")

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	cmd.SetArgs([]string{
		fmt.Sprintf("%d", customgovtypes.PermClaimValidator),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=%s", cli2.FlagAddr, addr.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, types3.NewCoins(types3.NewCoin(s.cfg.BondDenom, types3.NewInt(100))).String()),
	})

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	// Vote Proposal
	out.Reset()
	cmd = cli.GetTxVoteProposal()
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
}

func (s IntegrationTestSuite) TestCreateProposalUpsertDataRegistry() {
	// Query permissions for role Validator
	val := s.network.Validators[0]

	// We create some random address where we will give perms.
	addr, err := types3.AccAddressFromBech32("kira1alzyfq30zjsveet87jlg8jxetwqmr0a22c9uz9")
	s.Require().NoError(err)

	cmd := cli.GetTxProposalUpsertDataRegistry()
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
		fmt.Sprintf("--%s=%s", cli2.FlagAddr, addr.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, types3.NewCoins(types3.NewCoin(s.cfg.BondDenom, types3.NewInt(100))).String()),
	})

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	// Vote Proposal
	out.Reset()
	cmd = cli.GetTxVoteProposal()
	cmd.SetArgs([]string{
		fmt.Sprintf("%d", 2), // Proposal ID
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

func (s IntegrationTestSuite) TestCreateProposalSetNetworkProperty() {
	// Query permissions for role Validator
	val := s.network.Validators[0]

	cmd := cli.GetTxProposalSetNetworkProperty()
	_, out := testutil.ApplyMockIO(cmd)
	clientCtx := val.ClientCtx.WithOutput(out).WithOutputFormat("json")

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	cmd.SetArgs([]string{
		fmt.Sprintf("%s", "MIN_TX_FEE"),
		fmt.Sprintf("%d", 12345),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, types3.NewCoins(types3.NewCoin(s.cfg.BondDenom, types3.NewInt(100))).String()),
	})

	err := cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	// Vote Proposal
	out.Reset()
	cmd = cli.GetTxVoteProposal()
	cmd.SetArgs([]string{
		fmt.Sprintf("%d", 2), // Proposal ID
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
