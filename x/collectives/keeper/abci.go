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
		if bondsValue.LT(minCollectiveBond) && collective.Status == types.CollectiveInactive {
			collective.Status = types.CollectiveActive
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
}
