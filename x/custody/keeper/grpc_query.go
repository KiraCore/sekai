package keeper

import (
	"context"

	"github.com/KiraCore/sekai/x/custody/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Querier struct {
	keeper Keeper
}

func NewQuerier(keeper Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}

var _ types.QueryServer = Querier{}

func (q Querier) CustodyByAddress(c context.Context, request *types.CustodyByAddressRequest) (*types.CustodyByAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.CustodyByAddressResponse{
		CustodySettings: q.keeper.GetCustodyInfoByAddress(ctx, request.Addr),
	}, nil
}

func (q Querier) CustodyWhiteListByAddress(c context.Context, request *types.CustodyWhiteListByAddressRequest) (*types.CustodyWhiteListByAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.CustodyWhiteListByAddressResponse{
		CustodyWhiteList: q.keeper.GetCustodyWhiteListByAddress(ctx, request.Addr),
	}, nil
}
