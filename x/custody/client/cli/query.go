package cli

import (
	"context"
	"github.com/KiraCore/sekai/x/custody/types"
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
		Short: "query commands for the custody module",
	}

	queryCmd.AddCommand(GetCmdQueryCustodyByAddress())

	return queryCmd
}

// GetCmdQueryCustodyByAddress is the querier for custody by address.
func GetCmdQueryCustodyByAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [addr]",
		Short: "Query custody assigned to an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			accAddr, err := sdk.AccAddressFromBech32(args[0])

			if err != nil {
				return errors.Wrap(err, "invalid account address")
			}

			params := &types.CustodyByAddressRequest{Addr: accAddr}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.CustodyByAddress(context.Background(), params)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
