package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/KiraCore/cosmos-sdk/codec"

	"github.com/KiraCore/cosmos-sdk/client"
	"github.com/KiraCore/cosmos-sdk/client/flags"
	"github.com/KiraCore/cosmos-sdk/client/tx"
	"github.com/KiraCore/cosmos-sdk/server"
	"github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/cosmos-sdk/types/module"
	"github.com/KiraCore/cosmos-sdk/version"
	types2 "github.com/KiraCore/cosmos-sdk/x/bank/types"
	"github.com/KiraCore/cosmos-sdk/x/genutil"
	"github.com/KiraCore/cosmos-sdk/x/staking/client/cli"
	cumstomtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmtypes "github.com/tendermint/tendermint/types"
)

func GenTxClaimCmd(mbm module.BasicManager, txEncCfg client.TxEncodingConfig, genBalIterator types2.GenesisBalancesIterator, defaultNodeHome string) *cobra.Command {
	ipDefault, _ := server.ExternalIP()
	fsCreateValidator, defaultsDesc := cli.CreateValidatorMsgFlagSet(ipDefault)

	var cmd = &cobra.Command{
		Use:   "gentx-claim [key_name]",
		Short: "Generate a genesis tx to claim a validator seat",
		Args:  cobra.ExactArgs(1),
		Long: fmt.Sprintf(`Generate a genesis transaction that creates a validator with a self-delegation,
that is signed by the key in the Keyring referenced by a given name. A node ID and Bech32 consensus
pubkey may optionally be provided. If they are omitted, they will be retrieved from the priv_validator.json
file. The following default parameters are included: 
    %s
				
Example:
$ %s gentx my-key-name --home=/path/to/home/dir --keyring-backend=os --chain-id=test-chain-1 \
    --amount=1000000stake \
    --moniker="myvalidator" \
    --commission-max-change-rate=0.01 \
    --commission-max-rate=1.0 \
    --commission-rate=0.07 \
    --details="..." \
    --security-contact="..." \
    --website="..."
`, defaultsDesc, version.AppName,
		),
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
			err = genDoc.SaveAs(config.GenesisFile())
			if err != nil {
				return err
			}

			fmt.Printf("genesis state updated to include validator")

			return nil
		},
	}
	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.Flags().String(flags.FlagOutputDocument, "", "Write the genesis transaction JSON document to the given file instead of the default location")
	cmd.Flags().String(flags.FlagChainID, "", "The network chain ID")
	cmd.Flags().AddFlagSet(fsCreateValidator)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func makeOutputFilepath(rootDir, nodeID string) (string, error) {
	writePath := filepath.Join(rootDir, "config", "gentx")
	if err := tmos.EnsureDir(writePath, 0700); err != nil {
		return "", err
	}

	return filepath.Join(writePath, fmt.Sprintf("gentx-%v.json", nodeID)), nil
}

func readUnsignedGenTxFile(clientCtx client.Context, r io.Reader) (types.Tx, error) {
	bz, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	aTx, err := clientCtx.TxConfig.TxJSONDecoder()(bz)
	if err != nil {
		return nil, err
	}

	return aTx, err
}

func writeSignedGenTx(clientCtx client.Context, outputDocument string, tx types.Tx) error {
	outputFile, err := os.OpenFile(outputDocument, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	json, err := clientCtx.TxConfig.TxJSONEncoder()(tx)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(outputFile, "%s\n", json)

	return err
}

func BuildClaimValidatorMsg(clientCtx client.Context, config cli.TxCreateValidatorConfig, txBldr tx.Factory, generateOnly bool) (tx.Factory, types.Msg, error) {
	valAddr := clientCtx.GetFromAddress()
	pkStr := config.PubKey

	pk, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, pkStr)
	if err != nil {
		return txBldr, nil, err
	}

	commission, err := types.NewDecFromStr(config.CommissionRate)

	msg, err := cumstomtypes.NewMsgClaimValidator(
		config.Moniker,
		config.Website,
		"The SOCIAL, change please",
		config.Identity,
		commission,
		types.ValAddress(valAddr),
		pk,
	)

	if generateOnly {
		ip := config.IP
		nodeID := config.NodeID

		if nodeID != "" && ip != "" {
			txBldr = txBldr.WithMemo(fmt.Sprintf("%s@%s:26656", nodeID, ip))
		}
	}

	return txBldr, msg, nil
}
