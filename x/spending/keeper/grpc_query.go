package keeper

import (
	"context"

	"github.com/KiraCore/sekai/x/spending/types"
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
