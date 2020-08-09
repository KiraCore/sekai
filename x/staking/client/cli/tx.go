package cli

import (
	"bufio"
	"encoding/json"
	"fmt"

	"github.com/KiraCore/cosmos-sdk/crypto/keyring"
	authclient "github.com/KiraCore/cosmos-sdk/x/auth/client"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

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
)

const (
	FlagMoniker   = "moniker"
	FlagWebsite   = "website"
	FlagSocial    = "social"
	FlagIdentity  = "identity"
	FlagComission = "commission"
	FlagValKey    = "validator-key"
	FlagPubKey    = "public-key"
)

func GetTxClaimValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-validator-seat",
		Short: "Claim validator seat to become a Validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			moniker, _ := cmd.Flags().GetString(FlagMoniker)
			website, _ := cmd.Flags().GetString(FlagWebsite)
			social, _ := cmd.Flags().GetString(FlagSocial)
			identity, _ := cmd.Flags().GetString(FlagIdentity)
			comission, _ := cmd.Flags().GetString(FlagComission)
			valKeyStr, _ := cmd.Flags().GetString(FlagValKey)
			pubKeyStr, _ := cmd.Flags().GetString(FlagPubKey)

			comm, err := types.NewDecFromStr(comission)
			val, err := types.ValAddressFromBech32(valKeyStr)
			if err != nil {
				return err
			}

			valPubKey, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, pubKeyStr)
			if err != nil {
				return errors.Wrap(err, "failed to get validator public key")
			}

			msg, err := cumstomtypes.NewMsgClaimValidator(moniker, website, social, identity, comm, val, valPubKey)
			if err != nil {
				return fmt.Errorf("error creating tx: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagMoniker, "", "the Moniker")
	cmd.Flags().String(FlagWebsite, "", "the Website")
	cmd.Flags().String(FlagSocial, "", "the social")
	cmd.Flags().String(FlagIdentity, "", "the Identity")
	cmd.Flags().String(FlagComission, "", "the commission")
	cmd.Flags().String(FlagValKey, "", "the validator key")
	cmd.Flags().String(FlagPubKey, "", "the public key")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func GenTxClaimCmd(mbm module.BasicManager, txEncCfg client.TxEncodingConfig, genBalIterator types2.GenesisBalancesIterator, defaultNodeHome string) *cobra.Command {
	ipDefault, _ := server.ExternalIP()
	fsCreateValidator, defaultsDesc := cli.CreateValidatorMsgFlagSet(ipDefault)

	cmd := &cobra.Command{
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

			nodeID, valPubKey, err := genutil.InitializeNodeValidatorFiles(serverCtx.Config)
			if err != nil {
				return errors.Wrap(err, "failed to initialize node validator files")
			}

			// read --nodeID, if empty take it from priv_validator.json
			if nodeIDString, _ := cmd.Flags().GetString(cli.FlagNodeID); nodeIDString != "" {
				nodeID = nodeIDString
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

			inBuf := bufio.NewReader(cmd.InOrStdin())

			name := args[0]
			key, err := clientCtx.Keyring.Key(name)
			if err != nil {
				return errors.Wrapf(err, "failed to fetch '%s' from the keyring", name)
			}

			moniker := config.Moniker
			if m, _ := cmd.Flags().GetString(cli.FlagMoniker); m != "" {
				moniker = m
			}

			// set flags for creating a gentx
			createValCfg, err := cli.PrepareConfigForTxCreateValidator(cmd.Flags(), moniker, nodeID, genDoc.ChainID, valPubKey)
			if err != nil {
				return errors.Wrap(err, "error creating configuration to create validator msg")
			}

			amount, _ := cmd.Flags().GetString(cli.FlagAmount)
			coins, err := types.ParseCoins(amount)
			if err != nil {
				return errors.Wrap(err, "failed to parse coins")
			}
			//
			err = genutil.ValidateAccountInGenesis(genesisState, genBalIterator, key.GetAddress(), coins, cdc)
			if err != nil {
				return errors.Wrap(err, "failed to validate account in genesis")
			}

			txFactory := tx.NewFactoryCLI(clientCtx, cmd.Flags())

			clientCtx = clientCtx.WithInput(inBuf).WithFromAddress(key.GetAddress())

			txBldr, msg, err := BuildClaimValidatorMsg(clientCtx, createValCfg, txFactory, true)
			if err != nil {
				return errors.Wrap(err, "failed to build create-validator message")
			}

			if key.GetType() == keyring.TypeOffline || key.GetType() == keyring.TypeMulti {
				cmd.PrintErrln("Offline key passed in. Use `tx sign` command to sign.")
				return authclient.PrintUnsignedStdTx(txBldr, clientCtx, []types.Msg{msg})
			}

			//// write the unsigned transaction to the buffer
			//w := bytes.NewBuffer([]byte{})
			//clientCtx = clientCtx.WithOutput(w)
			//
			//if err = authclient.PrintUnsignedStdTx(txBldr, clientCtx, []sdk.Msg{msg}); err != nil {
			//	return errors.Wrap(err, "failed to print unsigned std tx")
			//}
			//
			//// read the transaction
			//stdTx, err := readUnsignedGenTxFile(clientCtx, w)
			//if err != nil {
			//	return errors.Wrap(err, "failed to read unsigned gen tx file")
			//}
			//
			//// sign the transaction and write it to the output file
			//txBuilder, err := clientCtx.TxConfig.WrapTxBuilder(stdTx)
			//if err != nil {
			//	return fmt.Errorf("error creating tx builder: %w", err)
			//}
			//
			//err = authclient.SignTx(txFactory, clientCtx, name, txBuilder, true)
			//if err != nil {
			//	return errors.Wrap(err, "failed to sign std tx")
			//}
			//
			//outputDocument, _ := cmd.Flags().GetString(flags.FlagOutputDocument)
			//if outputDocument == "" {
			//	outputDocument, err = makeOutputFilepath(config.RootDir, nodeID)
			//	if err != nil {
			//		return errors.Wrap(err, "failed to create output file path")
			//	}
			//}
			//
			//if err := writeSignedGenTx(clientCtx, outputDocument, stdTx); err != nil {
			//	return errors.Wrap(err, "failed to write signed gen tx")
			//}
			//
			//cmd.PrintErrf("Genesis transaction written to %q\n", outputDocument)
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

//func makeOutputFilepath(rootDir, nodeID string) (string, error) {
//	writePath := filepath.Join(rootDir, "config", "gentx")
//	if err := tmos.EnsureDir(writePath, 0700); err != nil {
//		return "", err
//	}
//
//	return filepath.Join(writePath, fmt.Sprintf("gentx-%v.json", nodeID)), nil
//}
//
//func readUnsignedGenTxFile(clientCtx client.Context, r io.Reader) (sdk.Tx, error) {
//	bz, err := ioutil.ReadAll(r)
//	if err != nil {
//		return nil, err
//	}
//
//	aTx, err := clientCtx.TxConfig.TxJSONDecoder()(bz)
//	if err != nil {
//		return nil, err
//	}
//
//	return aTx, err
//}
//
//func writeSignedGenTx(clientCtx client.Context, outputDocument string, tx sdk.Tx) error {
//	outputFile, err := os.OpenFile(outputDocument, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
//	if err != nil {
//		return err
//	}
//	defer outputFile.Close()
//
//	json, err := clientCtx.TxConfig.TxJSONEncoder()(tx)
//	if err != nil {
//		return err
//	}
//
//	_, err = fmt.Fprintf(outputFile, "%s\n", json)
//
//	return err
//}

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
