package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/KiraCore/sekai/x/layer2/types"
)

// NewQueryCmd returns a root CLI command handler for all x/layer2 transaction commands.
func NewQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   types.RouterKey,
		Short: "query commands for the layer2 module",
	}
	queryCmd.AddCommand(
		GetCmdQueryExecutionRegistrar(),
		GetCmdQueryAllDapps(),
		GetCmdQueryTransferDapp(),
	)

	return queryCmd
}

func GetCmdQueryExecutionRegistrar() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execution-registrar [dapp-name]",
		Short: "Queries a execution registrar for a dapp",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ExecutionRegistrar(context.Background(), &types.QueryExecutionRegistrarRequest{
				Identifier: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryAllDapps() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-dapps",
		Short: "Queries all dapps",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.AllDapps(context.Background(), &types.QueryAllDappsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryTransferDapp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-dapps",
		Short: "Queries transfer dapps",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.TransferDapps(context.Background(), &types.QueryTransferDappsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
