package signerkeys

import (
	"fmt"

	"github.com/KiraCore/cosmos-sdk/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/types"
	signerkey "github.com/KiraCore/sekai/x/kiraHub/transactions/upsertSignerKey"
	abci "github.com/tendermint/tendermint/abci/types"
)

func QueryListSignerKeys(ctx sdk.Context, path []string, req abci.RequestQuery, keeper signerkey.Keeper) ([]byte, error) {

	var queryOutput []types.SignerKey
	curator, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return []byte{}, fmt.Errorf("Invalid curator address %s: %+v", path[0], err)
	}

	queryOutput = keeper.GetSignerKeys(ctx, curator)

	res, marshalJSONIndentError := codec.MarshalJSONIndent(keeper.GetCodec(), queryOutput)
	if marshalJSONIndentError != nil {
		panic(marshalJSONIndentError)
	}

	return res, nil
}
