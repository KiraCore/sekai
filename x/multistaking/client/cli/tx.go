package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/KiraCore/sekai/x/multistaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// define flags
const (
	FlagEnabled    = "enabled"
	FlagCommission = "commission"
)

// NewTxCmd returns a root CLI command handler for all x/multistaking transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "multistaking sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetTxUpsertStakingPool(),
		GetTxDelegate(),
		GetTxUndelegate(),
		GetTxClaimRewards(),
		GetTxClaimUndelegation(),
		GetTxClaimMaturedUndelegations(),
		GetTxSetCompoundInfo(),
		GetTxRegisterDelegator(),
	)

	return txCmd
}

func GetTxUpsertStakingPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upsert-staking-pool [validator-key]",
		Short: "Submit a transaction to upsert staking pool.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			enabled, err := cmd.Flags().GetBool(FlagEnabled)
			if err != nil {
				return fmt.Errorf("invalid flag: %w", err)
			}

			commissionStr, err := cmd.Flags().GetString(FlagCommission)
			if err != nil {
				return fmt.Errorf("failed to parse commission: %s", err.Error())
			}
			commission, err := sdk.NewDecFromStr(commissionStr)
			if err != nil {
				return err
			}
			msg := types.NewMsgUpsertStakingPool(clientCtx.FromAddress.String(), args[0], enabled, commission)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().Bool(FlagEnabled, true, "enabled flag for the pool")
	cmd.Flags().String(FlagCommission, "", "the commission rate ranged from 0.01 to 0.5.(1%% to 50%%)")
	cmd.MarkFlagRequired(FlagCommission)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxDelegate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegate [validator] [coins]",
		Short: "Submit a transaction to delegate to a pool.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDelegate(clientCtx.FromAddress.String(), args[0], coins)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxUndelegate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "undelegate [validator] [coins]",
		Short: "Submit a transaction to start undelegation from a pool.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgUndelegate(clientCtx.FromAddress.String(), args[0], coins)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxClaimRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-rewards",
		Short: "Submit a transaction to claim rewards from a pool.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgClaimRewards(clientCtx.FromAddress.String())

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxClaimUndelegation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-undelegation [undelegationId]",
		Short: "Submit a transaction to claim matured undelegation.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			undelegationId, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			msg := types.NewMsgClaimUndelegation(clientCtx.FromAddress.String(), uint64(undelegationId))

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxClaimMaturedUndelegations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-matured-undelegations",
		Short: "Submit a transaction to claim all matured undelegations.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgClaimMaturedUndelegations(clientCtx.FromAddress.String())
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxSetCompoundInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-compound-info [all_denom] [specific_denoms]",
		Short: "Submit a transaction to set compound info.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			allDenom, err := strconv.ParseBool(args[0])
			if err != nil {
				return err
			}

			denoms := strings.Split(args[1], ",")
			msg := types.NewMsgSetCompoundInfo(clientCtx.FromAddress.String(), allDenom, denoms)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxRegisterDelegator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-delegator",
		Short: "Submit a transaction to register a delegator.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRegisterDelegator(clientCtx.FromAddress.String())

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
