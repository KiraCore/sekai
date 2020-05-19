package listOrderBooks

import (
	"github.com/KiraCore/cosmos-sdk/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/x/kiraHub/transactions/createOrderBook"
	"github.com/KiraCore/sekai/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func QueryGetOrderBooks(ctx sdk.Context, path []string, req abci.RequestQuery, keeper createOrderBook.Keeper) ([]byte, error) {

	var queryOutput []types.OrderBook

	if path[0] == "ID" {

		queryOutput = keeper.GetOrderBookByID(ctx, path[1])

	} else if path[0] == "Index" {

		queryOutput = keeper.GetOrderBookByIndex(ctx, path[1])

	} else if path[0] == "Quote" {

		queryOutput = keeper.GetOrderBookByQuote(ctx, path[1])

	} else if path[0] == "Base" {

		queryOutput = keeper.GetOrderBookByBase(ctx, path[1])

	} else if path[0] == "Trading_Pair" {

		queryOutput = keeper.GetOrderBookByTP(ctx, path[0], path[1])

	} else if path[0] == "Curator" {

		queryOutput = keeper.GetOrderBookByCurator(ctx, path[1])
	}

	res, marshalJSONIndentError := codec.MarshalJSONIndent(packageCodec, queryOutput)
	if marshalJSONIndentError != nil {
		panic(marshalJSONIndentError)
	}

	return res, nil


}