package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/custody/types"
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
