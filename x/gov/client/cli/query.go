package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/gov/types"
)

const (
	FlagRole = "role"
)

// GetCmdQueryPermissions the query delegation command.
func GetCmdQueryPermissions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "permissions addr",
		Short: "Get the permissions of an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			accAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return errors.Wrap(err, "invalid account address")
			}

			params := &types.PermissionsByAddressRequest{ValAddr: accAddr}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.PermissionsByAddress(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res.Permissions)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryRolesByAddress the query delegation command.
func GetCmdQueryRolesByAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "roles addr",
		Short: "Get the roles assigned to an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			accAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return errors.Wrap(err, "invalid account address")
			}

			params := &types.RolesByAddressRequest{ValAddr: accAddr}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.RolesByAddress(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryRolePermissions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role-permissions arg-num",
		Short: "Get the permissions of all the roles",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			roleNum, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid role number")
			}

			params := &types.RolePermissionsRequest{
				Role: roleNum,
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.RolePermissions(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res.Permissions)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryNetworkProperties implement query network properties
func GetCmdQueryNetworkProperties() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network-properties",
		Short: "Get the network properties",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			params := &types.NetworkPropertiesRequest{}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetNetworkProperties(context.Background(), params)
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryExecutionFee query for execution fee by execution name
func GetCmdQueryExecutionFee() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execution-fee",
		Short: "Get the execution fee by [transaction_type]",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			params := &types.ExecutionFeeRequest{
				TransactionType: args[0],
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetExecutionFee(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdQueryCouncilRegistry() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "council-registry [--addr || --flagMoniker]",
		Short: "Query the governance registry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			addr, err := cmd.Flags().GetString(FlagAddress)
			if err != nil {
				return err
			}

			moniker, err := cmd.Flags().GetString(FlagMoniker)
			if err != nil {
				return err
			}
			if addr == "" && moniker == "" {
				return fmt.Errorf("at least one flag (--flag or --moniker) is mandatory")
			}

			var res *types.CouncilorResponse
			if moniker != "" {
				params := &types.CouncilorByMonikerRequest{Moniker: moniker}

				queryClient := types.NewQueryClient(clientCtx)
				res, err = queryClient.CouncilorByMoniker(context.Background(), params)
				if err != nil {
					return err
				}
			} else {
				bech32, err := sdk.AccAddressFromBech32(addr)
				if err != nil {
					return fmt.Errorf("invalid address: %w", err)
				}

				params := &types.CouncilorByAddressRequest{ValAddr: bech32}

				queryClient := types.NewQueryClient(clientCtx)
				res, err = queryClient.CouncilorByAddress(context.Background(), params)
				if err != nil {
					return err
				}
			}

			return clientCtx.PrintOutput(&res.Councilor)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	cmd.Flags().String(FlagAddress, "", "the address you want to query information")
	cmd.Flags().String(FlagMoniker, "", "the moniker you want to query information")

	return cmd
}
