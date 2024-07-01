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

func (q Querier) GetTokenInfo(goCtx context.Context, request *types.TokenInfoRequest) (*types.TokenInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	info := q.keeper.GetTokenInfo(ctx, request.Denom)
	supply := q.keeper.bankKeeper.GetSupply(ctx, request.Denom)
	return &types.TokenInfoResponse{
		Data:   info,
		Supply: supply,
	}, nil
}

func (q Querier) GetTokenInfosByDenom(ctx context.Context, request *types.TokenInfosByDenomRequest) (*types.TokenInfosByDenomResponse, error) {
	infos := q.keeper.GetTokenInfosByDenom(sdk.UnwrapSDKContext(ctx), request.Denoms)
	return &types.TokenInfosByDenomResponse{Data: infos}, nil
}

func (q Querier) GetAllTokenInfos(goCtx context.Context, request *types.AllTokenInfosRequest) (*types.AllTokenInfosResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	infos := q.keeper.GetAllTokenInfos(ctx)
	data := []types.TokenInfoResponse{}
	for _, info := range infos {
		supply := q.keeper.bankKeeper.GetSupply(ctx, info.Denom)
		data = append(data, types.TokenInfoResponse{
			Data:   &info,
			Supply: supply,
		})
	}
	return &types.AllTokenInfosResponse{Data: data}, nil
}

func (q Querier) GetTokenBlackWhites(ctx context.Context, request *types.TokenBlackWhitesRequest) (*types.TokenBlackWhitesResponse, error) {
	data := q.keeper.GetTokenBlackWhites(sdk.UnwrapSDKContext(ctx))
	return &types.TokenBlackWhitesResponse{Data: data}, nil
}
