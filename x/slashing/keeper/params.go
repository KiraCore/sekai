package keeper

import (
	"time"

	"github.com/KiraCore/sekai/x/slashing/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MaxMischance - minimum blocks signed per window
func (k Keeper) MaxMischance(ctx sdk.Context) int64 {
	var maxMischance int64
	k.paramspace.Get(ctx, types.KeyMaxMischance, &maxMischance)
	return maxMischance
}

// DowntimeInactiveDuration - Downtime unbond duration
func (k Keeper) DowntimeInactiveDuration(ctx sdk.Context) (res time.Duration) {
	k.paramspace.Get(ctx, types.KeyDowntimeInactiveDuration, &res)
	return
}

// GetParams returns the total set of slashing parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramspace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the slashing parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramspace.SetParamSet(ctx, &params)
}
