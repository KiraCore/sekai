package cli

import (
	"fmt"

	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

// NewTxCmd returns a root CLI command handler for all x/slashing transaction commands.
func NewTxCmd() *cobra.Command {
	slashingTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Slashing transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	slashingTxCmd.AddCommand(
		NewActivateTxCmd(),
		NewPauseTxCmd(),
		NewUnpauseTxCmd(),
		GetTxProposalResetWholeValidatorRankCmd(),
	)
	return slashingTxCmd
}

// NewActivateTxCmd defines MsgActivate tx
func NewActivateTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activate",
		Args:  cobra.NoArgs,
		Short: "Activate a validator previously inactivated for downtime",
		Long: `activate an inactivated validator:

$ <appd> tx slashing activate --from mykey
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			valAddr := clientCtx.GetFromAddress()

			msg := types.NewMsgActivate(sdk.ValAddress(valAddr))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewPauseTxCmd defines MsgPause tx
func NewPauseTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause",
		Args:  cobra.NoArgs,
		Short: "Pause a validator",
		Long: `Pause a validator before stopping of a node to avoid automatic inactivation:

$ <appd> tx customslashing pause --from validator --chain-id=testing --keyring-backend=test --fees=100ukex --home=$HOME/.sekaid --yes
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			valAddr := clientCtx.GetFromAddress()

			msg := types.NewMsgPause(sdk.ValAddress(valAddr))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewUnpauseTxCmd defines MsgUnpause tx
func NewUnpauseTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unpause",
		Args:  cobra.NoArgs,
		Short: "Unpause a validator previously paused for downtime",
		Long: `Unpause a paused validator:

$ <appd> tx slashing unpause --from mykey
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			valAddr := clientCtx.GetFromAddress()

			msg := types.NewMsgUnpause(sdk.ValAddress(valAddr))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetTxProposalResetWholeValidatorRankCmd implement cli command for ProposalResetWholeValidatorRank
func GetTxProposalResetWholeValidatorRankCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-reset-whole-validator-rank",
		Short: "Create a proposal to unjail validator (the from address is the validator)",
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

			msg, err := govtypes.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				types.NewResetWholeValidatorRankProposal(clientCtx.FromAddress),
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
