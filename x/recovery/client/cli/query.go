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
		GetCmdQueryRRHolderRewards(),
		GetCmdQueryRRHolders(),
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

// GetCmdQueryRRHolderRewards implements the command to query rr holder rewards
func GetCmdQueryRRHolderRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rr-holder-rewards [address]",
		Short: "Query an account's rr holder rewards information",
		Long: strings.TrimSpace(`Query an account's rr holder rewards information:

$ <appd> query recovery rr-holder-rewards kira15nxzg5lrmyu42vuzlztdnlhq9sngerenu520ey
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryRRHolderRewardsRequest{Address: args[0]}
			res, err := queryClient.RRHolderRewards(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryRRHolders implements the command to query rr holders
func GetCmdQueryRRHolders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rr-holders [rr_token]",
		Short: "Query registered rr holders",
		Long: strings.TrimSpace(`Query registered rr holders:

$ <appd> query recovery rr-holders rr/val1
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryRegisteredRRTokenHoldersRequest{RecoveryToken: args[0]}
			res, err := queryClient.RegisteredRRTokenHolders(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
