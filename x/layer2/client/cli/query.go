package cli

import (
	"context"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/KiraCore/sekai/x/basket/types"
)

// NewQueryCmd returns a root CLI command handler for all x/basket transaction commands.
func NewQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   types.RouterKey,
		Short: "query commands for the basket module",
	}
	queryCmd.AddCommand(
		GetCmdQueryTokenBasketById(),
	)

	return queryCmd
}

func GetCmdQueryTokenBasketById() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-basket-by-id [id]",
		Short: "Queries a single basket by id",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.TokenBasketById(context.Background(), &types.QueryTokenBasketByIdRequest{
				Id: uint64(id),
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
