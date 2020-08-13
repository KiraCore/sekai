package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KiraCore/cosmos-sdk/client"
	"github.com/KiraCore/cosmos-sdk/client/flags"
	sdk "github.com/KiraCore/cosmos-sdk/types"

	cumstomtypes "github.com/KiraCore/sekai/x/staking/types"
)

// GetCmdQueryValidatorByAddress the query delegation command.
func GetCmdQueryValidatorByAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-by-address []",
		Short: "Query a validator based on address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := cumstomtypes.NewQueryClient(clientCtx)

			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := &cumstomtypes.ValidatorByAddressRequest{ValAddr: valAddr}

			res, err := queryClient.ValidatorByAddress(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res.Validator)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
