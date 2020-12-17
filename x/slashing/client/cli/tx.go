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

	slashingTxCmd.AddCommand(NewActivateTxCmd())
	return slashingTxCmd
}

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
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

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
