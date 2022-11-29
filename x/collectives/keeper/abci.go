package keeper

import (
	"github.com/KiraCore/sekai/x/collectives/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {

}

func (k Keeper) DistributeCollectiveRewards(ctx sdk.Context, collective types.Collective) error {
	delegator := collective.GetCollectiveAddress()
	k.mk.RegisterDelegator(ctx, delegator)
	coins := k.mk.ClaimRewards(ctx, delegator)

	// send to spending pools based on weight
	for _, pool := range collective.SpendingPools {
		portionCoins := calcPortion(coins, pool.Weight)
		pool := k.spk.GetSpendingPool(ctx, pool.Name)
		if pool == nil {
			continue
		}

		err := k.spk.DepositSpendingPoolFromAccount(ctx, delegator, pool.Name, portionCoins)
		if err != nil {
			return err
		}
	}

	delegator = collective.GetCollectiveDonationAddress()
	k.mk.RegisterDelegator(ctx, delegator)
	coins = k.mk.ClaimRewards(ctx, delegator)
	collective.Donations = sdk.Coins(collective.Donations).Add(coins...)
	err := k.bk.SendCoinsFromAccountToModule(ctx, delegator, types.ModuleName, coins)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) EndBlocker(ctx sdk.Context) {
	collectives := k.GetAllCollectives(ctx)
	properties := k.gk.GetNetworkProperties(ctx)

	collectives = k.GetAllCollectives(ctx)
	for _, collective := range collectives {
		if collective.Status != types.CollectiveActive {
			continue
		}

		// Do distribution per interval or just after claim start
		blockTime := uint64(ctx.BlockTime().Unix())
		if (collective.ClaimStart >= blockTime && collective.LastDistribution == 0) ||
			collective.LastDistribution+collective.ClaimPeriod <= blockTime {
			cacheCtx, write := ctx.CacheContext()
			err := k.DistributeCollectiveRewards(cacheCtx, collective)
			if err == nil {
				write()
			}
			collective.LastDistribution = uint64(ctx.BlockTime().Unix())
		}
	}

	for _, collective := range collectives {
		bondsValue := k.GetBondsValue(ctx, collective.Bonds)

		// For the collective to become activated a minimum bond amount of tokens will have to be committed to the collective pool,
		// the default `min_collective_bond` should be equivalent to 100’000 KEX
		// and configurable in the [Network Properties](https://www.notion.so/de74fe4b731a47df86683f2e9eefa793)
		minCollectiveBond := sdk.NewDec(int64(properties.MinCollectiveBond)).Mul(sdk.NewDec(1000_000))

		// To be `active`, ClaimStart time should pass
		if collective.ClaimStart <= uint64(ctx.BlockTime().Unix()) &&
			(collective.ClaimEnd == 0 || collective.ClaimEnd >= uint64(ctx.BlockTime().Unix())) &&
			collective.Status != types.CollectivePaused {
			if bondsValue.GTE(minCollectiveBond) {
				collective.Status = types.CollectiveActive
			} else {
				collective.Status = types.CollectiveInactive
			}
		}
		k.SetCollective(ctx, collective)

		// if minimum collective bonding time pass
		if int64(collective.CreationTime+properties.MinCollectiveBondingTime) <= ctx.BlockTime().Unix() {
			if bondsValue.LT(minCollectiveBond) {
				cacheCtx, write := ctx.CacheContext()
				err := k.ExecuteCollectiveRemove(cacheCtx, collective)
				if err == nil {
					write()
				}
			}
		}
	}

}
