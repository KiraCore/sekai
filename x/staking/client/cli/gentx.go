package cli

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/KiraCore/sekai/x/genutil"
	genutiltypes "github.com/KiraCore/sekai/x/genutil/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func GenTxClaimCmd(genBalIterator banktypes.GenesisBalancesIterator, defaultNodeHome string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "gentx-claim [key_name]",
		Short: "Adds validator into the genesis set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.JSONCodec

			config := serverCtx.Config
			config.SetRoot(clientCtx.HomeDir)

			_, valPubKey, err := genutil.InitializeNodeValidatorFiles(serverCtx.Config)
			if err != nil {
				return errors.Wrap(err, "failed to initialize node validator files")
			}

			if valPubKeyString, _ := cmd.Flags().GetString(cli.FlagPubKey); valPubKeyString != "" {
				if err := clientCtx.Codec.UnmarshalInterfaceJSON([]byte(valPubKeyString), &valPubKey); err != nil {
					return errors.Wrap(err, "failed to get consensus node public key")
				}
			}

			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(config.GenesisFile())

			name := args[0]
			key, err := clientCtx.Keyring.Key(name)
			if err != nil {
				return errors.Wrapf(err, "failed to fetch '%s' from the keyring", name)
			}

			validator, err := stakingtypes.NewValidator(
				types.ValAddress(key.GetAddress()),
				valPubKey,
			)
			if err != nil {
				return errors.Wrap(err, "failed to create new validator")
			}

			var stakingGenesisState stakingtypes.GenesisState
			stakingGenesisState.Validators = append(stakingGenesisState.Validators, validator)
			bzStakingGen := cdc.MustMarshalJSON(&stakingGenesisState)
			appState[stakingtypes.ModuleName] = bzStakingGen

			var customGovGenState govtypes.GenesisState
			cdc.MustUnmarshalJSON(appState[govtypes.ModuleName], &customGovGenState)

			// Only first validator is network actor
			networkActor := govtypes.NewNetworkActor(
				types.AccAddress(validator.ValKey),
				govtypes.Roles{
					uint64(govtypes.RoleSudo),
				},
				govtypes.Active,
				[]govtypes.VoteOption{
					govtypes.OptionYes,
					govtypes.OptionAbstain,
					govtypes.OptionNo,
					govtypes.OptionNoWithVeto,
				},
				govtypes.NewPermissions(nil, nil),
				1,
			)
			customGovGenState.NetworkActors = append(customGovGenState.NetworkActors, &networkActor)
			moniker := config.Moniker
			if m, _ := cmd.Flags().GetString(cli.FlagMoniker); m != "" {
				moniker = m
			}
			for _, record := range customGovGenState.IdentityRecords {
				if record.Key == "moniker" && record.Value == moniker {
					panic(fmt.Sprintf("same moniker exists, moniker = %s", moniker))
				}
			}
			customGovGenState.IdentityRecords = append(customGovGenState.IdentityRecords, govtypes.IdentityRecord{
				Id:        customGovGenState.LastIdentityRecordId + 1,
				Address:   types.AccAddress(validator.ValKey).String(),
				Key:       "moniker",
				Value:     moniker,
				Date:      time.Now().UTC(),
				Verifiers: []string{},
			})
			customGovGenState.LastIdentityRecordId++
			appState[govtypes.ModuleName] = cdc.MustMarshalJSON(&customGovGenState)

			appGenStateJSON, err := json.Marshal(appState)
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

			fmt.Printf("genesis state updated to include validator\n")

			return nil
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|kwallet|pass|test)")
	AddValidatorFlags(cmd)

	return cmd
}
