package cli

import (
	"strconv"

	"github.com/KiraCore/sekai/x/basket/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

// NewTxCmd returns a root CLI command handler for all x/bank transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Tokens sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetTxDisableBasketDepositsCmd(),
	)

	return txCmd
}

// GetTxDisableBasketDepositsCmd implement cli command for MsgDisableBasketDeposits
func GetTxDisableBasketDepositsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable-basket-deposits [basket_id] [disabled]",
		Short: "Emergency function & permission to disable one or all deposits of one or all token in the basket",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			basketId, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			disabled, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDisableBasketDeposits(
				clientCtx.FromAddress,
				uint64(basketId),
				disabled,
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
