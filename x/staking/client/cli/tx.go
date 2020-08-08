package cli

import (
	"fmt"

	"github.com/KiraCore/cosmos-sdk/client/flags"

	"github.com/KiraCore/cosmos-sdk/client"
	"github.com/KiraCore/cosmos-sdk/client/tx"
	"github.com/KiraCore/cosmos-sdk/types"
	cumstomtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/spf13/cobra"
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
			pub, err := types.AccAddressFromBech32(pubKeyStr)
			if err != nil {
				return err
			}

			msg, err := cumstomtypes.NewMsgClaimValidator(moniker, website, social, identity, comm, val, pub)
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
