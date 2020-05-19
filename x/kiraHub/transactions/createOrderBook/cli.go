package createOrderBook

import (
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
		Use:   "createOrderBook [base] [quote] [Mnemonic]",
		Short: "Create and sign transaction to mint an asset",
		Args:  cobra.ExactArgs(3),
		RunE: func(command *cobra.Command, args []string) error {
			bufioReader := bufio.NewReader(command.InOrStdin())
			transactionBuilder := auth.NewTxBuilderFromCLI(bufioReader).WithTxEncoder(auth.DefaultTxEncoder(codec))
			cliContext := context.NewCLIContext().WithCodec(codec)

			var curator = cliContext.GetFromAddress()
			message := Message {
				Base: args[0],
				Quote: args[1],
				Mnemonic: args[2],
				Curator: string(curator),
			}

			if err := message.ValidateBasic(); err != nil {
				return err
			}

			return client.GenerateOrBroadcastMsgs(cliContext, transactionBuilder, []sdk.Msg{message})
		},
	}
}