package cli

import (
	"context"
	"strconv"

	"github.com/KiraCore/sekai/x/kiraHub/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

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

func GetSignerKeysCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "getsignerkeys",
		Short: "List signer key(s) by curator address",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.GetSignerKeysRequest{}
			queryClient := types.NewQueryClient(clientCtx)
			_, err := queryClient.GetSignerKeys(context.Background(), params)
			// res, err := queryClient.GetSignerKeys(context.Background(), params)
			if err != nil {
				return err
			}
			// return clientCtx.PrintOutput(&res.XXXX)
			return nil
			// var owner = clientCtx.GetFromAddress()

			// res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/kiraHub/getsignerkeys/%s", owner.String()), nil)
			// if err != nil {
			// 	fmt.Printf("could not query. Searching By - %s \n", owner.String())
			// 	return nil
			// }

			// var out []types.SignerKey
			// cdc.MustUnmarshalJSON(res, &out)
			// return clientCtx.PrintOutput(out)
		},
	}
}
