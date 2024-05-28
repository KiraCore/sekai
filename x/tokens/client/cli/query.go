package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/KiraCore/sekai/x/tokens/types"
)

// NewQueryCmd returns a root CLI command handler for all x/tokens transaction commands.
func NewQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   types.RouterKey,
		Short: "query commands for the tokens module",
	}
	queryCmd.AddCommand(
		GetCmdQueryTokenInfo(),
		GetCmdQueryAllTokenInfos(),
		GetCmdQueryTokenInfosByDenom(),
		GetCmdQueryTokenBlackWhites(),
	)

	return queryCmd
}

// GetCmdQueryTokenInfo the query token rate command.
func GetCmdQueryTokenInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rate <denom>",
		Short: "Get the token rate by denom",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.TokenInfoRequest{Denom: args[0]}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetTokenInfo(context.Background(), params)
			if err != nil {
				return err
			}

			if res.Data == nil {
				return fmt.Errorf("%s denom does not exist", params.Denom)
			}

			return clientCtx.PrintProto(res.Data)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryAllTokenInfos the query all token rates command.
func GetCmdQueryAllTokenInfos() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-rates",
		Short: "Get all token rates",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.AllTokenInfosRequest{}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetAllTokenInfos(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryTokenInfosByDenom the query token aliases by denom command.
func GetCmdQueryTokenInfosByDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rates-by-denom",
		Short: "Get token rates by denom",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			denoms := strings.Split(args[0], ",")
			params := &types.TokenInfosByDenomRequest{
				Denoms: denoms,
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetTokenInfosByDenom(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryTokenBlackWhites the query token blacklists / whitelists
func GetCmdQueryTokenBlackWhites() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-black-whites",
		Short: "Get token black whites",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			params := &types.TokenBlackWhitesRequest{}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetTokenBlackWhites(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
