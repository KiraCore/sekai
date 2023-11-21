package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {

}

func (k Keeper) EndBlocker(ctx sdk.Context) {
	pools := k.GetAllSpendingPools(ctx)

	for _, pool := range pools {
		if !pool.DynamicRate {
			continue
		}
		if pool.DynamicRatePeriod+pool.LastDynamicRateCalcTime > uint64(ctx.BlockTime().Unix()) {
			continue
		}

		claimInfos := k.GetPoolClaimInfos(ctx, pool.Name)
		totalWeight := sdk.ZeroDec()
		for _, info := range claimInfos {
			addr := sdk.MustAccAddressFromBech32(info.Account)
			weight := k.GetBeneficiaryWeight(ctx, addr, *pool.Beneficiaries)
			totalWeight = totalWeight.Add(weight)
		}

		if totalWeight.IsZero() {
			continue
		}

		// If any of the weights is changed to value different then 1, then token rates should be recalculated accordingly, that is for each token `x`
		// in `token-deposits` table, the
		// `new_token_rate(x) = ( ( token_deposits(x) /  (dynamic_rate_period * weights_sum))`
		poolRates := sdk.DecCoins{}
		for _, deposit := range pool.Balances {
			rate := sdk.NewDecFromInt(deposit.Amount).Quo(sdk.NewDec(int64(pool.DynamicRatePeriod)).Mul(totalWeight))
			poolRates = poolRates.Add(sdk.NewDecCoinFromDec(deposit.Denom, rate))
		}
		pool.Rates = poolRates
		pool.LastDynamicRateCalcTime = uint64(ctx.BlockTime().Unix())
		k.SetSpendingPool(ctx, pool)
	}
}
