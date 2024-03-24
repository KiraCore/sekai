package cli

import (
	"github.com/KiraCore/sekai/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"strconv"
)

// NewTxCmd returns a root CLI command handler for all x/bank transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "bridge sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(TxChangeCosmosEthereum())
	txCmd.AddCommand(TxChangeEthereumCosmos())

	return txCmd
}

func TxChangeCosmosEthereum() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change_cosmos_ethereum",
		Short: "Create new change request from Cosmos to Ethereum",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			to := args[1]

			inAmount, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return err
			}

			outAmount, _ := strconv.ParseInt(args[3], 10, 64)

			msg := types.NewMsgChangeCosmosEthereum(
				clientCtx.FromAddress,
				to,
				inAmount,
				outAmount,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func TxChangeEthereumCosmos() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change_ethereum_cosmos",
		Short: "Create new change request from Ethereum to Cosmos",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := args[1]

			to, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			inAmount, _ := strconv.ParseInt(args[4], 10, 64)

			outAmount, err := sdk.ParseCoinsNormalized(args[3])
			if err != nil {
				return err
			}

			msg := types.NewMsgChangeEthereumCosmos(
				clientCtx.FromAddress,
				from,
				to,
				inAmount,
				outAmount,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
