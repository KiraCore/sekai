package cli

import (
	"fmt"
	"strconv"
	"strings"

	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/spending/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
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

			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return fmt.Errorf("invalid name: %w", err)
			}

			claimStart, err := cmd.Flags().GetInt32(FlagClaimStart)
			if err != nil {
				return fmt.Errorf("invalid claim start: %w", err)
			}

			claimEnd, err := cmd.Flags().GetInt32(FlagClaimEnd)
			if err != nil {
				return fmt.Errorf("invalid claim end: %w", err)
			}

			ratesStr, err := cmd.Flags().GetString(FlagRates)
			if err != nil {
				return fmt.Errorf("invalid rate: %w", err)
			}
			rates, err := sdk.ParseDecCoins(ratesStr)
			if err != nil {
				return fmt.Errorf("invalid rates: %w", err)
			}

			voteQuorum, err := cmd.Flags().GetInt32(FlagVoteQuorum)
			if err != nil {
				return fmt.Errorf("invalid vote quorum: %w", err)
			}

			votePeriod, err := cmd.Flags().GetInt32(FlagVotePeriod)
			if err != nil {
				return fmt.Errorf("invalid vote period: %w", err)
			}

			voteEnactment, err := cmd.Flags().GetInt32(FlagVoteEnactment)
			if err != nil {
				return fmt.Errorf("invalid vote enactment: %w", err)
			}

			ownerRolesStr, err := cmd.Flags().GetString(FlagOwnerRoles)
			if err != nil {
				return fmt.Errorf("invalid owner roles: %w", err)
			}
			ownerRoles := []uint64{}
			if len(ownerRolesStr) > 0 {
				ownerRoleStrArr := strings.Split(ownerRolesStr, ",")
				for _, roleStr := range ownerRoleStrArr {
					role, err := strconv.Atoi(roleStr)
					if err != nil {
						return err
					}
					ownerRoles = append(ownerRoles, uint64(role))
				}
			}

			ownerAccountsStr, err := cmd.Flags().GetString(FlagOwnerAccounts)
			if err != nil {
				return fmt.Errorf("invalid owner accounts: %w", err)
			}

			ownerAccounts := []string{}
			if len(ownerAccountsStr) > 0 {
				ownerAccounts = strings.Split(ownerAccountsStr, ",")
			}

			beneficiary, err := parseBeneficiaryRolesAndAccounts(cmd)
			if err != nil {
				return err
			}

			dynamicRate, err := cmd.Flags().GetBool(FlagDynamicRate)
			if err != nil {
				return fmt.Errorf("invalid dynamic rate: %w", err)
			}

			dynamicRatePeriod, err := cmd.Flags().GetUint64(FlagDynamicRatePeriod)
			if err != nil {
				return fmt.Errorf("invalid dynamic rate: %w", err)
			}

			msg := types.NewMsgCreateSpendingPool(
				name, uint64(claimStart), uint64(claimEnd), rates,
				sdk.NewDecWithPrec(int64(voteQuorum), 2), uint64(votePeriod), uint64(voteEnactment),
				types.PermInfo{
					OwnerRoles:    ownerRoles,
					OwnerAccounts: ownerAccounts,
				},
				beneficiary,
				clientCtx.GetFromAddress(),
				dynamicRate,
				dynamicRatePeriod,
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

	cmd.Flags().String(FlagName, "", "The name of the spending pool.")
	cmd.Flags().Int32(FlagClaimStart, 0, "The claim start timestamp of the spending pool.")
	cmd.Flags().Int32(FlagClaimEnd, 0, "The claim end timestamp of the spending pool.")
	cmd.Flags().Uint64(FlagClaimExpiry, 43200, "claim expiry time when the users' rewards cut.")
	cmd.Flags().String(FlagRates, "", "reward rates of the spending pool.")
	cmd.Flags().Int32(FlagVoteQuorum, 0, "vote quorum of the spending pool.")
	cmd.Flags().Int32(FlagVotePeriod, 0, "vote period of the spending pool.")
	cmd.Flags().Int32(FlagVoteEnactment, 0, "vote enactment period of the spending pool.")
	cmd.Flags().String(FlagOwnerRoles, "", "owner roles of the spending pool.")
	cmd.Flags().String(FlagOwnerAccounts, "", "owner accounts of the spending pool.")
	cmd.Flags().String(FlagBeneficiaryRoles, "", "beneficiary roles of the spending pool.")
	cmd.Flags().String(FlagBeneficiaryRoleWeights, "", "beneficiary role weights on the spending pool.")
	cmd.Flags().String(FlagBeneficiaryAccounts, "", "beneficiary accounts of the spending pool.")
	cmd.Flags().String(FlagBeneficiaryAccountWeights, "", "beneficiary account weights on the spending pool.")
	cmd.Flags().Bool(FlagDynamicRate, false, "flag to dynamically calculate rates on the spending pool.")
	cmd.Flags().Uint64(FlagDynamicRatePeriod, 0, "dynamic rate recalculation period on the spending pool.")

	return cmd
}

// GetTxDepositSpendingPoolCmd implement cli command for MsgDepositSpendingPool
func GetTxDepositSpendingPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-spending-pool",
		Short: "Deposit spending pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return fmt.Errorf("invalid name: %w", err)
			}

			amountStr, err := cmd.Flags().GetString(FlagAmount)
			if err != nil {
				return fmt.Errorf("invalid amount: %w", err)
			}

			amount, err := sdk.ParseCoinsNormalized(amountStr)
			if err != nil {
				return fmt.Errorf("invalid amount: %w", err)
			}

			// NewMsgDepositSpendingPool
			msg := types.NewMsgDepositSpendingPool(name, amount, clientCtx.GetFromAddress())

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	cmd.Flags().String(FlagName, "", "The name of the spending pool.")
	cmd.Flags().String(FlagAmount, "", "The amount of coins to deposit on the spending pool.")

	return cmd
}

// GetTxRegisterSpendingPoolBeneficiaryCmd implement cli command for MsgRegisterSpendingPoolBeneficiary
func GetTxRegisterSpendingPoolBeneficiaryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-spending-pool-beneficiary",
		Short: "Register spending pool beneficiary",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return fmt.Errorf("invalid name: %w", err)
			}

			msg := types.NewMsgRegisterSpendingPoolBeneficiary(name, clientCtx.GetFromAddress())

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	cmd.Flags().String(FlagName, "", "The name of the spending pool.")

	return cmd
}

// GetTxClaimSpendingPoolCmd implement cli command for MsgClaimSpendingPool
func GetTxClaimSpendingPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-spending-pool",
		Short: "Claim spending pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return fmt.Errorf("invalid name: %w", err)
			}

			msg := types.NewMsgClaimSpendingPool(name, clientCtx.GetFromAddress())

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	cmd.Flags().String(FlagName, "", "The name of the spending pool.")

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

			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return fmt.Errorf("invalid name: %w", err)
			}

			claimStart, err := cmd.Flags().GetInt32(FlagClaimStart)
			if err != nil {
				return fmt.Errorf("invalid claim start: %w", err)
			}

			claimEnd, err := cmd.Flags().GetInt32(FlagClaimEnd)
			if err != nil {
				return fmt.Errorf("invalid claim end: %w", err)
			}

			ratesStr, err := cmd.Flags().GetString(FlagRates)
			if err != nil {
				return fmt.Errorf("invalid rates: %w", err)
			}
			rates, err := sdk.ParseDecCoins(ratesStr)
			if err != nil {
				return fmt.Errorf("invalid rates: %w", err)
			}

			voteQuorum, err := cmd.Flags().GetInt32(FlagVoteQuorum)
			if err != nil {
				return fmt.Errorf("invalid vote quorum: %w", err)
			}

			votePeriod, err := cmd.Flags().GetInt32(FlagVotePeriod)
			if err != nil {
				return fmt.Errorf("invalid vote period: %w", err)
			}

			voteEnactment, err := cmd.Flags().GetInt32(FlagVoteEnactment)
			if err != nil {
				return fmt.Errorf("invalid vote enactment: %w", err)
			}

			ownerRolesStr, err := cmd.Flags().GetString(FlagOwnerRoles)
			if err != nil {
				return fmt.Errorf("invalid owner roles: %w", err)
			}
			ownerRoles := []uint64{}
			if len(ownerRolesStr) > 0 {
				ownerRoleStrArr := strings.Split(ownerRolesStr, ",")
				for _, roleStr := range ownerRoleStrArr {
					role, err := strconv.Atoi(roleStr)
					if err != nil {
						return err
					}
					ownerRoles = append(ownerRoles, uint64(role))
				}
			}

			ownerAccountsStr, err := cmd.Flags().GetString(FlagOwnerAccounts)
			if err != nil {
				return fmt.Errorf("invalid owner accounts: %w", err)
			}
			ownerAccounts := []string{}
			if len(ownerAccountsStr) > 0 {
				ownerAccounts = strings.Split(ownerAccountsStr, ",")
			}

			beneficiary, err := parseBeneficiaryRolesAndAccounts(cmd)
			if err != nil {
				return err
			}

			dynamicRate, err := cmd.Flags().GetBool(FlagDynamicRate)
			if err != nil {
				return fmt.Errorf("invalid dynamic rate: %w", err)
			}

			dynamicRatePeriod, err := cmd.Flags().GetUint64(FlagDynamicRatePeriod)
			if err != nil {
				return fmt.Errorf("invalid dynamic rate: %w", err)
			}

			msg, err := govtypes.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				types.NewUpdateSpendingPoolProposal(
					name, uint64(claimStart), uint64(claimEnd), rates,
					sdk.NewDecWithPrec(int64(voteQuorum), 2), uint64(votePeriod), uint64(voteEnactment),
					types.PermInfo{
						OwnerRoles:    ownerRoles,
						OwnerAccounts: ownerAccounts,
					},
					beneficiary,
					dynamicRate,
					dynamicRatePeriod,
				),
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

	cmd.Flags().String(FlagName, "", "The name of the spending pool.")
	cmd.Flags().Int32(FlagClaimStart, 0, "The claim start timestamp of the spending pool.")
	cmd.Flags().Int32(FlagClaimEnd, 0, "The claim end timestamp of the spending pool.")
	cmd.Flags().String(FlagRates, "", "reward rates of the spending pool.")
	cmd.Flags().Int32(FlagVoteQuorum, 0, "vote quorum of the spending pool.")
	cmd.Flags().Int32(FlagVotePeriod, 0, "vote period of the spending pool.")
	cmd.Flags().Int32(FlagVoteEnactment, 0, "vote enactment period of the spending pool.")
	cmd.Flags().String(FlagOwnerRoles, "", "owner roles of the spending pool.")
	cmd.Flags().String(FlagOwnerAccounts, "", "owner accounts of the spending pool.")
	cmd.Flags().String(FlagBeneficiaryRoles, "", "beneficiary roles of the spending pool.")
	cmd.Flags().String(FlagBeneficiaryAccounts, "", "beneficiary accounts of the spending pool.")
	cmd.Flags().String(FlagBeneficiaryRoleWeights, "", "beneficiary role weights on the spending pool.")
	cmd.Flags().String(FlagBeneficiaryAccountWeights, "", "beneficiary account weights on the spending pool.")
	cmd.Flags().Bool(FlagDynamicRate, false, "flag to dynamically calculate rates on the spending pool.")
	cmd.Flags().String(FlagDynamicRatePeriod, "", "dynamic rate recalculation period on the spending pool.")

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

			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return fmt.Errorf("invalid name: %w", err)
			}

			msg, err := govtypes.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				types.NewSpendingPoolDistributionProposal(name),
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

	cmd.Flags().String(FlagName, "", "The name of the spending pool.")

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

			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return fmt.Errorf("invalid name: %w", err)
			}

			beneficiaryAccountsStr, err := cmd.Flags().GetString(FlagBeneficiaryAccounts)
			if err != nil {
				return fmt.Errorf("invalid beneficiary accounts: %w", err)
			}
			beneficiaryAccounts := []string{}
			if len(beneficiaryAccountsStr) > 0 {
				beneficiaryAccounts = strings.Split(beneficiaryAccountsStr, ",")
			}

			amountStr, err := cmd.Flags().GetString(FlagAmount)
			if err != nil {
				return fmt.Errorf("invalid amount: %w", err)
			}

			amount, err := sdk.ParseCoinsNormalized(amountStr)
			if err != nil {
				return fmt.Errorf("invalid amount: %w", err)
			}

			msg, err := govtypes.NewMsgSubmitProposal(
				clientCtx.GetFromAddress(),
				title,
				description,
				types.NewSpendingPoolWithdrawProposal(
					name, beneficiaryAccounts, amount,
				),
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

	cmd.Flags().String(FlagName, "", "The name of the spending pool.")
	cmd.Flags().String(FlagBeneficiaryAccounts, "", "beneficiary accounts of the spending pool.")
	cmd.Flags().String(FlagAmount, "", "The amount of coins to deposit on the spending pool.")

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func parseBeneficiaryRolesAndAccounts(cmd *cobra.Command) (types.WeightedPermInfo, error) {

	beneficiaryRolesStr, err := cmd.Flags().GetString(FlagBeneficiaryRoles)
	if err != nil {
		return types.WeightedPermInfo{}, fmt.Errorf("invalid beneficiary roles: %w", err)
	}
	beneficiaryRoleWeightsStr, err := cmd.Flags().GetString(FlagBeneficiaryRoleWeights)
	if err != nil {
		return types.WeightedPermInfo{}, fmt.Errorf("invalid beneficiary role weights: %w", err)
	}
	beneficiaryRoles := []types.WeightedRole{}
	if len(beneficiaryRolesStr) > 0 {
		beneficiaryRoleStrArr := strings.Split(beneficiaryRolesStr, ",")
		beneficiaryRoleWeightStrArr := strings.Split(beneficiaryRoleWeightsStr, ",")
		if len(beneficiaryRoleStrArr) != len(beneficiaryRoleWeightStrArr) {
			return types.WeightedPermInfo{}, fmt.Errorf("beneficiary role and weight count mismatch")
		}
		for index, roleStr := range beneficiaryRoleStrArr {
			role, err := strconv.Atoi(roleStr)
			if err != nil {
				return types.WeightedPermInfo{}, err
			}
			weight, err := sdk.NewDecFromStr(beneficiaryRoleWeightStrArr[index])
			if err != nil {
				return types.WeightedPermInfo{}, err
			}
			beneficiaryRoles = append(beneficiaryRoles, types.WeightedRole{
				Role:   uint64(role),
				Weight: weight,
			})
		}
	}

	beneficiaryAccountsStr, err := cmd.Flags().GetString(FlagBeneficiaryAccounts)
	if err != nil {
		return types.WeightedPermInfo{}, fmt.Errorf("invalid beneficiary accounts: %w", err)
	}
	beneficiaryAccountWeightsStr, err := cmd.Flags().GetString(FlagBeneficiaryAccountWeights)
	if err != nil {
		return types.WeightedPermInfo{}, fmt.Errorf("invalid beneficiary accounts: %w", err)
	}
	beneficiaryAccounts := []types.WeightedAccount{}
	if len(beneficiaryAccountsStr) > 0 {
		beneficiaryAccountStrArr := strings.Split(beneficiaryAccountsStr, ",")
		beneficiaryAccountWeightStrArr := strings.Split(beneficiaryAccountWeightsStr, ",")
		if len(beneficiaryAccountStrArr) != len(beneficiaryAccountWeightStrArr) {
			return types.WeightedPermInfo{}, fmt.Errorf("beneficiary account and weight count mismatch")
		}
		for index, account := range beneficiaryAccountStrArr {
			weight, err := sdk.NewDecFromStr(beneficiaryAccountWeightStrArr[index])
			if err != nil {
				return types.WeightedPermInfo{}, err
			}
			beneficiaryAccounts = append(beneficiaryAccounts, types.WeightedAccount{
				Account: account,
				Weight:  weight,
			})
		}
	}
	return types.WeightedPermInfo{
		Roles:    beneficiaryRoles,
		Accounts: beneficiaryAccounts,
	}, nil
}
