package signerkeys

import (
	"fmt"

	"github.com/KiraCore/cosmos-sdk/client/context"
	"github.com/KiraCore/cosmos-sdk/codec"
	"github.com/KiraCore/sekai/types"

	"github.com/spf13/cobra"
)

func ListSignerKeysCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "listsignerkeys",
		Short: "List signer key(s) by curator address",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var owner = cliCtx.GetFromAddress()

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/kiraHub/listSignerKeys/%s", owner.String()), nil)
			if err != nil {
				fmt.Printf("could not query. Searching By - %s \n", owner.String())
				return nil
			}

			var out []types.SignerKey
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
