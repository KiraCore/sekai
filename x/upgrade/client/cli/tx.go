package cli

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/KiraCore/sekai/x/upgrade/types"
)

const (
	// FlagResourceId            = "resource-id"
	// FlagResourceGit           = "resource-git"
	// FlagResourceCheckout      = "resource-checkout"
	// FlagResourceChecksum      = "resource-checksum"
	FlagName                  = "name"
	FlagResources             = "resources"
	FlagHeight                = "height"
	FlagMinUpgradeTime        = "min-upgrade-time"
	FlagOldChainId            = "old-chain-id"
	FlagNewChainId            = "new-chain-id"
	FlagRollbackMemo          = "rollback-memo"
	FlagMaxEnrollmentDuration = "max-enrollment-duration"
	FlagUpgradeMemo           = "upgrade-memo"
	FlagInstateUpgrade        = "instate-upgrade"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Upgrade transaction subcommands",
	}

	cmd.AddCommand(
		GetTxSetUpgradePlan(),
	)

	return cmd
}

func GetTxSetUpgradePlan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-plan",
		Short: "Set upgrade plan",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			// resoureId, err := cmd.Flags().GetString(FlagResourceId)
			// if err != nil {
			// 	return fmt.Errorf("invalid resource id")
			// }

			// resourceGit, err := cmd.Flags().GetString(FlagResourceGit)
			// if err != nil {
			// 	return fmt.Errorf("invalid resource git")
			// }

			// resourceCheckout, err := cmd.Flags().GetString(FlagResourceCheckout)
			// if err != nil {
			// 	return fmt.Errorf("invalid resource checkout")
			// }

			// resourceChecksum, err := cmd.Flags().GetString(FlagResourceChecksum)
			// if err != nil {
			// 	return fmt.Errorf("invalid resource checksum")
			// }

			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return fmt.Errorf("invalid upgrade name")
			}

			resourcesJson, err := cmd.Flags().GetString(FlagResources)
			if err != nil {
				return fmt.Errorf("invalid resources json")
			}

			resources := []types.Resource{}
			err = json.Unmarshal([]byte(resourcesJson), &resources)
			if err != nil {
				return err
			}

			height, err := cmd.Flags().GetInt64(FlagHeight)
			if err != nil {
				return fmt.Errorf("invalid height")
			}

			minUpgradeTime, err := cmd.Flags().GetInt64(FlagMinUpgradeTime)
			if err != nil {
				return fmt.Errorf("invalid min halt time")
			}

			oldChainId, err := cmd.Flags().GetString(FlagOldChainId)
			if err != nil {
				return fmt.Errorf("invalid old chain id")
			}

			newChainId, err := cmd.Flags().GetString(FlagNewChainId)
			if err != nil {
				return fmt.Errorf("invalid new chain id")
			}

			rollBackMemo, err := cmd.Flags().GetString(FlagRollbackMemo)
			if err != nil {
				return fmt.Errorf("invalid rollback memo")
			}

			maxEnrollmentDuration, err := cmd.Flags().GetInt64(FlagMaxEnrollmentDuration)
			if err != nil {
				return fmt.Errorf("invalid max enrollment duration")
			}

			upgradeMemo, err := cmd.Flags().GetString(FlagUpgradeMemo)
			if err != nil {
				return fmt.Errorf("invalid upgrade memo")
			}

			instateUpgrade, err := cmd.Flags().GetBool(FlagInstateUpgrade)
			if err != nil {
				return fmt.Errorf("invalid instate upgrade flag")
			}

			msg := types.NewMsgProposalSoftwareUpgradeRequest(
				clientCtx.FromAddress,
				name,
				// resoureId, resourceGit, resourceCheckout, resourceChecksum,
				resources,
				height, minUpgradeTime, oldChainId, newChainId, rollBackMemo, maxEnrollmentDuration, upgradeMemo,
				instateUpgrade,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	// cmd.Flags().String(FlagResourceId, "", "id of resource")
	// cmd.Flags().String(FlagResourceGit, "", "git of resource")
	// cmd.Flags().String(FlagResourceCheckout, "", "checkout of resource")
	// cmd.Flags().String(FlagResourceChecksum, "", "checksum of resource")
	cmd.Flags().String(FlagName, "upgrade1", "upgrade name")
	cmd.Flags().String(FlagResources, "[]", "resource info")
	cmd.Flags().Int64(FlagHeight, 0, "upgrade height")
	cmd.Flags().Int64(FlagMinUpgradeTime, 0, "min halt time")
	cmd.Flags().String(FlagOldChainId, "", "old chain id")
	cmd.Flags().String(FlagNewChainId, "", "new chain id")
	cmd.Flags().String(FlagRollbackMemo, "", "rollback memo")
	cmd.Flags().Int64(FlagMaxEnrollmentDuration, 0, "max enrollment duration")
	cmd.Flags().String(FlagUpgradeMemo, "", "upgrade memo")
	cmd.Flags().Bool(FlagInstateUpgrade, true, "instate upgrade flag")

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
