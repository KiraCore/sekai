package queries

import (
	"strconv"

	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/x/kiraHub/keeper"
	"github.com/KiraCore/sekai/x/kiraHub/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func QueryGetOrders(ctx sdk.Context, path []string, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {

	var queryOutput []types.LimitOrder

	var int1, _ = strconv.Atoi(path[1])
	var int2, _ = strconv.Atoi(path[2])

	queryOutput = keeper.GetOrders(ctx, path[0], uint32(int1), uint32(int2))

	res, err := types.ModuleCdc.MarshalJSON(queryOutput)
	if err != nil {
		panic(err)
	}

	return res, nil
}
