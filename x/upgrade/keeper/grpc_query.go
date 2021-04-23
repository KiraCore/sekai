package keeper

import (
	"context"
	"github.com/KiraCore/sekai/x/upgrade/types"
)

type Querier struct {
	keeper Keeper
}

func (q Querier) CurrentPlan(ctx context.Context, request *types.QueryCurrentPlanRequest) (*types.QueryCurrentPlanResponse, error) {
	panic("implement me")
}

func NewQuerier(keeper Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}
