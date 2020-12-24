package cli_test

import (
	"context"
	"fmt"

	"github.com/KiraCore/sekai/x/gov/client/cli"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	cli2 "github.com/KiraCore/sekai/x/staking/client/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s IntegrationTestSuite) WhitelistPermissions(addr sdk.AccAddress, perm customgovtypes.PermValue) {
	val := s.network.Validators[0]
	cmd := cli.GetTxSetWhitelistPermissions()

	_, out := testutil.ApplyMockIO(cmd)

	clientCtx := val.ClientCtx.WithOutput(out).WithOutputFormat("json")

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	cmd.SetArgs(
		[]string{
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			fmt.Sprintf("--%s=%s", cli2.FlagAddr, addr.String()),
			fmt.Sprintf("--%s=%d", cli.FlagPermission, perm),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
		},
	)

	err := cmd.ExecuteContext(ctx)
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
	s.Require().True(perms.IsWhitelisted(perm))
}
