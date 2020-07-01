package createOrder

import (
	"fmt"
	"strconv"

	"bufio"
	"github.com/KiraCore/cosmos-sdk/client/context"
	"github.com/KiraCore/cosmos-sdk/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/cosmos-sdk/x/auth"
	"github.com/KiraCore/cosmos-sdk/x/auth/client"
	"github.com/spf13/cobra"
)

func TransactionCommand(codec *codec.Codec) *cobra.Command {

	return &cobra.Command{
		Use:   "createOrder [order_book_id] [type] [amount] [price]",
		Short: "Create Order",
		Long: "0 - Limit Buy, 1 - Limit Sell",
		Args:  cobra.ExactArgs(4),
		RunE: func(command *cobra.Command, args []string) error {
			bufioReader := bufio.NewReader(command.InOrStdin())
			transactionBuilder := auth.NewTxBuilderFromCLI(bufioReader).WithTxEncoder(auth.DefaultTxEncoder(codec))
			cliContext := context.NewCLIContext().WithCodec(codec)

			var curator = cliContext.GetFromAddress()
			var orderType, _ =  strconv.Atoi(args[1])

			// Limit Order

			if uint8(orderType) == 1 || uint8(orderType) == 2 {

				var amount, _ = strconv.Atoi(args[2])
				var limitPrice, _ = strconv.Atoi(args[3])

				var message = Message {
					OrderBookID: args[0],
					OrderType: uint8(orderType),
					Amount: int64(amount),
					LimitPrice: int64(limitPrice),
					Curator: curator,
				}

				if err := message.ValidateBasic(); err != nil {
					return err
				}

				return client.GenerateOrBroadcastMsgs(cliContext, transactionBuilder, []sdk.Msg{message})

			}

			fmt.Println("did not get in limit order")

			return client.GenerateOrBroadcastMsgs(cliContext, transactionBuilder, []sdk.Msg{})
		},
	}
}
