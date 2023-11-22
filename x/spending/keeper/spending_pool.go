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

func (k Keeper) RemoveClaimInfo(ctx sdk.Context, claimInfo types.ClaimInfo) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ClaimInfoKey(claimInfo.PoolName, claimInfo.Account))
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

func (k Keeper) GetPoolClaimInfos(ctx sdk.Context, poolName string) []types.ClaimInfo {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PoolClaimInfoPrefix(poolName))
	defer iterator.Close()

	claimInfos := []types.ClaimInfo{}
	for ; iterator.Valid(); iterator.Next() {
		claimInfo := types.ClaimInfo{}

		k.cdc.MustUnmarshal(iterator.Value(), &claimInfo)
		claimInfos = append(claimInfos, claimInfo)
	}

	return claimInfos
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

	weight := k.GetBeneficiaryWeight(ctx, sender, *pool.Beneficiaries)
	if weight.IsZero() {
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

	if claimStart >= claimEnd {
		return types.ErrNoMoreRewardsToClaim
	}

	if pool.DynamicRate { // dynamic rate case
		if claimStart < int64(pool.LastDynamicRateCalcTime) {
			claimStart = int64(pool.LastDynamicRateCalcTime)
		}
	}

	rewards := sdk.Coins{}
	for _, rate := range pool.Rates {
		duration := claimEnd - claimStart
		if duration > int64(pool.ClaimExpiry) {
			duration = int64(pool.ClaimExpiry)
		}
		amount := rate.Amount.Mul(sdk.NewDec(duration)).Mul(weight).RoundInt()
		rewards = rewards.Add(sdk.NewCoin(rate.Denom, amount))
	}

	// update pool to reduce pool's balance
	pool.Balances = sdk.Coins(pool.Balances).Sub(rewards...)
	k.SetSpendingPool(ctx, *pool)

	err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, rewards)
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

func (k Keeper) DepositSpendingPoolFromModule(ctx sdk.Context, moduleName, poolName string, amounts sdk.Coins) error {
	err := k.bk.SendCoinsFromModuleToModule(ctx, moduleName, types.ModuleName, amounts)
	if err != nil {
		return err
	}

	pool := k.GetSpendingPool(ctx, poolName)
	if pool == nil {
		return types.ErrPoolDoesNotExist
	}

	pool.Balances = sdk.Coins(pool.Balances).Add(amounts...)
	k.SetSpendingPool(ctx, *pool)
	return nil
}

func (k Keeper) DepositSpendingPoolFromAccount(ctx sdk.Context, addr sdk.AccAddress, poolName string, amounts sdk.Coins) error {
	err := k.bk.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, amounts)
	if err != nil {
		return err
	}

	pool := k.GetSpendingPool(ctx, poolName)
	if pool == nil {
		return types.ErrPoolDoesNotExist
	}

	pool.Balances = sdk.Coins(pool.Balances).Add(amounts...)
	k.SetSpendingPool(ctx, *pool)
	return nil
}

// https://www.notion.so/kira-network/KIP-83-Staking-Collectives-31f72aa35e184978a1a62eb0769ab909
// ### Dynamic Distribution Rates

// Since the amount of staking rewards deposited from the Staking Collective to the Spending Pool changes unpredictably every `claim-period`
// (staking rewards vary from block to block) we need to create a mechanism of dynamic token rates to ensure that all staking rewards
// are spent within some specific time period.

// A new boolean field `dynamic-rate` should be included in the spending pool structure to define whether or not dynamically adjusted
// token rates should be enforced. Once `dynamic-rate` is set to value “true” another field `dynamic-rate-period`
// should define every what period of time all the “token-rates” should change.
// If `dynamic-rate` is set to false, then the spending pool should operate in the same way as it would originally in KIP-71
// (with the exception that it can support claiming multiple tokens and not just one).

// In order to calculate new dynamic token rates we will utilize a "token-deposits” table in which all deposits will be collected
// for the time period equal to `dynamic-rate-period`.
// We will also need to calculate the total number of beneficiaries (users assigned by account and role) and
// divide their number by the total distribution rate we want to achieve.
// As the result, we would allow the fair distribution of all tokens deposited within `dynamic-rate-period` to a set of beneficiaries
// so that no tokens would be left at the end of the period.
// There is no need to take into account that the number of beneficiaries might be changing,
// the only thing that is important is the recalculation of token rates every `dynamic-rate-period`.
// To get the “total_beneficiaries_count” we will only take into account users in the `claims`
// table that registered their intent to withdraw tokens.

// In essence, for each token `x` in `token-deposits` table, the
// `new_token_rate(x) = ( token_deposits(x) / dynamic_rate_period ) / total_beneficiaries_count`.

// It might happen that the pool will run out of tokens because the number of beneficiaries suddenly increased.
// To mitigate that issue we should NOT allow for withdraws of tokens unless the beneficiary is eligible and registered to claim
// the tokens **BEFORE** the current `dynamic-rate-period` started.
// This means that new beneficiaries must await after registration for a new claim period before they can start receiving tokens from the pool.
// If `dynamic-rate` is set to true, then `claim-expire` should be ignored or set to value the same as `dynamic-rate-period` so that the
// maximum time period for which the user should be able to claim tokens is equal to `dynamic-rate-period`.
