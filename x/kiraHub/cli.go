package kiraHub

import (
	"github.com/KiraCore/sekai/x/kiraHub/client/cli"
	"github.com/KiraCore/sekai/x/kiraHub/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
)

func GetCLIRootTransactionCommand(codec *codec.LegacyAmino) *cobra.Command {
	rootTxCmd := &cobra.Command{
		Use:                        types.TransactionRoute,
		Short:                      "Asset root transaction command.",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	rootTxCmd.AddCommand(
		cli.CreateOrder(),
		cli.CreateOrderBook(),
		cli.UpsertSignerKey(),
	)

	rootTxCmd.PersistentFlags().String("node", "tcp://localhost:26657", "<host>:<port> to Tendermint RPC interface for this chain")
	rootTxCmd.PersistentFlags().String("keyring-backend", "os", "Select keyring's backend (os|file|test)")
	rootTxCmd.PersistentFlags().String("from", "", "Name or address of private key with which to sign")
	rootTxCmd.PersistentFlags().String("broadcast-mode", "sync", "Transaction broadcasting mode (sync|async|block)")

	return rootTxCmd
}

func GetCLIRootQueryCommand(codec *codec.LegacyAmino) *cobra.Command {
	rootQueryCmd := &cobra.Command{
		Use:                        types.QuerierRoute,
		Short:                      "Asset root query command.",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	rootQueryCmd.AddCommand(
		cli.GetOrderBooksCmd(),
		cli.GetOrderBooksByTPCmd(),
		cli.GetOrdersCmd(),
		cli.ListSignerKeysCmd(),
	)

	rootQueryCmd.PersistentFlags().String("node", "tcp://localhost:26657", "<host>:<port> to Tendermint RPC interface for this chain")
	rootQueryCmd.PersistentFlags().String("keyring-backend", "os", "Select keyring's backend (os|file|test)")
	rootQueryCmd.PersistentFlags().String("from", "", "Name or address of private key with which to sign")

	return rootQueryCmd
}
