package cli

import (
	"github.com/KiraCore/sekai/x/basket/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

// flags for basket module txs
const ()

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
		GetTxDisableBasketDepositsCmd(),
	)

	return txCmd
}

// GetTxDisableBasketDepositsCmd implement cli command for MsgDisableBasketDeposits
func GetTxDisableBasketDepositsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable-basket-deposits",
		Short: "Emergency function & permission to disable one or all deposits of one or all token in the basket",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			msg := types.NewMsgDisableBasketDeposits(
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

// GetTxDisableBasketWithdrawsCmd implement cli command for MsgDisableBasketWithdraws
func GetTxDisableBasketWithdrawsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable-basket-withdraws",
		Short: "Emergency function & permission to disable one or all withdraws of one or all token in the basket",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			msg := types.NewMsgDisableBasketWithdraws(
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

// GetTxDisableBasketSwapsCmd implement cli command for MsgDisableBasketSwaps
func GetTxDisableBasketSwapsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable-basket-swaps",
		Short: "Emergency function & permission to disable one or all swaps of one or all token in the basket",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			msg := types.NewMsgDisableBasketSwaps(
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

// GetTxBasketTokenMintCmd implement cli command for MsgBasketTokenMint
func GetTxBasketTokenMintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-basket-tokens",
		Short: "mint basket tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			msg := types.NewMsgBasketTokenMint(
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

// GetTxBasketTokenBurnCmd implement cli command for MsgBasketTokenBurn
func GetTxBasketTokenBurnCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-basket-tokens",
		Short: "burn basket tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			msg := types.NewMsgBasketTokenBurn(
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

// GetTxBasketTokenSwapCmd implement cli command for MsgBasketTokenSwap
func GetTxBasketTokenSwapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap-basket-tokens",
		Short: "swap one or many of the basket tokens for one or many others",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			msg := types.NewMsgBasketTokenSwap(
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

// GetTxBasketClaimRewardsCmd implement cli command for MsgBasketClaimRewards
func GetTxBasketClaimRewardsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "basket-claim-rewards",
		Short: "force staking derivative `SDB` basket to claim outstanding rewards of one all or many aggregate `V<ID>` tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			msg := types.NewMsgBasketClaimRewards(
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
