package signerkey

import (
	"fmt"
	"time"

	"bufio"

	"github.com/KiraCore/cosmos-sdk/client/context"
	"github.com/KiraCore/cosmos-sdk/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/cosmos-sdk/x/auth"
	"github.com/KiraCore/cosmos-sdk/x/auth/client"
	"github.com/KiraCore/sekai/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// TransactionCommand is a cli command to upsertSignerKey
func TransactionCommand(codec *codec.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "upsertSignerKey [pubKey] [keyType] [expiryTime] [enabled] [permissions]",
		Short: "upsert signer key",
		Long:  "Secp256k1 | Ed25519 for keyType",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			bufioReader := bufio.NewReader(command.InOrStdin())
			transactionBuilder := auth.NewTxBuilderFromCLI(bufioReader).WithTxEncoder(auth.DefaultTxEncoder(codec))
			cliContext := context.NewCLIContext().WithCodec(codec)

			var curator = cliContext.GetFromAddress()
			var pubKeyFlag = args[0]
			var pubKey = [4096]byte{}
			var keyTypeFlag = viper.GetString("key-type")
			var enabledFlag = viper.GetBool("enabled")
			var permissionsFlag = viper.GetIntSlice("permissions")
			var expiryTimeFlag = viper.GetInt64("expiry-time")

			var keyType = types.Secp256k1
			// TODO: check encoding types, because users will be providing strings
			// check encoding types, because users will be providing strings
			switch keyTypeFlag {
			case types.Secp256k1.String():
				// TODO: should set pubKey from pubKeyFlag for Secp256k1
				// tendermint/PubKeySecp256k1",
				// "AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w"
				// Library "github.com/tendermint/tendermint/crypto/secp256k1"
				// TODO: should check the byte length
				secp256k1PubKey := secp256k1.GenPrivKey().PubKey()
				err := codec.UnmarshalBinaryBare([]byte(pubKeyFlag), &secp256k1PubKey)
				if err != nil {
					fmt.Println("invalid secp256k1PubKey", err)
					return client.GenerateOrBroadcastMsgs(cliContext, transactionBuilder, []sdk.Msg{})
				}
			case types.Ed25519.String():
				// TODO: should set pubKey from pubKeyFlag for Ed25519
				// tendermint/PrivKeyEd25519"
				// "TXgDkmTYpPRwU/PvDbfbhbwiYA7jXMwQgNffHVey1dC644OBBI4OQdf4Tro6hzimT1dHYzPiGZB0aYWJBC2keQ=="
				// Library "github.com/tendermint/tendermint/crypto/ed25519"
				// TODO: should check the byte length
				keyType = types.Ed25519
				ed25519PubKey := ed25519.GenPrivKey().PubKey()
				err := codec.UnmarshalBinaryBare([]byte(pubKeyFlag), &ed25519PubKey)
				if err != nil {
					fmt.Println("invalid ed25519PubKey", err)
					return client.GenerateOrBroadcastMsgs(cliContext, transactionBuilder, []sdk.Msg{})
				}
			default:
				fmt.Println("invalid pubKey type")
				return client.GenerateOrBroadcastMsgs(cliContext, transactionBuilder, []sdk.Msg{})
			}

			var message = Message{
				PubKey:      pubKey,
				KeyType:     keyType,
				ExpiryTime:  expiryTimeFlag,
				Enabled:     enabledFlag,
				Permissions: permissionsFlag,
				Curator:     curator,
			}

			if err := message.ValidateBasic(); err != nil {
				return err
			}

			return client.GenerateOrBroadcastMsgs(cliContext, transactionBuilder, []sdk.Msg{message})
		},
	}
	cmd.Flags().String("key-type", "Secp256k1", "flag to set pubKey type; Secp256k1 | Ed25519")
	cmd.Flags().Bool("enabled", true, "flag to enable/disable pubKey")
	cmd.Flags().IntSlice("permissions", []int{}, "flag to set permissions set for the pubKey")
	cmd.Flags().Int64("expiry-time", time.Hour.Milliseconds()*24*10, "flag to set permissions set for the pubKey")
	return cmd
}
