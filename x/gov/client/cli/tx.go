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
	FlagMinCustodyReward  = "min_custody_reward"
	FlagTxType            = "transaction_type"
	FlagExecutionFee      = "execution_fee"
	FlagFailureFee        = "failure_fee"
	FlagTimeout           = "timeout"
	FlagDefaultParameters = "default_parameters"
	FlagMoniker           = "moniker"
	FlagAddr              = "addr"
	FlagWhitelistPerms    = "whitelist"
	FlagBlacklistPerms    = "blacklist"
	FlagInfosFile         = "infos-file"
	FlagInfosJson         = "infos-json"
	FlagKeys              = "keys"
	FlagVerifier          = "verifier"
	FlagRecordIds         = "record-ids"
	FlagVerifierTip       = "verifier-tip"
	FlagApprove           = "approve"
	FlagSlash             = "slash"
	FlagUsername          = "username"
	FlagSocial            = "social"
	FlagContact           = "contact"
	FlagAvatar            = "avatar"
	FlagPollRoles         = "poll-roles"
	FlagPollOptions       = "poll-options"
	FlagPollCount         = "poll-count"
	FlagPollType          = "poll-type"
	FlagPollChoices       = "poll-choices"
	FlagPollDuration      = "poll-duration"
	FlagPollReference     = "poll-reference"
	FlagPollChecksum      = "poll-checksum"
	FlagCustomPollValue   = "poll-custom-value"
	FlagTxTypes           = "tx-types"
	FlagExecutionFees     = "execution-fees"
	FlagFailureFees       = "failure-fees"
	FlagTimeouts          = "timeouts"
	FlagDefaultParams     = "default-params"
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
		NewTxPollCmds(),
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

// NewTxPollCmds returns the subcommands of poll related commands.
func NewTxPollCmds() *cobra.Command {
	pollCmd := &cobra.Command{
		Use:                        "poll",
		Short:                      "Governance poll management subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	pollCmd.AddCommand(GetTxPollCreate())
	pollCmd.AddCommand(GetTxVotePoll())

	return pollCmd
}

// NewTxProposalCmds returns the subcommands of proposal related commands.
func NewTxProposalCmds() *cobra.Command {
	proposalCmd := &cobra.Command{
		Use:                        "proposal",
		Short:                      "Governance proposals management subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	accountProposalCmd := &cobra.Command{
		Use:                        "account",
		Short:                      "Account proposals management subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	accountProposalCmd.AddCommand(GetTxProposalWhitelistAccountPermission())
	accountProposalCmd.AddCommand(GetTxProposalBlacklistAccountPermission())
	accountProposalCmd.AddCommand(GetTxProposalRemoveWhitelistedAccountPermission())
	accountProposalCmd.AddCommand(GetTxProposalRemoveBlacklistedAccountPermission())
	accountProposalCmd.AddCommand(GetTxProposalAssignRoleToAccount())
	accountProposalCmd.AddCommand(GetTxProposalUnassignRoleFromAccount())

	roleProposalCmd := &cobra.Command{
		Use:                        "role",
		Short:                      "Role proposals management subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	roleProposalCmd.AddCommand(GetTxProposalCreateRole())
	roleProposalCmd.AddCommand(GetTxProposalRemoveRole())
	roleProposalCmd.AddCommand(GetTxProposalWhitelistRolePermission())
	roleProposalCmd.AddCommand(GetTxProposalBlacklistRolePermission())
	roleProposalCmd.AddCommand(GetTxProposalRemoveWhitelistedRolePermission())
	roleProposalCmd.AddCommand(GetTxProposalRemoveBlacklistedRolePermission())

	proposalCmd.AddCommand(GetTxVoteProposal())
	proposalCmd.AddCommand(GetTxProposalSetNetworkProperty())
	proposalCmd.AddCommand(GetTxProposalSetPoorNetworkMessages())
	proposalCmd.AddCommand(GetTxProposalUpsertDataRegistry())
	proposalCmd.AddCommand(GetTxProposalSetProposalDurations())
	proposalCmd.AddCommand(GetTxProposalResetWholeCouncilorRankCmd())
	proposalCmd.AddCommand(GetTxProposalJailCouncilorCmd())
	proposalCmd.AddCommand(GetTxProposalSetExecutionFeesCmd())

	proposalCmd.AddCommand(accountProposalCmd)
	proposalCmd.AddCommand(roleProposalCmd)

	return proposalCmd
}

// NewTxRoleCmds returns the subcommands of role related commands.
func NewTxRoleCmds() *cobra.Command {
	roleCmd := &cobra.Command{
		Use:                        "role",
		Short:                      "Role management subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	roleCmd.AddCommand(GetTxCreateRole())
	roleCmd.AddCommand(GetTxAssignRole())
	roleCmd.AddCommand(GetTxUnassignRole())

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
		Short:                      "Permission management subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	permCmd.AddCommand(
		GetTxSetWhitelistPermissions(),
		GetTxRemoveWhitelistedPermissions(),
		GetTxSetBlacklistPermissions(),
		GetTxRemoveBlacklistedPermissions(),
	)

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

	councilor.AddCommand(
		GetTxClaimCouncilorSeatCmd(),
		GetTxCouncilorPauseCmd(),
		GetTxCouncilorUnpauseCmd(),
		GetTxCouncilorActivateCmd(),
	)

	return councilor
}

func GetTxSetWhitelistPermissions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist",
		Short: "Assign permission to a kira address whitelist",
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

func GetTxRemoveWhitelistedPermissions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-whitelisted",
		Short: "Remove whitelisted permission from an address",
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

			msg := types.NewMsgRemoveWhitelistedPermissions(
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
		Use:   "blacklist",
		Short: "Assign permission to a kira account blacklist",
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

func GetTxRemoveBlacklistedPermissions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-blacklisted",
		Short: "Remove blacklisted permission from an address",
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

			msg := types.NewMsgRemoveBlacklistedPermissions(
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
			minCustodyReward, err := cmd.Flags().GetUint64(FlagMinCustodyReward)
			if err != nil {
				return fmt.Errorf("invalid min custody reward")
			}

			// TODO: should set more by flags
			msg := types.NewMsgSetNetworkProperties(
				clientCtx.FromAddress,
				&types.NetworkProperties{
					MinTxFee:                        minTxFee,
					MaxTxFee:                        maxTxFee,
					VoteQuorum:                      sdk.NewDecWithPrec(33, 2), // 33%
					MinimumProposalEndTime:          300,                       // 5min
					ProposalEnactmentTime:           300,                       // 5min
					EnableForeignFeePayments:        true,
					MischanceRankDecreaseAmount:     10,
					InactiveRankDecreasePercent:     sdk.NewDecWithPrec(50, 2), // 50%
					PoorNetworkMaxBankSend:          1000000,                   // 1M ukex
					MinValidators:                   minValidators,
					MinCustodyReward:                minCustodyReward,
					MinProposalEndBlocks:            2,
					MinProposalEnactmentBlocks:      1,
					MaxMischance:                    1,
					MinIdentityApprovalTip:          200,
					UniqueIdentityKeys:              "moniker,username",
					UbiHardcap:                      6000_000,
					ValidatorsFeeShare:              sdk.NewDecWithPrec(50, 2), // 50%
					InflationRate:                   sdk.NewDecWithPrec(18, 2), // 18%
					InflationPeriod:                 31557600,                  // 1 year
					UnstakingPeriod:                 2629800,                   // 1 month
					MaxDelegators:                   100,
					MinDelegationPushout:            10,
					SlashingPeriod:                  2629800,
					MaxJailedPercentage:             sdk.NewDecWithPrec(25, 2),
					MaxSlashingPercentage:           sdk.NewDecWithPrec(5, 3), // 0.5%
					MaxCustodyBufferSize:            10,
					MaxCustodyTxSize:                8192,
					AbstentionRankDecreaseAmount:    1,
					MaxAbstention:                   2,
					MinCollectiveBond:               100_000, // in KEX
					MinCollectiveBondingTime:        86400,   // in seconds
					MaxCollectiveOutputs:            10,
					MinCollectiveClaimPeriod:        14400,                     // 4hrs
					ValidatorRecoveryBond:           300000,                    // 300k KEX
					MaxAnnualInflation:              sdk.NewDecWithPrec(35, 2), // 35%
					MaxProposalTitleSize:            128,
					MaxProposalDescriptionSize:      1024,
					MaxProposalPollOptionSize:       64,
					MaxProposalPollOptionCount:      128,
					MaxProposalReferenceSize:        512,
					MaxProposalChecksumSize:         128,
					MinDappBond:                     1000000,
					MaxDappBond:                     10000000,
					DappBondDuration:                604800,
					DappVerifierBond:                sdk.NewDecWithPrec(1, 3), //0.1%
					DappAutoDenounceTime:            60,                       // 60s
					DappMischanceRankDecreaseAmount: 1,
					DappMaxMischance:                10,
					DappInactiveRankDecreasePercent: sdk.NewDecWithPrec(10, 2), // 10%
					DappPoolSlippageDefault:         sdk.NewDecWithPrec(1, 1),  // 10%
					DappLiquidationThreshold:        100_000_000_000,           // default 100’000 KEX
					DappLiquidationPeriod:           2419200,                   // default 2419200, ~28d
					MintingFtFee:                    100_000_000_000_000,
					MintingNftFee:                   100_000_000_000_000,
					VetoThreshold:                   sdk.NewDecWithPrec(3340, 4), // 33.40%
					AutocompoundIntervalNumBlocks:   17280,
					DowntimeInactiveDuration:        600,
				},
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().Uint64(FlagMinTxFee, 1, "min tx fee")
	cmd.Flags().Uint64(FlagMaxTxFee, 10000, "max tx fee")
	cmd.Flags().Uint64(FlagMinValidators, 2, "min validators")
	cmd.Flags().Uint64(FlagMinCustodyReward, 200, "min custody reward")
	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxWhitelistRolePermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist-permission [role_sid] [permission_id]",
		Short: "Whitelist a permission to a role",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			permission, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid permission: %w", err)
			}

			msg := types.NewMsgWhitelistRolePermission(
				clientCtx.FromAddress,
				args[0],
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
		Use:   "blacklist-permission [role_sid] [permission_id]",
		Short: "Blacklist a permission for the governance role",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			permission, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid permission: %w", err)
			}

			msg := types.NewMsgBlacklistRolePermission(
				clientCtx.FromAddress,
				args[0],
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
		Use:   "remove-whitelisted-permission [role_sid] [permission_id]",
		Short: "Remove a whitelisted permission from a governance role",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			permission, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid permission: %w", err)
			}

			msg := types.NewMsgRemoveWhitelistRolePermission(
				clientCtx.FromAddress,
				args[0],
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
		Use:   "remove-blacklisted-permission [role_sid] [permission_id]",
		Short: "Remove a blacklisted permission from a governance role",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			permission, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid permission: %w", err)
			}

			msg := types.NewMsgRemoveBlacklistRolePermission(
				clientCtx.FromAddress,
				args[0],
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
		Use:   "create [role_sid] [role_description]",
		Short: "Create new role",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateRole(
				clientCtx.FromAddress,
				args[0],
				args[1],
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
		Use:   "assign [role_id]",
		Short: "Assign role to account",
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

func GetTxUnassignRole() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unassign role",
		Short: "Unassign a role from account",
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

			msg := types.NewMsgUnassignRole(
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
			MINIMUM_PROPOSAL_END_TIME
			PROPOSAL_ENACTMENT_TIME
			MIN_PROPOSAL_END_BLOCKS
			MIN_PROPOSAL_ENACTMENT_BLOCKS
			ENABLE_FOREIGN_FEE_PAYMENTS
			MISCHANCE_RANK_DECREASE_AMOUNT
			MAX_MISCHANCE
			MISCHANCE_CONFIDENCE
			INACTIVE_RANK_DECREASE_PERCENT
			POOR_NETWORK_MAX_BANK_SEND
			MIN_VALIDATORS
			UNJAIL_MAX_TIME
			ENABLE_TOKEN_WHITELIST
			ENABLE_TOKEN_BLACKLIST
			MIN_IDENTITY_APPROVAL_TIP
			UNIQUE_IDENTITY_KEYS
			UBI_HARDCAP
			VALIDATORS_FEE_SHARE
			INFLATION_RATE
			INFLATION_PERIOD
			UNSTAKING_PERIOD
			MAX_DELEGATORS
			MIN_DELEGATION_PUSHOUT
			SLASHING_PERIOD
			MAX_JAILED_PERCENTAGE
			MAX_SLASHING_PERCENTAGE
			MIN_CUSTODY_REWARD
			MAX_CUSTODY_BUFFER_SIZE
			MAX_CUSTODY_TX_SIZE
			ABSTENTION_RANK_DECREASE_AMOUNT
			MAX_ABSTENTION
			MIN_COLLECTIVE_BOND
			MIN_COLLECTIVE_BONDING_TIME
			MAX_COLLECTIVE_OUTPUTS
			MIN_COLLECTIVE_CLAIM_PERIOD
			VALIDATOR_RECOVERY_BOND
			MAX_ANNUAL_INFLATION
			MAX_PROPOSAL_TITLE_SIZE
			MAX_PROPOSAL_DESCRIPTION_SIZE
			MAX_PROPOSAL_POLL_OPTION_SIZE
			MAX_PROPOSAL_POLL_OPTION_COUNT
			MAX_PROPOSAL_REFERENCE_SIZE
			MAX_PROPOSAL_CHECKSUM_SIZE
			MIN_DAPP_BOND
			MAX_DAPP_BOND
			DAPP_LIQUIDATION_THRESHOLD
			DAPP_LIQUIDATION_PERIOD
			DAPP_BOND_DURATION
			DAPP_VERIFIER_BOND
			DAPP_AUTO_DENOUNCE_TIME
			DAPP_MISCHANCE_RANK_DECREASE_AMOUNT
			DAPP_MAX_MISCHANCE
			DAPP_INACTIVE_RANK_DECREASE_PERCENT
			DAPP_POOL_SLIPPAGE_DEFAULT
			MINTING_FT_FEE
			MINTING_NFT_FEE
			VETO_THRESHOLD
			AUTOCOMPOUND_INTERVAL_NUM_BLOCKS
			DOWNTIME_INACTIVE_DURATION
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
			switch types.NetworkProperty(property) {
			case types.InactiveRankDecreasePercent:
				fallthrough
			case types.UniqueIdentityKeys:
				fallthrough
			case types.ValidatorsFeeShare:
				fallthrough
			case types.InflationRate:
				fallthrough
			case types.MaxJailedPercentage:
				fallthrough
			case types.MaxSlashingPercentage:
				fallthrough
			case types.MaxAnnualInflation:
				fallthrough
			case types.DappVerifierBond:
				fallthrough
			case types.DappPoolSlippageDefault:
				fallthrough
			case types.DappInactiveRankDecreasePercent:
				fallthrough
			case types.VoteQuorum:
				fallthrough
			case types.VetoThreshold:
				value.StrValue = args[1]
			default:
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

func GetTxProposalAssignRoleToAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assign-role [role_identifier]",
		Short: "Create a proposal to assign a role to an address.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
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
				types.NewAssignRoleToAccountProposal(addr, args[0]),
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

func GetTxProposalUnassignRoleFromAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unassign-role [role]",
		Short: "Create a proposal to unassign a role from an address.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
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
				types.NewUnassignRoleFromAccountProposal(addr, args[0]),
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

func GetTxProposalWhitelistAccountPermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist-permission [permission_id]",
		Short: "Create a proposal to whitelist a permission to an address.",
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
				types.NewWhitelistAccountPermissionProposal(addr, types.PermValue(perm)),
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

func GetTxProposalBlacklistAccountPermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blacklist-permission [permission_id]",
		Short: "Create a proposal to blacklist a permission to an address.",
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
				types.NewBlacklistAccountPermissionProposal(addr, types.PermValue(perm)),
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

func GetTxProposalRemoveWhitelistedAccountPermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-whitelisted-permission [permission_id]",
		Short: "Create a proposal to remove a whitelisted permission from an address.",
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
				types.NewRemoveWhitelistedAccountPermissionProposal(addr, types.PermValue(perm)),
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

func GetTxProposalRemoveBlacklistedAccountPermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-blacklisted-permission [permission_id]",
		Short: "Create a proposal to remove a blacklisted permission from an address.",
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
				types.NewRemoveBlacklistedAccountPermissionProposal(addr, types.PermValue(perm)),
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

			slashStr, err := cmd.Flags().GetString(FlagSlash)
			if err != nil {
				return err
			}
			slash, err := sdk.NewDecFromStr(slashStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgVoteProposal(
				uint64(proposalID),
				clientCtx.FromAddress,
				types.VoteOption(voteOption),
				slash,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)
	cmd.Flags().String(FlagSlash, "0.01", "slash value on the proposal")

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
		return nil, fmt.Errorf("invalid Bech32 address")
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
			username, _ := cmd.Flags().GetString(FlagUsername)
			description, _ := cmd.Flags().GetString(FlagDescription)
			social, _ := cmd.Flags().GetString(FlagSocial)
			contact, _ := cmd.Flags().GetString(FlagContact)
			avatar, _ := cmd.Flags().GetString(FlagAvatar)

			msg := types.NewMsgClaimCouncilor(
				clientCtx.FromAddress,
				moniker,
				username,
				description,
				social,
				contact,
				avatar,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().String(FlagAddr, "", "the address")
	cmd.Flags().String(FlagMoniker, "", "the Moniker")
	cmd.Flags().String(FlagUsername, "", "the Username")
	cmd.Flags().String(FlagDescription, "", "the description")
	cmd.Flags().String(FlagSocial, "", "the social")
	cmd.Flags().String(FlagContact, "", "the contact")
	cmd.Flags().String(FlagAvatar, "", "the avatar")

	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// CouncilorPause - signal to the network that Councilor will NOT be present for a prolonged period of time
func GetTxCouncilorPauseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause",
		Short: "Pause councilor",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCouncilorPause(
				clientCtx.GetFromAddress(),
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// CouncilorUnpause - signal to the network that Councilor wishes to regain voting ability after planned absence
func GetTxCouncilorUnpauseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unpause",
		Short: "Unpause councilor",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCouncilorUnpause(
				clientCtx.GetFromAddress(),
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// CouncilorActivate - signal to the network that Councilor wishes to regain voting ability after planned absence
func GetTxCouncilorActivateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activate",
		Short: "Activate councilor",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCouncilorActivate(
				clientCtx.GetFromAddress(),
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxProposalCreateRole() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [role_sid] [role_description]",
		Short: "Raise governance proposal to create a new role.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
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
					args[0],
					args[1],
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

func GetTxProposalRemoveRole() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [role_sid]",
		Short: "Raise governance proposal to remove a role.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
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
				types.NewRemoveRoleProposal(
					args[0],
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
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxProposalWhitelistRolePermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist-permission [role_sid] [permission_id]",
		Short: "Raise governance proposal to whitelist a permission for a role.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			perm, err := strconv.Atoi(args[1])
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
				types.NewWhitelistRolePermissionProposal(
					args[0],
					types.PermValue(perm),
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
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxProposalBlacklistRolePermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blacklist-permission [role_sid] [role_description]",
		Short: "Raise governance proposal to blacklist a permission for a role.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			perm, err := strconv.Atoi(args[1])
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
				types.NewBlacklistRolePermissionProposal(
					args[0],
					types.PermValue(perm),
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
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxProposalRemoveWhitelistedRolePermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-whitelisted-permission [role_sid] [permission]",
		Short: "Raise governance proposal to remove whitelisted permission from a role.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			perm, err := strconv.Atoi(args[1])
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
				types.NewRemoveWhitelistedRolePermissionProposal(
					args[0],
					types.PermValue(perm),
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
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxProposalRemoveBlacklistedRolePermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-blacklisted-permission [role_sid] [permission_id]",
		Short: "Raise governance proposal to remove a blacklisted permission from a role.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			perm, err := strconv.Atoi(args[1])
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
				types.NewRemoveBlacklistedRolePermissionProposal(
					args[0],
					types.PermValue(perm),
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
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxProposalSetProposalDurations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-proposal-durations-proposal [proposal_types] [durations]",
		Short: "Create a proposal to set batch proposal durations.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposalTypes := strings.Split(args[0], ",")
			proposalDurationStrs := strings.Split(args[1], ",")
			proposalDurations := []uint64{}
			for _, durationStr := range proposalDurationStrs {
				duration, err := strconv.Atoi(durationStr)
				if err != nil {
					return err
				}
				proposalDurations = append(proposalDurations, uint64(duration))
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
				types.NewSetProposalDurationsProposal(
					proposalTypes,
					proposalDurations,
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
	cmd.Flags().String(FlagDescription, "", "The description of the proposal, it can be a url, some text, etc.")
	cmd.MarkFlagRequired(FlagTitle)
	cmd.MarkFlagRequired(FlagDescription)
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

			tipStr, err := cmd.Flags().GetString(FlagVerifierTip)
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

	cmd.Flags().String(FlagVerifierTip, "", "The tip to be given to the verifier.")
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

func GetTxPollCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a poll.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
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

			reference, err := cmd.Flags().GetString(FlagPollReference)
			if err != nil {
				return fmt.Errorf("invalid reference: %w", err)
			}

			checksum, err := cmd.Flags().GetString(FlagPollChecksum)
			if err != nil {
				return fmt.Errorf("invalid checksum: %w", err)
			}

			options, err := cmd.Flags().GetStringSlice(FlagPollOptions)
			if err != nil {
				return fmt.Errorf("invalid options: %w", err)
			}

			var filteredOptions []string
			for _, v := range options {
				filteredOptions = append(filteredOptions, strings.ToLower(strings.TrimSpace(v)))
			}

			roles, err := cmd.Flags().GetStringSlice(FlagPollRoles)
			if err != nil {
				return fmt.Errorf("invalid roles: %w", err)
			}

			valueCount, err := cmd.Flags().GetUint64(FlagPollCount)
			if err != nil {
				return fmt.Errorf("invalid count: %w", err)
			}

			valueType, err := cmd.Flags().GetString(FlagPollType)
			if err != nil {
				return fmt.Errorf("invalid type: %w", err)
			}

			possibleChoices, err := cmd.Flags().GetUint64(FlagPollChoices)
			if err != nil {
				return fmt.Errorf("invalid choices: %w", err)
			}

			duration, err := cmd.Flags().GetString(FlagPollDuration)
			if err != nil {
				return fmt.Errorf("invalid duration: %w", err)
			}

			msg := types.NewMsgPollCreate(
				clientCtx.FromAddress,
				title,
				description,
				reference,
				checksum,
				filteredOptions,
				roles,
				valueCount,
				valueType,
				possibleChoices,
				duration,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(FlagTitle, "", "The title of the poll.")
	cmd.MarkFlagRequired(FlagTitle)
	cmd.Flags().String(FlagDescription, "", "The description of the poll, it can be an url, some text, etc.")
	cmd.MarkFlagRequired(FlagDescription)
	cmd.Flags().String(FlagPollReference, "", "IPFS CID or URL reference to file describing poll and voting options in depth.")
	cmd.Flags().String(FlagPollChecksum, "", "Reference checksum.")
	cmd.Flags().StringSlice(FlagPollOptions, []string{}, "The options value in the format variant1,variant2.")
	cmd.MarkFlagRequired(FlagPollOptions)
	cmd.Flags().StringSlice(FlagPollRoles, []string{}, "List of roles that are allowed to take part in the poll vote in the format role1,role2.")
	cmd.MarkFlagRequired(FlagPollRoles)
	cmd.Flags().Uint64(FlagPollCount, 128, "Maximum number of voting options that poll can have.")
	cmd.Flags().String(FlagPollType, "", "Type of the options, all user supplied or predefined options must match its type.")
	cmd.MarkFlagRequired(FlagPollType)
	cmd.Flags().Uint64(FlagPollChoices, 1, "Should define maximum number of choices that voter can select.")
	cmd.Flags().String(FlagPollDuration, "", "The duration of the poll.")
	cmd.MarkFlagRequired(FlagPollDuration)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxVotePoll() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote [poll-id] [poll-option] ",
		Short: "Vote a poll.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pollID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid poll ID: %w", err)
			}

			optionID, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid option ID: %w", err)
			}

			value, err := cmd.Flags().GetString(FlagCustomPollValue)
			if err != nil {
				return fmt.Errorf("invalid custom value: %w", err)
			}

			msg := types.NewMsgVotePoll(
				uint64(pollID),
				clientCtx.FromAddress,
				types.PollVoteOption(optionID),
				value,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(FlagCustomPollValue, "", "The custom poll value.")
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

// GetTxProposalResetWholeCouncilorRankCmd implement cli command for ProposalResetWholeCouncilorRank
func GetTxProposalResetWholeCouncilorRankCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-reset-whole-councilor-rank",
		Short: "Create a proposal to reset whole councilor rank",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
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
				types.NewResetWholeCouncilorRankProposal(clientCtx.FromAddress),
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
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxProposalJailCouncilorCmd implement cli command for ProposalJailCouncilor
func GetTxProposalJailCouncilorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-jail-councilor [councilors]",
		Short: "Create a proposal to jail councilors",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
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

			councilors := strings.Split(args[0], ",")

			msg, err := types.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				types.NewJailCouncilorProposal(clientCtx.FromAddress, description, councilors),
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
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxProposalSetExecutionFeesCmd implement cli command for ProposalSetExecutionFees
func GetTxProposalSetExecutionFeesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "proposal-set-execution-fees",
		Short:   "Create a proposal to set execution fees",
		Example: `proposal-set-execution-fees --tx-types=[txTypes] --execution-fees=[executionFees] --failure-fees=[failureFees] --timeouts=[timeouts] --default-params=[defaultParams]`,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
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

			txTypesStr, err := cmd.Flags().GetString(FlagTxTypes)
			if err != nil {
				return fmt.Errorf("invalid tx types: %w", err)
			}

			execFeesStr, err := cmd.Flags().GetString(FlagExecutionFees)
			if err != nil {
				return fmt.Errorf("invalid execution fees: %w", err)
			}

			failureFeesStr, err := cmd.Flags().GetString(FlagFailureFees)
			if err != nil {
				return fmt.Errorf("invalid failure fees: %w", err)
			}

			timeoutsStr, err := cmd.Flags().GetString(FlagTimeouts)
			if err != nil {
				return fmt.Errorf("invalid timeouts: %w", err)
			}

			defaultParamsStr, err := cmd.Flags().GetString(FlagTimeouts)
			if err != nil {
				return fmt.Errorf("invalid default params: %w", err)
			}

			txTypes := strings.Split(txTypesStr, ",")
			execFeeStrs := strings.Split(execFeesStr, ",")
			failureFeeStrs := strings.Split(failureFeesStr, ",")
			timeoutStrs := strings.Split(timeoutsStr, ",")
			defaultParamStrs := strings.Split(defaultParamsStr, ",")
			executionFees := []types.ExecutionFee{}
			for i, txType := range txTypes {
				execFee, err := strconv.Atoi(execFeeStrs[i])
				if err != nil {
					return err
				}
				failureFee, err := strconv.Atoi(failureFeeStrs[i])
				if err != nil {
					return err
				}
				timeout, err := strconv.Atoi(timeoutStrs[i])
				if err != nil {
					return err
				}
				defaultParams, err := strconv.Atoi(defaultParamStrs[i])
				if err != nil {
					return err
				}
				executionFees = append(executionFees, types.ExecutionFee{
					TransactionType:   txType,
					ExecutionFee:      uint64(execFee),
					FailureFee:        uint64(failureFee),
					Timeout:           uint64(timeout),
					DefaultParameters: uint64(defaultParams),
				})
			}

			msg, err := types.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				types.NewSetExecutionFeesProposal(clientCtx.FromAddress, description, executionFees),
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
	cmd.Flags().String(FlagTxTypes, "", "Transaction types to set execution fees")
	cmd.Flags().String(FlagExecutionFees, "", "Execution fees")
	cmd.Flags().String(FlagFailureFees, "", "Failure fees")
	cmd.Flags().String(FlagTimeouts, "", "Timeouts")
	cmd.Flags().String(FlagDefaultParams, "", "Default params")

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// convertAsPermValues convert array of int32 to PermValue array.
func convertAsPermValues(values []int32) []types.PermValue {
	var v []types.PermValue
	for _, perm := range values {
		v = append(v, types.PermValue(perm))
	}

	return v
}

// convertAsPermValues convert array of int32 to PermValue array.
func convertAsOptionValues(values []int32) []types.PollVoteOption {
	var v []types.PollVoteOption
	for _, option := range values {
		v = append(v, types.PollVoteOption(option))
	}

	return v
}
