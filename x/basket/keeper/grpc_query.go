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
	_ = ctx
	return &types.QueryTokenBasketByIdResponse{}, nil
}

func (q Querier) TokenBasketByDenom(c context.Context, request *types.QueryTokenBasketByDenomRequest) (*types.QueryTokenBasketByDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_ = ctx
	return &types.QueryTokenBasketByDenomResponse{}, nil
}

func (q Querier) TokenBaskets(c context.Context, request *types.QueryTokenBasketsRequest) (*types.QueryTokenBasketsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_ = ctx
	return &types.QueryTokenBasketsResponse{}, nil
}
