package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {

}

func (k Keeper) EndBlocker(ctx sdk.Context) {
	dapps := k.GetAllDapps(ctx)
	properties := k.gk.GetNetworkProperties(ctx)

	for _, dapp := range dapps {
		minDappBond := sdk.NewInt(int64(properties.MinDappBond)).Mul(sdk.NewInt(1000_000))

		if int64(dapp.CreationTime+properties.DappBondDuration) <= ctx.BlockTime().Unix() {
			if dapp.TotalBond.Amount.LT(minDappBond) {
				cacheCtx, write := ctx.CacheContext()
				err := k.ExecuteDappRemove(cacheCtx, dapp)
				if err == nil {
					write()
				}
			}
		}
	}

}
