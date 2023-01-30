package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/KiraCore/sekai/x/recovery/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

// NewTxCmd returns a root CLI command handler for all x/recovery transaction commands.
func NewTxCmd() *cobra.Command {
	recoveryTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Recovery transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	recoveryTxCmd.AddCommand(
		NewRegisterRecoverySecretTxCmd(),
		NewRotateRecoveryAddressTxCmd(),
		NewIssueRecoveryTokensTxCmd(),
		NewBurnRecoveryTokensTxCmd(),
		NewGenerateRecoverySecretCmd(),
	)
	return recoveryTxCmd
}

// NewRegisterRecoverySecretTxCmd defines MsgRegisterRecoverySecret tx
func NewRegisterRecoverySecretTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-recovery-secret [challenge] [nonce] [proof]",
		Args:  cobra.ExactArgs(3),
		Short: "Register recovery secret",
		Long: `Register recovery secret:

$ <appd> tx recovery register-recovery-secret [challenge] [nonce] [proof] --from mykey
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRegisterRecoverySecret(clientCtx.GetFromAddress().String(), args[0], args[1], args[2])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewRotateRecoveryAddressTxCmd defines MsgRotateRecoveryAddress tx
func NewRotateRecoveryAddressTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rotate-recovery-address [recovery] [proof]",
		Args:  cobra.ExactArgs(2),
		Short: "Rotate an address to recovery address",
		Long: `Rotate an address to recovery address:

$ <appd> tx recovery rotate-recovery-address [recovery] [proof] --from validator --chain-id=testing --keyring-backend=test --fees=100ukex --home=$HOME/.sekaid --yes
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRotateRecoveryAddress(clientCtx.GetFromAddress().String(), args[0], args[1])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewIssueRecoveryTokensTxCmd defines MsgIssueRecoveryTokens tx
func NewIssueRecoveryTokensTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-recovery-tokens",
		Args:  cobra.NoArgs,
		Short: "Issue recovery tokens",
		Long: `Issue recovery tokens:

$ <appd> tx recovery issue-recovery-tokens --from mykey
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgIssueRecoveryTokens(clientCtx.GetFromAddress().String())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewBurnRecoveryTokensTxCmd defines MsgBurnRecoveryTokens tx
func NewBurnRecoveryTokensTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-recovery-tokens",
		Args:  cobra.NoArgs,
		Short: "Burn recovery tokens",
		Long: `Burn recovery tokens:

$ <appd> tx recovery burn-recovery-tokens --from mykey
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgBurnRecoveryTokens(clientCtx.GetFromAddress())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewGenerateRecoverySecretCmd generates recovery secret
func NewGenerateRecoverySecretCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-recovery-secret [key]",
		Args:  cobra.ExactArgs(1),
		Short: "Generate recovery secret",
		Long: `Generate recovery secret:

$ <appd> tx recovery generate-recovery-secret 10a0fbe01030000122300000000000
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			privKey, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			proof := sha256.Sum256(privKey)
			challenge := sha256.Sum256(proof[:])

			fmt.Println("nonce", "00")
			fmt.Println("proof", hex.EncodeToString(proof[:]))
			fmt.Println("challenge", hex.EncodeToString(challenge[:]))

			return nil
		},
	}

	return cmd
}
