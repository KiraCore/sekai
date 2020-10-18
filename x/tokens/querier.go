package tokens

import (
	"context"

	"github.com/KiraCore/sekai/x/tokens/keeper"
	"github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Querier struct {
	keeper keeper.Keeper
}

func NewQuerier(keeper keeper.Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}

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

	return &types.TokenRateResponse{Data: rate.ToHumanReadable()}, nil
}

func (q Querier) GetTokenRatesByDenom(ctx context.Context, request *types.TokenRatesByDenomRequest) (*types.TokenRatesByDenomResponse, error) {
	ratesRaw := q.keeper.GetTokenRatesByDenom(sdk.UnwrapSDKContext(ctx), request.Denoms)
	ratesHR := make(map[string]*types.TokenRateHumanReadable)

	for k, v := range ratesRaw {
		ratesHR[k] = v.ToHumanReadable()
	}

	return &types.TokenRatesByDenomResponse{Data: ratesHR}, nil
}

func (q Querier) GetAllTokenRates(ctx context.Context, request *types.AllTokenRatesRequest) (*types.AllTokenRatesResponse, error) {
	ratesRaw := q.keeper.ListTokenRate(sdk.UnwrapSDKContext(ctx))
	ratesHR := []*types.TokenRateHumanReadable{}

	for _, rate := range ratesRaw {
		ratesHR = append(ratesHR, rate.ToHumanReadable())
	}

	return &types.AllTokenRatesResponse{Data: ratesHR}, nil
}
