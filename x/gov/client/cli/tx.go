package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	types2 "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/staking/client/cli"
)

// define flags
const (
	FlagPermission        = "permission"
	FlagMinTxFee          = "min_tx_fee"
	FlagMaxTxFee          = "max_tx_fee"
	FlagExecName          = "execution_name"
	FlagTxType            = "transaction_type"
	FlagExecutionFee      = "execution_fee"
	FlagFailureFee        = "failure_fee"
	FlagTimeout           = "timeout"
	FlagDefaultParameters = "default_parameters"
	FlagWebsite           = "website"
	FlagMoniker           = "moniker"
	FlagSocial            = "social"
	FlagIdentity          = "identity"
	FlagAddress           = "address"
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

	txCmd.AddCommand(
		GetTxSetWhitelistPermissions(),
		GetTxSetNetworkProperties(),
		GetTxSetExecutionFee(),
		GetTxSetBlacklistPermissions(),
		GetTxClaimGovernanceCmd(),
		GetTxBlacklistRolePermission(),
		GetTxWhitelistRolePermission(),
		GetTxSetWhitelistPermissions(),
		GetTxSetBlacklistPermissions(),
		GetTxClaimGovernanceCmd(),
		GetTxCreateRole(),
		GetTxRemoveRole(),
		GetTxBlacklistRolePermission(),
		GetTxWhitelistRolePermission(),
		GetTxRemoveWhitelistRolePermission(),
		GetTxRemoveBlacklistRolePermission(),
	)

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

// GetTxSetNetworkProperties is a function to set network properties tx command
func GetTxSetNetworkProperties() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-network-properties",
		Short: "Set network properties",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			minTxFee, err := cmd.Flags().GetUint64(FlagMinTxFee)
			if err != nil {
				return fmt.Errorf("invalid minimum tx fee")
			}
			maxTxFee, err := cmd.Flags().GetUint64(FlagMaxTxFee)
			if err != nil {
				return fmt.Errorf("invalid maximum tx fee")
			}

			msg := types.NewMsgSetNetworkProperties(
				clientCtx.FromAddress,
				&types.NetworkProperties{
					MinTxFee: minTxFee,
					MaxTxFee: maxTxFee,
				},
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().Uint64(FlagMinTxFee, 1, "min tx fee")
	cmd.Flags().Uint64(FlagMaxTxFee, 10000, "max tx fee")
	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxWhitelistRolePermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist-role-permissions role permission",
		Short: "Whitelist role permissions",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			role, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid role: %w", err)
			}

			permission, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid permission: %w", err)
			}

			msg := types.NewMsgWhitelistRolePermission(
				clientCtx.FromAddress,
				uint32(role),
				uint32(permission),
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxSetExecutionFee is a function to set network properties tx command
func GetTxSetExecutionFee() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-execution-fee",
		Short: "Set execution fee",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			execName, err := cmd.Flags().GetString(FlagExecName)
			if err != nil {
				return fmt.Errorf("invalid execution name")
			}
			txType, err := cmd.Flags().GetString(FlagTxType)
			if err != nil {
				return fmt.Errorf("invalid transaction type")
			}

			execFee, err := cmd.Flags().GetUint64(FlagExecutionFee)
			if err != nil {
				return fmt.Errorf("invalid execution fee")
			}
			failureFee, err := cmd.Flags().GetUint64(FlagFailureFee)
			if err != nil {
				return fmt.Errorf("invalid failure fee")
			}
			timeout, err := cmd.Flags().GetUint64(FlagTimeout)
			if err != nil {
				return fmt.Errorf("invalid timeout")
			}
			defaultParams, err := cmd.Flags().GetUint64(FlagDefaultParameters)
			if err != nil {
				return fmt.Errorf("invalid default parameters")
			}

			msg := types.NewMsgSetExecutionFee(
				execName,
				txType,
				execFee,
				failureFee,
				timeout,
				defaultParams,
				clientCtx.FromAddress,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().String(FlagExecName, "", "execution name")
	cmd.Flags().String(FlagTxType, "", "execution type")
	cmd.Flags().Uint64(FlagExecutionFee, 10, "execution fee")
	cmd.Flags().Uint64(FlagFailureFee, 1, "failure fee")
	cmd.Flags().Uint64(FlagTimeout, 0, "timeout")
	cmd.Flags().Uint64(FlagDefaultParameters, 0, "default parameters")
	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxBlacklistRolePermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blacklist-role-permissions role permission",
		Short: "Blacklist role permissions",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			role, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid role: %w", err)
			}

			permission, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid permission: %w", err)
			}

			msg := types.NewMsgBlacklistRolePermission(
				clientCtx.FromAddress,
				uint32(role),
				uint32(permission),
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxRemoveWhitelistRolePermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-whitelist-role-permissions role permission",
		Short: "Remove whitelist role permissions",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			role, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid role: %w", err)
			}

			permission, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid permission: %w", err)
			}

			msg := types.NewMsgRemoveWhitelistRolePermission(
				clientCtx.FromAddress,
				uint32(role),
				uint32(permission),
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxRemoveBlacklistRolePermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-blacklist-role-permissions role permission",
		Short: "Remove blacklist role permissions",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			role, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid role: %w", err)
			}

			permission, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid permission: %w", err)
			}

			msg := types.NewMsgRemoveBlacklistRolePermission(
				clientCtx.FromAddress,
				uint32(role),
				uint32(permission),
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxCreateRole() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-role role",
		Short: "Create new role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			role, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid role: %w", err)
			}

			msg := types.NewMsgCreateRole(
				clientCtx.FromAddress,
				uint32(role),
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxAssignRole() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assign-role role",
		Short: "Assign new role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			role, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid role: %w", err)
			}

			addr, err := getAddressFromFlag(cmd)
			if err != nil {
				return fmt.Errorf("error getting address: %w", err)
			}

			msg := types.NewMsgAssignRole(
				clientCtx.FromAddress,
				addr,
				uint32(role),
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(cli.FlagAddr, "", "the address to set permissions")

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(cli.FlagAddr)

	return cmd
}

func GetTxRemoveRole() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-role role",
		Short: "Remove new role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			role, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid role: %w", err)
			}

			addr, err := getAddressFromFlag(cmd)
			if err != nil {
				return fmt.Errorf("error getting address: %w", err)
			}

			msg := types.NewMsgRemoveRole(
				clientCtx.FromAddress,
				addr,
				uint32(role),
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(cli.FlagAddr, "", "the address to set permissions")

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(cli.FlagAddr)

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

func GetTxClaimGovernanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-councilor-seat",
		Short: "Claim governance seat to become a Councilor",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			moniker, _ := cmd.Flags().GetString(FlagMoniker)
			website, _ := cmd.Flags().GetString(FlagWebsite)
			social, _ := cmd.Flags().GetString(FlagSocial)
			identity, _ := cmd.Flags().GetString(FlagIdentity)
			address, _ := cmd.Flags().GetString(FlagAddress)

			bech32, err := types2.AccAddressFromBech32(address)
			if err != nil {
				return err
			}

			msg := types.NewMsgClaimCouncilor(
				moniker,
				website,
				social,
				identity,
				bech32,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().String(FlagMoniker, "", "the Moniker")
	cmd.Flags().String(FlagWebsite, "", "the Website")
	cmd.Flags().String(FlagSocial, "", "the social")
	cmd.Flags().String(FlagIdentity, "", "the Identity")
	cmd.Flags().String(FlagAddress, "", "the address")

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
