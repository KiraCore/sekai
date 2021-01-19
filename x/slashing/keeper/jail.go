package keeper

import (
	"github.com/KiraCore/sekai/x/slashing/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Jail attempts to jail a validator. The slash is delegated to the staking module
// to make the necessary validator changes.
func (k Keeper) Jail(ctx sdk.Context, consAddr sdk.ConsAddress) {
	validator, err := k.sk.GetValidatorByConsAddr(ctx, consAddr)
	if err == nil && !validator.IsInactivated() {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeInactivate,
				sdk.NewAttribute(types.AttributeKeyInactivated, consAddr.String()),
			),
		)

		// TODO: should be modified to k.sk.Jail() function call
		k.sk.Inactivate(ctx, validator.ValKey)
	}
}
