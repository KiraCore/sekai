package cli

import (
	"github.com/KiraCore/sekai/x/ethereum/types"
	"github.com/spf13/cobra"
)

// NewQueryCmd returns a root CLI command handler for all x/ethereum transaction commands.
func NewQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   types.RouterKey,
		Short: "query commands for the ethereum module",
	}

	return queryCmd
}
