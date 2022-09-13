package keeper

import (
	"context"

	"github.com/KiraCore/sekai/x/basket/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Querier struct {
	keeper Keeper
}

func NewQuerier(keeper Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}

var _ types.QueryServer = Querier{}

func (q Querier) TokenBasketById(c context.Context, request *types.QueryTokenBasketByIdRequest) (*types.QueryTokenBasketByIdResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	basket, err := q.keeper.GetBasketById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &types.QueryTokenBasketByIdResponse{
		Basket: &basket,
	}, nil
}

func (q Querier) TokenBasketByDenom(c context.Context, request *types.QueryTokenBasketByDenomRequest) (*types.QueryTokenBasketByDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	basket, err := q.keeper.GetBasketByDenom(ctx, request.Denom)
	if err != nil {
		return nil, err
	}
	return &types.QueryTokenBasketByDenomResponse{
		Basket: &basket,
	}, nil
}

func (q Querier) TokenBaskets(c context.Context, request *types.QueryTokenBasketsRequest) (*types.QueryTokenBasketsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// if `tokens` flag is set (comma separated array of strings) return
	// list of all `id`'s of all the baskets that accept any of the specified `denom`’s as deposit
	baskets := q.keeper.GetAllBaskets(ctx)

	if len(request.Tokens) > 0 {
		filtered := []types.Basket{}
		for _, basket := range baskets {
			if basket.DenomExists(request.Tokens) {
				filtered = append(filtered, basket)
			}
		}
		return &types.QueryTokenBasketsResponse{
			Baskets: filtered,
		}, nil
	}

	return &types.QueryTokenBasketsResponse{
		Baskets: baskets,
	}, nil
}
