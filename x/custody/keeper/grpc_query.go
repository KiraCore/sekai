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

func (q Querier) CustodyCustodiansByAddress(c context.Context, request *types.CustodyCustodiansByAddressRequest) (*types.CustodyCustodiansByAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.CustodyCustodiansByAddressResponse{
		CustodyCustodians: q.keeper.GetCustodyCustodiansByAddress(ctx, request.Addr),
	}, nil
}

func (q Querier) CustodyPoolByAddress(c context.Context, request *types.CustodyPoolByAddressRequest) (*types.CustodyPoolByAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.CustodyPoolByAddressResponse{
		Transactions: q.keeper.GetCustodyPoolByAddress(ctx, request.Addr),
	}, nil
}

func (q Querier) CustodyWhiteListByAddress(c context.Context, request *types.CustodyWhiteListByAddressRequest) (*types.CustodyWhiteListByAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.CustodyWhiteListByAddressResponse{
		CustodyWhiteList: q.keeper.GetCustodyWhiteListByAddress(ctx, request.Addr),
	}, nil
}

func (q Querier) CustodyLimitsByAddress(c context.Context, request *types.CustodyLimitsByAddressRequest) (*types.CustodyLimitsByAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.CustodyLimitsByAddressResponse{
		CustodyLimits: q.keeper.GetCustodyLimitsByAddress(ctx, request.Addr),
	}, nil
}

func (q Querier) CustodyLimitsStatusByAddress(c context.Context, request *types.CustodyLimitsStatusByAddressRequest) (*types.CustodyLimitsStatusByAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.CustodyLimitsStatusByAddressResponse{
		CustodyStatuses: q.keeper.GetCustodyLimitsStatusByAddress(ctx, request.Addr),
	}, nil
}
