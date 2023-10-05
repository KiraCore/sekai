package cli

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	tmjson "github.com/cometbft/cometbft/libs/json"
	tmos "github.com/cometbft/cometbft/libs/os"
)

// GetStandardGenesisExportedCmd returns standard genesis from genesis with spaces
func GetStandardGenesisExportedCmd(mbm module.BasicManager, txEncCfg client.TxEncodingConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-minimized-genesis [path-to-exported.json] [path-to-standard.json]",
		Short: "Get minimized genesis from genesis with spaces",
		Args:  cobra.ExactArgs(2),
		Long: fmt.Sprintf(`Get minimized genesis from not standard genesis.
Example:
$ %s export-minimized-genesis exported-genesis.json new-genesis.json
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

			if err = mbm.ValidateGenesis(cdc, txEncCfg, genesisState); err != nil {
				return errors.Wrap(err, "failed to validate genesis state")
			}

			appState, err := json.MarshalIndent(genesisState, "", " ")
			if err != nil {
				return errors.Wrap(err, "Failed to marshall default genesis state")
			}

			genDoc.AppState = appState

			genDocBytes, err := tmjson.Marshal(genDoc)
			if err != nil {
				return err
			}
			return tmos.WriteFile(args[1], genDocBytes, 0644)
		},
	}

	return cmd
}
