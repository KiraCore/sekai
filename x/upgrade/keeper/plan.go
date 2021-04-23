package keeper

import (
	"fmt"
	"github.com/KiraCore/sekai/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) ApplyUpgradePlan(ctx sdk.Context, plan types.Plan) {
	if ctx.BlockHeight() > plan.MinHaltTime {
		panic(fmt.Sprintf("UPGRADE \"%s\" NEEDED at %s", plan.Name, plan.MinHaltTime))
	}
}
