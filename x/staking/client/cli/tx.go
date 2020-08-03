package cli

import (
	"github.com/KiraCore/cosmos-sdk/client"
	"github.com/KiraCore/cosmos-sdk/client/tx"
	"github.com/KiraCore/cosmos-sdk/types"
	types2 "github.com/KiraCore/sekai/x/staking/types"
	"github.com/spf13/cobra"
)

const (
	flagMoniker   = "moniker"
	flagWebsite   = "website"
	flagSocial    = "social"
	flagIdentity  = "identity"
	flagComission = "commission"
	flagValKey    = "validator-key"
	flagPubKey    = "public-key"
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

			moniker, _ := cmd.Flags().GetString(flagMoniker)
			website, _ := cmd.Flags().GetString(flagWebsite)
			social, _ := cmd.Flags().GetString(flagSocial)
			identity, _ := cmd.Flags().GetString(flagIdentity)
			comission, _ := cmd.Flags().GetString(flagComission)
			valKeyStr, _ := cmd.Flags().GetString(flagValKey)
			pubKeyStr, _ := cmd.Flags().GetString(flagPubKey)

			comm, err := types.NewDecFromStr(comission)
			val, err := types.ValAddressFromBech32(valKeyStr)
			if err != nil {
				return err
			}
			pub, err := types.AccAddressFromBech32(pubKeyStr)
			if err != nil {
				return err
			}

			msg, err := types2.NewMsgClaimValidator(moniker, website, social, identity, comm, val, pub)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(flagMoniker, "", "the Moniker")
	cmd.Flags().String(flagWebsite, "", "the Website")
	cmd.Flags().String(flagSocial, "", "the social")
	cmd.Flags().String(flagIdentity, "", "the Identity")
	cmd.Flags().String(flagComission, "", "the commission")
	cmd.Flags().String(flagValKey, "", "the validator key")
	cmd.Flags().String(flagPubKey, "", "the public key")

	return cmd
}
