package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/KiraCore/sekai/x/ubi/types"
)

// NewQueryCmd returns a root CLI command handler for all x/ubi transaction commands.
func NewQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   types.RouterKey,
		Short: "query commands for the ubi module",
	}
	queryCmd.AddCommand(
		GetCmdQueryUBIRecordByName(),
		GetCmdQueryUBIRecords(),
	)

	return queryCmd
}

func GetCmdQueryUBIRecordByName() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ubi-record-by-name",
		Short: "Get ubi record by name",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryUBIRecordByName(context.Background(), &types.QueryUBIRecordByNameRequest{
				Name: args[0],
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

func GetCmdQueryUBIRecords() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ubi-records",
		Short: "Get all ubi records",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryUBIRecords(context.Background(), &types.QueryUBIRecordsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
