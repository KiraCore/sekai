package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/KiraCore/sekai/x/spending/types"
)

// NewQueryCmd returns a root CLI command handler for all x/tokens transaction commands.
func NewQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   types.RouterKey,
		Short: "query commands for the tokens module",
	}
	queryCmd.AddCommand(
		GetCmdQueryPoolNames(),
		GetCmdQueryPoolByName(),
		GetCmdQueryPoolProposals(),
		GetCmdQueryPoolsByAccount(),
	)

	return queryCmd
}

// GetCmdQueryPoolNames the query list of pool names
func GetCmdQueryPoolNames() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool-names",
		Short: "Get the query pool names",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			req := &types.QueryPoolNamesRequest{}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryPoolNames(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryPoolByName the query pool by name.
func GetCmdQueryPoolByName() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool-by-name [name]",
		Short: "Query pool by name",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			req := &types.QueryPoolByNameRequest{Name: args[0]}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryPoolByName(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryPoolProposals the query pool proposals by name.
func GetCmdQueryPoolProposals() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool-proposals [pool-name]",
		Short: "Get proposals for the pool by name",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryPoolProposalsRequest{
				PoolName:   args[0],
				Pagination: pageReq,
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryPoolProposals(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "spending")

	return cmd
}

// GetCmdQueryPoolsByAccount the query pool proposals by name.
func GetCmdQueryPoolsByAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pools-by-account [addr]",
		Short: "Query list of pool names where specific kira account can register its claim or otherwise claim tokens from",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			req := &types.QueryPoolsByAccountRequest{Account: args[0]}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryPoolsByAccount(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
