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
	k.SetStakingPool(ctx, pool)
}
