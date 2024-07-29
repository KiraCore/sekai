package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	appparams "github.com/KiraCore/sekai/app/params"
	"github.com/KiraCore/sekai/x/genutil"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	v0317upgradetypes "github.com/KiraCore/sekai/x/upgrade/legacy/v0317"
	upgradetypes "github.com/KiraCore/sekai/x/upgrade/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	tmos "github.com/cometbft/cometbft/libs/os"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	FlagJsonMinimize = "json-minimize"
	FlagModulesOnly  = "modules-only"
)

func upgradedPlan(plan *v0317upgradetypes.PlanV0317) *upgradetypes.Plan {
	if plan == nil {
		return nil
	}

	return &upgradetypes.Plan{
		Name:                      plan.Name,
		Resources:                 upgradedResources(plan.Resources),
		UpgradeTime:               plan.UpgradeTime,
		OldChainId:                plan.OldChainId,
		NewChainId:                plan.NewChainId,
		RollbackChecksum:          plan.RollbackChecksum,
		MaxEnrolmentDuration:      plan.MaxEnrolmentDuration,
		InstateUpgrade:            plan.InstateUpgrade,
		RebootRequired:            plan.RebootRequired,
		SkipHandler:               plan.SkipHandler,
		ProposalID:                plan.ProposalID,
		ProcessedNoVoteValidators: plan.ProcessedNoVoteValidators,
	}
}

func upgradedResources(resources []v0317upgradetypes.ResourceV0317) []upgradetypes.Resource {
	upgraded := []upgradetypes.Resource{}
	for _, resource := range resources {
		upgraded = append(upgraded, upgradetypes.Resource{
			Id:       resource.Id,
			Url:      resource.Url,
			Version:  resource.Version,
			Checksum: resource.Checksum,
		})
	}
	return upgraded
}

// GetNewGenesisFromExportedCmd returns new genesis from exported genesis
func GetNewGenesisFromExportedCmd(mbm module.BasicManager, txEncCfg client.TxEncodingConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-genesis-from-exported [path-to-exported.json] [path-to-new.json]",
		Short: "Get new genesis from exported app state json",
		Args:  cobra.ExactArgs(2),
		Long: fmt.Sprintf(`Get new genesis from exported app state json.
- Change chain-id to new_chain_id as indicated by the upgrade plan
- Replace current upgrade plan in the app_state.upgrade with next plan and set next plan to null

Example:
$ %s new-genesis-from-exported exported-genesis.json new-genesis.json
`, version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			cdc := clientCtx.Codec

			genDoc, err := tmtypes.GenesisDocFromFile(args[0])
			if err != nil {
				return errors.Wrapf(err, "failed to read genesis file %s", args[0])
			}

			var genesisState map[string]json.RawMessage
			if err = json.Unmarshal(genDoc.AppState, &genesisState); err != nil {
				return errors.Wrap(err, "failed to unmarshal genesis state")
			}

			modulesCombined, err := cmd.Flags().GetString(FlagModulesOnly)
			if err != nil {
				return err
			}
			if modulesCombined != "" {
				newGenesis := mbm.DefaultGenesis(cdc)
				modules := strings.Split(modulesCombined, ",")
				for _, module := range modules {
					moduleGenesis, ok := genesisState[module]
					if !ok {
						return errors.New("invalid module")
					}
					newGenesis[module] = moduleGenesis
				}
				genesisState = newGenesis
			}

			upgradeGenesisV03123 := v0317upgradetypes.GenesisStateV0317{}
			err = cdc.UnmarshalJSON(genesisState[upgradetypes.ModuleName], &upgradeGenesisV03123)
			if err == nil { // which means old upgrade genesis
				upgradeGenesis := upgradetypes.GenesisState{
					Version:     "v0.3.45",
					CurrentPlan: upgradedPlan(upgradeGenesisV03123.CurrentPlan),
					NextPlan:    upgradedPlan(upgradeGenesisV03123.NextPlan),
				}
				genesisState[upgradetypes.ModuleName] = cdc.MustMarshalJSON(&upgradeGenesis)
			} else {
				fmt.Println("error exists v0.3.17 upgrade genesis parsing", err)
			}

			upgradeGenesis := upgradetypes.GenesisState{}
			cdc.MustUnmarshalJSON(genesisState[upgradetypes.ModuleName], &upgradeGenesis)
			if upgradeGenesis.Version == "" {
				upgradeGenesis.Version = "v0.3.45"
				fmt.Println("upgraded the upgrade module genesis to v0.3.45")
			}

			// if upgradeGenesis.NextPlan == nil {
			// 	return fmt.Errorf("next plan is not available")
			// }

			// if genDoc.ChainID != upgradeGenesis.NextPlan.OldChainId {
			// 	return fmt.Errorf("next plan has different oldchain id, current chain_id=%s, next_plan.old_chain_id=%s", genDoc.ChainID, upgradeGenesis.NextPlan.OldChainId)
			// }
			if upgradeGenesis.NextPlan != nil {
				genDoc.ChainID = upgradeGenesis.NextPlan.NewChainId
			}

			upgradeGenesis.CurrentPlan = upgradeGenesis.NextPlan
			upgradeGenesis.NextPlan = nil

			genesisState[upgradetypes.ModuleName] = cdc.MustMarshalJSON(&upgradeGenesis)

			// upgrade gov genesis for more role permissions
			govGenesis := govtypes.GenesisState{}
			err = cdc.UnmarshalJSON(genesisState[govtypes.ModuleName], &govGenesis)
			if err == nil {
				if govGenesis.DefaultDenom == "" {
					govGenesis.DefaultDenom = appparams.DefaultDenom
				}
				if govGenesis.Bech32Prefix == "" {
					govGenesis.Bech32Prefix = appparams.AccountAddressPrefix
				}
				govGenesis.NetworkProperties.VoteQuorum = sdk.NewDecWithPrec(33, 2)                     // 33%
				govGenesis.NetworkProperties.VetoThreshold = sdk.NewDecWithPrec(3340, 4)                // 33.4%
				govGenesis.NetworkProperties.DappInactiveRankDecreasePercent = sdk.NewDecWithPrec(1, 1) // 10%
				govGenesis.NetworkProperties.SlashingPeriod = 2629800
				genesisState[govtypes.ModuleName] = cdc.MustMarshalJSON(&govGenesis)
			} else {
				fmt.Println("parse error for latest gov genesis", err)
				fmt.Println("trying to parse v0.3.17 gov genesis for following error on genesis parsing")
				govGenesisV0317 := make(map[string]interface{})
				err = json.Unmarshal(genesisState[govtypes.ModuleName], &govGenesisV0317)
				if err != nil {
					panic(err)
				}

				fmt.Println("Setting default gov data", appparams.DefaultDenom, appparams.AccountAddressPrefix)
				govGenesisV0317["default_denom"] = appparams.DefaultDenom
				govGenesisV0317["bech32_prefix"] = appparams.AccountAddressPrefix
				bz, err := json.Marshal(&govGenesisV0317)
				if err != nil {
					panic(err)
				}
				genesisState[govtypes.ModuleName] = bz
			}

			appState, err := json.MarshalIndent(genesisState, "", " ")
			if err != nil {
				return errors.Wrap(err, "Failed to marshal default genesis state")
			}

			genDoc.AppState = appState
			genDoc.InitialHeight = 0

			if jsonMinimize, _ := cmd.Flags().GetBool(FlagJsonMinimize); jsonMinimize {
				genDocBytes, err := tmjson.Marshal(genDoc)
				if err != nil {
					return err
				}
				return tmos.WriteFile(args[1], genDocBytes, 0644)
			}

			if err = genutil.ExportGenesisFile(genDoc, args[1]); err != nil {
				return errors.Wrap(err, "Failed to export genesis file")
			}
			return nil
		},
	}
	cmd.Flags().Bool(FlagJsonMinimize, true, "flag to export genesis in minimized version")
	cmd.Flags().String(FlagModulesOnly, "", "flag to derive only specific modules - one of followings auth,bank,customstaking,customslashing,evidence,consensus,params,upgrade,recovery,customgov,spending,distributor,basket,ubi,tokens,custody,multistaking,collectives,layer2")

	return cmd
}
