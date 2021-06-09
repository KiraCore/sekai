package keeper

import (
	"fmt"
	"time"

	"github.com/KiraCore/sekai/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/golang/protobuf/proto"
)

func (k Keeper) GetUpgradePlan(ctx sdk.Context) (*types.Plan, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.KeyUpgradePlan) {
		return nil, nil
	}

	plan := types.Plan{}
	bz := store.Get(types.KeyUpgradePlan)
	err := proto.Unmarshal(bz, &plan)
	return &plan, err
}

func (k Keeper) SaveUpgradePlan(ctx sdk.Context, plan types.Plan) {
	store := ctx.KVStore(k.storeKey)
	bz, err := proto.Marshal(&plan)
	if err != nil {
		panic(err)
	}
	store.Set(types.KeyUpgradePlan, bz)
}

func (k Keeper) ApplyUpgradePlan(ctx sdk.Context, plan types.Plan) {
	// TODO: how do we check that upgrade is already done and not halt any more?
	if plan.ShouldExecute(ctx) {
		panic(fmt.Sprintf("UPGRADE \"%s\" NEEDED at %s", plan.Name, time.Unix(plan.MinHaltTime, 0).String()))
	}
}
