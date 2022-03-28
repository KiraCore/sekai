package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/KiraCore/sekai/x/distributor/types"
)

// NewQueryCmd returns a root CLI command handler for all x/distributor transaction commands.
func NewQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   types.RouterKey,
		Short: "query commands for the distributor module",
	}
	queryCmd.AddCommand(
		GetCmdQuerydistributorRecordByName(),
		GetCmdQuerydistributorRecords(),
	)

	return queryCmd
}

func GetCmdQuerydistributorRecordByName() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "distributor-record-by-name",
		Short: "Get distributor record by name",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QuerydistributorRecordByName(context.Background(), &types.QuerydistributorRecordByNameRequest{
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

func GetCmdQuerydistributorRecords() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "distributor-records",
		Short: "Get all distributor records",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QuerydistributorRecords(context.Background(), &types.QuerydistributorRecordsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
