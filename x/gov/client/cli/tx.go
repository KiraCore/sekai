package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/staking/client/cli"
)

// define flags
const (
	FlagTitle             = "title"
	FlagDescription       = "description"
	FlagPermission        = "permission"
	FlagMinTxFee          = "min_tx_fee"
	FlagMaxTxFee          = "max_tx_fee"
	FlagMinValidators     = "min_validators"
	FlagExecName          = "execution_name"
	FlagTxType            = "transaction_type"
	FlagExecutionFee      = "execution_fee"
	FlagFailureFee        = "failure_fee"
	FlagTimeout           = "timeout"
	FlagDefaultParameters = "default_parameters"
	FlagMoniker           = "moniker"
	FlagAddress           = "address"
	FlagWhitelistPerms    = "whitelist"
	FlagBlacklistPerms    = "blacklist"
	FlagInfosFile         = "infos-file"
	FlagInfosJson         = "infos-json"
	FlagKeys              = "keys"
	FlagVerifier          = "verifier"
	FlagRecordIds         = "record-ids"
	FlagTip               = "tip"
	FlagApprove           = "approve"
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
		NewTxCouncilorCmds(),
		NewTxProposalCmds(),
		NewTxRoleCmds(),
		NewTxPermissionCmds(),
		NewTxSetNetworkProperties(),
		NewTxSetExecutionFee(),
		GetTxRegisterIdentityRecords(),
		GetTxDeleteIdentityRecords(),
		GetTxRequestIdentityRecordsVerify(),
		GetTxHandleIdentityRecordsVerifyRequest(),
		GetTxCancelIdentityRecordsVerifyRequest(),
	)

	return txCmd
}

// NewTxProposalCmds returns the subcommands of proposal related commands.
func NewTxProposalCmds() *cobra.Command {
	proposalCmd := &cobra.Command{
		Use:                        "proposal",
		Short:                      "Proposal subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	proposalCmd.AddCommand(GetTxProposalAssignPermission())
	proposalCmd.AddCommand(GetTxVoteProposal())
	proposalCmd.AddCommand(GetTxProposalSetNetworkProperty())
	proposalCmd.AddCommand(GetTxProposalSetPoorNetworkMessages())
	proposalCmd.AddCommand(GetTxProposalCreateRole())
	proposalCmd.AddCommand(GetTxProposalUpsertDataRegistry())

	return proposalCmd
}

// NewTxRoleCmds returns the subcommands of role related commands.
func NewTxRoleCmds() *cobra.Command {
	roleCmd := &cobra.Command{
		Use:                        "role",
		Short:                      "Role subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	roleCmd.AddCommand(GetTxCreateRole())
	roleCmd.AddCommand(GetTxRemoveRole())

	roleCmd.AddCommand(GetTxBlacklistRolePermission())
	roleCmd.AddCommand(GetTxWhitelistRolePermission())
	roleCmd.AddCommand(GetTxRemoveWhitelistRolePermission())
	roleCmd.AddCommand(GetTxRemoveBlacklistRolePermission())

	return roleCmd
}

// NewTxPermissionCmds returns the subcommands of permission related commands.
func NewTxPermissionCmds() *cobra.Command {
	permCmd := &cobra.Command{
		Use:                        "permission",
		Short:                      "Permission subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	permCmd.AddCommand(GetTxSetWhitelistPermissions())
	permCmd.AddCommand(GetTxSetBlacklistPermissions())

	return permCmd
}

func NewTxCouncilorCmds() *cobra.Command {
	councilor := &cobra.Command{
		Use:                        "councilor",
		Short:                      "Councilor subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	councilor.AddCommand(GetTxClaimCouncilorSeatCmd())

	return councilor
}

func GetTxSetWhitelistPermissions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist-permission",
		Short: "Whitelists permission into an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

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
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxSetBlacklistPermissions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blacklist-permission",
		Short: "Blacklist permission into an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
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
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// NewTxSetNetworkProperties is a function to set network properties tx command
func NewTxSetNetworkProperties() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-network-properties",
		Short: "Submit a transaction to set network properties",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
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
			minValidators, err := cmd.Flags().GetUint64(FlagMinValidators)
			if err != nil {
				return fmt.Errorf("invalid min validators")
			}

			// TODO: should set more by flags
			msg := types.NewMsgSetNetworkProperties(
				clientCtx.FromAddress,
				&types.NetworkProperties{
					MinTxFee:                    minTxFee,
					MaxTxFee:                    maxTxFee,
					VoteQuorum:                  33,
					ProposalEndTime:             1, // 1min
					ProposalEnactmentTime:       2, // 2min
					EnableForeignFeePayments:    true,
					MischanceRankDecreaseAmount: 10,
					InactiveRankDecreasePercent: 50,      // 50%
					PoorNetworkMaxBankSend:      1000000, // 1M ukex
					MinValidators:               minValidators,
				},
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().Uint64(FlagMinTxFee, 1, "min tx fee")
	cmd.Flags().Uint64(FlagMaxTxFee, 10000, "max tx fee")
	cmd.Flags().Uint64(FlagMinValidators, 2, "min validators")
	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxWhitelistRolePermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist-permission role permission",
		Short: "Whitelist a permission to a role",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
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
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// NewTxSetExecutionFee is a function to set network properties tx command
func NewTxSetExecutionFee() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-execution-fee",
		Short: "Submit a transaction to set execution fee",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
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
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxBlacklistRolePermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blacklist-permission role permission",
		Short: "Blacklist a permission on a role",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
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
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxRemoveWhitelistRolePermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-whitelist-permission role permission",
		Short: "Remove a whitelisted permission from a role",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
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
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxRemoveBlacklistRolePermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-blacklist-permission role permission",
		Short: "Remove a blacklisted permission from a role",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
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
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxCreateRole() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create role",
		Short: "Create new role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
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
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxAssignRole() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assign-role role",
		Short: "Assign new role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
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
				clientCtx.GetFromAddress(),
				addr,
				uint32(role),
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(cli.FlagAddr, "", "the address to set permissions")

	cmd.MarkFlagRequired(flags.FlagFrom)
	cmd.MarkFlagRequired(cli.FlagAddr)

	return cmd
}

func GetTxRemoveRole() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove role",
		Short: "Remove role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
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

	cmd.MarkFlagRequired(flags.FlagFrom)
	cmd.MarkFlagRequired(cli.FlagAddr)

	return cmd
}

// GetTxProposalSetPoorNetworkMessages defines command to send proposal tx to modify poor network messages
func GetTxProposalSetPoorNetworkMessages() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-poor-network-msgs <messages>",
		Short: "Create a proposal to set a value on a network property.",
		Long: `
		$ %s tx customgov proposal set-poor-network-msgs XXXX,YYY --from=<key_or_address>

		All the message types supported could be added here
			create-role
			assign-role
			remove-role
			...
		`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			messages := strings.Split(args[0], ",")

			title, err := cmd.Flags().GetString(FlagTitle)
			if err != nil {
				return fmt.Errorf("invalid title: %w", err)
			}

			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return fmt.Errorf("invalid description: %w", err)
			}

			msg, err := types.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				types.NewSetPoorNetworkMessagesProposal(messages),
			)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagTitle, "", "The title of the proposal.")
	cmd.MarkFlagRequired(FlagTitle)
	cmd.Flags().String(FlagDescription, "", "The description of the proposal, it can be a url, some text, etc.")
	cmd.MarkFlagRequired(FlagDescription)

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxProposalSetNetworkProperty defines command to send proposal tx to modify a network property
func GetTxProposalSetNetworkProperty() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-network-property <property> <value> [flags]",
		Short: "Create a proposal to set a value on a network property.",
		Long: `
		$ %s tx customgov proposal set-network-property MIN_TX_FEE 100 --from=<key_or_address>

		Available properties:
			MIN_TX_FEE
			MAX_TX_FEE
			VOTE_QUORUM
			PROPOSAL_END_TIME
			PROPOSAL_ENACTMENT_TIME
			ENABLE_FOREIGN_TX_FEE_PAYMENTS
		`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			property, ok := types.NetworkProperty_value[args[0]]
			if !ok {
				return fmt.Errorf("invalid network property name: %s", args[0])
			}

			value := types.NetworkPropertyValue{}
			if property == int32(types.UniqueIdentityKeys) {
				value.StrValue = args[1]
			} else {
				numVal, err := strconv.Atoi(args[1])
				if err != nil {
					return fmt.Errorf("invalid network property value: %w", err)
				}
				value.Value = uint64(numVal)
			}

			title, err := cmd.Flags().GetString(FlagTitle)
			if err != nil {
				return fmt.Errorf("invalid title: %w", err)
			}

			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return fmt.Errorf("invalid description: %w", err)
			}

			msg, err := types.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				types.NewSetNetworkPropertyProposal(types.NetworkProperty(property), value),
			)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagTitle, "", "The title of the proposal.")
	cmd.MarkFlagRequired(FlagTitle)
	cmd.Flags().String(FlagDescription, "", "The description of the proposal, it can be a url, some text, etc.")
	cmd.MarkFlagRequired(FlagDescription)

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxProposalAssignPermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assign-permission permission",
		Short: "Create a proposal to assign a permission to an address.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			perm, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid perm: %w", err)
			}

			addr, err := getAddressFromFlag(cmd)
			if err != nil {
				return fmt.Errorf("error getting address: %w", err)
			}

			title, err := cmd.Flags().GetString(FlagTitle)
			if err != nil {
				return fmt.Errorf("invalid title: %w", err)
			}

			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return fmt.Errorf("invalid description: %w", err)
			}

			msg, err := types.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				types.NewAssignPermissionProposal(addr, types.PermValue(perm)),
			)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagTitle, "", "The title of the proposal.")
	cmd.MarkFlagRequired(FlagTitle)
	cmd.Flags().String(FlagDescription, "", "The description of the proposal, it can be a url, some text, etc.")
	cmd.MarkFlagRequired(FlagDescription)

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(cli.FlagAddr, "", "the address to set permissions")

	cmd.MarkFlagRequired(flags.FlagFrom)
	cmd.MarkFlagRequired(cli.FlagAddr)

	return cmd
}

func GetTxProposalUpsertDataRegistry() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upsert-data-registry [key] [hash] [reference] [encoding] [size] [flags]",
		Short: "Create a proposal to upsert a key in the data registry",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			key := args[0]
			hash := args[1]
			reference := args[2]
			encoding := args[3]
			size, err := strconv.Atoi(args[4])
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(FlagTitle)
			if err != nil {
				return fmt.Errorf("invalid title: %w", err)
			}

			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return fmt.Errorf("invalid description: %w", err)
			}

			msg, err := types.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				types.NewUpsertDataRegistryProposal(
					key,
					hash,
					reference,
					encoding,
					uint64(size),
				),
			)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagTitle, "", "The title of the proposal.")
	cmd.Flags().String(FlagDescription, "", "The description of the proposal, it can be a url, some text, etc.")
	flags.AddTxFlagsToCmd(cmd)

	cmd.MarkFlagRequired(FlagTitle)
	cmd.MarkFlagRequired(FlagDescription)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxVoteProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote proposal-id vote-option",
		Short: "Vote a proposal.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposalID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid proposal ID: %w", err)
			}

			voteOption, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid vote option: %w", err)
			}

			msg := types.NewMsgVoteProposal(
				uint64(proposalID),
				clientCtx.FromAddress,
				types.VoteOption(voteOption),
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// setPermissionFlags sets the flags needed for set blacklist and set whitelist permission
// commands.
func setPermissionFlags(cmd *cobra.Command) {
	cmd.Flags().String(cli.FlagAddr, "", "the address to set permissions")
	cmd.Flags().Uint32(FlagPermission, 0, "the permission")
}

// getAddressFromFlag returns the AccAddress from FlagAddr in Command.
func getAddressFromFlag(cmd *cobra.Command) (sdk.AccAddress, error) {
	addr, err := cmd.Flags().GetString(cli.FlagAddr)
	if err != nil {
		return nil, fmt.Errorf("error getting address")
	}

	bech, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid address")
	}

	return bech, nil
}

func GetTxClaimCouncilorSeatCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-seat",
		Short: "Claim councilor seat",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			moniker, _ := cmd.Flags().GetString(FlagMoniker)
			address, _ := cmd.Flags().GetString(FlagAddress)

			bech32, err := sdk.AccAddressFromBech32(address)
			if err != nil {
				return err
			}

			msg := types.NewMsgClaimCouncilor(
				moniker,
				bech32,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().String(FlagMoniker, "", "the Moniker")
	cmd.Flags().String(FlagAddress, "", "the address")

	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxProposalCreateRole() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-role role",
		Short: "Create a proposal to add a new role.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			role, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid perm: %w", err)
			}

			wAsInts, err := cmd.Flags().GetInt32Slice(FlagWhitelistPerms)
			if err != nil {
				return fmt.Errorf("invalid whitelist perms: %w", err)
			}
			whitelistPerms := convertAsPermValues(wAsInts)

			bAsInts, err := cmd.Flags().GetInt32Slice(FlagBlacklistPerms)
			if err != nil {
				return fmt.Errorf("invalid blacklist perms: %w", err)
			}
			blacklistPerms := convertAsPermValues(bAsInts)

			title, err := cmd.Flags().GetString(FlagTitle)
			if err != nil {
				return fmt.Errorf("invalid title: %w", err)
			}

			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return fmt.Errorf("invalid description: %w", err)
			}

			msg, err := types.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				types.NewCreateRoleProposal(
					types.Role(role),
					whitelistPerms,
					blacklistPerms,
				),
			)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().String(FlagTitle, "", "The title of the proposal.")
	cmd.MarkFlagRequired(FlagTitle)
	cmd.Flags().String(FlagDescription, "", "The description of the proposal, it can be a url, some text, etc.")
	cmd.MarkFlagRequired(FlagDescription)
	cmd.Flags().Int32Slice(FlagWhitelistPerms, []int32{}, "the whitelist value in format 1,2,3")
	cmd.Flags().Int32Slice(FlagBlacklistPerms, []int32{}, "the blacklist values in format 1,2,3")
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxRegisterIdentityRecords() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-identity-records",
		Short: "Submit a transaction to create an identity record.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			infos, err := parseIdInfoJSON(cmd.Flags())
			if err != nil {
				return err
			}

			msg := types.NewMsgRegisterIdentityRecords(
				clientCtx.FromAddress,
				infos,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().String(FlagInfosFile, "", "The infos file for identity request.")
	cmd.Flags().String(FlagInfosJson, "", "The infos json for identity request.")
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxDeleteIdentityRecords() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-identity-records",
		Short: "Submit a transaction to delete an identity records.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			keysStr, err := cmd.Flags().GetString(FlagKeys)
			if err != nil {
				return err
			}

			keys := strings.Split(keysStr, ",")
			if keysStr == "" {
				keys = []string{}
			}

			msg := types.NewMsgDeleteIdentityRecords(
				clientCtx.FromAddress,
				keys,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().String(FlagKeys, "", "The keys to remove.")
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxRequestIdentityRecordsVerify() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request-identity-record-verify",
		Short: "Submit a transaction to request an identity verify record.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			verifierText, err := cmd.Flags().GetString(FlagVerifier)
			if err != nil {
				return err
			}
			verifier, err := sdk.AccAddressFromBech32(verifierText)

			recordIdsStr, err := cmd.Flags().GetString(FlagRecordIds)
			if err != nil {
				return err
			}

			recordIdsSplit := strings.Split(recordIdsStr, ",")
			recordIds := []uint64{}
			for _, str := range recordIdsSplit {
				id, err := strconv.ParseUint(str, 10, 64)
				if err != nil {
					return err
				}
				recordIds = append(recordIds, id)
			}

			tipStr, err := cmd.Flags().GetString(FlagTip)
			if err != nil {
				return err
			}
			coin, err := sdk.ParseCoinNormalized(tipStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgRequestIdentityRecordsVerify(
				clientCtx.FromAddress,
				verifier,
				recordIds,
				coin,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().String(FlagTip, "", "The tip to be given to the verifier.")
	cmd.Flags().String(FlagRecordIds, "", "Concatenated identity record ids array. e.g. 1,2")
	cmd.Flags().String(FlagVerifier, "", "The verifier of the record ids")
	cmd.MarkFlagRequired(FlagRecordIds)
	cmd.MarkFlagRequired(FlagVerifier)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxHandleIdentityRecordsVerifyRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "handle-identity-records-verify-request [id]",
		Short: "Submit a transaction to approve or reject identity records verify request.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			requestId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			isApprove, err := cmd.Flags().GetBool(FlagApprove)
			if err != nil {
				return err
			}

			msg := types.NewMsgHandleIdentityRecordsVerifyRequest(
				clientCtx.FromAddress,
				requestId,
				isApprove,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Bool(FlagApprove, true, "The flag to approve or reject")
	cmd.MarkFlagRequired(FlagApprove)
	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxCancelIdentityRecordsVerifyRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-identity-records-verify-request [id]",
		Short: "Submit a transaction to cancel identity records verification request.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			requestId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelIdentityRecordsVerifyRequest(
				clientCtx.FromAddress,
				requestId,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func parseIdInfoJSON(fs *pflag.FlagSet) ([]types.IdentityInfoEntry, error) {
	var err error
	infos := make(map[string]string)
	infosFile, _ := fs.GetString(FlagInfosFile)
	infosJson, _ := fs.GetString(FlagInfosJson)

	if infosFile == "" && infosJson == "" {
		return nil, fmt.Errorf("should input infos file json using the --%s flag or infos json using the --%s flag", FlagInfosFile, FlagInfosJson)
	}

	if infosFile != "" && infosJson != "" {
		return nil, fmt.Errorf("should only set one of --%s flag or --%s flag", FlagInfosFile, FlagInfosJson)
	}

	contents := []byte(infosJson)

	if infosFile != "" {
		contents, err = ioutil.ReadFile(infosFile)
		if err != nil {
			return nil, err
		}
	}

	// make exception if unknown field exists
	err = json.Unmarshal(contents, &infos)
	if err != nil {
		return nil, err
	}

	return types.WrapInfos(infos), nil
}

// convertAsPermValues convert array of int32 to PermValue array.
func convertAsPermValues(values []int32) []types.PermValue {
	var v []types.PermValue
	for _, perm := range values {
		v = append(v, types.PermValue(perm))
	}

	return v
}
