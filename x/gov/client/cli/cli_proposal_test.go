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

	cmd := cli.GetTxProposalAssignPermission()
	_, out := testutil.ApplyMockIO(cmd)
	clientCtx := val.ClientCtx.WithOutput(out).WithOutputFormat("json")

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	cmd.SetArgs([]string{
		fmt.Sprintf("%d", customgovtypes.PermClaimValidator),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=%s", cli2.FlagAddr, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, types3.NewCoins(types3.NewCoin(s.cfg.BondDenom, types3.NewInt(10))).String()),
	})

	err := cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	s.T().Logf("%s", out.String())
}
