package cli

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/KiraCore/sekai/x/genutil"
	upgradetypes "github.com/KiraCore/sekai/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
)

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
			cdc := clientCtx.JSONCodec

			genDoc, err := tmtypes.GenesisDocFromFile(args[0])
			if err != nil {
				return errors.Wrapf(err, "failed to read genesis file %s", args[0])
			}

			var genesisState map[string]json.RawMessage
			if err = json.Unmarshal(genDoc.AppState, &genesisState); err != nil {
				return errors.Wrap(err, "failed to unmarshal genesis state")
			}

			if err = mbm.ValidateGenesis(cdc, txEncCfg, genesisState); err != nil {
				return errors.Wrap(err, "failed to validate genesis state")
			}

			upgradeGenesis := upgradetypes.GenesisState{}
			cdc.MustUnmarshalJSON(genesisState[upgradetypes.ModuleName], &upgradeGenesis)

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

			appState, err := json.MarshalIndent(genesisState, "", " ")
			if err != nil {
				return errors.Wrap(err, "Failed to marshall default genesis state")
			}

			genDoc.AppState = appState
			if err = genutil.ExportGenesisFile(genDoc, args[1]); err != nil {
				return errors.Wrap(err, "Failed to export gensis file")
			}
			return nil
		},
	}

	return cmd
}
