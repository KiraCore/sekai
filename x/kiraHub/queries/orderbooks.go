package queries

import (
	"strconv"

	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/x/kiraHub/keeper"
	"github.com/KiraCore/sekai/x/kiraHub/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func QueryGetOrderBooks(ctx sdk.Context, path []string, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {

	var queryOutput []types.OrderBook

	if path[0] == "ID" {

		queryOutput = keeper.GetOrderBookByID(ctx, path[1])

	} else if path[0] == "Index" {

		var int1, _ = strconv.Atoi(path[1])
		queryOutput = keeper.GetOrderBookByIndex(ctx, uint32(int1))

	} else if path[0] == "Quote" {

		queryOutput = keeper.GetOrderBookByQuote(ctx, path[1])

	} else if path[0] == "Base" {

		queryOutput = keeper.GetOrderBookByBase(ctx, path[1])

	} else if path[0] == "tp" {

		queryOutput = keeper.GetOrderBookByTP(ctx, path[1], path[2])

	} else if path[0] == "Curator" {

		queryOutput = keeper.GetOrderBookByCurator(ctx, path[1])
	}

	res, err := types.ModuleCdc.MarshalJSON(queryOutput)
	if err != nil {
		panic(err)
	}

	return res, nil
}
