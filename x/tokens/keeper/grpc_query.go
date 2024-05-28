package keeper

import (
	"context"

	"github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Querier struct {
	keeper Keeper
}

func NewQuerier(keeper Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}

var _ types.QueryServer = Querier{}

func (q Querier) GetTokenInfo(ctx context.Context, request *types.TokenInfoRequest) (*types.TokenInfoResponse, error) {
	rate := q.keeper.GetTokenInfo(sdk.UnwrapSDKContext(ctx), request.Denom)

	if rate == nil {
		return &types.TokenInfoResponse{Data: nil}, nil
	}
	return &types.TokenInfoResponse{Data: rate}, nil
}

func (q Querier) GetTokenInfosByDenom(ctx context.Context, request *types.TokenInfosByDenomRequest) (*types.TokenInfosByDenomResponse, error) {
	rates := q.keeper.GetTokenInfosByDenom(sdk.UnwrapSDKContext(ctx), request.Denoms)
	return &types.TokenInfosByDenomResponse{Data: rates}, nil
}

func (q Querier) GetAllTokenInfos(ctx context.Context, request *types.AllTokenInfosRequest) (*types.AllTokenInfosResponse, error) {
	rates := q.keeper.GetAllTokenInfos(sdk.UnwrapSDKContext(ctx))
	return &types.AllTokenInfosResponse{Data: rates}, nil
}

func (q Querier) GetTokenBlackWhites(ctx context.Context, request *types.TokenBlackWhitesRequest) (*types.TokenBlackWhitesResponse, error) {
	data := q.keeper.GetTokenBlackWhites(sdk.UnwrapSDKContext(ctx))
	return &types.TokenBlackWhitesResponse{Data: data}, nil
}
