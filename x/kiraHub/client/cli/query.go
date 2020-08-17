package cli

import (
	"fmt"

	"github.com/KiraCore/cosmos-sdk/client"
	"github.com/KiraCore/sekai/x/kiraHub/types"
	"github.com/spf13/cobra"
)

func GetOrderBooksCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "listorderbooks [by] [value]",
		Short: "List order book(s) by ID, Index, Quote, Base, or Curator",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.JSONMarshaler

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/kiraHub/listOrderBooks/%s/%s", args[0], args[1]), nil)
			if err != nil {
				fmt.Printf("could not query. Searching By - %s & Value - %s is invalid. \n", args[0], args[1])
				return nil
			}

			var out []types.OrderBook
			cdc.MustUnmarshalJSON(res, &out)
			return clientCtx.PrintOutput(out)
		},
	}
}

func GetOrderBooksByTPCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "listorderbooksTP [base] [quote]",
		Short: "List order book(s) by Trading_Pair",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.JSONMarshaler

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/kiraHub/listOrderBooks/tp/%s/%s", args[0], args[1]), nil)
			if err != nil {
				fmt.Printf("could not query. Searching By - %s & Value - %s is invalid. \n", args[0], args[1])
				return nil
			}

			var out []types.OrderBook
			cdc.MustUnmarshalJSON(res, &out)
			return clientCtx.PrintOutput(out)
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
			cdc := clientCtx.JSONMarshaler

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/kiraHub/listOrders/%s/%s/%s", args[0], args[1], args[2]), nil)
			if err != nil {
				fmt.Printf("could not query. Searching By - %s with max_orders - %s and min_amount - %s \n", args[0], args[1], args[2])
				return nil
			}

			var out []types.LimitOrder
			cdc.MustUnmarshalJSON(res, &out)
			return clientCtx.PrintOutput(out)
		},
	}
}

func ListSignerKeysCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "listsignerkeys",
		Short: "List signer key(s) by curator address",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.JSONMarshaler

			var owner = clientCtx.GetFromAddress()

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/kiraHub/listSignerKeys/%s", owner.String()), nil)
			if err != nil {
				fmt.Printf("could not query. Searching By - %s \n", owner.String())
				return nil
			}

			var out []types.SignerKey
			cdc.MustUnmarshalJSON(res, &out)
			return clientCtx.PrintOutput(out)
		},
	}
}
