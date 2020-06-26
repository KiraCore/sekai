package listOrders

import (
	"github.com/KiraCore/cosmos-sdk/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/kiraHub/transactions/createOrder"
	abci "github.com/tendermint/tendermint/abci/types"
	"strconv"
)

func QueryGetOrders(ctx sdk.Context, path []string, req abci.RequestQuery, keeper createOrder.Keeper) ([]byte, error) {

	var queryOutput []types.LimitOrder

	var int1, _ = strconv.Atoi(path[1])
	var int2, _ = strconv.Atoi(path[2])

	queryOutput = keeper.GetOrders(ctx, path[0], int1, int2)

	res, marshalJSONIndentError := codec.MarshalJSONIndent(packageCodec, queryOutput)
	if marshalJSONIndentError != nil {
		panic(marshalJSONIndentError)
	}

	return res, nil
}