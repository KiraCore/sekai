package kiraHub

import (
	"context"

	"github.com/KiraCore/sekai/x/kiraHub/keeper"
	"github.com/KiraCore/sekai/x/kiraHub/types"
)

type Querier struct {
	keeper keeper.Keeper
}

func NewQuerier(keeper keeper.Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}

func (q Querier) GetOrderBooks(ctx context.Context, request *types.HubRequest) (*types.HubResponse, error) {
	// c := sdk.UnwrapSDKContext(ctx)
	return &types.HubResponse{}, nil

	// var queryOutput []types.OrderBook

	// if path[0] == "ID" {

	// 	queryOutput = keeper.GetOrderBookByID(ctx, path[1])

	// } else if path[0] == "Index" {

	// 	var int1, _ = strconv.Atoi(path[1])
	// 	queryOutput = keeper.GetOrderBookByIndex(ctx, uint32(int1))

	// } else if path[0] == "Quote" {

	// 	queryOutput = keeper.GetOrderBookByQuote(ctx, path[1])

	// } else if path[0] == "Base" {

	// 	queryOutput = keeper.GetOrderBookByBase(ctx, path[1])

	// } else if path[0] == "tp" {

	// 	queryOutput = keeper.GetOrderBookByTP(ctx, path[1], path[2])

	// } else if path[0] == "Curator" {

	// 	queryOutput = keeper.GetOrderBookByCurator(ctx, path[1])
	// }

	// res, err := types.ModuleCdc.MarshalJSON(queryOutput)
	// if err != nil {
	// 	panic(err)
	// }

	// return res, nil
}

func (q Querier) GetOrderBooksByTP(ctx context.Context, request *types.HubRequest) (*types.HubResponse, error) {
	// c := sdk.UnwrapSDKContext(ctx)
	return &types.HubResponse{}, nil
}

func (q Querier) GetOrders(ctx context.Context, request *types.HubRequest) (*types.HubResponse, error) {
	// c := sdk.UnwrapSDKContext(ctx)
	return &types.HubResponse{}, nil

	// var queryOutput []types.LimitOrder

	// var int1, _ = strconv.Atoi(path[1])
	// var int2, _ = strconv.Atoi(path[2])

	// queryOutput = keeper.GetOrders(ctx, path[0], uint32(int1), uint32(int2))

	// res, err := types.ModuleCdc.MarshalJSON(queryOutput)
	// if err != nil {
	// 	panic(err)
	// }

	// return res, nil
}

func (q Querier) ListSignerKeys(ctx context.Context, request *types.HubRequest) (*types.HubResponse, error) {
	// c := sdk.UnwrapSDKContext(ctx)
	return &types.HubResponse{}, nil

	// var queryOutput []types.SignerKey
	// curator, err := sdk.AccAddressFromBech32(path[0])
	// if err != nil {
	// 	return []byte{}, fmt.Errorf("Invalid curator address %s: %+v", path[0], err)
	// }

	// queryOutput = keeper.GetSignerKeys(ctx, curator)

	// res, err := types.ModuleCdc.MarshalJSON(queryOutput)
	// if err != nil {
	// 	panic(err)
	// }

	// return res, nil
}
