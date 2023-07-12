package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/KiraCore/sekai/x/basket/types"
	govcli "github.com/KiraCore/sekai/x/gov/client/cli"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

// flags for basket module txs
const (
	FlagBasketId          = "basket-id"
	FlagBasketSuffix      = "basket-suffix"
	FlagBasketDescription = "basket-description"
	FlagSwapFee           = "swap-fee"
	FlagSlippageFeeMin    = "slippage-fee-min"
	FlagTokensCap         = "tokens-cap"
	FlagLimitsPeriod      = "limits-period"
	FlagMintsMin          = "mints-min"
	FlagMintsMax          = "mints-max"
	FlagMintsDisabled     = "mints-disabled"
	FlagBurnsMin          = "burns-min"
	FlagBurnsMax          = "burns-max"
	FlagBurnsDisabled     = "burns-disabled"
	FlagSwapsMin          = "swaps-min"
	FlagSwapsMax          = "swaps-max"
	FlagSwapsDisabled     = "swaps-disabled"
	FlagBasketTokens      = "basket-tokens"
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
		GetTxDisableBasketDepositsCmd(),
		GetTxDisableBasketWithdrawsCmd(),
		GetTxDisableBasketSwapsCmd(),
		GetTxBasketTokenMintCmd(),
		GetTxBasketTokenBurnCmd(),
		GetTxBasketTokenSwapCmd(),
		GetTxBasketClaimRewardsCmd(),
		GetTxProposalCreateBasketCmd(),
		GetTxProposalEditBasketCmd(),
		GetTxProposalBasketWithdrawSurplusCmd(),
	)

	return txCmd
}

// GetTxDisableBasketDepositsCmd implement cli command for MsgDisableBasketDeposits
func GetTxDisableBasketDepositsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable-basket-deposits [basket_id] [disabled]",
		Short: "Emergency function & permission to disable one or all deposits of one or all token in the basket",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			basketId, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			disabled, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDisableBasketDeposits(
				clientCtx.FromAddress,
				uint64(basketId),
				disabled,
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
		Use:   "disable-basket-withdraws [basket_id] [disabled]",
		Short: "Emergency function & permission to disable one or all withdraws of one or all token in the basket",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			basketId, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			disabled, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDisableBasketWithdraws(
				clientCtx.FromAddress,
				uint64(basketId),
				disabled,
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
		Use:   "disable-basket-swaps [basket_id] [disabled]",
		Short: "Emergency function & permission to disable one or all swaps of one or all token in the basket",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			basketId, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			disabled, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDisableBasketSwaps(
				clientCtx.FromAddress,
				uint64(basketId),
				disabled,
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
		Use:   "mint-basket-tokens [basket_id] [deposit_coins]",
		Short: "mint basket tokens",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			basketId, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgBasketTokenMint(
				clientCtx.FromAddress,
				uint64(basketId),
				coins,
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
		Use:   "burn-basket-tokens [basket_id] [burn_coin]",
		Short: "burn basket tokens",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			basketId, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			coin, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgBasketTokenBurn(
				clientCtx.FromAddress,
				uint64(basketId),
				coin,
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
		Use:   "swap-basket-tokens [basket_id] [in_amount1] [out_token1] [in_amount2] [out_token2] ...",
		Short: "swap one or many of the basket tokens for one or many others",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			basketId, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			pairs := []types.SwapPair{}
			for i := 1; i < len(args); i += 2 {
				coin, err := sdk.ParseCoinNormalized(args[i])
				if err != nil {
					return err
				}
				if i+1 == len(args) {
					return fmt.Errorf("out token not set for %s", args[i])
				}
				pairs = append(pairs, types.SwapPair{
					InAmount: coin,
					OutToken: args[i+1],
				})
			}

			msg := types.NewMsgBasketTokenSwap(
				clientCtx.FromAddress,
				uint64(basketId),
				pairs,
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
		Use:   "basket-claim-rewards [basket_tokens]",
		Short: "force staking derivative `SDB` basket to claim outstanding rewards of one all or many aggregate `V<ID>` tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			basketTokens, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgBasketClaimRewards(
				clientCtx.FromAddress,
				basketTokens,
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

// GetTxProposalCreateBasketCmd implement cli command for ProposalCreateBasket
func GetTxProposalCreateBasketCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-create-basket",
		Short: "Create a proposal to create a basket",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			title, err := cmd.Flags().GetString(govcli.FlagTitle)
			if err != nil {
				return fmt.Errorf("invalid title: %w", err)
			}
			description, err := cmd.Flags().GetString(govcli.FlagDescription)
			if err != nil {
				return fmt.Errorf("invalid description: %w", err)
			}

			suffix, err := cmd.Flags().GetString(FlagBasketSuffix)
			if err != nil {
				return fmt.Errorf("invalid basket suffix: %w", err)
			}

			basketDescription, err := cmd.Flags().GetString(FlagBasketDescription)
			if err != nil {
				return fmt.Errorf("invalid basket description: %w", err)
			}

			swapFeeStr, err := cmd.Flags().GetString(FlagSwapFee)
			if err != nil {
				return fmt.Errorf("invalid basket swap fee: %w", err)
			}

			swapFee, err := sdk.NewDecFromStr(swapFeeStr)
			if err != nil {
				return fmt.Errorf("invalid basket swap fee: %w", err)
			}

			slippageFeeMinStr, err := cmd.Flags().GetString(FlagSlippageFeeMin)
			if err != nil {
				return fmt.Errorf("invalid basket slippage fee min: %w", err)
			}

			slipppageFeeMin, err := sdk.NewDecFromStr(slippageFeeMinStr)
			if err != nil {
				return fmt.Errorf("invalid basket slippage fee min: %w", err)
			}

			tokensCapStr, err := cmd.Flags().GetString(FlagTokensCap)
			if err != nil {
				return fmt.Errorf("invalid basket tokens cap: %w", err)
			}

			tokensCap, err := sdk.NewDecFromStr(tokensCapStr)
			if err != nil {
				return fmt.Errorf("invalid basket tokens cap: %w", err)
			}

			limitsPeriod, err := cmd.Flags().GetUint64(FlagLimitsPeriod)
			if err != nil {
				return fmt.Errorf("invalid basket limits cap: %w", err)
			}

			mintsMinStr, err := cmd.Flags().GetString(FlagMintsMin)
			if err != nil {
				return fmt.Errorf("invalid basket minimum mints: %w", err)
			}

			mintsMin, ok := sdk.NewIntFromString(mintsMinStr)
			if !ok {
				return fmt.Errorf("invalid basket minimum mints: %s", mintsMinStr)
			}

			mintsMaxStr, err := cmd.Flags().GetString(FlagMintsMax)
			if err != nil {
				return fmt.Errorf("invalid basket maximum mints: %w", err)
			}

			mintsMax, ok := sdk.NewIntFromString(mintsMaxStr)
			if !ok {
				return fmt.Errorf("invalid basket maximum mints: %s", mintsMaxStr)
			}

			mintsDisabled, err := cmd.Flags().GetBool(FlagMintsDisabled)
			if err != nil {
				return fmt.Errorf("invalid basket mints disabled flag: %w", err)
			}

			burnsMinStr, err := cmd.Flags().GetString(FlagBurnsMin)
			if err != nil {
				return fmt.Errorf("invalid basket minimum burns: %w", err)
			}

			burnsMin, ok := sdk.NewIntFromString(burnsMinStr)
			if !ok {
				return fmt.Errorf("invalid basket minimum burns: %s", burnsMinStr)
			}

			burnsMaxStr, err := cmd.Flags().GetString(FlagBurnsMax)
			if err != nil {
				return fmt.Errorf("invalid basket maximum burns: %w", err)
			}

			burnsMax, ok := sdk.NewIntFromString(burnsMaxStr)
			if !ok {
				return fmt.Errorf("invalid basket maximum burns: %s", burnsMaxStr)
			}

			burnsDisabled, err := cmd.Flags().GetBool(FlagBurnsDisabled)
			if err != nil {
				return fmt.Errorf("invalid basket burns disabled flag: %w", err)
			}

			swapsMinStr, err := cmd.Flags().GetString(FlagSwapsMin)
			if err != nil {
				return fmt.Errorf("invalid basket minimum swaps: %w", err)
			}

			swapsMin, ok := sdk.NewIntFromString(swapsMinStr)
			if !ok {
				return fmt.Errorf("invalid basket minimum mints: %s", swapsMinStr)
			}

			swapsMaxStr, err := cmd.Flags().GetString(FlagSwapsMax)
			if err != nil {
				return fmt.Errorf("invalid basket maximum swaps: %w", err)
			}

			swapsMax, ok := sdk.NewIntFromString(swapsMaxStr)
			if !ok {
				return fmt.Errorf("invalid basket maximum swaps: %s", swapsMaxStr)
			}

			swapsDisabled, err := cmd.Flags().GetBool(FlagSwapsDisabled)
			if err != nil {
				return fmt.Errorf("invalid basket swaps disabled flag: %w", err)
			}

			basketTokensStr, err := cmd.Flags().GetString(FlagBasketTokens)
			if err != nil {
				return fmt.Errorf("invalid basket tokens: %w", err)
			}

			basketTokens, err := parseBasketTokens(basketTokensStr)
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				types.NewProposalCreateBasket(types.Basket{
					Suffix:          suffix,
					Description:     basketDescription,
					SwapFee:         swapFee,
					SlipppageFeeMin: slipppageFeeMin,
					TokensCap:       tokensCap,
					LimitsPeriod:    limitsPeriod,
					MintsMin:        mintsMin,
					MintsMax:        mintsMax,
					MintsDisabled:   mintsDisabled,
					BurnsMin:        burnsMin,
					BurnsMax:        burnsMax,
					BurnsDisabled:   burnsDisabled,
					SwapsMin:        swapsMin,
					SwapsMax:        swapsMax,
					SwapsDisabled:   swapsDisabled,
					Tokens:          basketTokens,
				}),
			)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(govcli.FlagTitle, "", "The title of the proposal.")
	cmd.MarkFlagRequired(govcli.FlagTitle)
	cmd.Flags().String(govcli.FlagDescription, "", "The description of the proposal, it can be a url, some text, etc.")
	cmd.MarkFlagRequired(govcli.FlagDescription)
	cmd.Flags().String(FlagBasketSuffix, "", "The suffix of the basket.")
	cmd.Flags().String(FlagBasketDescription, "", "The description of the basket.")
	cmd.Flags().String(FlagSwapFee, "", "Swap fee on the basket.")
	cmd.Flags().String(FlagSlippageFeeMin, "", "Minimum slippage fee on the basket.")
	cmd.Flags().String(FlagTokensCap, "", "Tokens cap on the basket.")
	cmd.Flags().Uint64(FlagLimitsPeriod, 0, "Limits period on the basket.")
	cmd.Flags().String(FlagMintsMin, "", "Min mint amount on the basket.")
	cmd.Flags().String(FlagMintsMax, "", "Max mint amount on a day on the basket.")
	cmd.Flags().Bool(FlagMintsDisabled, false, "Mints enabled flag on the basket.")
	cmd.Flags().String(FlagBurnsMin, "", "Min burn amount on the basket.")
	cmd.Flags().String(FlagBurnsMax, "", "Max burn amount on a day on the basket.")
	cmd.Flags().Bool(FlagBurnsDisabled, false, "Burns enabled flag on the basket.")
	cmd.Flags().String(FlagSwapsMin, "", "Min swap amount on the basket.")
	cmd.Flags().String(FlagSwapsMax, "", "Max swap amount on a day on the basket.")
	cmd.Flags().Bool(FlagSwapsDisabled, false, "Swap disabled flag on the basket.")
	cmd.Flags().String(FlagBasketTokens, "", "Underlying tokens with rates on the basket.")

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxProposalEditBasketCmd implement cli command for ProposalEditBasket
func GetTxProposalEditBasketCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-edit-basket",
		Short: "Create a proposal to edit a basket",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			title, err := cmd.Flags().GetString(govcli.FlagTitle)
			if err != nil {
				return fmt.Errorf("invalid title: %w", err)
			}
			description, err := cmd.Flags().GetString(govcli.FlagDescription)
			if err != nil {
				return fmt.Errorf("invalid description: %w", err)
			}

			basketId, err := cmd.Flags().GetUint64(FlagBasketId)
			if err != nil {
				return fmt.Errorf("invalid basket id: %w", err)
			}

			suffix, err := cmd.Flags().GetString(FlagBasketSuffix)
			if err != nil {
				return fmt.Errorf("invalid basket suffix: %w", err)
			}

			basketDescription, err := cmd.Flags().GetString(FlagBasketDescription)
			if err != nil {
				return fmt.Errorf("invalid basket description: %w", err)
			}

			swapFeeStr, err := cmd.Flags().GetString(FlagSwapFee)
			if err != nil {
				return fmt.Errorf("invalid basket swap fee: %w", err)
			}

			swapFee, err := sdk.NewDecFromStr(swapFeeStr)
			if err != nil {
				return fmt.Errorf("invalid basket swap fee: %w", err)
			}

			slippageFeeMinStr, err := cmd.Flags().GetString(FlagSlippageFeeMin)
			if err != nil {
				return fmt.Errorf("invalid basket slippage fee min: %w", err)
			}

			slipppageFeeMin, err := sdk.NewDecFromStr(slippageFeeMinStr)
			if err != nil {
				return fmt.Errorf("invalid basket slippage fee min: %w", err)
			}

			tokensCapStr, err := cmd.Flags().GetString(FlagTokensCap)
			if err != nil {
				return fmt.Errorf("invalid basket tokens cap: %w", err)
			}

			tokensCap, err := sdk.NewDecFromStr(tokensCapStr)
			if err != nil {
				return fmt.Errorf("invalid basket tokens cap: %w", err)
			}

			limitsPeriod, err := cmd.Flags().GetUint64(FlagLimitsPeriod)
			if err != nil {
				return fmt.Errorf("invalid basket limits cap: %w", err)
			}

			mintsMinStr, err := cmd.Flags().GetString(FlagMintsMin)
			if err != nil {
				return fmt.Errorf("invalid basket minimum mints: %w", err)
			}

			mintsMin, ok := sdk.NewIntFromString(mintsMinStr)
			if !ok {
				return fmt.Errorf("invalid basket minimum mints: %s", mintsMinStr)
			}

			mintsMaxStr, err := cmd.Flags().GetString(FlagMintsMax)
			if err != nil {
				return fmt.Errorf("invalid basket maximum mints: %w", err)
			}

			mintsMax, ok := sdk.NewIntFromString(mintsMaxStr)
			if !ok {
				return fmt.Errorf("invalid basket maximum mints: %s", mintsMaxStr)
			}

			mintsDisabled, err := cmd.Flags().GetBool(FlagMintsDisabled)
			if err != nil {
				return fmt.Errorf("invalid basket mints disabled flag: %w", err)
			}

			burnsMinStr, err := cmd.Flags().GetString(FlagBurnsMin)
			if err != nil {
				return fmt.Errorf("invalid basket minimum burns: %w", err)
			}

			burnsMin, ok := sdk.NewIntFromString(burnsMinStr)
			if !ok {
				return fmt.Errorf("invalid basket minimum burns: %s", burnsMinStr)
			}

			burnsMaxStr, err := cmd.Flags().GetString(FlagBurnsMax)
			if err != nil {
				return fmt.Errorf("invalid basket maximum burns: %w", err)
			}

			burnsMax, ok := sdk.NewIntFromString(burnsMaxStr)
			if !ok {
				return fmt.Errorf("invalid basket maximum burns: %s", burnsMaxStr)
			}

			burnsDisabled, err := cmd.Flags().GetBool(FlagBurnsDisabled)
			if err != nil {
				return fmt.Errorf("invalid basket burns disabled flag: %w", err)
			}

			swapsMinStr, err := cmd.Flags().GetString(FlagSwapsMin)
			if err != nil {
				return fmt.Errorf("invalid basket minimum swaps: %w", err)
			}

			swapsMin, ok := sdk.NewIntFromString(swapsMinStr)
			if !ok {
				return fmt.Errorf("invalid basket minimum mints: %s", swapsMinStr)
			}

			swapsMaxStr, err := cmd.Flags().GetString(FlagSwapsMax)
			if err != nil {
				return fmt.Errorf("invalid basket maximum swaps: %w", err)
			}

			swapsMax, ok := sdk.NewIntFromString(swapsMaxStr)
			if !ok {
				return fmt.Errorf("invalid basket maximum swaps: %s", swapsMaxStr)
			}

			swapsDisabled, err := cmd.Flags().GetBool(FlagSwapsDisabled)
			if err != nil {
				return fmt.Errorf("invalid basket swaps disabled flag: %w", err)
			}

			basketTokensStr, err := cmd.Flags().GetString(FlagBasketTokens)
			if err != nil {
				return fmt.Errorf("invalid basket tokens: %w", err)
			}

			basketTokens, err := parseBasketTokens(basketTokensStr)
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				types.NewProposalEditBasket(types.Basket{
					Id:              basketId,
					Suffix:          suffix,
					Description:     basketDescription,
					SwapFee:         swapFee,
					SlipppageFeeMin: slipppageFeeMin,
					TokensCap:       tokensCap,
					LimitsPeriod:    limitsPeriod,
					MintsMin:        mintsMin,
					MintsMax:        mintsMax,
					MintsDisabled:   mintsDisabled,
					BurnsMin:        burnsMin,
					BurnsMax:        burnsMax,
					BurnsDisabled:   burnsDisabled,
					SwapsMin:        swapsMin,
					SwapsMax:        swapsMax,
					SwapsDisabled:   swapsDisabled,
					Tokens:          basketTokens,
				}),
			)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(govcli.FlagTitle, "", "The title of the proposal.")
	cmd.MarkFlagRequired(govcli.FlagTitle)
	cmd.Flags().String(govcli.FlagDescription, "", "The description of the proposal, it can be a url, some text, etc.")
	cmd.MarkFlagRequired(govcli.FlagDescription)
	cmd.Flags().Uint64(FlagBasketId, 0, "The id of the basket.")
	cmd.Flags().String(FlagBasketSuffix, "", "The suffix of the basket.")
	cmd.Flags().String(FlagBasketDescription, "", "The description of the basket.")
	cmd.Flags().String(FlagSwapFee, "", "Swap fee on the basket.")
	cmd.Flags().String(FlagSlippageFeeMin, "", "Minimum slippage fee on the basket.")
	cmd.Flags().String(FlagTokensCap, "", "Tokens cap on the basket.")
	cmd.Flags().Uint64(FlagLimitsPeriod, 0, "Limits period on the basket.")
	cmd.Flags().String(FlagMintsMin, "", "Min mint amount on the basket.")
	cmd.Flags().String(FlagMintsMax, "", "Max mint amount on a day on the basket.")
	cmd.Flags().Bool(FlagMintsDisabled, false, "Mints enabled flag on the basket.")
	cmd.Flags().String(FlagBurnsMin, "", "Min burn amount on the basket.")
	cmd.Flags().String(FlagBurnsMax, "", "Max burn amount on a day on the basket.")
	cmd.Flags().Bool(FlagBurnsDisabled, false, "Burns enabled flag on the basket.")
	cmd.Flags().String(FlagSwapsMin, "", "Min swap amount on the basket.")
	cmd.Flags().String(FlagSwapsMax, "", "Max swap amount on a day on the basket.")
	cmd.Flags().Bool(FlagSwapsDisabled, false, "Swap disabled flag on the basket.")
	cmd.Flags().String(FlagBasketTokens, "", "Underlying tokens with rates on the basket.")

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxProposalBasketWithdrawSurplusCmd implement cli command for ProposalBasketWithdrawSurplus
func GetTxProposalBasketWithdrawSurplusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-basket-withdraw-surplus [basket_ids] [withdraw_target]",
		Short: "Create a proposal to withdraw surplus from the basket",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			title, err := cmd.Flags().GetString(govcli.FlagTitle)
			if err != nil {
				return fmt.Errorf("invalid title: %w", err)
			}
			description, err := cmd.Flags().GetString(govcli.FlagDescription)
			if err != nil {
				return fmt.Errorf("invalid description: %w", err)
			}

			basketIdsArr := strings.Split(args[0], ",")
			basketIds := []uint64{}
			for _, basketIdStr := range basketIdsArr {
				basketId, err := strconv.Atoi(basketIdStr)
				if err != nil {
					return err
				}
				basketIds = append(basketIds, uint64(basketId))
			}

			msg, err := govtypes.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				types.NewProposalBasketWithdrawSurplus(basketIds, args[1]),
			)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(govcli.FlagTitle, "", "The title of the proposal.")
	cmd.MarkFlagRequired(govcli.FlagTitle)
	cmd.Flags().String(govcli.FlagDescription, "", "The description of the proposal, it can be a url, some text, etc.")
	cmd.MarkFlagRequired(govcli.FlagDescription)

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func parseBasketTokens(basketTokensStr string) ([]types.BasketToken, error) {
	basketTokens := []types.BasketToken{}
	basketTokensArr := strings.Split(basketTokensStr, ",")
	for _, basketTokenStr := range basketTokensArr {
		split := strings.Split(basketTokenStr, "#")
		if len(split) != 5 {
			return basketTokens, fmt.Errorf("invalid basket token info: %s", basketTokenStr)
		}
		weight, err := sdk.NewDecFromStr(split[1])
		if err != nil {
			return basketTokens, fmt.Errorf("invalid basket token weight: %w", err)
		}
		deposits, err := strconv.ParseBool(split[2])
		if err != nil {
			return basketTokens, fmt.Errorf("invalid basket token deposits: %w", err)
		}
		withdraws, _ := strconv.ParseBool(split[3])
		if err != nil {
			return basketTokens, fmt.Errorf("invalid basket token withdraws: %w", err)
		}
		swaps, _ := strconv.ParseBool(split[4])
		if err != nil {
			return basketTokens, fmt.Errorf("invalid basket token swaps: %w", err)
		}
		basketTokens = append(basketTokens, types.BasketToken{
			Denom:     split[0],
			Weight:    weight,
			Deposits:  deposits,
			Withdraws: withdraws,
			Swaps:     swaps,
		})
	}
	return basketTokens, nil
}
