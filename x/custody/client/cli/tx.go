package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/KiraCore/sekai/x/custody/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

const (
	OldKey = "okey"
	NewKey = "nkey"
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
	txCmd.AddCommand(NewCustodiansTxCmd())
	txCmd.AddCommand(NewWhiteListTxCmd())
	txCmd.AddCommand(NewLimitsTxCmd())
	txCmd.AddCommand(NewBankTxCmd())

	return txCmd
}

func NewLimitsTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        "limits",
		Short:                      "custody limits sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(GetTxAddToCustodyLimits())
	txCmd.AddCommand(GetTxRemoveFromCustodyLimits())
	txCmd.AddCommand(GetTxDropCustodyLimits())

	return txCmd
}

func NewCustodiansTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        "custodians",
		Short:                      "custody custodians sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(GetTxAddToCustodyCustodians())
	txCmd.AddCommand(GetTxRemoveFromCustodyCustodians())
	txCmd.AddCommand(GetTxDropCustodyCustodians())
	txCmd.AddCommand(GetTxAproveCustodyTransaction())
	txCmd.AddCommand(GetTxDeclineCustodyTransaction())

	return txCmd
}

func GetTxAddToCustodyCustodians() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [addr]",
		Short: "Add new address to the custody custodians",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var newAddr []sdk.AccAddress

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			addresses := strings.Split(args[0], ",")

			for _, addr := range addresses {
				accAddr, err := sdk.AccAddressFromBech32(addr)
				if err != nil {
					return err
				}

				newAddr = append(newAddr, accAddr)
			}

			oldKey, err := cmd.Flags().GetString(OldKey)
			if err != nil {
				return fmt.Errorf("invalid old key: %w", err)
			}

			newKey, err := cmd.Flags().GetString(NewKey)
			if err != nil {
				return fmt.Errorf("invalid new key: %w", err)
			}

			msg := types.NewMsgAddToCustodyCustodians(
				clientCtx.FromAddress,
				newAddr,
				oldKey,
				newKey,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(OldKey, "", "Previous hash string.")
	cmd.MarkFlagRequired(OldKey)

	cmd.Flags().String(NewKey, "", "Next hash string.")
	cmd.MarkFlagRequired(NewKey)

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxRemoveFromCustodyCustodians() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [addr]",
		Short: "Remove address from the custody custodians",
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

			oldKey, err := cmd.Flags().GetString(OldKey)
			if err != nil {
				return fmt.Errorf("invalid old key: %w", err)
			}

			newKey, err := cmd.Flags().GetString(NewKey)
			if err != nil {
				return fmt.Errorf("invalid new key: %w", err)
			}

			msg := types.NewMsgRemoveFromCustodyCustodians(
				clientCtx.FromAddress,
				accAddr,
				oldKey,
				newKey,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(OldKey, "", "Previous hash string.")
	cmd.MarkFlagRequired(OldKey)

	cmd.Flags().String(NewKey, "", "Next hash string.")
	cmd.MarkFlagRequired(NewKey)

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxDropCustodyCustodians() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drop",
		Short: "Drop the custody custodians",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			oldKey, err := cmd.Flags().GetString(OldKey)
			if err != nil {
				return fmt.Errorf("invalid old key: %w", err)
			}

			newKey, err := cmd.Flags().GetString(NewKey)
			if err != nil {
				return fmt.Errorf("invalid new key: %w", err)
			}

			msg := types.NewMsgDropCustodyCustodians(
				clientCtx.FromAddress,
				oldKey,
				newKey,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(OldKey, "", "Previous hash string.")
	cmd.MarkFlagRequired(OldKey)

	cmd.Flags().String(NewKey, "", "Next hash string.")
	cmd.MarkFlagRequired(NewKey)

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxAproveCustodyTransaction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve [hash]",
		Short: "Approve custody transaction by hash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgApproveCustodyTransaction(
				clientCtx.FromAddress,
				args[0],
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxDeclineCustodyTransaction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decline [hash]",
		Short: "Decline custody transaction by hash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeclineCustodyTransaction(
				clientCtx.FromAddress,
				args[0],
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
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
			var newAddr []sdk.AccAddress

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			addresses := strings.Split(args[0], "_")

			for _, addr := range addresses {
				accAddr, err := sdk.AccAddressFromBech32(addr)
				if err != nil {
					return err
				}

				newAddr = append(newAddr, accAddr)
			}

			oldKey, err := cmd.Flags().GetString(OldKey)
			if err != nil {
				return fmt.Errorf("invalid old key: %w", err)
			}

			newKey, err := cmd.Flags().GetString(NewKey)
			if err != nil {
				return fmt.Errorf("invalid new key: %w", err)
			}

			msg := types.NewMsgAddToCustodyWhiteList(
				clientCtx.FromAddress,
				newAddr,
				oldKey,
				newKey,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(OldKey, "", "Previous hash string.")
	cmd.MarkFlagRequired(OldKey)

	cmd.Flags().String(NewKey, "", "Next hash string.")
	cmd.MarkFlagRequired(NewKey)

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

			oldKey, err := cmd.Flags().GetString(OldKey)
			if err != nil {
				return fmt.Errorf("invalid old key: %w", err)
			}

			newKey, err := cmd.Flags().GetString(NewKey)
			if err != nil {
				return fmt.Errorf("invalid new key: %w", err)
			}

			msg := types.NewMsgRemoveFromCustodyWhiteList(
				clientCtx.FromAddress,
				accAddr,
				oldKey,
				newKey,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(OldKey, "", "Previous hash string.")
	cmd.MarkFlagRequired(OldKey)

	cmd.Flags().String(NewKey, "", "Next hash string.")
	cmd.MarkFlagRequired(NewKey)

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

			oldKey, err := cmd.Flags().GetString(OldKey)
			if err != nil {
				return fmt.Errorf("invalid old key: %w", err)
			}

			newKey, err := cmd.Flags().GetString(NewKey)
			if err != nil {
				return fmt.Errorf("invalid new key: %w", err)
			}

			msg := types.NewMsgDropCustodyWhiteList(
				clientCtx.FromAddress,
				oldKey,
				newKey,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(OldKey, "", "Previous hash string.")
	cmd.MarkFlagRequired(OldKey)

	cmd.Flags().String(NewKey, "", "Next hash string.")
	cmd.MarkFlagRequired(NewKey)

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxCreateCustody() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create ",
		Short: "Create new custody settings",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			custodyMode, _ := strconv.Atoi(args[0])
			usePassword, _ := strconv.ParseBool(args[1])
			useLimits, _ := strconv.ParseBool(args[2])
			useWhiteList, _ := strconv.ParseBool(args[3])

			oldKey, err := cmd.Flags().GetString(OldKey)
			if err != nil {
				return fmt.Errorf("invalid old key: %w", err)
			}

			newKey, err := cmd.Flags().GetString(NewKey)
			if err != nil {
				return fmt.Errorf("invalid new key: %w", err)
			}

			msg := types.NewMsgCreateCustody(
				clientCtx.FromAddress,
				types.CustodySettings{
					CustodyEnabled: true,
					CustodyMode:    uint64(custodyMode),
					UsePassword:    usePassword,
					UseLimits:      useLimits,
					UseWhiteList:   useWhiteList,
				},
				oldKey,
				newKey,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(OldKey, "", "Previous hash string.")
	cmd.MarkFlagRequired(OldKey)

	cmd.Flags().String(NewKey, "", "Next hash string.")
	cmd.MarkFlagRequired(NewKey)

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxAddToCustodyLimits() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [addr]",
		Short: "Add new address to the custody limits",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			oldKey, err := cmd.Flags().GetString(OldKey)
			if err != nil {
				return fmt.Errorf("invalid old key: %w", err)
			}

			newKey, err := cmd.Flags().GetString(NewKey)
			if err != nil {
				return fmt.Errorf("invalid new key: %w", err)
			}

			msg := types.NewMsgAddToCustodyLimits(
				clientCtx.FromAddress,
				args[0],
				uint64(amount),
				args[2],
				oldKey,
				newKey,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(OldKey, "", "Previous hash string.")
	cmd.MarkFlagRequired(OldKey)

	cmd.Flags().String(NewKey, "", "Next hash string.")
	cmd.MarkFlagRequired(NewKey)

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxRemoveFromCustodyLimits() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [addr]",
		Short: "Remove address from the custody limits",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			oldKey, err := cmd.Flags().GetString(OldKey)
			if err != nil {
				return fmt.Errorf("invalid old key: %w", err)
			}

			newKey, err := cmd.Flags().GetString(NewKey)
			if err != nil {
				return fmt.Errorf("invalid new key: %w", err)
			}

			msg := types.NewMsgRemoveFromCustodyLimits(
				clientCtx.FromAddress,
				args[0],
				oldKey,
				newKey,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(OldKey, "", "Previous hash string.")
	cmd.MarkFlagRequired(OldKey)

	cmd.Flags().String(NewKey, "", "Next hash string.")
	cmd.MarkFlagRequired(NewKey)

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxDropCustodyLimits() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drop",
		Short: "Drop the custody limits",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			oldKey, err := cmd.Flags().GetString(OldKey)
			if err != nil {
				return fmt.Errorf("invalid old key: %w", err)
			}

			newKey, err := cmd.Flags().GetString(NewKey)
			if err != nil {
				return fmt.Errorf("invalid new key: %w", err)
			}

			msg := types.NewMsgDropCustodyLimits(
				clientCtx.FromAddress,
				oldKey,
				newKey,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(OldKey, "", "Previous hash string.")
	cmd.MarkFlagRequired(OldKey)

	cmd.Flags().String(NewKey, "", "Next hash string.")
	cmd.MarkFlagRequired(NewKey)

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func NewBankTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "send [from_key_or_address] [to_address] [amount]",
		Short: `Send funds from one account to another. Note, the'--from' flag is
ignored as it is implied from [from_key_or_address].`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[0])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			toAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgSend(clientCtx.GetFromAddress(), toAddr, coins, args[3])

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
