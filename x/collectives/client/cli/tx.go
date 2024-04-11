package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/KiraCore/sekai/x/collectives/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

const (
	FlagTitle                 = "title"
	FlagDescription           = "description"
	FlagCollectiveName        = "collective-name"
	FlagCollectiveDescription = "collective-description"
	FlagCollectiveStatus      = "collective-status"
	FlagBonds                 = "bonds"
	FlagDepositAny            = "deposit-any"
	FlagDepositRoles          = "deposit-roles"
	FlagDepositAccounts       = "deposit-accounts"
	FlagOwnerRoles            = "owner-roles"
	FlagOwnerAccounts         = "owner-accounts"
	FlagWeightedSpendingPools = "weighted-spending-pools"
	FlagClaimStart            = "claim-start"
	FlagClaimPeriod           = "claim-period"
	FlagClaimEnd              = "claim-end"
	FlagVoteQuorum            = "vote-quorum"
	FlagVotePeriod            = "vote-period"
	FlagVoteEnactment         = "vote-enactment"
	FlagLocking               = "locking"
	FlagDonation              = "donation"
	FlagDonationLock          = "donation-lock"
	FlagAddr                  = "addr"
	FlagAmounts               = "amounts"
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
		GetTxProposalCollectiveSendDonationCmd(),
		GetTxProposalCollectiveUpdateCmd(),
		GetTxProposalCollectiveRemoveCmd(),
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

			name, err := cmd.Flags().GetString(FlagCollectiveName)
			if err != nil {
				return fmt.Errorf("invalid name: %w", err)
			}
			description, err := cmd.Flags().GetString(FlagCollectiveDescription)
			if err != nil {
				return fmt.Errorf("invalid description: %w", err)
			}
			bondsStr, err := cmd.Flags().GetString(FlagBonds)
			if err != nil {
				return fmt.Errorf("invalid bonds: %w", err)
			}
			bonds, err := sdk.ParseCoinsNormalized(bondsStr)
			if err != nil {
				return fmt.Errorf("invalid bonds: %w", err)
			}

			depositAny, err := cmd.Flags().GetBool(FlagDepositAny)
			if err != nil {
				return fmt.Errorf("invalid deposit any: %w", err)
			}
			depositRolesStr, err := cmd.Flags().GetString(FlagDepositRoles)
			if err != nil {
				return fmt.Errorf("invalid deposit roles: %w", err)
			}
			depositAccountsStr, err := cmd.Flags().GetString(FlagDepositAccounts)
			if err != nil {
				return fmt.Errorf("invalid deposit accounts: %w", err)
			}

			depositRoles, depositAccounts, err := parseRolesAccounts(depositRolesStr, depositAccountsStr)
			if err != nil {
				return fmt.Errorf("invalid deposit role/accounts: %w", err)
			}

			depositWhitelist := types.DepositWhitelist{
				Any:      depositAny,
				Accounts: depositAccounts,
				Roles:    depositRoles,
			}

			ownerRolesStr, err := cmd.Flags().GetString(FlagOwnerRoles)
			if err != nil {
				return fmt.Errorf("invalid deposit roles: %w", err)
			}
			ownerAccountsStr, err := cmd.Flags().GetString(FlagOwnerAccounts)
			if err != nil {
				return fmt.Errorf("invalid deposit accounts: %w", err)
			}

			ownerRoles, ownerAccounts, err := parseRolesAccounts(ownerRolesStr, ownerAccountsStr)
			if err != nil {
				return fmt.Errorf("invalid owner roles/accounts: %w", err)
			}

			ownerWhitelist := types.OwnersWhitelist{
				Accounts: ownerAccounts,
				Roles:    ownerRoles,
			}

			weightedSpendingPoolsStr, err := cmd.Flags().GetString(FlagWeightedSpendingPools)
			if err != nil {
				return fmt.Errorf("invalid weighted spending pools: %w", err)
			}

			weightedSpendingPools, err := parseWeightedSpendingPools(weightedSpendingPoolsStr)
			if err != nil {
				return fmt.Errorf("invalid weighted spending pools: %w", err)
			}

			claimStart, err := cmd.Flags().GetUint64(FlagClaimStart)
			if err != nil {
				return fmt.Errorf("invalid claim start: %w", err)
			}
			claimPeriod, err := cmd.Flags().GetUint64(FlagClaimPeriod)
			if err != nil {
				return fmt.Errorf("invalid claim period: %w", err)
			}
			claimEnd, err := cmd.Flags().GetUint64(FlagClaimEnd)
			if err != nil {
				return fmt.Errorf("invalid claim end: %w", err)
			}
			voteQuorum, err := cmd.Flags().GetUint64(FlagVoteQuorum)
			if err != nil {
				return fmt.Errorf("invalid vote quorum: %w", err)
			}
			votePeriod, err := cmd.Flags().GetUint64(FlagVotePeriod)
			if err != nil {
				return fmt.Errorf("invalid vote period: %w", err)
			}
			voteEnactment, err := cmd.Flags().GetUint64(FlagVoteEnactment)
			if err != nil {
				return fmt.Errorf("invalid vote period: %w", err)
			}

			msg := types.NewMsgCreateCollective(
				clientCtx.FromAddress,
				name, description,
				bonds,
				depositWhitelist,
				ownerWhitelist,
				weightedSpendingPools,
				claimStart, claimPeriod, claimEnd,
				sdk.NewDecWithPrec(int64(voteQuorum), 2), votePeriod, voteEnactment,
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

	cmd.Flags().String(FlagCollectiveName, "", "the name of the collective.")
	cmd.Flags().String(FlagCollectiveDescription, "", "the description of the collective.")
	cmd.Flags().String(FlagBonds, "", "the bonds to put on the collective.")
	cmd.Flags().Bool(FlagDepositAny, false, "flag to enable anyone to deposit on the collective.")
	cmd.Flags().String(FlagDepositRoles, "", "roles to deposit on the collective.")
	cmd.Flags().String(FlagDepositAccounts, "", "accounts to deposit on the collective.")
	cmd.Flags().String(FlagOwnerRoles, "", "owner roles on the collective.")
	cmd.Flags().String(FlagOwnerAccounts, "", "owner accounts on the collective.")
	cmd.Flags().String(FlagWeightedSpendingPools, "", "weighted spending pools configuration for collective. e.g. pool1#0.1,pool2#0.9")
	cmd.Flags().Uint64(FlagClaimStart, 0, "claim start timestamp of the collective")
	cmd.Flags().Uint64(FlagClaimPeriod, 0, "claim period of the collective")
	cmd.Flags().Uint64(FlagClaimEnd, 0, "claim end timestamp of the collective")
	cmd.Flags().Uint64(FlagVoteQuorum, 0, "vote quorum of the collective")
	cmd.Flags().Uint64(FlagVotePeriod, 0, "vote period of the collective")
	cmd.Flags().Uint64(FlagVoteEnactment, 0, "vote enactment of the collective")

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

			name, err := cmd.Flags().GetString(FlagCollectiveName)
			if err != nil {
				return fmt.Errorf("invalid name: %w", err)
			}
			bondsStr, err := cmd.Flags().GetString(FlagBonds)
			if err != nil {
				return fmt.Errorf("invalid bonds: %w", err)
			}
			bonds, err := sdk.ParseCoinsNormalized(bondsStr)
			if err != nil {
				return fmt.Errorf("invalid bonds: %w", err)
			}

			msg := types.NewMsgBondCollective(
				clientCtx.FromAddress, name, bonds,
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

	cmd.Flags().String(FlagCollectiveName, "", "the name of the collective.")
	cmd.Flags().String(FlagBonds, "", "the bonds to put on the collective.")

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

			name, err := cmd.Flags().GetString(FlagCollectiveName)
			if err != nil {
				return fmt.Errorf("invalid name: %w", err)
			}

			locking, err := cmd.Flags().GetUint64(FlagLocking)
			if err != nil {
				return fmt.Errorf("invalid locking: %w", err)
			}

			donationStr, err := cmd.Flags().GetString(FlagDonation)
			if err != nil {
				return fmt.Errorf("invalid donation: %w", err)
			}
			donation, err := sdk.NewDecFromStr(donationStr)
			if err != nil {
				return fmt.Errorf("invalid donation: %w", err)
			}

			donationLock, err := cmd.Flags().GetBool(FlagDonationLock)
			if err != nil {
				return fmt.Errorf("invalid donation lock: %w", err)
			}
			msg := types.NewMsgDonateCollective(
				clientCtx.FromAddress,
				name, locking, donation, donationLock,
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

	cmd.Flags().String(FlagCollectiveName, "", "the name of the collective.")
	cmd.Flags().Uint64(FlagLocking, 0, "the locking duration of the collective.")
	cmd.Flags().String(FlagDonation, "", "the donation percentage on the collective contributor rewards.")
	cmd.Flags().Bool(FlagDonationLock, false, "flag to lock contribution on the collective.")

	return cmd
}

// GetTxWithdrawCollectiveCmd can be sent by any whitelisted “contributor” to withdraw
// their tokens (unless locking is enabled)
func GetTxWithdrawCollectiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-collective",
		Short: "sent by any whitelisted “contributor” to withdraw their tokens (unless locking is enabled)",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			name, err := cmd.Flags().GetString(FlagCollectiveName)
			if err != nil {
				return fmt.Errorf("invalid name: %w", err)
			}
			msg := types.NewMsgWithdrawCollective(
				clientCtx.FromAddress,
				name,
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
	cmd.Flags().String(FlagCollectiveName, "", "the name of the collective.")

	return cmd
}

func parseRolesAccounts(rolesStr, accountsStr string) ([]uint64, []string, error) {
	roleStrArr := strings.Split(rolesStr, ",")
	if rolesStr == "" {
		roleStrArr = []string{}
	}
	roles := []uint64{}
	for _, roleStr := range roleStrArr {
		role, err := strconv.Atoi(roleStr)
		if err != nil {
			return []uint64{}, []string{}, fmt.Errorf("invalid role: %w", err)
		}
		roles = append(roles, uint64(role))
	}
	accounts := strings.Split(accountsStr, ",")
	if accountsStr == "" {
		accounts = []string{}
	}
	return roles, accounts, nil
}

func parseWeightedSpendingPools(str string) ([]types.WeightedSpendingPool, error) {
	pools := []types.WeightedSpendingPool{}
	weightedPoolsStr := strings.Split(str, ",")
	for _, weightedPool := range weightedPoolsStr {
		split := strings.Split(weightedPool, "#")
		if len(split) != 2 {
			return []types.WeightedSpendingPool{}, fmt.Errorf("invalid sub weighted pool")
		}
		weight, err := sdk.NewDecFromStr(split[1])
		if err != nil {
			return []types.WeightedSpendingPool{}, err
		}
		pools = append(pools, types.WeightedSpendingPool{
			Name:   split[0],
			Weight: weight,
		})
	}
	return pools, nil
}

// GetTxProposalCollectiveSendDonationCmd implement cli command for ProposalCollectiveSendDonation
func GetTxProposalCollectiveSendDonationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-send-donation",
		Short: "Create a proposal to withdraw donation",
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

			name, err := cmd.Flags().GetString(FlagCollectiveName)
			if err != nil {
				return fmt.Errorf("invalid name: %w", err)
			}
			address, err := cmd.Flags().GetString(FlagAddr)
			if err != nil {
				return fmt.Errorf("invalid address: %w", err)
			}
			amountsStr, err := cmd.Flags().GetString(FlagAmounts)
			if err != nil {
				return fmt.Errorf("invalid amounts: %w", err)
			}
			amounts, err := sdk.ParseCoinsNormalized(amountsStr)
			if err != nil {
				return fmt.Errorf("invalid amounts: %w", err)
			}

			msg, err := govtypes.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				types.NewProposalCollectiveSendDonation(
					name, address, amounts,
				),
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

	cmd.Flags().String(FlagCollectiveName, "", "the name of the collective.")
	cmd.Flags().String(FlagAddr, "", "The address to receive from donation pool.")
	cmd.Flags().String(FlagAmounts, "", "The amounts to receive from donation pool.")

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxProposalCollectiveUpdateCmd implement cli command for ProposalCollectiveUpdate
func GetTxProposalCollectiveUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-collective-update",
		Short: "Create a proposal to update collective",
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

			name, err := cmd.Flags().GetString(FlagCollectiveName)
			if err != nil {
				return fmt.Errorf("invalid collective name: %w", err)
			}

			collectiveDesc, err := cmd.Flags().GetString(FlagCollectiveDescription)
			if err != nil {
				return fmt.Errorf("invalid collective description: %w", err)
			}

			statusStr, err := cmd.Flags().GetString(FlagCollectiveStatus)
			if err != nil {
				return fmt.Errorf("invalid collective description: %w", err)
			}
			status := types.CollectiveStatus_value[statusStr]

			depositAny, err := cmd.Flags().GetBool(FlagDepositAny)
			if err != nil {
				return fmt.Errorf("invalid deposit any: %w", err)
			}
			depositRolesStr, err := cmd.Flags().GetString(FlagDepositRoles)
			if err != nil {
				return fmt.Errorf("invalid deposit roles: %w", err)
			}
			depositAccountsStr, err := cmd.Flags().GetString(FlagDepositAccounts)
			if err != nil {
				return fmt.Errorf("invalid deposit accounts: %w", err)
			}

			depositRoles, depositAccounts, err := parseRolesAccounts(depositRolesStr, depositAccountsStr)
			if err != nil {
				return fmt.Errorf("invalid deposit role/accounts: %w", err)
			}

			depositWhitelist := types.DepositWhitelist{
				Any:      depositAny,
				Accounts: depositAccounts,
				Roles:    depositRoles,
			}

			ownerRolesStr, err := cmd.Flags().GetString(FlagOwnerRoles)
			if err != nil {
				return fmt.Errorf("invalid deposit roles: %w", err)
			}
			ownerAccountsStr, err := cmd.Flags().GetString(FlagOwnerAccounts)
			if err != nil {
				return fmt.Errorf("invalid deposit accounts: %w", err)
			}

			ownerRoles, ownerAccounts, err := parseRolesAccounts(ownerRolesStr, ownerAccountsStr)
			if err != nil {
				return fmt.Errorf("invalid owner roles/accounts: %w", err)
			}

			ownerWhitelist := types.OwnersWhitelist{
				Accounts: ownerAccounts,
				Roles:    ownerRoles,
			}

			weightedSpendingPoolsStr, err := cmd.Flags().GetString(FlagWeightedSpendingPools)
			if err != nil {
				return fmt.Errorf("invalid weighted spending pools: %w", err)
			}

			weightedSpendingPools, err := parseWeightedSpendingPools(weightedSpendingPoolsStr)
			if err != nil {
				return fmt.Errorf("invalid weighted spending pools: %w", err)
			}

			claimStart, err := cmd.Flags().GetUint64(FlagClaimStart)
			if err != nil {
				return fmt.Errorf("invalid claim start: %w", err)
			}
			claimPeriod, err := cmd.Flags().GetUint64(FlagClaimPeriod)
			if err != nil {
				return fmt.Errorf("invalid claim period: %w", err)
			}
			claimEnd, err := cmd.Flags().GetUint64(FlagClaimEnd)
			if err != nil {
				return fmt.Errorf("invalid claim end: %w", err)
			}
			voteQuorum, err := cmd.Flags().GetUint64(FlagVoteQuorum)
			if err != nil {
				return fmt.Errorf("invalid vote quorum: %w", err)
			}
			votePeriod, err := cmd.Flags().GetUint64(FlagVotePeriod)
			if err != nil {
				return fmt.Errorf("invalid vote period: %w", err)
			}
			voteEnactment, err := cmd.Flags().GetUint64(FlagVoteEnactment)
			if err != nil {
				return fmt.Errorf("invalid vote period: %w", err)
			}

			msg, err := govtypes.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				types.NewProposalCollectiveUpdate(
					name, collectiveDesc, types.CollectiveStatus(status),
					depositWhitelist,
					ownerWhitelist,
					weightedSpendingPools,
					claimStart, claimPeriod, claimEnd,
					sdk.NewDecWithPrec(int64(voteQuorum), 2), votePeriod, voteEnactment,
				),
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

	cmd.Flags().String(FlagCollectiveName, "", "the name of the collective.")
	cmd.Flags().String(FlagCollectiveDescription, "", "the description of the collective.")
	cmd.Flags().String(FlagCollectiveStatus, "", "the status of the collective.")
	cmd.Flags().Bool(FlagDepositAny, false, "flag to enable anyone to deposit on the collective.")
	cmd.Flags().String(FlagDepositRoles, "", "roles to deposit on the collective.")
	cmd.Flags().String(FlagDepositAccounts, "", "accounts to deposit on the collective.")
	cmd.Flags().String(FlagOwnerRoles, "", "owner roles on the collective.")
	cmd.Flags().String(FlagOwnerAccounts, "", "owner accounts on the collective.")
	cmd.Flags().String(FlagWeightedSpendingPools, "", "weighted spending pools configuration for collective. e.g. pool1#0.1,pool2#0.9")
	cmd.Flags().Uint64(FlagClaimStart, 0, "claim start timestamp of the collective")
	cmd.Flags().Uint64(FlagClaimPeriod, 0, "claim period of the collective")
	cmd.Flags().Uint64(FlagClaimEnd, 0, "claim end timestamp of the collective")
	cmd.Flags().Uint64(FlagVoteQuorum, 0, "vote quorum of the collective")
	cmd.Flags().Uint64(FlagVotePeriod, 0, "vote period of the collective")
	cmd.Flags().Uint64(FlagVoteEnactment, 0, "vote enactment of the collective")

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxProposalCollectiveRemoveCmd implement cli command for ProposalCollectiveRemove
func GetTxProposalCollectiveRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-remove-collective",
		Short: "Create a proposal to withdraw collective",
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

			name, err := cmd.Flags().GetString(FlagCollectiveName)
			if err != nil {
				return fmt.Errorf("invalid name: %w", err)
			}

			msg, err := govtypes.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				types.NewProposalCollectiveRemove(
					name,
				),
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

	cmd.Flags().String(FlagCollectiveName, "", "the name of the collective.")

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
