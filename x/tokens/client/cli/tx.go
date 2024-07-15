package cli

import (
	"fmt"

	appparams "github.com/KiraCore/sekai/app/params"
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
	FlagSymbol            = "symbol"
	FlagName              = "name"
	FlagIcon              = "icon"
	FlagDecimals          = "decimals"
	FlagDenoms            = "denoms"
	FlagDenom             = "denom"
	FlagFeeRate           = "fee_rate"
	FlagStakeCap          = "stake_cap"
	FlagStakeToken        = "stake_token"
	FlagStakeMin          = "stake_min"
	FlagFeeEnabled        = "fee_payments"
	FlagIsBlacklist       = "is_blacklist"
	FlagIsAdd             = "is_add"
	FlagTokens            = "tokens"
	FlagTitle             = "title"
	FlagDescription       = "description"
	FlagInvalidated       = "invalidated"
	FlagSupply            = "supply"
	FlagSupplyCap         = "supply_cap"
	FlagWebsite           = "website"
	FlagSocial            = "social"
	FlagMintingFee        = "minting_fee"
	FlagOwner             = "owner"
	FlagOwnerEditDisabled = "owner_edit_disabled"
	FlagNftMetadata       = "nft_metadata"
	FlagNftHash           = "nft_hash"
	FlagTokenType         = "token_type"
	FlagTokenRate         = "token_rate"
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
		GetTxUpsertTokenInfoCmd(),
		GetTxProposalUpsertTokenInfoCmd(),
		GetTxProposalTokensBlackWhiteChangeCmd(),
	)

	return txCmd
}

// GetTxProposalUpsertTokenInfoCmd implement cli command for MsgUpsertTokenInfos
func GetTxProposalUpsertTokenInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-upsert-rate",
		Short: "Create a proposal to upsert token rate",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return fmt.Errorf("invalid is_blacklist flag: %w", err)
			}

			denom, err := cmd.Flags().GetString(FlagDenom)
			if err != nil {
				return fmt.Errorf("invalid denom")
			}
			if denom == appparams.DefaultDenom {
				return fmt.Errorf("bond denom rate is read-only")
			}

			tokenType, err := cmd.Flags().GetString(FlagTokenType)
			if err != nil {
				return fmt.Errorf("invalid tokenType")
			}

			rateString, err := cmd.Flags().GetString(FlagTokenRate)
			if err != nil {
				return fmt.Errorf("invalid rate")
			}

			feeRate, err := sdk.NewDecFromStr(rateString)
			if err != nil {
				return err
			}

			feeEnabled, err := cmd.Flags().GetBool(FlagFeeEnabled)
			if err != nil {
				return fmt.Errorf("invalid fee enabled flag")
			}

			supplyString, err := cmd.Flags().GetString(FlagSupply)
			if err != nil {
				return fmt.Errorf("invalid supply: %w", err)
			}

			supply, ok := sdk.NewIntFromString(supplyString)
			if !ok {
				return fmt.Errorf("invalid supply: %s", supplyString)
			}

			supplyCapString, err := cmd.Flags().GetString(FlagSupplyCap)
			if err != nil {
				return fmt.Errorf("invalid supply cap: %w", err)
			}

			supplyCap, ok := sdk.NewIntFromString(supplyCapString)
			if !ok {
				return fmt.Errorf("invalid supply cap: %s", supplyCapString)
			}

			title, err := cmd.Flags().GetString(FlagTitle)
			if err != nil {
				return fmt.Errorf("invalid title: %w", err)
			}

			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return fmt.Errorf("invalid description: %w", err)
			}

			website, err := cmd.Flags().GetString(FlagWebsite)
			if err != nil {
				return fmt.Errorf("invalid website: %w", err)
			}

			social, err := cmd.Flags().GetString(FlagSocial)
			if err != nil {
				return fmt.Errorf("invalid social: %w", err)
			}

			mintingFeeStr, err := cmd.Flags().GetString(FlagMintingFee)
			if err != nil {
				return fmt.Errorf("invalid mintingFee: %w", err)
			}

			mintingFee, ok := sdk.NewIntFromString(mintingFeeStr)
			if !ok {
				return fmt.Errorf("invalid minting fee: %s", mintingFeeStr)
			}

			owner, err := cmd.Flags().GetString(FlagOwner)
			if err != nil {
				return fmt.Errorf("invalid owner: %w", err)
			}

			ownerEditDisabled, err := cmd.Flags().GetBool(FlagOwnerEditDisabled)
			if err != nil {
				return fmt.Errorf("invalid ownerEditDisabled: %w", err)
			}

			nftMetadata, err := cmd.Flags().GetString(FlagNftMetadata)
			if err != nil {
				return fmt.Errorf("invalid nftMetadata: %w", err)
			}

			nftHash, err := cmd.Flags().GetString(FlagNftHash)
			if err != nil {
				return fmt.Errorf("invalid nftHash: %w", err)
			}

			stakeToken, err := cmd.Flags().GetBool(FlagStakeToken)
			if err != nil {
				return fmt.Errorf("invalid stake token flag")
			}

			stakeCapStr, err := cmd.Flags().GetString(FlagStakeCap)
			if err != nil {
				return fmt.Errorf("invalid stake cap: %w", err)
			}

			stakeCap, err := sdk.NewDecFromStr(stakeCapStr)
			if err != nil {
				return fmt.Errorf("invalid stake cap: %w", err)
			}

			stakeMinStr, err := cmd.Flags().GetString(FlagStakeMin)
			if err != nil {
				return fmt.Errorf("invalid stake min: %w", err)
			}

			stakeMin, ok := sdk.NewIntFromString(stakeMinStr)
			if !ok {
				return fmt.Errorf("invalid stake min: %s", stakeMinStr)
			}

			isInvalidated, err := cmd.Flags().GetBool(FlagInvalidated)
			if err != nil {
				return fmt.Errorf("invalid invalidated flag: %w", err)
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

			msg, err := govtypes.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				types.NewUpsertTokenInfosProposal(
					denom,
					tokenType,
					feeRate,
					feeEnabled,
					supply,
					supplyCap,
					stakeCap,
					stakeMin,
					stakeToken,
					isInvalidated,
					symbol,
					name,
					icon,
					decimals,
					description,
					website,
					social,
					0,
					mintingFee,
					owner,
					ownerEditDisabled,
					nftMetadata,
					nftHash,
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
	cmd.Flags().String(FlagFeeRate, "1.0", "rate to register, max decimal 9, max value 10^10")
	cmd.MarkFlagRequired(FlagFeeRate)
	cmd.Flags().Bool(FlagFeeEnabled, true, "use registry as fee payment")
	cmd.MarkFlagRequired(FlagFeeEnabled)
	cmd.Flags().String(FlagTitle, "", "The title of a proposal.")
	cmd.MarkFlagRequired(FlagTitle)
	cmd.Flags().String(FlagDescription, "", "The description of the proposal, it can be a url, some text, etc.")
	cmd.MarkFlagRequired(FlagDescription)
	cmd.Flags().String(FlagStakeCap, "0.1", "rewards to be allocated for the token.")
	cmd.Flags().String(FlagStakeMin, "1", "min amount to stake at a time.")
	cmd.Flags().Bool(FlagStakeToken, false, "flag of if staking token or not.")
	cmd.Flags().Bool(FlagInvalidated, false, "Flag to show token rate is invalidated or not")
	cmd.Flags().String(FlagSymbol, "KEX", "Ticker (eg. ATOM, KEX, BTC)")
	cmd.Flags().String(FlagName, "Kira", "Token Name (e.g. Cosmos, Kira, Bitcoin)")
	cmd.Flags().String(FlagIcon, "", "Graphical Symbol (url link to graphics)")
	cmd.Flags().Uint32(FlagDecimals, 6, "Integer number of max decimals")
	cmd.Flags().String(FlagSupply, "", "Supply of token")
	cmd.Flags().String(FlagSupplyCap, "", "Supply cap of token")
	cmd.Flags().String(FlagWebsite, "", "Website")
	cmd.Flags().String(FlagSocial, "", "Social")
	cmd.Flags().String(FlagMintingFee, "", "Minting fee")
	cmd.Flags().String(FlagOwner, "", "Owner")
	cmd.Flags().Bool(FlagOwnerEditDisabled, false, "Owner edit disabled flag")
	cmd.Flags().String(FlagNftMetadata, "", "Nft metadata")
	cmd.Flags().String(FlagNftHash, "", "Nft hash")
	cmd.Flags().String(FlagTokenType, "", "Token type")
	cmd.Flags().String(FlagTokenRate, "", "Token rate")

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxUpsertTokenInfoCmd implement cli command for MsgUpsertTokenInfo
func GetTxUpsertTokenInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upsert-rate",
		Short: "Upsert token rate",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			denom, err := cmd.Flags().GetString(FlagDenom)
			if err != nil {
				return fmt.Errorf("invalid denom")
			}
			if denom == appparams.DefaultDenom {
				return fmt.Errorf("bond denom rate is read-only")
			}

			tokenType, err := cmd.Flags().GetString(FlagTokenType)
			if err != nil {
				return fmt.Errorf("invalid tokenType")
			}

			rateString, err := cmd.Flags().GetString(FlagFeeRate)
			if err != nil {
				return fmt.Errorf("invalid rate")
			}

			feeRate, err := sdk.NewDecFromStr(rateString)
			if err != nil {
				return err
			}

			feeEnabled, err := cmd.Flags().GetBool(FlagFeeEnabled)
			if err != nil {
				return fmt.Errorf("invalid fee payments")
			}

			supplyString, err := cmd.Flags().GetString(FlagSupply)
			if err != nil {
				return fmt.Errorf("invalid supply: %w", err)
			}

			supply, ok := sdk.NewIntFromString(supplyString)
			if !ok {
				return fmt.Errorf("invalid supply: %s", supplyString)
			}

			supplyCapString, err := cmd.Flags().GetString(FlagSupplyCap)
			if err != nil {
				return fmt.Errorf("invalid supply cap: %w", err)
			}

			supplyCap, ok := sdk.NewIntFromString(supplyCapString)
			if !ok {
				return fmt.Errorf("invalid supply cap: %s", supplyCapString)
			}

			stakeToken, err := cmd.Flags().GetBool(FlagStakeToken)
			if err != nil {
				return fmt.Errorf("invalid stake token flag")
			}

			stakeCapStr, err := cmd.Flags().GetString(FlagStakeCap)
			if err != nil {
				return fmt.Errorf("invalid stake cap: %w", err)
			}

			stakeCap, err := sdk.NewDecFromStr(stakeCapStr)
			if err != nil {
				return fmt.Errorf("invalid stake cap: %w", err)
			}

			stakeMinStr, err := cmd.Flags().GetString(FlagStakeMin)
			if err != nil {
				return fmt.Errorf("invalid stake min: %w", err)
			}

			stakeMin, ok := sdk.NewIntFromString(stakeMinStr)
			if !ok {
				return fmt.Errorf("invalid stake min: %s", stakeMinStr)
			}

			isInvalidated, err := cmd.Flags().GetBool(FlagInvalidated)
			if err != nil {
				return fmt.Errorf("invalid invalidated flag: %w", err)
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

			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return fmt.Errorf("invalid description: %w", err)
			}

			website, err := cmd.Flags().GetString(FlagWebsite)
			if err != nil {
				return fmt.Errorf("invalid website: %w", err)
			}

			social, err := cmd.Flags().GetString(FlagSocial)
			if err != nil {
				return fmt.Errorf("invalid social: %w", err)
			}

			mintingFeeStr, err := cmd.Flags().GetString(FlagMintingFee)
			if err != nil {
				return fmt.Errorf("invalid mintingFee: %w", err)
			}

			mintingFee, ok := sdk.NewIntFromString(mintingFeeStr)
			if !ok {
				return fmt.Errorf("invalid minting fee: %s", mintingFeeStr)
			}

			owner, err := cmd.Flags().GetString(FlagOwner)
			if err != nil {
				return fmt.Errorf("invalid owner: %w", err)
			}

			ownerEditDisabled, err := cmd.Flags().GetBool(FlagOwnerEditDisabled)
			if err != nil {
				return fmt.Errorf("invalid ownerEditDisabled: %w", err)
			}

			nftMetadata, err := cmd.Flags().GetString(FlagNftMetadata)
			if err != nil {
				return fmt.Errorf("invalid nftMetadata: %w", err)
			}

			nftHash, err := cmd.Flags().GetString(FlagNftHash)
			if err != nil {
				return fmt.Errorf("invalid nftHash: %w", err)
			}

			msg := types.NewMsgUpsertTokenInfo(
				clientCtx.FromAddress,
				denom,
				tokenType,
				feeRate,
				feeEnabled,
				supply,
				supplyCap,
				stakeCap,
				stakeMin,
				stakeToken,
				isInvalidated,
				symbol,
				name,
				icon,
				decimals,
				description,
				website,
				social,
				0,
				mintingFee,
				owner,
				ownerEditDisabled,
				nftMetadata,
				nftHash,
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
	cmd.Flags().String(FlagFeeRate, "1.0", "rate to register, max decimal 9, max value 10^10")
	cmd.MarkFlagRequired(FlagFeeRate)
	cmd.Flags().Bool(FlagFeeEnabled, true, "use registry as fee payment")
	cmd.MarkFlagRequired(FlagFeeEnabled)
	cmd.Flags().String(FlagStakeCap, "0.1", "rewards to be allocated for the token.")
	cmd.Flags().String(FlagStakeMin, "1", "min amount to stake at a time.")
	cmd.Flags().Bool(FlagStakeToken, false, "flag of if staking token or not.")
	cmd.Flags().Bool(FlagInvalidated, false, "Flag to show token rate is invalidated or not")
	cmd.Flags().String(FlagSymbol, "KEX", "Ticker (eg. ATOM, KEX, BTC)")
	cmd.Flags().String(FlagName, "Kira", "Token Name (e.g. Cosmos, Kira, Bitcoin)")
	cmd.Flags().String(FlagIcon, "", "Graphical Symbol (url link to graphics)")
	cmd.Flags().Uint32(FlagDecimals, 6, "Integer number of max decimals")
	cmd.Flags().String(FlagSupply, "", "Supply of token")
	cmd.Flags().String(FlagSupplyCap, "", "Supply cap of token")
	cmd.Flags().String(FlagWebsite, "", "Website")
	cmd.Flags().String(FlagSocial, "", "Social")
	cmd.Flags().String(FlagMintingFee, "", "Minting fee")
	cmd.Flags().String(FlagOwner, "", "Owner")
	cmd.Flags().Bool(FlagOwnerEditDisabled, false, "Owner edit disabled flag")
	cmd.Flags().String(FlagNftMetadata, "", "Nft metadata")
	cmd.Flags().String(FlagNftHash, "", "Nft hash")
	cmd.Flags().String(FlagTokenType, "", "Token type")
	cmd.Flags().String(FlagDescription, "", "Token description")
	cmd.Flags().String(FlagTokenRate, "", "Token rate")

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
			if err != nil {
				return err
			}

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
