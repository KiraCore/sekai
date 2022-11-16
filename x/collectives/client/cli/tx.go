package cli

import (
	"github.com/KiraCore/sekai/x/collectives/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

// NewTxCmd returns a root CLI command handler for all x/bank transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Collectives sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetTxCreateCollectiveCmd(),
		GetTxContributeCollectiveCmd(),
		GetTxDonateCollectiveCmd(),
		GetTxWithdrawCollectiveCmd(),
	)

	return txCmd
}

// GetTxCreateCollectiveCmd defines a method for creating collective.
func GetTxCreateCollectiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-collective",
		Short: "a method to create collective",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateCollective(
				clientCtx.FromAddress,
			)

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

// GetTxContributeCollectiveCmd defines a method for putting bonds on collective.
func GetTxContributeCollectiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contribute-collective",
		Short: "a method to put bonds on collective",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgBondCollective(
				clientCtx.FromAddress,
			)

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

// GetTxDonateCollectiveCmd defines a method for putting bonds on collective.
func GetTxDonateCollectiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "donate-collective",
		Short: "a method to set lock and donation for bonds on the collection",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDonateCollective(
				clientCtx.FromAddress,
			)

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

// GetTxWithdrawCollectiveCmd can be sent by any whitelisted “contributor” to withdraw
// their tokens (unless locking is enabled)
func GetTxWithdrawCollectiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "donate-collective",
		Short: "sent by any whitelisted “contributor” to withdraw their tokens (unless locking is enabled)",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawCollective(
				clientCtx.FromAddress,
			)

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
