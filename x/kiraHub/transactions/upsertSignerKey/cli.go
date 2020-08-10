package signerkey

import (
	"fmt"

	"bufio"

	"github.com/KiraCore/cosmos-sdk/client/context"
	"github.com/KiraCore/cosmos-sdk/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/cosmos-sdk/x/auth"
	"github.com/KiraCore/cosmos-sdk/x/auth/client"
	"github.com/KiraCore/sekai/types"
	"github.com/spf13/cobra"
)

// TransactionCommand is a cli command to upsertSignerKey
func TransactionCommand(codec *codec.Codec) *cobra.Command {

	return &cobra.Command{
		Use:   "upsertSignerKey [pubKey] [keyType] [expiryTime] [enabled] [permissions]",
		Short: "upsert signer key",
		Long:  "Secp256k1 | Ed25519 for keyType",
		Args:  cobra.ExactArgs(4),
		RunE: func(command *cobra.Command, args []string) error {
			bufioReader := bufio.NewReader(command.InOrStdin())
			transactionBuilder := auth.NewTxBuilderFromCLI(bufioReader).WithTxEncoder(auth.DefaultTxEncoder(codec))
			cliContext := context.NewCLIContext().WithCodec(codec)

			var curator = cliContext.GetFromAddress()
			var pubKey = [4096]byte{}
			var keyTypeString = args[1]
			var keyType = types.Secp256k1
			var enabled = false       // TODO: should set from args[2]
			var permissions = []int{} // TODO: should set from args[3]

			switch keyTypeString {
			case types.Secp256k1.String():
				// TODO: should set pubKey from args[0] for Secp256k1
				// tendermint/PubKeySecp256k1",
				// "AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w"
				// Library "github.com/tendermint/tendermint/crypto/secp256k1"
			case types.Ed25519.String():
				// TODO: should set pubKey from args[0] for Ed25519
				// tendermint/PrivKeyEd25519"
				// "TXgDkmTYpPRwU/PvDbfbhbwiYA7jXMwQgNffHVey1dC644OBBI4OQdf4Tro6hzimT1dHYzPiGZB0aYWJBC2keQ=="
				// Library "github.com/tendermint/tendermint/crypto/ed25519"
				keyType = types.Ed25519
			default:
				fmt.Println("invalid pubKey type")
				return client.GenerateOrBroadcastMsgs(cliContext, transactionBuilder, []sdk.Msg{})
			}

			var message = Message{
				PubKey:      pubKey,
				KeyType:     keyType,
				ExpiryTime:  0, // TODO: should discuss if it should be set here
				Enabled:     enabled,
				Permissions: permissions,
				Curator:     curator,
			}

			if err := message.ValidateBasic(); err != nil {
				return err
			}

			// TODO: should add readme for sample upsertSignerKey cli and rest
			// TODO: should add readme for both Secp256k1 and Ed25519
			return client.GenerateOrBroadcastMsgs(cliContext, transactionBuilder, []sdk.Msg{message})
		},
	}
}
