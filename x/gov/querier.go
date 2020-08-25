package gov

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/gov/types"
)

type Querier struct {
	keeper keeper.Keeper
}

func NewQuerier(keeper keeper.Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}

func (q Querier) PermissionsByAddress(ctx context.Context, request *types.PermissionsByAddressRequest) (*types.PermissionsResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)

	networkActor, err := q.keeper.GetNetworkActorByAddress(sdkContext, request.ValAddr)
	if err != nil {
		return nil, err
	}

	return &types.PermissionsResponse{Permissions: networkActor.Permissions}, nil
}
