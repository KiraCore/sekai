package cli_test

import (
	"context"
	"testing"

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
