package cli

import (
	"context"
	"fmt"

	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        stakingtypes.ModuleName,
		Short:                      "Querying commands for the staking module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdQueryValidator(),
		GetCmdQueryValidators(),
	)

	return queryCmd
}

// GetCmdQueryValidator the query delegation command.
func GetCmdQueryValidators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validators",
		Short: "Query a validators by specific filter  [--addr || --val-addr || --moniker || --status || --pubkey || --proposer ]",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pubKey, _ := cmd.Flags().GetString(FlagPubKey)
			proposer, _ := cmd.Flags().GetString(FlagProposer)
			status, _ := cmd.Flags().GetString(FlagStatus)
			valAddr, _ := cmd.Flags().GetString(FlagValAddr)
			addr, _ := cmd.Flags().GetString(FlagAddr)
			moniker, _ := cmd.Flags().GetString(FlagMoniker)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			request := &stakingtypes.ValidatorsRequest{
				Address:    addr,
				Valkey:     valAddr,
				Pubkey:     pubKey,
				Moniker:    moniker,
				Status:     status,
				Proposer:   proposer,
				Pagination: pageReq,
			}

			queryClient := stakingtypes.NewQueryClient(clientCtx)
			res, err := queryClient.Validators(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)

		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	cmd.Flags().String(FlagAddr, "", "the addres in AccAddress format.")
	cmd.Flags().String(FlagValAddr, "", "the addres in ValAddress format.")
	cmd.Flags().String(FlagMoniker, "", "the moniker")
	cmd.Flags().String(FlagPubKey, "", "the pubKey")
	cmd.Flags().String(FlagProposer, "", "the proposer")
	cmd.Flags().String(FlagStatus, "", "the status")

	flags.AddPaginationFlagsToCmd(cmd, "customstaking")

	return cmd
}

// GetCmdQueryValidator the query delegation command.
func GetCmdQueryValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator [--addr || --val-addr || --moniker] ",
		Short: "Query a validator based on address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			err := validateQueryValidatorFlags(cmd.Flags())
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

				params := &stakingtypes.ValidatorByAddressRequest{ValAddr: valAddr}

				queryClient := stakingtypes.NewQueryClient(clientCtx)
				res, err := queryClient.ValidatorByAddress(context.Background(), params)
				if err != nil {
					return err
				}

				return clientCtx.PrintProto(&res.Validator)
			}

			moniker, _ := cmd.Flags().GetString(FlagMoniker)
			if moniker != "" {
				params := &stakingtypes.ValidatorByMonikerRequest{Moniker: moniker}

				queryClient := stakingtypes.NewQueryClient(clientCtx)
				res, err := queryClient.ValidatorByMoniker(context.Background(), params)
				if err != nil {
					return err
				}

				return clientCtx.PrintProto(&res.Validator)
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
