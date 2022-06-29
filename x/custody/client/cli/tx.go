package cli

import (
	"github.com/KiraCore/sekai/x/custody/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"strconv"
)

// NewTxCmd returns a root CLI command handler for all x/bank transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "custody sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(GetTxCreateCustody())

	return txCmd
}

func GetTxCreateCustody() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create new custody settings",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			propagating, _ := strconv.ParseBool(args[0])
			password, _ := strconv.ParseBool(args[1])
			useLimit, _ := strconv.ParseBool(args[2])
			useWhitelist, _ := strconv.ParseBool(args[3])

			msg := types.NewMsgCreateCustody(
				clientCtx.FromAddress,
				types.CustodySettings{
					Propagating:  propagating,
					Password:     password,
					UseLimit:     useLimit,
					UseWhitelist: useWhitelist,
				},
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
