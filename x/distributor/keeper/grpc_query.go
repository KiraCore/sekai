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

	return &types.QueryFeesTreasuryResponse{
		Coins: q.keeper.GetFeesTreasury(ctx),
	}, nil
}

func (q Querier) SnapshotPeriod(c context.Context, request *types.QuerySnapshotPeriodRequest) (*types.QuerySnapshotPeriodResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QuerySnapshotPeriodResponse{
		SnapshotPeriod: q.keeper.GetSnapPeriod(ctx),
	}, nil
}

func (q Querier) SnapshotPeriodPerformance(c context.Context, request *types.QuerySnapshotPeriodPerformanceRequest) (*types.QuerySnapshotPeriodPerformanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	valAddr, err := sdk.ValAddressFromBech32(request.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	performance, err := q.keeper.GetValidatorPerformance(ctx, valAddr)
	return &types.QuerySnapshotPeriodPerformanceResponse{
		SnapshotPeriod: q.keeper.GetSnapPeriod(ctx),
		Performance:    performance,
	}, nil
}

func (q Querier) YearStartSnapshot(c context.Context, request *types.QueryYearStartSnapshotRequest) (*types.QueryYearStartSnapshotResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryYearStartSnapshotResponse{
		Snapshot: q.keeper.GetYearStartSnapshot(ctx),
	}, nil
}

func (q Querier) PeriodicSnapshot(c context.Context, request *types.QueryPeriodicSnapshotRequest) (*types.QueryPeriodicSnapshotResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryPeriodicSnapshotResponse{
		Snapshot: q.keeper.GetPeriodicSnapshot(ctx),
	}, nil
}
