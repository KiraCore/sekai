package cli

import (
	"fmt"
	"strings"

	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/tokens/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		Short: "Create a proposal to upsert token alias",
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
				types.NewUpsertTokenAliasProposal(
					symbol,
					name,
					icon,
					decimals,
					strings.Split(denoms, ","),
				),
			)
			if err != nil {
				return err
			}

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

	cmd.Flags().String(FlagTitle, "", "The title of the proposal.")
	cmd.MarkFlagRequired(FlagTitle)
	cmd.Flags().String(FlagDescription, "", "The description of the proposal, it can be a url, some text, etc.")
	cmd.MarkFlagRequired(FlagDescription)

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxProposalUpsertTokenRatesCmd implement cli command for MsgUpsertTokenAlias
func GetTxProposalUpsertTokenRatesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-upsert-rate",
		Short: "Create a proposal to upsert token rate",
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
				types.NewUpsertTokenRatesProposal(
					denom,
					rate,
					feePayments,
				),
			)
			if err != nil {
				return err
			}

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
	cmd.Flags().String(FlagTitle, "", "The title of a proposal.")
	cmd.MarkFlagRequired(FlagTitle)
	cmd.Flags().String(FlagDescription, "", "The description of the proposal, it can be a url, some text, etc.")
	cmd.MarkFlagRequired(FlagDescription)

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
		Use:   "proposal-update-tokens-blackwhite",
		Short: "Create a proposal to update whitelisted and blacklisted tokens",
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
				types.NewTokensWhiteBlackChangeProposal(
					isBlacklist,
					isAdd,
					tokens,
				),
			)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Bool(FlagIsBlacklist, true, "true to modify blacklist otherwise false")
	cmd.Flags().Bool(FlagIsAdd, true, "true to add otherwise false")
	cmd.Flags().StringArray(FlagTokens, []string{}, "tokens array (eg. ATOM, KEX, BTC)")
	cmd.Flags().String(FlagTitle, "", "The title of a proposal.")
	cmd.MarkFlagRequired(FlagTitle)
	cmd.Flags().String(FlagDescription, "", "The description of the proposal, it can be a url, some text, etc.")
	cmd.MarkFlagRequired(FlagDescription)

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
