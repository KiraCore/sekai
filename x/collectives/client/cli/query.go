package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/KiraCore/sekai/x/collectives/types"
)

// NewQueryCmd returns a root CLI command handler for all x/basket transaction commands.
func NewQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   types.RouterKey,
		Short: "query commands for the collectives module",
	}
	queryCmd.AddCommand(
		GetCmdQueryCollective(),
		GetCmdQueryCollectives(),
		GetCmdQueryCollectivesProposals(),
		GetCmdQueryCollectivesByAccount(),
	)

	return queryCmd
}

func GetCmdQueryCollective() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collective [name]",
		Short: "Queries a collective by name",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Collective(context.Background(), &types.CollectiveRequest{
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

func GetCmdQueryCollectives() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collectives",
		Short: "Collectives query list of all staking collectives.",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Collectives(context.Background(), &types.CollectivesRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryCollectivesProposals() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collectives-proposals",
		Short: "list id of all proposals in regards to staking collectives,\n(or proposals in regards to a specific collective if `name` / `id` is specified in the query)",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.CollectivesProposals(context.Background(), &types.CollectivesProposalsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryCollectivesByAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collectives-by-account [account]",
		Short: "query list of staking collectives by an individual KIRA address",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.CollectivesByAccount(context.Background(), &types.CollectivesByAccountRequest{
				Account: args[0],
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
