package cli

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/spf13/pflag"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cumstomtypes "github.com/KiraCore/sekai/x/staking/types"
)

const (
	FlagValAddr = "val-addr"
	FlagAddr    = "addr"
)

// GetCmdQueryValidator the query delegation command.
func GetCmdQueryValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator [--addr || --val-addr || --flagMoniker] ",
		Short: "Query a validator based on address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			err = validateQueryValidatorFlags(cmd.Flags())
			if err != nil {
				return err
			}

			valAddrStr, _ := cmd.Flags().GetString(FlagValAddr)
			addr, _ := cmd.Flags().GetString(FlagAddr)
			if valAddrStr != "" || addr != "" {
				var valAddr sdk.ValAddress
				if addr != "" {
					bechAddr, err := sdk.AccAddressFromBech32(addr)
					if err != nil {
						return errors.Wrap(err, "invalid account address")
					}
					valAddr = sdk.ValAddress(bechAddr)
				} else {
					valAddr, err = sdk.ValAddressFromBech32(valAddrStr)
					if err != nil {
						return errors.Wrap(err, "invalid validator address")
					}
				}

				params := &cumstomtypes.ValidatorByAddressRequest{ValAddr: valAddr}

				queryClient := cumstomtypes.NewQueryClient(clientCtx)
				res, err := queryClient.ValidatorByAddress(context.Background(), params)
				if err != nil {
					return err
				}

				return clientCtx.PrintOutput(&res.Validator)
			}

			moniker, _ := cmd.Flags().GetString(FlagMoniker)
			if moniker != "" {
				params := &cumstomtypes.ValidatorByMonikerRequest{Moniker: moniker}

				queryClient := cumstomtypes.NewQueryClient(clientCtx)
				res, err := queryClient.ValidatorByMoniker(context.Background(), params)
				if err != nil {
					return err
				}

				return clientCtx.PrintOutput(&res.Validator)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	cmd.Flags().String(FlagAddr, "", "the addres in AccAddress format.")
	cmd.Flags().String(FlagValAddr, "", "the addres in ValAddress format.")
	cmd.Flags().String(FlagMoniker, "", "the moniker")

	return cmd
}

// validateQueryValidatorFlags return the validator flags.
func validateQueryValidatorFlags(flagSet *pflag.FlagSet) error {
	moniker, err := flagSet.GetString(FlagMoniker)
	if err != nil {
		return err
	}
	addr, err := flagSet.GetString(FlagAddr)
	if err != nil {
		return err
	}
	valAddr, err := flagSet.GetString(FlagValAddr)
	if err != nil {
		return err
	}

	if moniker == "" && addr == "" && valAddr == "" {
		return fmt.Errorf("at least one of flags (--moniker, --val-addr, --addr) needs to be set")
	}

	return nil
}
