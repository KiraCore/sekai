package keeper

import (
	"fmt"
	"time"

	"github.com/KiraCore/sekai/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/golang/protobuf/proto"
)

func (k Keeper) GetCurrentPlan(ctx sdk.Context) (*types.Plan, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.KeyCurrentPlan) {
		return nil, nil
	}

	plan := types.Plan{}
	bz := store.Get(types.KeyCurrentPlan)
	err := proto.Unmarshal(bz, &plan)
	return &plan, err
}

func (k Keeper) SaveCurrentPlan(ctx sdk.Context, plan types.Plan) {
	store := ctx.KVStore(k.storeKey)
	bz, err := proto.Marshal(&plan)
	if err != nil {
		panic(err)
	}
	store.Set(types.KeyCurrentPlan, bz)
}

func (k Keeper) GetNextPlan(ctx sdk.Context) (*types.Plan, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.KeyNextPlan) {
		return nil, nil
	}

	plan := types.Plan{}
	bz := store.Get(types.KeyNextPlan)
	err := proto.Unmarshal(bz, &plan)
	return &plan, err
}

func (k Keeper) SaveNextPlan(ctx sdk.Context, plan types.Plan) error {
	if plan.UpgradeTime <= ctx.BlockTime().Unix() {
		return types.ErrInvalidUpgradeTime
	}
	store := ctx.KVStore(k.storeKey)
	bz, err := proto.Marshal(&plan)
	if err != nil {
		panic(err)
	}
	store.Set(types.KeyNextPlan, bz)
	return nil
}

func (k Keeper) ClearNextPlan(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyNextPlan)
}

func (k Keeper) ApplyUpgradePlan(ctx sdk.Context, plan types.Plan) {
	if plan.ShouldExecute(ctx) {
		k.SaveCurrentPlan(ctx, plan)
		k.ClearNextPlan(ctx)

		if plan.InstateUpgrade == false {
			panic(fmt.Sprintf("UPGRADE \"%s\" NEEDED at upgrade_time=%s", plan.Name, time.Unix(plan.UpgradeTime, 0).String()))
		}

		handler := k.upgradeHandlers[plan.Name]
		if plan.SkipHandler == true {
			fmt.Printf("Skipping handler for \"%s\" upgrade", plan.Name)
		} else if handler == nil {
			panic(fmt.Sprintf("Handler for \"%s\" instate upgrade is not set", plan.Name))
		} else {
			handler(ctx, plan)
		}
	}
}
