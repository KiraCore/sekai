package keeper

import (
	"context"

	"github.com/KiraCore/sekai/x/layer2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) ExecutionRegistrar(goCtx context.Context, request *types.QueryExecutionRegistrarRequest) (*types.QueryExecutionRegistrarResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	dapp := k.GetDapp(ctx, request.Identifier)
	registarar := k.GetDappSession(ctx, request.Identifier)

	return &types.QueryExecutionRegistrarResponse{
		Dapp:               &dapp,
		ExecutionRegistrar: &registarar,
	}, nil
}

func (k Keeper) AllDapps(goCtx context.Context, request *types.QueryAllDappsRequest) (*types.QueryAllDappsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.QueryAllDappsResponse{
		Dapps: k.GetAllDapps(ctx),
	}, nil
}

func (k Keeper) TransferDapps(goCtx context.Context, request *types.QueryTransferDappsRequest) (*types.QueryTransferDappsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryTransferDappsResponse{
		XAMs: k.GetXAMs(ctx),
	}, nil
}
