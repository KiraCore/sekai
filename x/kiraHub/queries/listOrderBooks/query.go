package listOrderBooks

import (
	"github.com/KiraCore/cosmos-sdk/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	sdkerrors "github.com/KiraCore/cosmos-sdk/types/errors"
	"github.com/KiraCore/cosmos-sdk/x/mint"
	abci "github.com/tendermint/tendermint/abci/types"
)

func QueryListOrderBooks(ctx sdk.Context, path []string, req abci.RequestQuery, keeper mint.Keeper) ([]byte, error) {
	var owner, err = sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, err
	}

	coin := keeper.GetCoin(ctx, owner)

	if coin.Owner == nil {
		return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "could not get coin - unknown owner")
	}

	res, marshalJSONIndentError := codec.MarshalJSONIndent(packageCodec, coin)
	if marshalJSONIndentError != nil {
		panic(marshalJSONIndentError)
	}

	return res, nil
}
