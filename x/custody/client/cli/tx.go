package cli

import (
	"strconv"

	"github.com/KiraCore/sekai/x/custody/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
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
	txCmd.AddCommand(NewWhiteListTxCmd())

	return txCmd
}

func NewWhiteListTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        "whitelist",
		Short:                      "custody whitelist sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(GetTxAddToCustodyWhiteList())
	txCmd.AddCommand(GetTxRemoveFromCustodyWhiteList())
	txCmd.AddCommand(GetTxDropCustodyWhiteList())

	return txCmd
}

func GetTxAddToCustodyWhiteList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [addr]",
		Short: "Add new address to the custody whitelist",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			accAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgAddToCustodyWhiteList(
				clientCtx.FromAddress,
				accAddr,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxRemoveFromCustodyWhiteList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [addr]",
		Short: "Remove address from the custody whitelist",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			accAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgRemoveFromCustodyWhiteList(
				clientCtx.FromAddress,
				accAddr,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxDropCustodyWhiteList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drop",
		Short: "Drop the custody whitelist",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDropCustodyWhiteList(
				clientCtx.FromAddress,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
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
			useWhiteList, _ := strconv.ParseBool(args[3])

			msg := types.NewMsgCreateCustody(
				clientCtx.FromAddress,
				types.CustodySettings{
					Propagating:  propagating,
					Password:     password,
					UseLimit:     useLimit,
					UseWhiteList: useWhiteList,
				},
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
