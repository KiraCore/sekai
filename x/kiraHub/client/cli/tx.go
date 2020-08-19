package cli

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/KiraCore/sekai/x/kiraHub/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateOrder() *cobra.Command {

	return &cobra.Command{
		Use:   "createOrder [order_book_id] [type] [amount] [price] [expiry_time]",
		Short: "Create Order",
		Long:  "0 - Limit Buy, 1 - Limit Sell",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			var curator = clientCtx.GetFromAddress()
			var orderType, _ = strconv.Atoi(args[1])

			// Limit Order
			if uint8(orderType) == 0 || uint8(orderType) == 1 {

				var amount, _ = strconv.Atoi(args[2])
				var limitPrice, _ = strconv.Atoi(args[3])

				message, _ := types.NewMsgCreateOrder(
					args[0],
					types.LimitOrderType(orderType),
					int64(amount),
					int64(limitPrice),
					curator,
				)

				if err := message.ValidateBasic(); err != nil {
					return err
				}
				return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), message)
			}

			fmt.Println("did not get in limit order")

			return errors.New("Invalid limit order type")
		},
	}
}

func CreateOrderBook() *cobra.Command {

	return &cobra.Command{
		Use:   "createOrderBook [base] [quote] [mnemonic]",
		Short: "Create OrderBook",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			var curator = clientCtx.GetFromAddress()

			message, _ := types.NewMsgCreateOrderBook(
				args[0],
				args[1],
				args[2],
				curator,
			)

			if err := message.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), message)
		},
	}
}

// UpsertSignerKey is a cli command to upsertSignerKey
func UpsertSignerKey() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "upsertSignerKey [pubKey] [keyType] [expiryTime] [enabled] [permissions]",
		Short: "upsert signer key",
		Long:  "Secp256k1 | Ed25519 for keyType",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			var curator = clientCtx.GetFromAddress()
			var pubKeyFlag = args[0]
			var keyTypeFlag = viper.GetString("key-type")
			var enabledFlag = viper.GetBool("enabled")
			var permissionsFlag = viper.GetIntSlice("permissions")
			var expiryTimeFlag = viper.GetInt64("expiry-time")

			var permissions = []int64{}

			for _, p := range permissionsFlag {
				permissions = append(permissions, int64(p))
			}

			var keyType = types.SignerKeyType_Secp256k1
			// TODO: check encoding types, because users will be providing strings
			// check encoding types, because users will be providing strings
			switch keyTypeFlag {
			case types.SignerKeyType_Secp256k1.String():
				// TODO: should set pubKey from pubKeyFlag for SignerKeyType_Secp256k1
				// tendermint/PubKeySecp256k1",
				// "AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w"
				// Library "github.com/tendermint/tendermint/crypto/secp256k1"
				// TODO: should check the byte length
				// secp256k1PubKey := secp256k1.GenPrivKey().PubKey()
				// err := cdc.UnmarshalBinaryBare([]byte(pubKeyFlag), &secp256k1PubKey)
				// if err != nil {
				// 	fmt.Println("invalid secp256k1PubKey", err)
				// 	return fmt.Errorf("Invalid secp256k1PubKey: %s", err.Error())
				// }
			case types.SignerKeyType_Ed25519.String():
				// TODO: should set pubKey from pubKeyFlag for SignerKeyType_Ed25519
				// tendermint/PrivKeyEd25519"
				// "TXgDkmTYpPRwU/PvDbfbhbwiYA7jXMwQgNffHVey1dC644OBBI4OQdf4Tro6hzimT1dHYzPiGZB0aYWJBC2keQ=="
				// Library "github.com/tendermint/tendermint/crypto/ed25519"
				// TODO: should check the byte length
				keyType = types.SignerKeyType_Ed25519
				// ed25519PubKey := ed25519.GenPrivKey().PubKey()
				// err := cdc.UnmarshalBinaryBare([]byte(pubKeyFlag), &ed25519PubKey)
				// if err != nil {
				// 	fmt.Println("invalid ed25519PubKey", err)
				// 	return fmt.Errorf("Invalid ed25519PubKey: %s", err.Error())
				// }
			default:
				fmt.Println("invalid pubKey type")
				return fmt.Errorf("invalid pubKey type: %s", keyTypeFlag)
			}

			message, _ := types.NewMsgUpsertSignerKey(
				pubKeyFlag,
				keyType,
				expiryTimeFlag,
				enabledFlag,
				permissions,
				curator,
			)

			if err := message.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), message)
		},
	}
	cmd.Flags().String("key-type", "Secp256k1", "flag to set pubKey type; Secp256k1 | Ed25519")
	cmd.Flags().Bool("enabled", true, "flag to enable/disable pubKey")
	cmd.Flags().IntSlice("permissions", []int{}, "flag to set permissions set for the pubKey")
	cmd.Flags().Int64("expiry-time", time.Now().Add(time.Hour*24*10).Unix(), "flag to set permissions set for the pubKey")
	return cmd
}
