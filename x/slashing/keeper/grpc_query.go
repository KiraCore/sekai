package keeper

import (
	"context"

	kiratypes "github.com/KiraCore/sekai/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	multistakingtypes "github.com/KiraCore/sekai/x/multistaking/types"
	"github.com/KiraCore/sekai/x/slashing/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

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

func (k Keeper) SlashProposals(c context.Context, request *types.QuerySlashProposalsRequest) (*types.QuerySlashProposalsResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	proposals, _ := k.gk.GetProposals(ctx)
	slashProposals := []govtypes.Proposal{}
	for _, proposal := range proposals {
		if proposal.GetContent().ProposalType() == kiratypes.ProposalTypeSlashValidator {
			slashProposals = append(slashProposals, proposal)
		}
	}
	return &types.QuerySlashProposalsResponse{
		Proposals: slashProposals,
	}, nil
}

func (k Keeper) SlashedStakingPools(c context.Context, request *types.QuerySlashedStakingPoolsRequest) (*types.QuerySlashedStakingPoolsResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	pools := k.msk.GetAllStakingPools(ctx)
	slashedPools := []multistakingtypes.StakingPool{}
	for _, pool := range pools {
		if pool.Slashed.IsPositive() {
			slashedPools = append(slashedPools, pool)
		}
	}
	return &types.QuerySlashedStakingPoolsResponse{
		Pools: slashedPools,
	}, nil
}

func (k Keeper) ActiveStakingPools(c context.Context, request *types.QueryActiveStakingPoolsRequest) (*types.QueryActiveStakingPoolsResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	pools := k.msk.GetAllStakingPools(ctx)
	activePools := []multistakingtypes.StakingPool{}
	for _, pool := range pools {
		if pool.Slashed.IsZero() && pool.Enabled {
			activePools = append(activePools, pool)
		}
	}
	return &types.QueryActiveStakingPoolsResponse{
		Pools: activePools,
	}, nil
}

func (k Keeper) InactiveStakingPools(c context.Context, request *types.QueryInactiveStakingPoolsRequest) (*types.QueryInactiveStakingPoolsResponse, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	pools := k.msk.GetAllStakingPools(ctx)
	inactivePools := []multistakingtypes.StakingPool{}
	for _, pool := range pools {
		if pool.Slashed.IsZero() && pool.Enabled {
			inactivePools = append(inactivePools, pool)
		}
	}
	return &types.QueryInactiveStakingPoolsResponse{
		Pools: inactivePools,
	}, nil
}
