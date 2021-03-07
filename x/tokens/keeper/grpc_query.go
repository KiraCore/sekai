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

func (q Querier) GetTokenAlias(ctx context.Context, request *types.TokenAliasRequest) (*types.TokenAliasResponse, error) {
	alias := q.keeper.GetTokenAlias(sdk.UnwrapSDKContext(ctx), request.Symbol)

	return &types.TokenAliasResponse{Data: alias}, nil
}

func (q Querier) GetTokenAliasesByDenom(ctx context.Context, request *types.TokenAliasesByDenomRequest) (*types.TokenAliasesByDenomResponse, error) {
	aliases := q.keeper.GetTokenAliasesByDenom(sdk.UnwrapSDKContext(ctx), request.Denoms)

	return &types.TokenAliasesByDenomResponse{Data: aliases}, nil
}

func (q Querier) GetAllTokenAliases(ctx context.Context, request *types.AllTokenAliasesRequest) (*types.AllTokenAliasesResponse, error) {
	aliases := q.keeper.ListTokenAlias(sdk.UnwrapSDKContext(ctx))

	return &types.AllTokenAliasesResponse{Data: aliases}, nil
}

func (q Querier) GetTokenRate(ctx context.Context, request *types.TokenRateRequest) (*types.TokenRateResponse, error) {
	rate := q.keeper.GetTokenRate(sdk.UnwrapSDKContext(ctx), request.Denom)

	if rate == nil {
		return &types.TokenRateResponse{Data: nil}, nil
	}
	return &types.TokenRateResponse{Data: rate}, nil
}

func (q Querier) GetTokenRatesByDenom(ctx context.Context, request *types.TokenRatesByDenomRequest) (*types.TokenRatesByDenomResponse, error) {
	rates := q.keeper.GetTokenRatesByDenom(sdk.UnwrapSDKContext(ctx), request.Denoms)
	return &types.TokenRatesByDenomResponse{Data: rates}, nil
}

func (q Querier) GetAllTokenRates(ctx context.Context, request *types.AllTokenRatesRequest) (*types.AllTokenRatesResponse, error) {
	rates := q.keeper.ListTokenRate(sdk.UnwrapSDKContext(ctx))
	return &types.AllTokenRatesResponse{Data: rates}, nil
}

func (q Querier) GetTokenBlackWhites(ctx context.Context, request *types.TokenBlackWhitesRequest) (*types.TokenBlackWhitesResponse, error) {
	data := q.keeper.GetTokenBlackWhites(sdk.UnwrapSDKContext(ctx))
	return &types.TokenBlackWhitesResponse{Data: data}, nil
}
