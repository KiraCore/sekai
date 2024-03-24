package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/KiraCore/sekai/x/genutil"
	v01228govtypes "github.com/KiraCore/sekai/x/gov/legacy/v01228"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	v03123upgradetypes "github.com/KiraCore/sekai/x/upgrade/legacy/v03123"
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

func upgradedPlan(plan *v03123upgradetypes.PlanV03123) *upgradetypes.Plan {
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

func upgradedResources(resources []v03123upgradetypes.ResourceV03123) []upgradetypes.Resource {
	upgraded := []upgradetypes.Resource{}
	for _, resource := range resources {
		upgraded = append(upgraded, upgradetypes.Resource{
			Id:       resource.Id,
			Url:      resource.Git,
			Version:  resource.Checkout,
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
			} else {
				if err = mbm.ValidateGenesis(cdc, txEncCfg, genesisState); err != nil {
					return errors.Wrap(err, "failed to validate genesis state")
				}

				upgradeGenesisV03123 := v03123upgradetypes.GenesisStateV03123{}
				err = cdc.UnmarshalJSON(genesisState[upgradetypes.ModuleName], &upgradeGenesisV03123)
				if err == nil { // which means old upgrade genesis
					upgradeGenesis := upgradetypes.GenesisState{
						Version:     "v0.3.1.24",
						CurrentPlan: upgradedPlan(upgradeGenesisV03123.CurrentPlan),
						NextPlan:    upgradedPlan(upgradeGenesisV03123.NextPlan),
					}
					genesisState[upgradetypes.ModuleName] = cdc.MustMarshalJSON(&upgradeGenesis)
				} else {
					fmt.Println("error exists v0.3.1.23 upgrade genesis parsing", err)
				}

				upgradeGenesis := upgradetypes.GenesisState{}
				cdc.MustUnmarshalJSON(genesisState[upgradetypes.ModuleName], &upgradeGenesis)
				oldVersion := upgradeGenesis.Version
				if upgradeGenesis.Version == "" {
					upgradeGenesis.Version = "v0.1.22.11"
					fmt.Println("upgraded the upgrade module genesis to v0.1.22.11")
				}

				if upgradeGenesis.NextPlan == nil {
					return fmt.Errorf("next plan is not available")
				}

				if genDoc.ChainID != upgradeGenesis.NextPlan.OldChainId {
					return fmt.Errorf("next plan has different oldchain id, current chain_id=%s, next_plan.old_chain_id=%s", genDoc.ChainID, upgradeGenesis.NextPlan.OldChainId)
				}

				genDoc.ChainID = upgradeGenesis.NextPlan.NewChainId
				upgradeGenesis.CurrentPlan = upgradeGenesis.NextPlan
				upgradeGenesis.NextPlan = nil

				genesisState[upgradetypes.ModuleName] = cdc.MustMarshalJSON(&upgradeGenesis)

				govGenesisV01228 := v01228govtypes.GenesisStateV01228{}
				err = cdc.UnmarshalJSON(genesisState[govtypes.ModuleName], &govGenesisV01228)

				// we are referencing oldPlan.name to determine upgrade genesis or not
				if err == nil && oldVersion == "" { // it means v0.1.22.8 gov genesis
					govGenesis := govtypes.GenesisState{
						StartingProposalId: govGenesisV01228.StartingProposalId,
						NextRoleId:         govtypes.DefaultGenesis().NextRoleId,
						Roles:              govtypes.DefaultGenesis().Roles,
						RolePermissions:    govGenesisV01228.Permissions,
						NetworkActors:      govGenesisV01228.NetworkActors,
						NetworkProperties: &govtypes.NetworkProperties{
							MinTxFee:                         govGenesisV01228.NetworkProperties.MinTxFee,
							MaxTxFee:                         govGenesisV01228.NetworkProperties.MaxTxFee,
							VoteQuorum:                       govGenesisV01228.NetworkProperties.VoteQuorum,
							MinimumProposalEndTime:           govGenesisV01228.NetworkProperties.ProposalEndTime,
							ProposalEnactmentTime:            govGenesisV01228.NetworkProperties.ProposalEnactmentTime,
							MinProposalEndBlocks:             govGenesisV01228.NetworkProperties.MinProposalEndBlocks,
							MinProposalEnactmentBlocks:       govGenesisV01228.NetworkProperties.MinProposalEnactmentBlocks,
							EnableForeignFeePayments:         govGenesisV01228.NetworkProperties.EnableForeignFeePayments,
							MischanceRankDecreaseAmount:      govGenesisV01228.NetworkProperties.MischanceRankDecreaseAmount,
							MaxMischance:                     govGenesisV01228.NetworkProperties.MaxMischance,
							MischanceConfidence:              govGenesisV01228.NetworkProperties.MischanceConfidence,
							InactiveRankDecreasePercent:      sdk.NewDecWithPrec(int64(govGenesisV01228.NetworkProperties.InactiveRankDecreasePercent), 2),
							MinValidators:                    govGenesisV01228.NetworkProperties.MinValidators,
							PoorNetworkMaxBankSend:           govGenesisV01228.NetworkProperties.PoorNetworkMaxBankSend,
							UnjailMaxTime:                    govGenesisV01228.NetworkProperties.JailMaxTime,
							EnableTokenWhitelist:             govGenesisV01228.NetworkProperties.EnableTokenWhitelist,
							EnableTokenBlacklist:             govGenesisV01228.NetworkProperties.EnableTokenBlacklist,
							MinIdentityApprovalTip:           govGenesisV01228.NetworkProperties.MinIdentityApprovalTip,
							UniqueIdentityKeys:               govGenesisV01228.NetworkProperties.UniqueIdentityKeys,
							UbiHardcap:                       6000_000,
							ValidatorsFeeShare:               sdk.NewDecWithPrec(50, 2), // 50%
							InflationRate:                    sdk.NewDecWithPrec(18, 2), // 18%
							InflationPeriod:                  31557600,                  // 1 year
							UnstakingPeriod:                  2629800,                   // 1 month
							MaxDelegators:                    100,
							MinDelegationPushout:             10,
							SlashingPeriod:                   3600,
							MaxJailedPercentage:              sdk.NewDecWithPrec(25, 2),
							MaxSlashingPercentage:            sdk.NewDecWithPrec(1, 2),
							MinCustodyReward:                 200,
							MaxCustodyTxSize:                 8192,
							MaxCustodyBufferSize:             10,
							AbstentionRankDecreaseAmount:     1,
							MaxAbstention:                    2,
							MinCollectiveBond:                100_000, // in KEX
							MinCollectiveBondingTime:         86400,   // in seconds
							MaxCollectiveOutputs:             10,
							MinCollectiveClaimPeriod:         14400,  // 4hrs
							ValidatorRecoveryBond:            300000, // 300k KEX
							MaxAnnualInflation:               sdk.NewDecWithPrec(35, 2),
							MaxProposalTitleSize:             128,
							MaxProposalDescriptionSize:       1024,
							MaxProposalPollOptionSize:        64,
							MaxProposalPollOptionCount:       128,
							MinDappBond:                      1000000,
							MaxDappBond:                      10000000,
							DappBondDuration:                 604800,
							DappVerifierBond:                 sdk.NewDecWithPrec(1, 3), //0.1%
							DappAutoDenounceTime:             60,                       // 60s
							DappMischanceRankDecreaseAmount:  1,
							DappMaxMischance:                 10,
							DappInactiveRankDecreasePercent:  10,
							DappPoolSlippageDefault:          sdk.NewDecWithPrec(1, 1), // 10%
							MintingFtFee:                     100_000_000_000_000,
							MintingNftFee:                    100_000_000_000_000,
							VetoThreshold:                    sdk.NewDecWithPrec(3340, 2), //33.40%
							AutocompoundIntervalNumBlocks:    17280,
							BridgeAddress:                    "test",
							BridgeCosmosEthereumExchangeRate: sdk.NewDec(10),
							BridgeEthereumCosmosExchangeRate: sdk.NewDecWithPrec(1, 1),
						},
						ExecutionFees:               govGenesisV01228.ExecutionFees,
						PoorNetworkMessages:         govGenesisV01228.PoorNetworkMessages,
						Proposals:                   []govtypes.Proposal{}, // govGenesisV01228.Proposals,
						Votes:                       []govtypes.Vote{},     // govGenesisV01228.Votes,
						DataRegistry:                govGenesisV01228.DataRegistry,
						IdentityRecords:             govGenesisV01228.IdentityRecords,
						LastIdentityRecordId:        govGenesisV01228.LastIdentityRecordId,
						IdRecordsVerifyRequests:     govGenesisV01228.IdRecordsVerifyRequests,
						LastIdRecordVerifyRequestId: govGenesisV01228.LastIdRecordVerifyRequestId,
						ProposalDurations:           make(map[string]uint64),
					}

					genesisState[govtypes.ModuleName] = cdc.MustMarshalJSON(&govGenesis)
				} else {
					fmt.Println("GovGenesis01228 unmarshal test: ", err)
					fmt.Println("Skipping governance module upgrade since it is not v0.1.22.8 genesis")
				}

				// upgrade gov genesis for more role permissions
				govGenesis := govtypes.GenesisState{}
				err = cdc.UnmarshalJSON(genesisState[govtypes.ModuleName], &govGenesis)
				if err == nil {
					govGenesis.RolePermissions[govtypes.RoleSudo] = govtypes.DefaultGenesis().RolePermissions[govtypes.RoleSudo]
					genesisState[govtypes.ModuleName] = cdc.MustMarshalJSON(&govGenesis)
				} else {
					fmt.Println("parse error for latest gov genesis", err)
					fmt.Println("trying to parse v03123 gov genesis for following error on genesis parsing")
					govGenesisV03123 := make(map[string]interface{})
					err = json.Unmarshal(genesisState[govtypes.ModuleName], &govGenesisV03123)
					if err != nil {
						panic(err)
					}
					govGenesisV03123["proposals"] = []govtypes.Proposal{}
					govGenesisV03123["votes"] = []govtypes.Vote{}
					bz, err := json.Marshal(&govGenesisV03123)
					if err != nil {
						panic(err)
					}
					genesisState[govtypes.ModuleName] = bz
				}
			}

			appState, err := json.MarshalIndent(genesisState, "", " ")
			if err != nil {
				return errors.Wrap(err, "Failed to marshal default genesis state")
			}

			genDoc.AppState = appState

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
