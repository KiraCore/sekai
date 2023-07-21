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
		GetCmdQueryCompoundInfo(),
		GetCmdQueryStakingPoolDelegators(),
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
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.QueryOutstandingRewardsRequest{
				Delegator: args[0],
			}

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
		Use:   "undelegations [delegator] [val-addr]",
		Short: "Query all undelegations",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.QueryUndelegationsRequest{
				Delegator:  args[0],
				ValAddress: args[1],
			}

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

func GetCmdQueryCompoundInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compound-info [delegator]",
		Short: "Query compound information of a delegator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.QueryCompoundInfoRequest{
				Delegator: args[0],
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.CompoundInfo(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryStakingPoolDelegators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "staking-pool-delegators [validator]",
		Short: "Query staking pool delegators",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.QueryStakingPoolDelegatorsRequest{
				Validator: args[0],
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.StakingPoolDelegators(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
