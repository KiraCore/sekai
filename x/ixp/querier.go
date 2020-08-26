package ixp

import (
	"context"
	"strconv"

	"github.com/KiraCore/sekai/x/ixp/keeper"
	"github.com/KiraCore/sekai/x/ixp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Querier struct {
	keeper keeper.Keeper
}

func NewQuerier(keeper keeper.Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}

func (q Querier) GetOrderBooks(ctx context.Context, request *types.GetOrderBooksRequest) (*types.GetOrderBooksResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)

	var queryOutput []types.OrderBook

	switch request.QueryType {
	case "ID":
		queryOutput = q.keeper.GetOrderBookByID(c, request.QueryValue)
	case "Index":
		var int1, _ = strconv.Atoi(request.QueryValue)
		queryOutput = q.keeper.GetOrderBookByIndex(c, uint32(int1))
	case "Quote":
		queryOutput = q.keeper.GetOrderBookByQuote(c, request.QueryValue)
	case "Base":
		queryOutput = q.keeper.GetOrderBookByBase(c, request.QueryValue)
	case "TradingPair":
		queryOutput = q.keeper.GetOrderBookByTradingPair(c, request.QueryValue, request.QueryValue2)
	case "Curator":
		queryOutput = q.keeper.GetOrderBookByCurator(c, request.QueryValue)
	}

	return &types.GetOrderBooksResponse{
		Orderbooks: queryOutput,
	}, nil
}

func (q Querier) GetOrderBooksByTradingPair(ctx context.Context, request *types.GetOrderBooksByTradingPairRequest) (*types.GetOrderBooksResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return &types.GetOrderBooksResponse{
		Orderbooks: q.keeper.GetOrderBookByTradingPair(c, request.Base, request.Quote),
	}, nil
}

func (q Querier) GetOrders(ctx context.Context, request *types.GetOrdersRequest) (*types.GetOrdersResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)

	var queryOutput []types.LimitOrder
	queryOutput = q.keeper.GetOrders(c, request.OrderBookID, request.MaxOrders, request.MinAmount)

	return &types.GetOrdersResponse{
		Orders: queryOutput,
	}, nil
}

func (q Querier) GetSignerKeys(ctx context.Context, request *types.GetSignerKeysRequest) (*types.GetSignerKeysResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)

	var queryOutput []types.SignerKey
	// curator, err := sdk.AccAddressFromBech32()
	// if err != nil {
	// 	return &types.GetSignerKeysResponse{}, fmt.Errorf("Invalid curator address %s: %+v", path[0], err)
	// }

	queryOutput = q.keeper.GetSignerKeys(c, request.Curator)

	return &types.GetSignerKeysResponse{
		Signerkeys: queryOutput,
	}, nil
}
