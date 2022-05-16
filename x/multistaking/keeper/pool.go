package keeper

import (
	"fmt"

	"github.com/KiraCore/sekai/x/multistaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: set id and validator mapping

func (k Keeper) GetLastPoolId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyLastPoolId)
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLastPoolId(ctx sdk.Context, id uint64) {
	idBz := sdk.Uint64ToBigEndian(id)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyLastPoolId, idBz)
}

func (k Keeper) GetStakingPoolByValidator(ctx sdk.Context, validator string) (pool types.StakingPool, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.KeyPrefixStakingPool), []byte(validator)...)
	bz := store.Get(key)
	if bz == nil {
		return pool, false
	}
	k.cdc.MustUnmarshal(bz, &pool)
	return pool, true
}

func (k Keeper) SetStakingPool(ctx sdk.Context, pool types.StakingPool) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.KeyPrefixStakingPool), []byte(pool.Validator)...)
	store.Set(key, k.cdc.MustMarshal(&pool))
}

func getPoolCoins(poolID uint64, coins sdk.Coins) sdk.Coins {
	prefix := fmt.Sprintf("v%d_", poolID)
	poolCoins := sdk.Coins{}
	for _, coin := range coins {
		poolCoins = poolCoins.Add(sdk.NewCoin(prefix+coin.Denom, coin.Amount))
	}
	return poolCoins
}
