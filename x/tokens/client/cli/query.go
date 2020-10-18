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

// GetCmdQueryTokenAlias the query token alias command.
func GetCmdQueryTokenAlias() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alias <symbol>",
		Short: "Get the token alias by symbol",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.TokenAliasRequest{Symbol: args[0]}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetTokenAlias(context.Background(), params)
			if err != nil {
				return err
			}

			if res.Data == nil {
				return fmt.Errorf("%s symbol does not exist", params.Symbol)
			}

			return clientCtx.PrintOutput(res.Data)
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
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.AllTokenAliasesRequest{}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetAllTokenAliases(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryTokenAliasesByDenom the query token aliases by denom command.
func GetCmdQueryTokenAliasesByDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aliases-by-denom",
		Short: "Get token aliases by denom",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			denoms := strings.Split(args[0], ",")
			params := &types.TokenAliasesByDenomRequest{
				Denoms: denoms,
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetTokenAliasesByDenom(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res)
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
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.TokenRateRequest{Denom: args[0]}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetTokenRate(context.Background(), params)
			if err != nil {
				return err
			}

			if res.Data == nil {
				return fmt.Errorf("%s denom does not exist", params.Denom)
			}

			return clientCtx.PrintOutput(res.Data)
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
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.AllTokenRatesRequest{}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetAllTokenRates(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res)
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
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			denoms := strings.Split(args[0], ",")
			params := &types.TokenRatesByDenomRequest{
				Denoms: denoms,
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetTokenRatesByDenom(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
