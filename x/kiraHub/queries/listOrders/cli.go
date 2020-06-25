package listOrders

import (
	"fmt"
	"github.com/KiraCore/cosmos-sdk/client/context"
	"github.com/KiraCore/cosmos-sdk/codec"
	"github.com/KiraCore/sekai/types"

	"github.com/spf13/cobra"
)

func GetOrdesCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "listorders [id] [max_orders] [min_amount]",
		Short: "List order(s) by ID",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			//var owner = cliCtx.GetFromAddress()

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/kiraHub/listOrders/%s/%s/%s", args[0], args[1], args[2]), nil)
			if err != nil {
				fmt.Printf("could not query. Searching By - %s with max_orders - %s and min_amount - %s \n", args[0], args[1], args[2])
				return nil
			}

			var out []types.LimitOrder
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
