package cli

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto"

	"github.com/KiraCore/cosmos-sdk/server"

	"github.com/KiraCore/cosmos-sdk/x/genutil"
	"github.com/KiraCore/cosmos-sdk/x/staking/client/cli"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/KiraCore/cosmos-sdk/client"
	"github.com/KiraCore/cosmos-sdk/client/flags"
	"github.com/KiraCore/cosmos-sdk/client/tx"
	"github.com/KiraCore/cosmos-sdk/types"
	cumstomtypes "github.com/KiraCore/sekai/x/staking/types"
)

const (
	FlagMoniker   = "moniker"
	FlagWebsite   = "website"
	FlagSocial    = "social"
	FlagIdentity  = "identity"
	FlagComission = "commission"
	FlagValKey    = "validator-key"
)

func GetTxClaimValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-validator-seat",
		Short: "Claim validator seat to become a Validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)
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

			// read --pubkey, if empty take it from priv_validator.json
			var valPubKey crypto.PubKey
			if valPubKeyString, _ := cmd.Flags().GetString(cli.FlagPubKey); valPubKeyString != "" {
				valPubKey, err = types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, valPubKeyString)
				if err != nil {
					return errors.Wrap(err, "failed to get consensus node public key")
				}
			} else {
				_, valPubKey, err = genutil.InitializeNodeValidatorFiles(serverCtx.Config)
				if err != nil {
					return errors.Wrap(err, "failed to initialize node validator files")
				}
			}

			comm, err := types.NewDecFromStr(comission)
			val, err := types.ValAddressFromBech32(valKeyStr)
			if err != nil {
				return err
			}

			msg, err := cumstomtypes.NewMsgClaimValidator(moniker, website, social, identity, comm, val, valPubKey)
			if err != nil {
				return fmt.Errorf("error creating tx: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	AddValidatorFlags(cmd)
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
