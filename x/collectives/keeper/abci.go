package keeper

import (
	"github.com/KiraCore/sekai/x/collectives/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {

}

func (k Keeper) EndBlocker(ctx sdk.Context) {
	collectives := k.GetAllCollectives(ctx)
	properties := k.gk.GetNetworkProperties(ctx)

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
		if collective.CreationTime <= ctx.BlockTime().Unix()+int64(properties.MinCollectiveBondingTime) {
			if bondsValue.LT(minCollectiveBond) {
				for _, contributer := range k.GetAllCollectiveContributers(ctx, collective.Name) {
					addr := sdk.MustAccAddressFromBech32(contributer.Address)
					err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, contributer.Bonds)
					if err != nil {
						panic(err)
					}
					k.DeleteCollectiveContributer(ctx, collective.Name, contributer.Address)
				}
				k.DeleteCollective(ctx, collective.Name)
			}
		}
	}

	// TODO: claim staking rewards and distribute them to specified spending pool (only for `active` status)
	// TODO: All donations should be subtracted from the amounts being sent to the spending pools.
	collectives = k.GetAllCollectives(ctx)
	for _, collective := range collectives {
		if collective.Status != types.CollectiveActive {
			continue
		}

		delegator := authtypes.GetModuleAccount(types.ModuleName, collective.Name)
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
				panic(err)
			}
		}

		delegator = authtypes.GetModuleAccount(types.DonationModuleAccount, collective.Name)
		k.mk.RegisterDelegator(ctx, delegator)
		coins = k.mk.ClaimRewards(ctx, delegator)
		collective.Donations = sdk.Coins(collective.Donations).Add(coins...)
		err := k.bk.SendCoinsFromAccountToModule(ctx, delegator, types.ModuleName, coins)
		if err != nil {
			panic(err)
		}
	}
}
