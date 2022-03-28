package keeper

import (
	"context"

	"github.com/KiraCore/sekai/x/distributor/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Querier struct {
	keeper Keeper
}

func NewQuerier(keeper Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}

var _ types.QueryServer = Querier{}

func (q Querier) FeesTreasury(c context.Context, request *types.QueryFeesTreasuryRequest) (*types.QueryFeesTreasuryResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_ = ctx

	return &types.QueryFeesTreasuryResponse{}, nil
}

func (q Querier) FeesCollected(c context.Context, request *types.QueryFeesCollectedRequest) (*types.QueryFeesCollectedResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_ = ctx

	return &types.QueryFeesCollectedResponse{}, nil
}

func (q Querier) SnapshotPeriod(c context.Context, request *types.QuerySnapshotPeriodRequest) (*types.QuerySnapshotPeriodResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_ = ctx

	return &types.QuerySnapshotPeriodResponse{}, nil
}
