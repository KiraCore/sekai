package cli

import (
	"context"

	"github.com/KiraCore/sekai/x/multistaking/types"
	stakingtypes "github.com/KiraCore/sekai/x/multistaking/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        stakingtypes.ModuleName,
		Short:                      "Querying commands for the multistaking module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdQueryStakingPools(),
		GetCmdQueryOutstandingRewards(),
		GetCmdQueryUndelegations(),
	)

	return queryCmd
}

// GetCmdQueryStakingPools the query available staking pools
func GetCmdQueryStakingPools() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pools",
		Short: "Query all staking pools",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.QueryStakingPoolsRequest{}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.StakingPools(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryOutstandingRewards the query outstanding rewards for a delegator
func GetCmdQueryOutstandingRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "outstanding-rewards [delegator]",
		Short: "Query outstanding rewards for a delegator",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.QueryOutstandingRewardsRequest{}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.OutstandingRewards(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryUndelegations the query all undelegations
func GetCmdQueryUndelegations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "undelegations",
		Short: "Query all undelegations",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.QueryUndelegationsRequest{}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Undelegations(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
