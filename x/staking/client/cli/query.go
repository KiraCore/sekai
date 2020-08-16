package cli

import (
	"context"
	"fmt"

	"github.com/spf13/pflag"

	"github.com/spf13/cobra"

	"github.com/KiraCore/cosmos-sdk/client"
	"github.com/KiraCore/cosmos-sdk/client/flags"
	sdk "github.com/KiraCore/cosmos-sdk/types"

	cumstomtypes "github.com/KiraCore/sekai/x/staking/types"
)

const (
	FlagValAddr = "val-addr"
	flagMoniker = "flagMoniker"
	flagAddr    = "addr"
)

// GetCmdQueryValidatorByAddress the query delegation command.
func GetCmdQueryValidatorByAddress() *cobra.Command {
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
			if valAddrStr != "" {
				valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
				if err != nil {
					return err
				}

				params := &cumstomtypes.ValidatorByAddressRequest{ValAddr: valAddr}

				queryClient := cumstomtypes.NewQueryClient(clientCtx)
				res, err := queryClient.ValidatorByAddress(context.Background(), params)
				if err != nil {
					return err
				}

				return clientCtx.PrintOutput(res.Validator)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	cmd.Flags().String(flagAddr, "", "the addres in AccAddress format.")
	cmd.Flags().String(FlagValAddr, "", "the addres in ValAddress format.")
	cmd.Flags().String(flagMoniker, "", "the moniker")

	return cmd
}

// validateQueryValidatorFlags return the validator flags.
func validateQueryValidatorFlags(flagSet *pflag.FlagSet) error {
	moniker, err := flagSet.GetString(FlagMoniker)
	if err != nil {
		return err
	}
	addr, err := flagSet.GetString(flagAddr)
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
