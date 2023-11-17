package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// combine multiple staking hooks, all hook functions are run in array sequence
type Hooks []MultistakingHooks

func NewMultiStakingHooks(hooks ...MultistakingHooks) Hooks {
	return hooks
}

func (h Hooks) AfterUpsertStakingPool(ctx sdk.Context, valAddr sdk.ValAddress, pool StakingPool) {
	for i := range h {
		h[i].AfterUpsertStakingPool(ctx, valAddr, pool)
	}
}

func (h Hooks) AfterSlashStakingPool(ctx sdk.Context, valAddr sdk.ValAddress, pool StakingPool, slash sdk.Dec) {
	for i := range h {
		h[i].AfterSlashStakingPool(ctx, valAddr, pool, slash)
	}
}
