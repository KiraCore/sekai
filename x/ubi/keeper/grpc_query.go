package keeper

import (
	"context"

	"github.com/KiraCore/sekai/x/ubi/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Querier struct {
	keeper Keeper
}

func NewQuerier(keeper Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}

var _ types.QueryServer = Querier{}

// QueryUBIRecords - query names of all UBIRecords
func (q Querier) QueryUBIRecords(c context.Context, request *types.QueryUBIRecordsRequest) (*types.QueryUBIRecordsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryUBIRecordsResponse{
		Records: q.keeper.GetUBIRecords(ctx),
	}, nil
}

// QueryUBIRecordByName - query specific UBIRecord by name
func (q Querier) QueryUBIRecordByName(c context.Context, request *types.QueryUBIRecordByNameRequest) (*types.QueryUBIRecordByNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryUBIRecordByNameResponse{
		Record: q.keeper.GetUBIRecordByName(ctx, request.Name),
	}, nil
}
