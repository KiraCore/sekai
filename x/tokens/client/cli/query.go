package cli

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/tokens/types"
)

const (
	FlagRole = "role"
)

// GetCmdQueryPermissions the query delegation command.
func GetCmdQueryPermissions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "permissions addr",
		Short: "Get the permissions of an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			accAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return errors.Wrap(err, "invalid account address")
			}

			params := &types.PermissionsByAddressRequest{ValAddr: accAddr}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.PermissionsByAddress(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res.Permissions)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
