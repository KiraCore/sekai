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
	queryCmd.AddCommand(NewCustodiansCmd())
	queryCmd.AddCommand(NewWhiteListCmd())
	queryCmd.AddCommand(NewLimitsCmd())

	return queryCmd
}

func NewCustodiansCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   "custodians",
		Short: "query commands for the custody custodians",
	}

	queryCmd.AddCommand(GetCmdQueryCustodiansByAddress())
	queryCmd.AddCommand(GetCmdQueryCustodyPoolByAddress())

	return queryCmd
}

func NewWhiteListCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   "whitelist",
		Short: "query commands for the custody whitelist",
	}

	queryCmd.AddCommand(GetCmdQueryWhiteListByAddress())

	return queryCmd
}

func NewLimitsCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   "limits",
		Short: "query commands for the custody limits",
	}

	queryCmd.AddCommand(GetCmdQueryLimitsByAddress())

	return queryCmd
}

func GetCmdQueryCustodiansByAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [addr]",
		Short: "Query custody custodians assigned to an address",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			accAddr, err := sdk.AccAddressFromBech32(args[0])

			if err != nil {
				return errors.Wrap(err, "invalid account address")
			}

			params := &types.CustodyCustodiansByAddressRequest{Addr: accAddr}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.CustodyCustodiansByAddress(context.Background(), params)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryCustodyPoolByAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool [addr]",
		Short: "Query custody pool",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			accAddr, err := sdk.AccAddressFromBech32(args[0])

			if err != nil {
				return errors.Wrap(err, "invalid account address")
			}

			params := &types.CustodyPoolByAddressRequest{Addr: accAddr}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.CustodyPoolByAddress(context.Background(), params)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryWhiteListByAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [addr]",
		Short: "Query custody whitelist assigned to an address",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			accAddr, err := sdk.AccAddressFromBech32(args[0])

			if err != nil {
				return errors.Wrap(err, "invalid account address")
			}

			params := &types.CustodyWhiteListByAddressRequest{Addr: accAddr}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.CustodyWhiteListByAddress(context.Background(), params)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryCustodyByAddress is the querier for custody by address.
func GetCmdQueryCustodyByAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [addr]",
		Short: "Query custody assigned to an address",
		Args:  cobra.MinimumNArgs(1),
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

func GetCmdQueryLimitsByAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [addr]",
		Short: "Query custody limits assigned to an address",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			accAddr, err := sdk.AccAddressFromBech32(args[0])

			if err != nil {
				return errors.Wrap(err, "invalid account address")
			}

			params := &types.CustodyLimitsByAddressRequest{Addr: accAddr}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.CustodyLimitsByAddress(context.Background(), params)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
