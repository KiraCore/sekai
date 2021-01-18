package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/KiraCore/sekai/x/staking/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Querier struct {
	keeper Keeper
}

func NewQuerier(keeper Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}

var _ types.QueryServer = Querier{}

func (q Querier) ValidatorByAddress(ctx context.Context, request *types.ValidatorByAddressRequest) (*types.ValidatorResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	val, err := q.keeper.GetValidator(c, request.ValAddr)
	if err != nil {
		return nil, errors.Wrap(errors.ErrKeyNotFound, err.Error())
	}
	return &types.ValidatorResponse{
		Validator: val,
	}, nil
}

func (q Querier) ValidatorByMoniker(ctx context.Context, request *types.ValidatorByMonikerRequest) (*types.ValidatorResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)

	val, err := q.keeper.GetValidatorByMoniker(c, request.Moniker)
	if err != nil {
		return nil, errors.Wrap(errors.ErrKeyNotFound, err.Error())
	}

	return &types.ValidatorResponse{
		Validator: val,
	}, nil
}

// Validators implements the Query all validators gRPC method
func (q Querier) Validators(ctx context.Context, request *types.ValidatorsRequest) (*types.ValidatorsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	store := sdk.UnwrapSDKContext(ctx).KVStore(q.keeper.storeKey)
	validatorStore := prefix.NewStore(store, ValidatorsKey)

	var validators []types.QueryValidator

	pageRes, err := query.Paginate(validatorStore, request.Pagination, func(key []byte, value []byte) error {
		var validator types.Validator
		q.keeper.cdc.MustUnmarshalBinaryBare(value, &validator)
		consPubkey, _ := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, validator.GetConsPubKey())
		fmt.Println("GetConsPubKey", validator.GetConsPubKey().String())
		fmt.Println("GetConsAddr", validator.GetConsAddr().String())

		fmt.Println("Rank", validator.Rank)
		validators = append(validators, types.QueryValidator{
			Address:    sdk.AccAddress(validator.ValKey).String(),
			ValKey:     validator.ValKey.String(),
			PubKey:     consPubkey,
			Moniker:    validator.Moniker,
			Website:    validator.Website,
			Social:     validator.Social,
			Identity:   validator.Identity,
			Commission: validator.Commission.String(),
			Status:     validator.Status.String(),
			Rank:       validator.Rank,
			Streak:     validator.Streak,
		})
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	fmt.Println("validators ========>", validators)
	response := types.ValidatorsResponse{Validators: validators, Pagination: pageRes}

	return &response, nil
}
