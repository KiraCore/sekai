package keeper

import (
	"github.com/KiraCore/sekai/x/distributor/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetFeesCollected(ctx sdk.Context, coins sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PrefixKeyFeesCollected, []byte(coins.String()))
}

func (k Keeper) GetFeesCollected(ctx sdk.Context) sdk.Coins {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PrefixKeyFeesCollected)
	if bz == nil {
		return sdk.Coins{}
	}
	coins, err := sdk.ParseCoinsNormalized(string(bz))
	if err != nil {
		panic(err)
	}
	return coins
}

func (k Keeper) SetFeesTreasury(ctx sdk.Context, coins sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PrefixKeyFeesTreasury, []byte(coins.String()))
}

func (k Keeper) GetFeesTreasury(ctx sdk.Context) sdk.Coins {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PrefixKeyFeesTreasury)
	if bz == nil {
		return sdk.Coins{}
	}
	coins, err := sdk.ParseCoinsNormalized(string(bz))
	if err != nil {
		panic(err)
	}
	return coins
}

func (k Keeper) SetSnapPeriod(ctx sdk.Context, period int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PrefixKeySnapPeriod, sdk.Uint64ToBigEndian(uint64(period)))
}

func (k Keeper) GetSnapPeriod(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PrefixKeySnapPeriod)
	if bz == nil {
		return 1
	}
	return int64(sdk.BigEndianToUint64(bz))
}
