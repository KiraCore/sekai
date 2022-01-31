package cli

import (
	"fmt"

	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/spending/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

// flags for tokens module txs
const (
	FlagSymbol      = "symbol"
	FlagName        = "name"
	FlagIcon        = "icon"
	FlagDecimals    = "decimals"
	FlagDenoms      = "denoms"
	FlagDenom       = "denom"
	FlagRate        = "rate"
	FlagFeePayments = "fee_payments"
	FlagIsBlacklist = "is_blacklist"
	FlagIsAdd       = "is_add"
	FlagTokens      = "tokens"
	FlagTitle       = "title"
	FlagDescription = "description"
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
		GetTxCreateSpendingPoolCmd(),
		GetTxDepositSpendingPoolCmd(),
		GetTxRegisterSpendingPoolBeneficiaryCmd(),
		GetTxClaimSpendingPoolCmd(),
		GetTxUpdateSpendingPoolProposalCmd(),
		GetTxSpendingPoolDistributionProposalCmd(),
		GetTxSpendingPoolWithdrawProposalCmd(),
	)
	return txCmd
}

// GetTxCreateSpendingPoolCmd implement cli command for MsgCreateSpendingPool
func GetTxCreateSpendingPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-spending-pool",
		Short: "Create spending pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			// NewMsgCreateSpendingPool
			msg := &types.MsgCreateSpendingPool{}

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

// GetTxDepositSpendingPoolCmd implement cli command for MsgDepositSpendingPool
func GetTxDepositSpendingPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-spending-pool",
		Short: "Deposit spending pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			// NewMsgDepositSpendingPool
			msg := &types.MsgDepositSpendingPool{}

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

// GetTxRegisterSpendingPoolBeneficiaryCmd implement cli command for MsgRegisterSpendingPoolBeneficiary
func GetTxRegisterSpendingPoolBeneficiaryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-spending-pool-beneficiary",
		Short: "Register spending pool beneficiary",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			// NewMsgRegisterSpendingPoolBeneficiary
			msg := &types.MsgRegisterSpendingPoolBeneficiary{}

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

// GetTxClaimSpendingPoolCmd implement cli command for MsgClaimSpendingPool
func GetTxClaimSpendingPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-spending-pool",
		Short: "Claim spending pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			// NewMsgClaimSpendingPool
			msg := &types.MsgClaimSpendingPool{}

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

// GetTxUpdateSpendingPoolProposalCmd implement cli command for UpdateSpendingPoolProposal
func GetTxUpdateSpendingPoolProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-update-spending-pool",
		Short: "Create a proposal to update spending pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			title, err := cmd.Flags().GetString(FlagTitle)
			if err != nil {
				return fmt.Errorf("invalid title: %w", err)
			}

			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return fmt.Errorf("invalid description: %w", err)
			}

			msg, err := govtypes.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				&types.UpdateSpendingPoolProposal{},
				// types.NewUpdateSpendingPoolProposal(
				// ),
			)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagTitle, "", "The title of a proposal.")
	cmd.MarkFlagRequired(FlagTitle)
	cmd.Flags().String(FlagDescription, "", "The description of the proposal, it can be a url, some text, etc.")
	cmd.MarkFlagRequired(FlagDescription)

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxSpendingPoolDistributionProposalCmd implement cli command for SpendingPoolDistributionProposal
func GetTxSpendingPoolDistributionProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-spending-pool-distribution",
		Short: "Create a proposal to distribute the spending pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			title, err := cmd.Flags().GetString(FlagTitle)
			if err != nil {
				return fmt.Errorf("invalid title: %w", err)
			}

			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return fmt.Errorf("invalid description: %w", err)
			}

			msg, err := govtypes.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				&types.SpendingPoolDistributionProposal{},
				// types.NewSpendingPoolDistributionProposal(
				// ),
			)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagTitle, "", "The title of a proposal.")
	cmd.MarkFlagRequired(FlagTitle)
	cmd.Flags().String(FlagDescription, "", "The description of the proposal, it can be a url, some text, etc.")
	cmd.MarkFlagRequired(FlagDescription)

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxSpendingPoolWithdrawProposalCmd implement cli command for SpendingPoolWithdrawProposal
func GetTxSpendingPoolWithdrawProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-spending-pool-withdraw",
		Short: "Create a proposal to withdraw spending pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			title, err := cmd.Flags().GetString(FlagTitle)
			if err != nil {
				return fmt.Errorf("invalid title: %w", err)
			}

			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return fmt.Errorf("invalid description: %w", err)
			}

			msg, err := govtypes.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				&types.SpendingPoolWithdrawProposal{},
				// types.NewSpendingPoolWithdrawProposal(
				// ),
			)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagTitle, "", "The title of a proposal.")
	cmd.MarkFlagRequired(FlagTitle)
	cmd.Flags().String(FlagDescription, "", "The description of the proposal, it can be a url, some text, etc.")
	cmd.MarkFlagRequired(FlagDescription)

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
