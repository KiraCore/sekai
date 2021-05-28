package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// GetQueryCmd returns the parent command for all x/upgrade CLi query commands.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the upgrade module",
	}

	cmd.AddCommand()

	return cmd
}
