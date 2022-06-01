package keeper

import (
	"github.com/KiraCore/sekai/x/multistaking/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetLastUndelegationId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyLastUndelegationId)
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLastUndelegationId(ctx sdk.Context, id uint64) {
	idBz := sdk.Uint64ToBigEndian(id)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyLastUndelegationId, idBz)
}

func (k Keeper) GetUndelegationById(ctx sdk.Context, id uint64) (undelegation types.Undelegation, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.KeyPrefixUndelegation), sdk.Uint64ToBigEndian(id)...)
	bz := store.Get(key)
	if bz == nil {
		return undelegation, false
	}
	k.cdc.MustUnmarshal(bz, &undelegation)
	return undelegation, true
}

func (k Keeper) GetAllUndelegations(ctx sdk.Context) []types.Undelegation {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixUndelegation)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	undelegations := []types.Undelegation{}
	for ; iterator.Valid(); iterator.Next() {
		undelegation := types.Undelegation{}
		k.cdc.MustUnmarshal(iterator.Value(), &undelegation)
		undelegations = append(undelegations, undelegation)
	}
	return undelegations
}

func (k Keeper) SetUndelegation(ctx sdk.Context, undelegation types.Undelegation) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixUndelegation, sdk.Uint64ToBigEndian(undelegation.Id)...)
	store.Set(key, k.cdc.MustMarshal(&undelegation))
}

func (k Keeper) SetPoolDelegator(ctx sdk.Context, poolId uint64, delegator sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := append(append(types.KeyPrefixPoolDelegator, sdk.Uint64ToBigEndian(poolId)...), delegator...)
	store.Set(key, delegator)
}

func (k Keeper) RemovePoolDelegator(ctx sdk.Context, poolId uint64, delegator sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := append(append(types.KeyPrefixPoolDelegator, sdk.Uint64ToBigEndian(poolId)...), delegator...)
	store.Delete(key)
}

func (k Keeper) GetPoolDelegators(ctx sdk.Context, poolId uint64) []sdk.AccAddress {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.KeyPrefixPoolDelegator, sdk.Uint64ToBigEndian(poolId)...))

	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	delegators := []sdk.AccAddress{}
	for ; iterator.Valid(); iterator.Next() {
		delegators = append(delegators, sdk.AccAddress(iterator.Value()))
	}
	return delegators
}

func (k Keeper) IncreaseDelegatorRewards(ctx sdk.Context, delegator sdk.AccAddress, amounts sdk.Coins) {
	rewards := k.GetDelegatorRewards(ctx, delegator)
	rewards = rewards.Add(amounts...)

	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixRewards, delegator...)
	store.Set(key, []byte(rewards.String()))
}

func (k Keeper) GetDelegatorRewards(ctx sdk.Context, delegator sdk.AccAddress) sdk.Coins {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixRewards, delegator...)
	bz := store.Get(key)
	if bz == nil {
		return sdk.Coins{}
	}

	coinStr := string(bz)
	coins, err := sdk.ParseCoinsNormalized(coinStr)
	if err != nil {
		panic(err)
	}

	return coins
}

func (k Keeper) IncreasePoolRewards(ctx sdk.Context, pool types.StakingPool, rewards sdk.Coins) {
	totalWeight := sdk.ZeroDec()
	for _, shareToken := range pool.TotalShareTokens {
		nativeDenom := getNativeDenom(pool.Id, shareToken.Denom)
		rate := k.tokenKeeper.GetTokenRate(ctx, nativeDenom)
		if rate == nil {
			continue
		}

		totalWeight = totalWeight.Add(shareToken.Amount.ToDec().Mul(rate.FeeRate))
	}

	if totalWeight.IsZero() {
		return
	}

	delegators := k.GetPoolDelegators(ctx, pool.Id)
	for _, delegator := range delegators {
		weight := sdk.ZeroDec()
		balances := k.bankKeeper.GetAllBalances(ctx, delegator)
		for _, shareToken := range pool.TotalShareTokens {
			nativeDenom := getNativeDenom(pool.Id, shareToken.Denom)
			rate := k.tokenKeeper.GetTokenRate(ctx, nativeDenom)
			balance := balances.AmountOf(shareToken.Denom)
			if rate == nil {
				continue
			}
			weight = weight.Add(balance.ToDec().Mul(rate.FeeRate))
		}

		delegatorRewards := sdk.Coins{}
		for _, reward := range rewards {
			delegatorRewards = delegatorRewards.Add(sdk.NewCoin(reward.Denom, reward.Amount.Mul(weight.Quo(totalWeight).RoundInt())))
		}

		k.IncreaseDelegatorRewards(ctx, delegator, delegatorRewards)
	}
}
