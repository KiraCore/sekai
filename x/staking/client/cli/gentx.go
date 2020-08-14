package cli

import (
	"encoding/json"
	"fmt"

	"github.com/KiraCore/cosmos-sdk/codec"

	"github.com/KiraCore/cosmos-sdk/client"
	"github.com/KiraCore/cosmos-sdk/client/flags"
	"github.com/KiraCore/cosmos-sdk/server"
	"github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/cosmos-sdk/types/module"
	types2 "github.com/KiraCore/cosmos-sdk/x/bank/types"
	"github.com/KiraCore/cosmos-sdk/x/genutil"
	"github.com/KiraCore/cosmos-sdk/x/staking/client/cli"
	cumstomtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	tmtypes "github.com/tendermint/tendermint/types"
)

func GenTxClaimCmd(mbm module.BasicManager, txEncCfg client.TxEncodingConfig, genBalIterator types2.GenesisBalancesIterator, defaultNodeHome string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "gentx-claim [key_name]",
		Short: "Adds validator into the genesis set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.JSONMarshaler

			config := serverCtx.Config
			config.SetRoot(clientCtx.HomeDir)

			_, valPubKey, err := genutil.InitializeNodeValidatorFiles(serverCtx.Config)
			if err != nil {
				return errors.Wrap(err, "failed to initialize node validator files")
			}

			// read --pubkey, if empty take it from priv_validator.json
			if valPubKeyString, _ := cmd.Flags().GetString(cli.FlagPubKey); valPubKeyString != "" {
				valPubKey, err = types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, valPubKeyString)
				if err != nil {
					return errors.Wrap(err, "failed to get consensus node public key")
				}
			}

			genDoc, err := tmtypes.GenesisDocFromFile(config.GenesisFile())
			if err != nil {
				return errors.Wrapf(err, "failed to read genesis doc file %s", config.GenesisFile())
			}

			var genesisState map[string]json.RawMessage
			if err = cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
				return errors.Wrap(err, "failed to unmarshal genesis state")
			}

			if err = mbm.ValidateGenesis(cdc, txEncCfg, genesisState); err != nil {
				return errors.Wrap(err, "failed to validate genesis state")
			}

			name := args[0]
			key, err := clientCtx.Keyring.Key(name)
			if err != nil {
				return errors.Wrapf(err, "failed to fetch '%s' from the keyring", name)
			}

			moniker := config.Moniker
			if m, _ := cmd.Flags().GetString(cli.FlagMoniker); m != "" {
				moniker = m
			}

			amount, _ := cmd.Flags().GetString(cli.FlagAmount)
			coins, err := types.ParseCoins(amount)
			if err != nil {
				return errors.Wrap(err, "failed to parse coins")
			}

			err = genutil.ValidateAccountInGenesis(genesisState, genBalIterator, key.GetAddress(), coins, cdc)
			if err != nil {
				return errors.Wrap(err, "failed to validate account in genesis")
			}

			website, _ := cmd.Flags().GetString(FlagWebsite)
			identity, _ := cmd.Flags().GetString(FlagIdentity)
			validator, err := cumstomtypes.NewValidator(
				moniker,
				website,
				"social",
				identity,
				types.NewDec(1),
				types.ValAddress(key.GetAddress()),
				valPubKey,
			)
			if err != nil {
				return errors.Wrap(err, "failed to create new validator")
			}

			var stakingGenesisState cumstomtypes.GenesisState
			stakingGenesisState.Validators = append(stakingGenesisState.Validators, validator)

			genesisState[cumstomtypes.ModuleName] = cdc.MustMarshalJSON(stakingGenesisState)
			appGenStateJSON, err := codec.MarshalJSONIndent(clientCtx.JSONMarshaler, genesisState)
			if err != nil {
				return err
			}

			genDoc.AppState = appGenStateJSON

			err = genDoc.ValidateAndComplete()
			if err != nil {
				return err
			}

			err = genDoc.SaveAs(config.GenesisFile())
			if err != nil {
				return err
			}

			fmt.Printf("genesis state updated to include validator")

			return nil
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|kwallet|pass|test)")
	AddValidatorFlags(cmd)

	return cmd
}
