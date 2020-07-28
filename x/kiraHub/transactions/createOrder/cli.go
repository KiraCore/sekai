package createOrder

import (
	"github.com/KiraCore/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func TransactionCommand(codec *codec.Codec) *cobra.Command {

	return &cobra.Command{
		Use:   "createOrder [order_book_id] [type] [amount] [price]",
		Short: "Create Order",
		Long:  "0 - Limit Buy, 1 - Limit Sell",
		Args:  cobra.ExactArgs(4),
		RunE: func(command *cobra.Command, args []string) error {
			//bufioReader := bufio.NewReader(command.InOrStdin())
			//transactionBuilder := authtypes.NewTxBuilderFromCLI(bufioReader).WithTxEncoder(auth.DefaultTxEncoder(codec))
			//cliContext := context.NewCLIContext().WithCodec(codec)
			//
			//var curator = cliContext.GetFromAddress()
			//var orderType, _ =  strconv.Atoi(args[1])
			//
			//// Limit Order
			//
			//if uint8(orderType) == 0 || uint8(orderType) == 1 {
			//
			//	var amount, _ = strconv.Atoi(args[2])
			//	var limitPrice, _ = strconv.Atoi(args[3])
			//
			//	var message = Message {
			//		OrderBookID: args[0],
			//		OrderType: uint8(orderType),
			//		Amount: int64(amount),
			//		LimitPrice: int64(limitPrice),
			//		Curator: curator,
			//	}
			//
			//	if err := message.ValidateBasic(); err != nil {
			//		return err
			//	}
			//
			//	return client.GenerateOrBroadcastMsgs(cliContext, transactionBuilder, []sdk.Msg{message})
			//
			//}
			//
			//fmt.Println("did not get in limit order")
			//
			//return client.GenerateOrBroadcastMsgs(cliContext, transactionBuilder, []sdk.Msg{})
			panic("implement me")
		},
	}
}
