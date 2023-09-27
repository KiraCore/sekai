package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/distributor/types"
)

func (k Keeper) GetPeriodicSnapshot(ctx sdk.Context) types.SupplySnapshot {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPeriodicSnapshot)
	if bz == nil {
		return types.SupplySnapshot{}
	}

	snapshot := types.SupplySnapshot{}
	k.cdc.MustUnmarshal(bz, &snapshot)
	return snapshot
}

func (k Keeper) SetPeriodicSnapshot(ctx sdk.Context, snapshot types.SupplySnapshot) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyPeriodicSnapshot, k.cdc.MustMarshal(&snapshot))
}
