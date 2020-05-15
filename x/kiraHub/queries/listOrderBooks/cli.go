package listOrderBooks

import (
	"fmt"
	"github.com/KiraCore/sekai/types"
	"github.com/KiraCore/cosmos-sdk/client/context"
	"github.com/KiraCore/cosmos-sdk/codec"

	"github.com/spf13/cobra"
)

func GetCoinCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "listorderbooks [by] [value]",
		Short: "List order book(s) by ID, Index, Quote, Base, Trading_Pair, or Curator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			//var owner = cliCtx.GetFromAddress()

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/kiraHub/listOrderBook/%s/%s", args[0], args[1]), nil)
			if err != nil {
				fmt.Printf("could not getcoin. owner - %s is invalid \n", args[0])
				return nil
			}

			var out types.OrderBook
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}