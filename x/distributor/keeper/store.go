package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetFeesCollected(ctx sdk.Context, coins sdk.Coins) {

}

func (k Keeper) GetFeesCollected(ctx sdk.Context) sdk.Coins {
	return sdk.Coins{}
}

func (k Keeper) SetFeesTreasury(ctx sdk.Context, coins sdk.Coins) {

}

func (k Keeper) GetFeesTreasury(ctx sdk.Context) sdk.Coins {
	return sdk.Coins{}

}

func (k Keeper) SetSnapPeriod(ctx sdk.Context, period int64) {
}

func (k Keeper) GetSnapPeriod(ctx sdk.Context) int64 {
	return 1
}
