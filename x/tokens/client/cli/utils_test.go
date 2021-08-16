package cli_test

import (
	"fmt"

	"github.com/KiraCore/sekai/x/gov/client/cli"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	customstakingcli "github.com/KiraCore/sekai/x/staking/client/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s IntegrationTestSuite) WhitelistPermissions(addr sdk.AccAddress, perm govtypes.PermValue) {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := cli.GetTxSetWhitelistPermissions()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=%s", customstakingcli.FlagAddr, addr.String()),
		fmt.Sprintf("--%s=%d", cli.FlagPermission, perm),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	// We check if the user has the permissions
	cmd = cli.GetCmdQueryPermissions()
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		addr.String(),
	})
	s.Require().NoError(err)

	var perms govtypes.Permissions
	clientCtx.JSONCodec.MustUnmarshalJSON(out.Bytes(), &perms)

	// Validator 1 has permission to Add Permissions.
	s.Require().True(perms.IsWhitelisted(perm))
}
