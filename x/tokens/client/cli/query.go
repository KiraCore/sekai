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
		GetCmdQueryTokenAlias(),
		GetCmdQueryAllTokenAliases(),
		GetCmdQueryTokenAliasesByDenom(),
		GetCmdQueryTokenRate(),
		GetCmdQueryAllTokenRates(),
		GetCmdQueryTokenRatesByDenom(),
		GetCmdQueryTokenBlackWhites(),
	)

	return queryCmd
}

// GetCmdQueryTokenAlias the query token alias command.
func GetCmdQueryTokenAlias() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alias <symbol>",
		Short: "Get the token alias by symbol",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.TokenAliasRequest{Symbol: args[0]}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetTokenAlias(context.Background(), params)
			if err != nil {
				return err
			}

			if res.Data == nil {
				return fmt.Errorf("%s symbol does not exist", params.Symbol)
			}

			return clientCtx.PrintProto(res.Data)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryAllTokenAliases the query all token aliases command.
func GetCmdQueryAllTokenAliases() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-aliases",
		Short: "Get all token aliases",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.AllTokenAliasesRequest{}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetAllTokenAliases(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryTokenAliasesByDenom the query token aliases by denom command.
func GetCmdQueryTokenAliasesByDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aliases-by-denom [aliases]",
		Short: "Get token aliases by denom",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			denoms := strings.Split(args[0], ",")
			params := &types.TokenAliasesByDenomRequest{
				Denoms: denoms,
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetTokenAliasesByDenom(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryTokenRate the query token rate command.
func GetCmdQueryTokenRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rate <denom>",
		Short: "Get the token rate by denom",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.TokenRateRequest{Denom: args[0]}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetTokenRate(context.Background(), params)
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

// GetCmdQueryAllTokenRates the query all token rates command.
func GetCmdQueryAllTokenRates() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-rates",
		Short: "Get all token rates",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.AllTokenRatesRequest{}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetAllTokenRates(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryTokenRatesByDenom the query token aliases by denom command.
func GetCmdQueryTokenRatesByDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rates-by-denom",
		Short: "Get token rates by denom",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			denoms := strings.Split(args[0], ",")
			params := &types.TokenRatesByDenomRequest{
				Denoms: denoms,
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetTokenRatesByDenom(context.Background(), params)
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
