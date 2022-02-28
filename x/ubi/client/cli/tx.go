package cli

import (
	"fmt"

	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/ubi/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

// flags for ubi module txs
const (
	FlagTitle             = "title"
	FlagDescription       = "description"
	FlagName              = "name"
	FlagDistributionStart = "distr-start"
	FlagDistributionEnd   = "distr-end"
	FlagAmount            = "amount"
	FlagPeriod            = "period"
	FlagPoolName          = "pool-name"
)

// NewTxCmd returns a root CLI command handler for all x/bank transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "UBI sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetTxProposalUpsertUBICmd(),
		GetTxProposalRemoveUBICmd(),
	)

	return txCmd
}

func GetTxProposalUpsertUBICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-upsert-ubi",
		Short: "Create a proposal to upsert ubi",
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

			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return fmt.Errorf("invalid ubi name: %w", err)
			}

			distrStart, err := cmd.Flags().GetUint64(FlagDistributionStart)
			if err != nil {
				return fmt.Errorf("invalid ubi distribution start: %w", err)
			}

			distrEnd, err := cmd.Flags().GetUint64(FlagDistributionEnd)
			if err != nil {
				return fmt.Errorf("invalid ubi distribution end: %w", err)
			}

			amount, err := cmd.Flags().GetUint64(FlagAmount)
			if err != nil {
				return fmt.Errorf("invalid ubi amount: %w", err)
			}

			period, err := cmd.Flags().GetUint64(FlagPeriod)
			if err != nil {
				return fmt.Errorf("invalid ubi period: %w", err)
			}

			poolName, err := cmd.Flags().GetString(FlagPoolName)
			if err != nil {
				return fmt.Errorf("invalid ubi pool name: %w", err)
			}

			msg, err := govtypes.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				types.NewUpsertUBIProposal(name, distrStart, distrEnd, amount, period, poolName),
			)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagName, "", "The name of ubi.")
	cmd.Flags().Uint64(FlagDistributionStart, 0, "The distribution start time of ubi.")
	cmd.Flags().Uint64(FlagDistributionEnd, 0, "The distribution end time of ubi.")
	cmd.Flags().Uint64(FlagAmount, 0, "The amount of tokens to be minted per period.")
	cmd.Flags().Uint64(FlagPeriod, 0, "The duration to to mint tokens.")
	cmd.Flags().String(FlagPoolName, "", "The target pool name to receive minted tokens.")
	cmd.Flags().String(FlagTitle, "", "The title of a proposal.")
	cmd.MarkFlagRequired(FlagTitle)
	cmd.Flags().String(FlagDescription, "", "The description of the proposal, it can be a url, some text, etc.")
	cmd.MarkFlagRequired(FlagDescription)

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GetTxProposalRemoveUBICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-remove-ubi",
		Short: "Create a proposal to remove ubi",
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

			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return fmt.Errorf("invalid ubi name: %w", err)
			}

			msg, err := govtypes.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				&types.RemoveUBIProposal{
					UbiName: name,
				},
			)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagName, "", "The name of ubi.")
	cmd.Flags().String(FlagTitle, "", "The title of a proposal.")
	cmd.MarkFlagRequired(FlagTitle)
	cmd.Flags().String(FlagDescription, "", "The description of the proposal, it can be a url, some text, etc.")
	cmd.MarkFlagRequired(FlagDescription)

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
