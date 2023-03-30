package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/distributor/types"
)

func (k Keeper) GetYearStartSnapshot(ctx sdk.Context) types.YearStartSnapshot {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyYearStartSnapshot)
	if bz == nil {
		return types.YearStartSnapshot{}
	}

	snapshot := types.YearStartSnapshot{}
	k.cdc.MustUnmarshal(bz, &snapshot)
	return snapshot
}

func (k Keeper) SetYearStartSnapshot(ctx sdk.Context, snapshot types.YearStartSnapshot) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyYearStartSnapshot, k.cdc.MustMarshal(&snapshot))
}

func (k Keeper) InflationPossible(ctx sdk.Context) bool {
	snapshot := k.GetYearStartSnapshot(ctx)
	if snapshot.SnapshotAmount.IsNil() || snapshot.SnapshotAmount.IsZero() {
		return true
	}
	yearlyInflation := k.gk.GetNetworkProperties(ctx).MaxAnnualInflation
	currSupply := k.bk.GetSupply(ctx, k.BondDenom(ctx))

	month := int64(86400 * 30)
	currTimeGone := ctx.BlockTime().Unix() - snapshot.SnapshotTime
	monthIndex := (currTimeGone + month - 1) / month
	currInflation := currSupply.Amount.ToDec().Quo(sdk.Dec(snapshot.SnapshotAmount.ToDec())).Sub(sdk.OneDec())
	fmt.Println("currInflation", currInflation.String())
	fmt.Println("monthIndex", monthIndex)
	fmt.Println("est.Inflation", yearlyInflation.Mul(sdk.NewDec(monthIndex)).Quo(sdk.NewDec(12)).String())
	if currInflation.GTE(yearlyInflation.Mul(sdk.NewDec(monthIndex)).Quo(sdk.NewDec(12))) {
		return false
	}
	return true
}