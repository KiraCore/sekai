package keeper

import (
	"context"
	"github.com/KiraCore/sekai/x/ethereum/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Querier struct {
	keeper Keeper
}

func NewQuerier(keeper Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}

var _ types.QueryServer = Querier{}

func (q Querier) RelayByAddress(goCtx context.Context, request *types.RelayByAddressRequest) (*types.RelayByAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.RelayByAddressResponse{
		MsgRelay: q.keeper.GetRelayByAddress(ctx, request.Addr),
	}, nil
}
