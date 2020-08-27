package ixp

import (
	"context"
	"strconv"

	"github.com/KiraCore/sekai/x/ixp/keeper"
	"github.com/KiraCore/sekai/x/ixp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Querier describe utilities for querying from node
type Querier struct {
	keeper keeper.Keeper
}

// NewQuerier returns an Querier instance
func NewQuerier(keeper keeper.Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}

// GetOrderBooks returns orderbooks by query params
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

// GetOrderBooksByTradingPair returns orderbooks by trading pair query params
func (q Querier) GetOrderBooksByTradingPair(ctx context.Context, request *types.GetOrderBooksByTradingPairRequest) (*types.GetOrderBooksResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return &types.GetOrderBooksResponse{
		Orderbooks: q.keeper.GetOrderBookByTradingPair(c, request.Base, request.Quote),
	}, nil
}

// GetOrders returns orders by by query params
func (q Querier) GetOrders(ctx context.Context, request *types.GetOrdersRequest) (*types.GetOrdersResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)

	var queryOutput []types.LimitOrder
	queryOutput = q.keeper.GetOrders(c, request.OrderBookID, request.MaxOrders, request.MinAmount)

	return &types.GetOrdersResponse{
		Orders: queryOutput,
	}, nil
}

// GetSignerKeys returns signer keys for curators
func (q Querier) GetSignerKeys(ctx context.Context, request *types.GetSignerKeysRequest) (*types.GetSignerKeysResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)

	var queryOutput []types.SignerKey
	queryOutput = q.keeper.GetSignerKeys(c, request.Curator)

	return &types.GetSignerKeysResponse{
		Signerkeys: queryOutput,
	}, nil
}
