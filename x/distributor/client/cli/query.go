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
		GetCmdSnapshotPeriod(),
		GetCmdValidatorSnapshotPerformance(),
		GetCmdYearStartSnapshot(),
		GetCmdPeriodicSnapshot(),
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

func GetCmdValidatorSnapshotPerformance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot-period-performance [val_addr]",
		Short: "Get snapshot period performance for a validator",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.SnapshotPeriodPerformance(context.Background(), &types.QuerySnapshotPeriodPerformanceRequest{
				ValidatorAddress: args[0],
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

func GetCmdYearStartSnapshot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "year-start-snapshot",
		Short: "Year start snapshot",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.YearStartSnapshot(context.Background(), &types.QueryYearStartSnapshotRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdPeriodicSnapshot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "periodic-snapshot",
		Short: "Periodic snapshot",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.PeriodicSnapshot(context.Background(), &types.QueryPeriodicSnapshotRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
