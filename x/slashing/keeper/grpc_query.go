package keeper

import (
	"context"

	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/slashing/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
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

	validator := stakingtypes.QueryValidator{}
	if req.IncludeValidator {
		val, err := k.sk.GetValidatorByConsAddr(ctx, consAddr)
		if err != nil {
			return nil, err
		}

		records := k.sk.GetIdRecordsByAddress(ctx, sdk.AccAddress(val.ValKey))
		validator = stakingtypes.QueryValidator{
			Address:  sdk.AccAddress(val.ValKey).String(),
			Valkey:   val.ValKey.String(),
			Pubkey:   val.GetConsPubKey().String(),
			Proposer: val.GetConsPubKey().Address().String(),
			Status:   val.Status.String(),
			Rank:     val.Rank,
			Streak:   val.Streak,
			Identity: records,
		}
	}

	return &types.QuerySigningInfoResponse{
		ValSigningInfo: signingInfo,
		Validator:      validator,
	}, nil
}

func (k Keeper) SigningInfos(c context.Context, request *types.QuerySigningInfosRequest) (*types.QuerySigningInfosResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	var signInfos []types.ValidatorSigningInfo
	var validators []stakingtypes.QueryValidator
	var pageRes *query.PageResponse
	var err error

	onResult := func(key []byte, value []byte, accumulate bool) (bool, error) {
		var info types.ValidatorSigningInfo
		err := k.cdc.Unmarshal(value, &info)
		if err != nil {
			return false, err
		}

		if accumulate {
			signInfos = append(signInfos, info)
			if request.IncludeValidator {
				consAddr, err := sdk.ConsAddressFromBech32(info.Address)
				if err != nil {
					return false, err
				}
				val, err := k.sk.GetValidatorByConsAddr(ctx, consAddr)
				if err != nil {
					return false, err
				}
				validators = append(validators, stakingtypes.QueryValidator{
					Address:  sdk.AccAddress(val.ValKey).String(),
					Valkey:   val.ValKey.String(),
					Pubkey:   val.GetConsPubKey().String(),
					Proposer: val.GetConsPubKey().Address().String(),
					Status:   val.Status.String(),
					Rank:     val.Rank,
					Streak:   val.Streak,
					Identity: k.sk.GetIdRecordsByAddress(ctx, sdk.AccAddress(val.ValKey)),
				})
			}
		}
		return true, nil
	}

	// we set maximum limit for safety of iteration
	if request.Pagination != nil && request.Pagination.Limit > kiratypes.PageIterationLimit {
		request.Pagination.Limit = kiratypes.PageIterationLimit
	}

	sigInfoStore := prefix.NewStore(store, types.ValidatorSigningInfoKeyPrefix)
	pageRes, err = query.FilteredPaginate(sigInfoStore, request.Pagination, onResult)

	if err != nil {
		return nil, err
	}

	return &types.QuerySigningInfosResponse{
		Info:       signInfos,
		Validators: validators,
		Pagination: pageRes,
	}, nil
}
