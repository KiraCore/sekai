package upgrade

import (
	"fmt"
	"time"

	"github.com/KiraCore/sekai/x/upgrade/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlock will check if there is a scheduled plan and if it is ready to be executed.
// If it is ready, it will execute it if the handler is installed, and panic/abort otherwise.
// If the plan is not ready, it will ensure the handler is not registered too early (and abort otherwise).
//
// The purpose is to ensure the binary is switched EXACTLY at the desired block, and to allow
// a migration to be executed if needed upon this switch (migration defined in the new binary)
func BeginBlocker(k keeper.Keeper, ctx sdk.Context, _ abci.RequestBeginBlock) {
	plan, err := k.GetNextPlan(ctx)
	if err != nil || plan == nil {
		return
	}

	// check if it's time to halt and halt the chain
	if plan.ShouldExecute(ctx) {
		ctx.Logger().Info(fmt.Sprintf("Applying the upgrade for \"%s\" at %s", plan.Name, time.Unix(plan.UpgradeTime, 0).String()))
		ctx = ctx.WithBlockGasMeter(sdk.NewInfiniteGasMeter())
		k.ApplyUpgradePlan(ctx, *plan)
		return
	}
}
