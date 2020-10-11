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
