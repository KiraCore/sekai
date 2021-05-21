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
	plan, err := q.keeper.GetUpgradePlan(ctx)
	return &types.QueryCurrentPlanResponse{
		Plan: plan,
	}, err
}

func NewQuerier(keeper Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}
