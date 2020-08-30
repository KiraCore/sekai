package cli

import (
	"fmt"

	types2 "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/staking/client/cli"

	"github.com/KiraCore/sekai/x/gov/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

const (
	FlagPermission = "permission"
)

func GetTxSetWhitelistPermissions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-whitelist-permissions",
		Short: "Whitelists permissions into an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			perm, err := cmd.Flags().GetUint32(FlagPermission)
			if err != nil {
				return fmt.Errorf("invalid permissions")
			}

			addr, err := cmd.Flags().GetString(cli.FlagAddr)
			if err != nil {
				return fmt.Errorf("error getting address")
			}

			bech, err := types2.AccAddressFromBech32(addr)
			if err != nil {
				return fmt.Errorf("invalid address")
			}

			msg := types.NewMsgWhitelistPermissions(
				clientCtx.FromAddress,
				bech,
				perm,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagAddr, "", "the address to set permissions")
	cmd.Flags().Uint32(FlagPermission, 0, "the list of permissions")

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func uintToUint32Slice(slice []uint) []uint32 {
	var converted []uint32
	for _, val := range slice {
		converted = append(converted, uint32(val))
	}

	return converted
}
