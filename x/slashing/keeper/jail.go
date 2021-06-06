package keeper

import (
	"github.com/KiraCore/sekai/x/slashing/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Jail attempts to jail a validator. The slash is delegated to the staking module
// to make the necessary validator changes.
func (k Keeper) Jail(ctx sdk.Context, consAddr sdk.ConsAddress) {
	validator, err := k.sk.GetValidatorByConsAddr(ctx, consAddr)
	if err == nil && !validator.IsJailed() {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeJail,
				sdk.NewAttribute(types.AttributeKeyJailed, consAddr.String()),
			),
		)

		k.sk.Jail(ctx, validator.ValKey)
	}
}
