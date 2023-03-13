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

	return &types.QueryExecutionRegistrarResponse{
		Dapp: &dapp,
	}, nil
}

func (k Keeper) AllDapps(goCtx context.Context, request *types.QueryAllDappsRequest) (*types.QueryAllDappsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.QueryAllDappsResponse{
		Dapps: k.GetAllDapps(ctx),
	}, nil
}

func (k Keeper) TransferDapp(goCtx context.Context, request *types.QueryTransferDappRequest) (*types.QueryTransferDappResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.QueryTransferDappResponse{}, nil
}

func (k Keeper) GlobalTokens(goCtx context.Context, request *types.QueryGlobalTokensRequest) (*types.QueryGlobalTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.QueryGlobalTokensResponse{}, nil
}
