package keeper

import (
	"context"

	kiratypes "github.com/KiraCore/sekai/types"
	kiraquery "github.com/KiraCore/sekai/types/query"
	"github.com/KiraCore/sekai/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Keeper) SigningInfo(c context.Context, req *types.QuerySigningInfoRequest) (*types.QuerySigningInfoResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.ConsAddress == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	consAddr, err := sdk.ConsAddressFromBech32(req.ConsAddress)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	signingInfo, found := k.GetValidatorSigningInfo(ctx, consAddr)
	if !found {
		return nil, status.Errorf(codes.NotFound, "SigningInfo not found for validator %s", req.ConsAddress)
	}

	return &types.QuerySigningInfoResponse{ValSigningInfo: signingInfo}, nil
}

func (k Keeper) SigningInfos(c context.Context, request *types.QuerySigningInfosRequest) (*types.QuerySigningInfosResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	var signInfos []types.ValidatorSigningInfo
	var pageRes *query.PageResponse
	var err error

	onResult := func(key []byte, value []byte, accumulate bool) (bool, error) {
		var info types.ValidatorSigningInfo
		err := k.cdc.UnmarshalBinaryBare(value, &info)
		if err != nil {
			return false, err
		}

		if accumulate {
			signInfos = append(signInfos, info)
		}
		return true, nil
	}

	// we set maximum limit for safety of iteration
	if request.Pagination.Limit > kiratypes.PageIterationLimit {
		request.Pagination.Limit = kiratypes.PageIterationLimit
	}

	pagination := types.SDKQueryPageReqFromCustomPageReq(request.Pagination)
	sigInfoStore := prefix.NewStore(store, types.ValidatorSigningInfoKeyPrefix)
	if request.All {
		pageRes, err = kiraquery.IterateAll(sigInfoStore, pagination, onResult)
	} else {
		pageRes, err = query.FilteredPaginate(sigInfoStore, pagination, onResult)
	}

	if err != nil {
		return nil, err
	}

	return &types.QuerySigningInfosResponse{Info: signInfos, Pagination: &types.PageResponse{
		NextKey: pageRes.NextKey,
		Total:   pageRes.Total,
	}}, nil
}
