package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/x/staking/client/cli"
)

const (
	FlagMoniker  = "moniker"
	FlagValAddr  = "val-addr"
	FlagAddr     = "addr"
	FlagPubKey   = "pubkey"
	FlagProposer = "proposer"
	FlagStatus   = "status"
)

// AddValidatorFlags adds the flags needed to create a validator.
func AddValidatorFlags(cmd *cobra.Command) {
	cmd.Flags().String(FlagMoniker, "", "the Moniker")
	cmd.Flags().String(FlagValKey, "", "the validator key")
	cmd.Flags().String(cli.FlagPubKey, "", "the public key")
}
