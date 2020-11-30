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
	)

	return txCmd
}

// GetTxUpsertTokenAliasCmd implement cli command for MsgUpsertTokenAlias
func GetTxUpsertTokenAliasCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upsert-alias",
		Short: "Upsert token alias",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

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

// GetTxUpsertTokenAliasCmd implement cli command for MsgUpsertTokenAlias
func GetTxProposalUpsertTokenAliasCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-upsert-alias",
		Short: "Creates an Upsert token alias",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

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

			denoms := args[4]

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
	cmd.Flags().String(FlagName, "Kira", "Token Name (e.g. Cosmos, Kira, Bitcoin)")
	cmd.Flags().String(FlagIcon, "", "Graphical Symbol (url link to graphics)")
	cmd.Flags().Uint32(FlagDecimals, 6, "Integer number of max decimals")
	cmd.Flags().String(FlagDenoms, "ukex,mkex", "An array of token denoms to be aliased")

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxUpsertTokenAliasCmd implement cli command for MsgUpsertTokenAlias
func GetTxProposalUpsertTokenRatesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-upsert-rate",
		Short: "Creates an Upsert token rate",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

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
	cmd.Flags().String(FlagRate, "1.0", "rate to register, max decimal 9, max value 10^10")
	cmd.Flags().Bool(FlagFeePayments, true, "use registry as fee payment")

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
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

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
	cmd.Flags().String(FlagRate, "1.0", "rate to register, max decimal 9, max value 10^10")
	cmd.Flags().Bool(FlagFeePayments, true, "use registry as fee payment")

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
