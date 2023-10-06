package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/distributor/types"
)

func (k Keeper) GetYearStartSnapshot(ctx sdk.Context) types.SupplySnapshot {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyYearStartSnapshot)
	if bz == nil {
		return types.SupplySnapshot{}
	}

	snapshot := types.SupplySnapshot{}
	k.cdc.MustUnmarshal(bz, &snapshot)
	return snapshot
}

func (k Keeper) SetYearStartSnapshot(ctx sdk.Context, snapshot types.SupplySnapshot) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyYearStartSnapshot, k.cdc.MustMarshal(&snapshot))
}

func (k Keeper) InflationPossible(ctx sdk.Context) bool {
	snapshot := k.GetYearStartSnapshot(ctx)
	if snapshot.SnapshotAmount.IsNil() || snapshot.SnapshotAmount.IsZero() {
		return true
	}
	yearlyInflation := k.gk.GetNetworkProperties(ctx).MaxAnnualInflation
	currSupply := k.bk.GetSupply(ctx, k.DefaultDenom(ctx))

	month := int64(86400 * 30)
	currTimeGone := ctx.BlockTime().Unix() - snapshot.SnapshotTime
	monthIndex := (currTimeGone + month - 1) / month
	currInflation := sdk.NewDecFromInt(currSupply.Amount).Quo(sdk.NewDecFromInt(snapshot.SnapshotAmount)).Sub(sdk.OneDec())
	if currInflation.GTE(yearlyInflation.Mul(sdk.NewDec(monthIndex)).Quo(sdk.NewDec(12))) {
		return false
	}
	return true
}
