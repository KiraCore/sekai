package keeper

import (
	"github.com/KiraCore/sekai/x/spending/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetSpendingPool stores spending pool record
func (k Keeper) SetSpendingPool(ctx sdk.Context, pool types.SpendingPool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&pool)
	store.Set(types.SpendingPoolKey(pool.Name), bz)
}

// GetSpendingPool returns SpendingPool stored inside keeper
func (k Keeper) GetSpendingPool(ctx sdk.Context, name string) *types.SpendingPool {
	var pool types.SpendingPool
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.SpendingPoolKey(name))
	if bz == nil {
		return nil
	}

	k.cdc.MustUnmarshal(bz, &pool)

	return &pool
}

func (k Keeper) GetAllSpendingPools(ctx sdk.Context) []types.SpendingPool {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.KeyPrefixSpendingPool))
	defer iterator.Close()

	pools := []types.SpendingPool{}
	for ; iterator.Valid(); iterator.Next() {
		pool := types.SpendingPool{}

		k.cdc.MustUnmarshal(iterator.Value(), &pool)
		pools = append(pools, pool)
	}
	return pools
}

func (k Keeper) CreateSpendingPool(ctx sdk.Context, pool types.SpendingPool) error {
	oldPool := k.GetSpendingPool(ctx, pool.Name)
	if oldPool != nil {
		return types.ErrAlreadyRegisteredPoolName
	}

	k.SetSpendingPool(ctx, pool)
	return nil
}

func (k Keeper) SetClaimInfo(ctx sdk.Context, claimInfo types.ClaimInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&claimInfo)
	store.Set(types.ClaimInfoKey(claimInfo.PoolName, claimInfo.Account), bz)
}

func (k Keeper) GetClaimInfo(ctx sdk.Context, poolName string, address sdk.AccAddress) *types.ClaimInfo {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ClaimInfoKey(poolName, address.String()))

	if bz == nil {
		return nil
	}

	claimInfo := types.ClaimInfo{}
	k.cdc.MustUnmarshal(bz, &claimInfo)
	return &claimInfo
}

func (k Keeper) GetAllClaimInfos(ctx sdk.Context) []types.ClaimInfo {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.KeyPrefixClaimInfo))
	defer iterator.Close()

	claimInfos := []types.ClaimInfo{}
	for ; iterator.Valid(); iterator.Next() {
		claimInfo := types.ClaimInfo{}

		k.cdc.MustUnmarshal(iterator.Value(), &claimInfo)
		claimInfos = append(claimInfos, claimInfo)
	}

	return claimInfos
}
