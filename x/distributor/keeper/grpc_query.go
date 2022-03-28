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

// QuerydistributorRecords - query names of all distributorRecords
func (q Querier) QuerydistributorRecords(c context.Context, request *types.QuerydistributorRecordsRequest) (*types.QuerydistributorRecordsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QuerydistributorRecordsResponse{
		Records: q.keeper.GetdistributorRecords(ctx),
	}, nil
}

// QuerydistributorRecordByName - query specific distributorRecord by name
func (q Querier) QuerydistributorRecordByName(c context.Context, request *types.QuerydistributorRecordByNameRequest) (*types.QuerydistributorRecordByNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QuerydistributorRecordByNameResponse{
		Record: q.keeper.GetdistributorRecordByName(ctx, request.Name),
	}, nil
}
