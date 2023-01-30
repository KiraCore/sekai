package keeper

import (
	"context"

	"github.com/KiraCore/sekai/x/recovery/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

// given KIRA public address as parameter return data from the recovery registrar
func (k Keeper) RecoveryRecord(c context.Context, req *types.QueryRecoveryRecordRequest) (*types.QueryRecoveryRecordResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	record, err := k.GetRecoveryRecord(ctx, req.Address)
	if err != nil {
		return nil, err
	}
	return &types.QueryRecoveryRecordResponse{Record: record}, nil
}

func (k Keeper) RecoveryToken(c context.Context, req *types.QueryRecoveryTokenRequest) (*types.QueryRecoveryTokenResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	token, err := k.GetRecoveryToken(ctx, req.Address)
	if err != nil {
		return nil, err
	}
	return &types.QueryRecoveryTokenResponse{Token: token}, nil
}
