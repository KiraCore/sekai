package cli

import (
	"github.com/KiraCore/cosmos-sdk/client/tx"
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

			return tx.GenerateOrBroadcastTx(cliCtx, txf, msg)
		},
	}

	cmd.Flags().String(flagMoniker, "", "")
	cmd.Flags().String(flagWebsite, "", "")
	cmd.Flags().String(flagSocial, "", "")
	cmd.Flags().String(flagIdentity, "", "")
	cmd.Flags().String(flagComission, "", "")
	cmd.Flags().String(flagValKey, "", "")
	cmd.Flags().String(flagPubKey, "", "")

	return cmd
}
