package cli

import (
	"context"
	"github.com/KiraCore/sekai/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// NewQueryCmd returns a root CLI command handler for all x/distributor transaction commands.
func NewQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   types.RouterKey,
		Short: "query commands for the bridge module",
	}

	queryCmd.AddCommand(GetCmdQueryChangeCosmosEthereumByAddress())
	queryCmd.AddCommand(GetCmdQueryChangeEthereumCosmosByAddress())

	return queryCmd
}

// GetCmdQueryChangeCosmosEthereumByAddress is the querier for change by address.
func GetCmdQueryChangeCosmosEthereumByAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get_cosmos_ethereum [addr]",
		Short: "Query change from Cosmos to Ethereum assigned to an address",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			accAddr, err := sdk.AccAddressFromBech32(args[0])

			if err != nil {
				return errors.Wrap(err, "invalid account address")
			}

			params := &types.ChangeCosmosEthereumByAddressRequest{Addr: accAddr}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ChangeCosmosEthereumByAddress(context.Background(), params)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryChangeEthereumCosmosByAddress is the querier for change by address.
func GetCmdQueryChangeEthereumCosmosByAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get_ethereum_cosmos [addr]",
		Short: "Query change from Ethereum to Cosmos assigned to an address",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			accAddr, err := sdk.AccAddressFromBech32(args[0])

			if err != nil {
				return errors.Wrap(err, "invalid account address")
			}

			params := &types.ChangeEthereumCosmosByAddressRequest{Addr: accAddr}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ChangeEthereumCosmosByAddress(context.Background(), params)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
