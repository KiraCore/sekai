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

			token, err := cmd.Flags().GetString(FlagToken)
			if err != nil {
				return fmt.Errorf("invalid token: %w", err)
			}

			rateStr, err := cmd.Flags().GetString(FlagRate)
			if err != nil {
				return fmt.Errorf("invalid rate: %w", err)
			}
			rate, err := sdk.NewDecFromStr(rateStr)
			if err != nil {
				return fmt.Errorf("invalid rate: %w", err)
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

			beneficiaryRolesStr, err := cmd.Flags().GetString(FlagBeneficiaryRoles)
			if err != nil {
				return fmt.Errorf("invalid beneficiary roles: %w", err)
			}
			beneficiaryRoles := []uint64{}
			if len(beneficiaryRolesStr) > 0 {
				beneficiaryRoleStrArr := strings.Split(beneficiaryRolesStr, ",")
				for _, roleStr := range beneficiaryRoleStrArr {
					role, err := strconv.Atoi(roleStr)
					if err != nil {
						return err
					}
					beneficiaryRoles = append(beneficiaryRoles, uint64(role))
				}
			}

			beneficiaryAccountsStr, err := cmd.Flags().GetString(FlagBeneficiaryAccounts)
			if err != nil {
				return fmt.Errorf("invalid beneficiary accounts: %w", err)
			}
			beneficiaryAccounts := []string{}
			if len(beneficiaryAccountsStr) > 0 {
				beneficiaryAccounts = strings.Split(beneficiaryAccountsStr, ",")
			}

			msg := types.NewMsgCreateSpendingPool(
				name, uint64(claimStart), uint64(claimEnd), token, rate,
				uint64(voteQuorum), uint64(votePeriod), uint64(voteEnactment),
				types.PermInfo{
					OwnerRoles:    ownerRoles,
					OwnerAccounts: ownerAccounts,
				},
				types.PermInfo{
					OwnerRoles:    beneficiaryRoles,
					OwnerAccounts: beneficiaryAccounts,
				},
				clientCtx.GetFromAddress(),
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
	cmd.Flags().String(FlagToken, "", "The reward token of the spending pool.")
	cmd.Flags().String(FlagRate, "", "reward rate of the spending pool.")
	cmd.Flags().Int32(FlagVoteQuorum, 0, "vote quorum of the spending pool.")
	cmd.Flags().Int32(FlagVotePeriod, 0, "vote period of the spending pool.")
	cmd.Flags().Int32(FlagVoteEnactment, 0, "vote enactment period of the spending pool.")
	cmd.Flags().String(FlagOwnerRoles, "", "owner roles of the spending pool.")
	cmd.Flags().String(FlagOwnerAccounts, "", "owner accounts of the spending pool.")
	cmd.Flags().String(FlagBeneficiaryRoles, "", "beneficiary roles of the spending pool.")
	cmd.Flags().String(FlagBeneficiaryAccounts, "", "beneficiary accounts of the spending pool.")

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

			beneficiaryRolesStr, err := cmd.Flags().GetString(FlagBeneficiaryRoles)
			if err != nil {
				return fmt.Errorf("invalid beneficiary roles: %w", err)
			}
			beneficiaryRoles := []uint64{}
			if len(beneficiaryRolesStr) > 0 {
				beneficiaryRoleStrArr := strings.Split(beneficiaryRolesStr, ",")
				for _, roleStr := range beneficiaryRoleStrArr {
					role, err := strconv.Atoi(roleStr)
					if err != nil {
						return err
					}
					beneficiaryRoles = append(beneficiaryRoles, uint64(role))
				}
			}
			beneficiaryAccountsStr, err := cmd.Flags().GetString(FlagBeneficiaryAccounts)
			if err != nil {
				return fmt.Errorf("invalid beneficiary accounts: %w", err)
			}
			beneficiaryAccounts := []string{}
			if len(beneficiaryAccountsStr) > 0 {
				beneficiaryAccounts = strings.Split(beneficiaryAccountsStr, ",")
			}

			msg := types.NewMsgRegisterSpendingPoolBeneficiary(name, types.PermInfo{
				OwnerRoles:    beneficiaryRoles,
				OwnerAccounts: beneficiaryAccounts,
			}, clientCtx.GetFromAddress())

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
	cmd.Flags().String(FlagBeneficiaryRoles, "", "beneficiary roles of the spending pool.")
	cmd.Flags().String(FlagBeneficiaryAccounts, "", "beneficiary accounts of the spending pool.")

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

			token, err := cmd.Flags().GetString(FlagToken)
			if err != nil {
				return fmt.Errorf("invalid token: %w", err)
			}

			rateStr, err := cmd.Flags().GetString(FlagRate)
			if err != nil {
				return fmt.Errorf("invalid rate: %w", err)
			}
			rate, err := sdk.NewDecFromStr(rateStr)
			if err != nil {
				return fmt.Errorf("invalid rate: %w", err)
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

			beneficiaryRolesStr, err := cmd.Flags().GetString(FlagBeneficiaryRoles)
			if err != nil {
				return fmt.Errorf("invalid beneficiary roles: %w", err)
			}
			beneficiaryRoles := []uint64{}
			if len(beneficiaryRolesStr) > 0 {
				beneficiaryRoleStrArr := strings.Split(beneficiaryRolesStr, ",")
				for _, roleStr := range beneficiaryRoleStrArr {
					role, err := strconv.Atoi(roleStr)
					if err != nil {
						return err
					}
					beneficiaryRoles = append(beneficiaryRoles, uint64(role))
				}
			}

			beneficiaryAccountsStr, err := cmd.Flags().GetString(FlagBeneficiaryAccounts)
			if err != nil {
				return fmt.Errorf("invalid beneficiary accounts: %w", err)
			}
			beneficiaryAccounts := []string{}
			if len(beneficiaryAccountsStr) > 0 {
				beneficiaryAccounts = strings.Split(beneficiaryAccountsStr, ",")
			}

			msg, err := govtypes.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				types.NewUpdateSpendingPoolProposal(
					name, uint64(claimStart), uint64(claimEnd), token, rate,
					uint64(voteQuorum), uint64(votePeriod), uint64(voteEnactment),
					types.PermInfo{
						OwnerRoles:    ownerRoles,
						OwnerAccounts: ownerAccounts,
					},
					types.PermInfo{
						OwnerRoles:    beneficiaryRoles,
						OwnerAccounts: beneficiaryAccounts,
					},
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
	cmd.Flags().String(FlagToken, "", "The reward token of the spending pool.")
	cmd.Flags().String(FlagRate, "", "reward rate of the spending pool.")
	cmd.Flags().Int32(FlagVoteQuorum, 0, "vote quorum of the spending pool.")
	cmd.Flags().Int32(FlagVotePeriod, 0, "vote period of the spending pool.")
	cmd.Flags().Int32(FlagVoteEnactment, 0, "vote enactment period of the spending pool.")
	cmd.Flags().String(FlagOwnerRoles, "", "owner roles of the spending pool.")
	cmd.Flags().String(FlagOwnerAccounts, "", "owner accounts of the spending pool.")
	cmd.Flags().String(FlagBeneficiaryRoles, "", "beneficiary roles of the spending pool.")
	cmd.Flags().String(FlagBeneficiaryAccounts, "", "beneficiary accounts of the spending pool.")

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
