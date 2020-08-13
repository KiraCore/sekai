package staking

import (
	"context"

	sdk "github.com/KiraCore/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/staking/keeper"
	"github.com/KiraCore/sekai/x/staking/types"
)

type Querier struct {
	keeper keeper.Keeper
}

func NewQuerier(keeper keeper.Keeper) *Querier {
	return &Querier{keeper: keeper}
}

func (q Querier) ValidatorByAddress(ctx context.Context, request *types.ValidatorByAddressRequest) (*types.ValidatorByAddressResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)

	validator := q.keeper.GetValidator(c, request.ValAddr)

	return &types.ValidatorByAddressResponse{Validator: validator}, nil
}
