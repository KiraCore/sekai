package staking

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/staking/keeper"
	"github.com/KiraCore/sekai/x/staking/types"
)

type Querier struct {
	keeper keeper.Keeper
}

func NewQuerier(keeper keeper.Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}

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
