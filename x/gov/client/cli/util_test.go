package cli_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/stretchr/testify/require"

	"github.com/KiraCore/sekai/testutil/network"

	"github.com/KiraCore/sekai/x/gov/client/cli"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/testutil"
	types3 "github.com/cosmos/cosmos-sdk/types"
)

// GetRolesByAddress calls the CLI command GetCmdQueryRolesByAddress and returns the roles.
func GetRolesByAddress(t *testing.T, network *network.Network, address types3.AccAddress) []uint64 {
	val := network.Validators[0]

	cmd := cli.GetCmdQueryRolesByAddress()
	_, out := testutil.ApplyMockIO(cmd)
	clientCtx := val.ClientCtx.WithOutput(out).WithOutputFormat("json")

	ctx := context.WithValue(context.Background(), client.ClientContextKey, &clientCtx)

	cmd.SetArgs([]string{
		address.String(),
	})

	err := cmd.ExecuteContext(ctx)

	var roles customgovtypes.RolesByAddressResponse
	err = val.ClientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &roles)
	require.NoError(t, err)

	return roles.Roles
}

// SetCouncilor calls CLI to set address in the Councilor Registry. The Validator 1 is the caller.
func (s IntegrationTestSuite) SetCouncilor(address types3.Address) {
	val := s.network.Validators[0]

	cmd := cli.GetTxClaimCouncilorSeatCmd()
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
			fmt.Sprintf("--%s=%s", cli.FlagAddress, address.String()),
			fmt.Sprintf("--%s=%s", cli.FlagMoniker, val.Moniker),
		},
	)

	err := cmd.ExecuteContext(ctx)
	s.Require().NoError(err)
}
