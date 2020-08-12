package kiraHub

import (
	constants "github.com/KiraCore/sekai/x/kiraHub/constants"
	"github.com/KiraCore/sekai/x/kiraHub/queries/listOrderBooks"
	"github.com/KiraCore/sekai/x/kiraHub/queries/listOrders"
	signerkeys "github.com/KiraCore/sekai/x/kiraHub/queries/listSignerKeys"
	"github.com/KiraCore/sekai/x/kiraHub/transactions/createOrder"
	"github.com/KiraCore/sekai/x/kiraHub/transactions/createOrderBook"
	signerkey "github.com/KiraCore/sekai/x/kiraHub/transactions/upsertSignerKey"
	"github.com/spf13/cobra"

	"github.com/KiraCore/cosmos-sdk/client"
	"github.com/KiraCore/cosmos-sdk/codec"
)

func GetCLIRootTransactionCommand(codec *codec.Codec) *cobra.Command {
	rootTransactionCommand := &cobra.Command{
		Use:                        constants.TransactionRoute,
		Short:                      "Asset root transaction command.",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	rootTransactionCommand.AddCommand(
		createOrderBook.TransactionCommand(codec),
		createOrder.TransactionCommand(codec),
		signerkey.TransactionCommand(codec),
	)...)

	rootTransactionCommand.PersistentFlags().String("node", "tcp://localhost:26657", "<host>:<port> to Tendermint RPC interface for this chain")
	rootTransactionCommand.PersistentFlags().String("keyring-backend", "os", "Select keyring's backend (os|file|test)")
	rootTransactionCommand.PersistentFlags().String("from", "", "Name or address of private key with which to sign")
	rootTransactionCommand.PersistentFlags().String("broadcast-mode", "sync", "Transaction broadcasting mode (sync|async|block)")


	return rootTransactionCommand
}

func GetCLIRootQueryCommand(codec *codec.Codec) *cobra.Command {
	rootQueryCommand := &cobra.Command{
		Use:                        constants.QuerierRoute,
		Short:                      "Asset root query command.",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	rootQueryCommand.AddCommand(
		listOrderBooks.GetOrderBooksCmd(codec),
		listOrderBooks.GetOrderBooksByTPCmd(codec),
		listOrders.GetOrdersCmd(codec),
		signerkeys.ListSignerKeysCmd(codec),
	)...)

	rootQueryCommand.PersistentFlags().String("node", "tcp://localhost:26657", "<host>:<port> to Tendermint RPC interface for this chain")
	rootQueryCommand.PersistentFlags().String("keyring-backend", "os", "Select keyring's backend (os|file|test)")
	rootQueryCommand.PersistentFlags().String("from", "", "Name or address of private key with which to sign")

	return rootQueryCommand
}
