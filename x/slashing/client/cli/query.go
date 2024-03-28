package cli

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/KiraCore/sekai/x/slashing/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group slashing queries under a subcommand
	slashingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the slashing module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	slashingQueryCmd.AddCommand(
		GetCmdQuerySigningInfo(),
		GetCmdQuerySigningInfos(),
		GetCmdQuerySlashProposals(),
		GetCmdQuerySlashedStakingPools(),
		GetCmdQueryActiveStakingPools(),
		GetCmdQueryInactiveStakingPools(),
	)

	return slashingQueryCmd

}

// GetCmdQuerySigningInfo implements the command to query signing info.
func GetCmdQuerySigningInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signing-info [valconsaddr]",
		Short: "Query a validator's signing information",
		Long: strings.TrimSpace(`Use a validators' consensus public key to find the signing-info for that validator:

$ <appd> query customslashing signing-info kiravalcons15nxzg5lrmyu42vuzlztdnlhq9sngerenu520ey
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QuerySigningInfoRequest{ConsAddress: args[0]}
			res, err := queryClient.SigningInfo(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.ValSigningInfo)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQuerySigningInfos implements the command to query signing infos.
func GetCmdQuerySigningInfos() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signing-infos",
		Short: "Query signing information of all validators",
		Long: strings.TrimSpace(`signing infos of validators:

$ <appd> query slashing signing-infos
`),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QuerySigningInfosRequest{Pagination: pageReq}
			res, err := queryClient.SigningInfos(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "signing infos")

	return cmd
}

// GetCmdQuerySlashProposals implements a command to fetch slash proposals.
func GetCmdQuerySlashProposals() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "slash-proposals",
		Short: "Query slash proposals",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(`Query the slash proposals:

$ <appd> query slashing slash-proposals
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QuerySlashProposalsRequest{}
			res, err := queryClient.SlashProposals(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQuerySlashedStakingPools implements a command to fetch slashed staking pools.
func GetCmdQuerySlashedStakingPools() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "slashed-staking-pools",
		Short: "Query slashed staking pools",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(`Query slashed staking pools:

$ <appd> query slashing slashed-staking-pools
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QuerySlashedStakingPoolsRequest{}
			res, err := queryClient.SlashedStakingPools(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryActiveStakingPools implements a command to fetch active staking pools.
func GetCmdQueryActiveStakingPools() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "active-staking-pools",
		Short: "Query active staking pools",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(`Query active staking pools:

$ <appd> query slashing active-staking-pools
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryActiveStakingPoolsRequest{}
			res, err := queryClient.ActiveStakingPools(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryInactiveStakingPools implements a command to fetch inactive staking pools.
func GetCmdQueryInactiveStakingPools() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inactive-staking-pools",
		Short: "Query inactive staking pools",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(`Query inactive staking pools:

$ <appd> query slashing inactive-staking-pools
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryInactiveStakingPoolsRequest{}
			res, err := queryClient.InactiveStakingPools(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
