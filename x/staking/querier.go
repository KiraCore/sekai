package staking

import (
	"context"

	"github.com/KiraCore/sekai/x/staking/keeper"
	"github.com/KiraCore/sekai/x/staking/types"
)

type Querier struct {
	keeper keeper.Keeper
}

func (q Querier) ValidatorByAddress(ctx context.Context, request *types.ValidatorByAddressRequest) (*types.ValidatorByAddressResponse, error) {
	panic("implement me")
}
