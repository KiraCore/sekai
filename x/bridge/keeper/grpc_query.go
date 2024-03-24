package keeper

import (
	"context"
	"github.com/KiraCore/sekai/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Querier struct {
	keeper Keeper
}

func NewQuerier(keeper Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}

var _ types.QueryServer = Querier{}

func (q Querier) ChangeCosmosEthereumByAddress(c context.Context, request *types.ChangeCosmosEthereumByAddressRequest) (*types.ChangeCosmosEthereumByAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	record := q.keeper.GetChangeCosmosEthereumRecord(ctx, request.Addr)

	return &types.ChangeCosmosEthereumByAddressResponse{
		From:      record.From,
		To:        record.To,
		InAmount:  record.InAmount,
		OutAmount: record.OutAmount,
	}, nil
}

func (q Querier) ChangeEthereumCosmosByAddress(c context.Context, request *types.ChangeEthereumCosmosByAddressRequest) (*types.ChangeEthereumCosmosByAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	record := q.keeper.GetChangeEthereumCosmosRecord(ctx, request.Addr)

	return &types.ChangeEthereumCosmosByAddressResponse{
		From:      record.From,
		To:        record.To,
		InAmount:  record.InAmount,
		OutAmount: record.OutAmount,
	}, nil
}
