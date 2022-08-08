package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SlashStakingPool(ctx sdk.Context, validator string, slash uint64) {
	pool, found := k.GetStakingPoolByValidator(ctx, validator)
	if !found {
		return
	}
	pool.Slashed = slash
	pool.Enabled = false

	totalStakingTokens := sdk.Coins{}
	for _, stakingToken := range pool.TotalStakingTokens {
		totalStakingTokens = totalStakingTokens.Add(sdk.NewCoin(stakingToken.Denom, stakingToken.Amount.Mul(sdk.NewInt(int64(100-slash))).Quo(sdk.NewInt(100))))
	}
	pool.TotalStakingTokens = totalStakingTokens
	k.SetStakingPool(ctx, pool)
}
