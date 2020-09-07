package cli

import (
	"fmt"

	types2 "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/staking/client/cli"

	"github.com/KiraCore/sekai/x/gov/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

const (
	FlagPermission = "permission"
)

// NewTxCmd returns a root CLI command handler for all x/bank transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Custom gov sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(GetTxSetWhitelistPermissions())

	return txCmd
}

func GetTxSetWhitelistPermissions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-whitelist-permissions",
		Short: "Whitelists permissions into an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			perm, err := cmd.Flags().GetUint32(FlagPermission)
			if err != nil {
				return fmt.Errorf("invalid permissions")
			}

			addr, err := getAddressFromFlag(cmd)
			if err != nil {
				return fmt.Errorf("error getting address: %w", err)
			}

			msg := types.NewMsgWhitelistPermissions(
				clientCtx.FromAddress,
				addr,
				perm,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	setPermissionFlags(cmd)

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxSetBlacklistPermissions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-blacklist-permissions",
		Short: "Blacklist permissions into an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			perm, err := cmd.Flags().GetUint32(FlagPermission)
			if err != nil {
				return fmt.Errorf("invalid permissions")
			}

			addr, err := getAddressFromFlag(cmd)
			if err != nil {
				return fmt.Errorf("error getting address: %w", err)
			}

			msg := types.NewMsgBlacklistPermissions(
				clientCtx.FromAddress,
				addr,
				perm,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	setPermissionFlags(cmd)

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// setPermissionFlags sets the flags needed for set blacklist and set whitelist permission
// commands.
func setPermissionFlags(cmd *cobra.Command) {
	cmd.Flags().String(cli.FlagAddr, "", "the address to set permissions")
	cmd.Flags().Uint32(FlagPermission, 0, "the permission")
}

// getAddressFromFlag returns the AccAddress from FlagAddr in Command.
func getAddressFromFlag(cmd *cobra.Command) (types2.AccAddress, error) {
	addr, err := cmd.Flags().GetString(cli.FlagAddr)
	if err != nil {
		return nil, fmt.Errorf("error getting address")
	}

	bech, err := types2.AccAddressFromBech32(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid address")
	}

	return bech, nil
}
