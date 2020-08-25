package cli

import (
	"context"
	"strconv"

	"github.com/KiraCore/sekai/x/kiraHub/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

// flags
const (
	FlagCurator = "curator"
)

// GetOrderBooksCmd is a command to query orderbooks
func GetOrderBooksCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "listorderbooks [by] [value]",
		Short: "List order book(s) by ID, Index, Quote, Base, or Curator",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.GetOrderBooksRequest{
				QueryType:  args[0],
				QueryValue: args[1],
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetOrderBooks(context.Background(), params)
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(res)
		},
	}
}

// GetOrderBooksByTradingPairCmd is a command to query orderbooks by trading pair
func GetOrderBooksByTradingPairCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "listorderbooks_tradingpair [base] [quote]",
		Short: "List order book(s) by trading pair",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.GetOrderBooksByTradingPairRequest{
				Base:  args[0],
				Quote: args[1],
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetOrderBooksByTradingPair(context.Background(), params)
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(res)
		},
	}
}

// GetOrdersCmd is a command to query orders
func GetOrdersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "listorders [order_book_id] [max_orders] [min_amount]",
		Short: "List order(s) by Order Book ID",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			maxOrders, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			minAmount, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}

			params := &types.GetOrdersRequest{
				OrderBookID: args[0],
				MaxOrders:   uint32(maxOrders),
				MinAmount:   uint32(minAmount),
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetOrders(context.Background(), params)
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(res)
		},
	}
}

// GetSignerKeysCmd is a command to query signer keys
func GetSignerKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getsignerkeys",
		Short: "List signer key(s) by curator address",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			curatorText, _ := cmd.Flags().GetString(FlagCurator)
			curator, err := sdk.AccAddressFromBech32(curatorText)
			if err != nil {
				return err
			}

			params := &types.GetSignerKeysRequest{
				Curator: curator,
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetSignerKeys(context.Background(), params)
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(res)
		},
	}
	cmd.Flags().String(FlagCurator, "", "address to query signer keys")
	return cmd
}
