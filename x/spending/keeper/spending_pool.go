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

func (k Keeper) ClaimSpendingPool(ctx sdk.Context, poolName string, sender sdk.AccAddress) error {

	pool := k.GetSpendingPool(ctx, poolName)
	if pool == nil {
		return types.ErrPoolDoesNotExist
	}

	if !k.IsAllowedAddress(ctx, sender, *pool.Beneficiaries) {
		return types.ErrNotPoolBeneficiary
	}

	claimInfo := k.GetClaimInfo(ctx, pool.Name, sender)
	if claimInfo == nil {
		return types.ErrNotRegisteredForRewards
	}

	claimStart := int64(pool.ClaimStart)
	if claimStart < int64(claimInfo.LastClaim) {
		claimStart = int64(claimInfo.LastClaim)
	}

	claimEnd := ctx.BlockTime().Unix()
	if pool.ClaimEnd != 0 && claimEnd > int64(pool.ClaimEnd) {
		claimEnd = int64(pool.ClaimEnd)
	}

	if claimStart > claimEnd {
		return types.ErrNoMoreRewardsToClaim
	}

	rewards := pool.Rate.Mul(sdk.NewDec(claimEnd - int64(claimStart))).TruncateInt()

	// update pool to reduce pool's balance
	pool.Balance = pool.Balance.Sub(rewards)
	k.SetSpendingPool(ctx, *pool)

	coins := sdk.Coins{sdk.NewCoin(pool.Token, rewards)}
	err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, coins)
	if err != nil {
		return err
	}

	k.SetClaimInfo(ctx, types.ClaimInfo{
		PoolName:  pool.Name,
		Account:   sender.String(),
		LastClaim: uint64(ctx.BlockTime().Unix()),
	})
	return nil
}

func (k Keeper) DepositSpendingPoolFromModule(ctx sdk.Context, moduleName, poolName string, amount sdk.Coin) error {
	err := k.bk.SendCoinsFromModuleToModule(ctx, moduleName, types.ModuleName, sdk.Coins{amount})
	if err != nil {
		return err
	}

	pool := k.GetSpendingPool(ctx, poolName)
	if pool == nil {
		return types.ErrPoolDoesNotExist
	}

	pool.Balance = pool.Balance.Add(sdk.Coins{amount}.AmountOf(pool.Token))
	k.SetSpendingPool(ctx, *pool)
	return nil
}
