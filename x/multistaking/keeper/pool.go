package keeper

import (
	"github.com/KiraCore/sekai/x/multistaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
