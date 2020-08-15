package cli

import (
	"github.com/spf13/cobra"

	"github.com/KiraCore/cosmos-sdk/x/staking/client/cli"
)

// AddValidatorFlags adds the flags needed to create a validator.
func AddValidatorFlags(cmd *cobra.Command) {
	cmd.Flags().String(FlagMoniker, "", "the Moniker")
	cmd.Flags().String(FlagWebsite, "", "the Website")
	cmd.Flags().String(FlagSocial, "", "the social")
	cmd.Flags().String(FlagIdentity, "", "the Identity")
	cmd.Flags().String(FlagComission, "", "the commission")
	cmd.Flags().String(FlagValKey, "", "the validator key")
	cmd.Flags().String(cli.FlagPubKey, "", "the public key")
}
