package keeper

import (
	"context"

	"github.com/KiraCore/sekai/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Querier struct {
	keeper Keeper
}

func (q Querier) CurrentPlan(goCtx context.Context, request *types.QueryCurrentPlanRequest) (*types.QueryCurrentPlanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	plan, err := q.keeper.GetCurrentPlan(ctx)
	return &types.QueryCurrentPlanResponse{
		Plan: plan,
	}, err
}

func (q Querier) NextPlan(goCtx context.Context, request *types.QueryNextPlanRequest) (*types.QueryNextPlanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	plan, err := q.keeper.GetNextPlan(ctx)
	return &types.QueryNextPlanResponse{
		Plan: plan,
	}, err
}

func NewQuerier(keeper Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}
