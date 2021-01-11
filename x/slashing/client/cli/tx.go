package cli

import (
	"github.com/spf13/cobra"

	"github.com/KiraCore/sekai/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewTxCmd returns a root CLI command handler for all x/slashing transaction commands.
func NewTxCmd() *cobra.Command {
	slashingTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Slashing transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	slashingTxCmd.AddCommand(
		NewActivateTxCmd(),
		NewPauseTxCmd(),
		NewUnpauseTxCmd(),
	)
	return slashingTxCmd
}

// NewActivateTxCmd defines MsgActivate tx
func NewActivateTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activate",
		Args:  cobra.NoArgs,
		Short: "activate validator previously inactivated for downtime",
		Long: `activate an inactivated validator:

$ <appd> tx slashing activate --from mykey
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			valAddr := clientCtx.GetFromAddress()

			msg := types.NewMsgActivate(sdk.ValAddress(valAddr))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewPauseTxCmd defines MsgPause tx
func NewPauseTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause",
		Args:  cobra.NoArgs,
		Short: "pause validator",
		Long: `pause a validator before stopping of a node to avoid automatic inactivation:

$ <appd> tx slashing pause --from mykey
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			valAddr := clientCtx.GetFromAddress()

			msg := types.NewMsgPause(sdk.ValAddress(valAddr))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewUnpauseTxCmd defines MsgUnpause tx
func NewUnpauseTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unpause",
		Args:  cobra.NoArgs,
		Short: "unpause validator previously paused for downtime",
		Long: `unpause a paused validator:

$ <appd> tx slashing unpause --from mykey
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			valAddr := clientCtx.GetFromAddress()

			msg := types.NewMsgUnpause(sdk.ValAddress(valAddr))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
