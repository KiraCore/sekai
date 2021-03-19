package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		GetTxUpsertTokenAliasCmd(),
		GetTxUpsertTokenRateCmd(),
		GetTxProposalUpsertTokenAliasCmd(),
		GetTxProposalUpsertTokenRatesCmd(),
		GetTxProposalTokensBlackWhiteChangeCmd(),
	)

	return txCmd
}

// GetTxUpsertTokenAliasCmd implement cli command for MsgUpsertTokenAlias
func GetTxUpsertTokenAliasCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upsert-alias",
		Short: "Upsert token alias",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			symbol, err := cmd.Flags().GetString(FlagSymbol)
			if err != nil {
				return fmt.Errorf("invalid symbol")
			}

			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return fmt.Errorf("invalid name")
			}

			icon, err := cmd.Flags().GetString(FlagIcon)
			if err != nil {
				return fmt.Errorf("invalid icon")
			}

			decimals, err := cmd.Flags().GetUint32(FlagDecimals)
			if err != nil {
				return fmt.Errorf("invalid decimals")
			}

			denomsString, err := cmd.Flags().GetString(FlagDenoms)
			if err != nil {
				return fmt.Errorf("invalid denoms")
			}

			denoms := strings.Split(denomsString, ",")
			for _, denom := range denoms {
				if err = sdk.ValidateDenom(denom); err != nil {
					return err
				}
			}

			msg := types.NewMsgUpsertTokenAlias(
				clientCtx.FromAddress,
				symbol,
				name,
				icon,
				decimals,
				denoms,
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagSymbol, "KEX", "Ticker (eg. ATOM, KEX, BTC)")
	cmd.Flags().String(FlagName, "Kira", "Token Name (e.g. Cosmos, Kira, Bitcoin)")
	cmd.Flags().String(FlagIcon, "", "Graphical Symbol (url link to graphics)")
	cmd.Flags().Uint32(FlagDecimals, 6, "Integer number of max decimals")
	cmd.Flags().String(FlagDenoms, "ukex,mkex", "An array of token denoms to be aliased")

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxProposalUpsertTokenAliasCmd implement cli command for MsgUpsertTokenAlias
func GetTxProposalUpsertTokenAliasCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-upsert-alias",
		Short: "Creates an Upsert token alias",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			symbol, err := cmd.Flags().GetString(FlagSymbol)
			if err != nil {
				return fmt.Errorf("invalid symbol: %w", err)
			}

			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return fmt.Errorf("invalid name: %w", err)
			}

			icon, err := cmd.Flags().GetString(FlagIcon)
			if err != nil {
				return fmt.Errorf("invalid icon: %w", err)
			}

			decimals, err := cmd.Flags().GetUint32(FlagDecimals)
			if err != nil {
				return fmt.Errorf("invalid decimals: %w", err)
			}

			denoms, err := cmd.Flags().GetString(FlagDenoms)
			if err != nil {
				return fmt.Errorf("invalid denoms: %w", err)
			}

			msg := types.NewMsgProposalUpsertTokenAlias(
				clientCtx.FromAddress,
				symbol,
				name,
				icon,
				decimals,
				strings.Split(denoms, ","),
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagSymbol, "KEX", "Ticker (eg. ATOM, KEX, BTC)")
	cmd.MarkFlagRequired(FlagSymbol)
	cmd.Flags().String(FlagName, "Kira", "Token Name (e.g. Cosmos, Kira, Bitcoin)")
	cmd.MarkFlagRequired(FlagName)
	cmd.Flags().String(FlagIcon, "", "Graphical Symbol (url link to graphics)")
	cmd.MarkFlagRequired(FlagIcon)
	cmd.Flags().Uint32(FlagDecimals, 6, "Integer number of max decimals")
	cmd.MarkFlagRequired(FlagDecimals)
	cmd.Flags().String(FlagDenoms, "ukex,mkex", "An array of token denoms to be aliased")
	cmd.MarkFlagRequired(FlagDenoms)

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxProposalUpsertTokenRatesCmd implement cli command for MsgUpsertTokenAlias
func GetTxProposalUpsertTokenRatesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-upsert-rate",
		Short: "Creates an Upsert token rate",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			denom, err := cmd.Flags().GetString(FlagDenom)
			if err != nil {
				return fmt.Errorf("invalid denom")
			}
			if denom == "ukex" {
				return fmt.Errorf("bond denom rate is read-only")
			}

			rateString, err := cmd.Flags().GetString(FlagRate)
			if err != nil {
				return fmt.Errorf("invalid rate")
			}

			rate, err := sdk.NewDecFromStr(rateString)
			if err != nil {
				return err
			}

			feePayments, err := cmd.Flags().GetBool(FlagFeePayments)
			if err != nil {
				return fmt.Errorf("invalid fee payments")
			}

			msg := types.NewMsgProposalUpsertTokenRates(
				clientCtx.FromAddress,
				denom,
				rate,
				feePayments,
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagDenom, "tbtc", "denom - identifier for token rates")
	cmd.MarkFlagRequired(FlagDenom)
	cmd.Flags().String(FlagRate, "1.0", "rate to register, max decimal 9, max value 10^10")
	cmd.MarkFlagRequired(FlagRate)
	cmd.Flags().Bool(FlagFeePayments, true, "use registry as fee payment")
	cmd.MarkFlagRequired(FlagFeePayments)

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxUpsertTokenRateCmd implement cli command for MsgUpsertTokenRate
func GetTxUpsertTokenRateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upsert-rate",
		Short: "Upsert token rate",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			denom, err := cmd.Flags().GetString(FlagDenom)
			if err != nil {
				return fmt.Errorf("invalid denom")
			}
			if denom == "ukex" {
				return fmt.Errorf("bond denom rate is read-only")
			}

			rateString, err := cmd.Flags().GetString(FlagRate)
			if err != nil {
				return fmt.Errorf("invalid rate")
			}

			rate, err := sdk.NewDecFromStr(rateString)
			if err != nil {
				return err
			}

			feePayments, err := cmd.Flags().GetBool(FlagFeePayments)
			if err != nil {
				return fmt.Errorf("invalid fee payments")
			}

			msg := types.NewMsgUpsertTokenRate(
				clientCtx.FromAddress,
				denom,
				rate,
				feePayments,
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagDenom, "tbtc", "denom - identifier for token rates")
	cmd.MarkFlagRequired(FlagDenom)
	cmd.Flags().String(FlagRate, "1.0", "rate to register, max decimal 9, max value 10^10")
	cmd.MarkFlagRequired(FlagRate)
	cmd.Flags().Bool(FlagFeePayments, true, "use registry as fee payment")
	cmd.MarkFlagRequired(FlagFeePayments)

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxProposalTokensBlackWhiteChangeCmd implement cli command for proposing tokens blacklist / whitelist update
func GetTxProposalTokensBlackWhiteChangeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "propose-update-tokens-blackwhite",
		Short: "Propose update whitelisted and blacklisted tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			isBlacklist, err := cmd.Flags().GetBool(FlagIsBlacklist)
			if err != nil {
				return fmt.Errorf("invalid is_blacklist flag: %w", err)
			}

			isAdd, err := cmd.Flags().GetBool(FlagIsAdd)
			if err != nil {
				return fmt.Errorf("invalid is_add flag: %w", err)
			}

			tokens, err := cmd.Flags().GetStringArray(FlagTokens)
			if err != nil {
				return fmt.Errorf("invalid tokens flag: %w", err)
			}

			msg := types.NewMsgProposalTokensWhiteBlackChange(
				clientCtx.FromAddress,
				isBlacklist,
				isAdd,
				tokens,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Bool(FlagIsBlacklist, true, "true to modify blacklist otherwise false")
	cmd.Flags().Bool(FlagIsAdd, true, "true to add otherwise false")
	cmd.Flags().StringArray(FlagTokens, []string{}, "tokens array (eg. ATOM, KEX, BTC)")

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
