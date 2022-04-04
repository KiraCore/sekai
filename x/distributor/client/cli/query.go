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
		GetCmdQueryFeesTreasury(),
		GetCmdQueryFeesCollected(),
		GetCmdSnapshotPeriod(),
	)

	return queryCmd
}

func GetCmdQueryFeesTreasury() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fees-treasury",
		Short: "Get fees treasury",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.FeesTreasury(context.Background(), &types.QueryFeesTreasuryRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryFeesCollected() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fees-collected",
		Short: "Get fees collected",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.FeesCollected(context.Background(), &types.QueryFeesCollectedRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdSnapshotPeriod() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot-period",
		Short: "Get snapshot period",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.SnapshotPeriod(context.Background(), &types.QuerySnapshotPeriodRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
