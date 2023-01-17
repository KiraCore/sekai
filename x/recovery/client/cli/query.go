package cli

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/KiraCore/sekai/x/recovery/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group recovery queries under a subcommand
	recoveryQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the recovery module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	recoveryQueryCmd.AddCommand(
		GetCmdQueryRecoveryRecord(),
		GetCmdQueryRecoveryToken(),
	)

	return recoveryQueryCmd

}

// GetCmdQueryRecoveryRecord implements the command to query a recovery record for an address.
func GetCmdQueryRecoveryRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recovery-record [address]",
		Short: "Query an account's recovery information",
		Long: strings.TrimSpace(`Query an account's recovery information:

$ <appd> query recovery recovery-record kira15nxzg5lrmyu42vuzlztdnlhq9sngerenu520ey
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryRecoveryRecordRequest{Address: args[0]}
			res, err := queryClient.RecoveryRecord(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryRecoveryToken implements the command to query a recovery token for an address.
func GetCmdQueryRecoveryToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recovery-token [address]",
		Short: "Query an account's recovery token information",
		Long: strings.TrimSpace(`Query an account's recovery token information:

$ <appd> query recovery recovery-token kira15nxzg5lrmyu42vuzlztdnlhq9sngerenu520ey
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryRecoveryTokenRequest{Address: args[0]}
			res, err := queryClient.RecoveryToken(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
