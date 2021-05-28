package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Upgrade transaction subcommands",
	}

	return cmd
}
