package cli

import (
	"context"

	"github.com/KiraCore/sekai/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the parent command for all x/upgrade CLi query commands.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the upgrade module",
	}

	cmd.AddCommand(
		GetCmdQueryCurrentPlan(),
	)

	return cmd
}

// GetCmdQueryCurrentPlan the query current plan.
func GetCmdQueryCurrentPlan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query current upgrade plan",
		Short: "Get the upgrade plan",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.QueryCurrentPlanRequest{}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.CurrentPlan(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
