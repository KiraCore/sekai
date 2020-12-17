package keeper

import (
	"github.com/KiraCore/sekai/x/slashing/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Activate calls the staking Activate function to activate a validator if the
// inactivated period has concluded
func (k Keeper) Activate(ctx sdk.Context, validatorAddr sdk.ValAddress) error {
	validator := k.sk.Validator(ctx, validatorAddr)
	if validator == nil {
		return types.ErrNoValidatorForAddress
	}

	// cannot be activated if not inactivated
	if !validator.IsInactivated() {
		return types.ErrValidatorNotInactivated
	}

	consAddr := validator.GetConsAddr()

	// If the validator has a ValidatorSigningInfo object that signals that the
	// validator was bonded and so we must check that the validator is not tombstoned
	// and can be activated at the current block.
	//
	// A validator that is inactivated but has no ValidatorSigningInfo object signals
	// that the validator was never bonded and must've been inactivated due to falling
	// below their minimum self-delegation. The validator can activate at any point
	// assuming they've now bonded above their minimum self-delegation.
	info, found := k.GetValidatorSigningInfo(ctx, consAddr)
	if found {
		// cannot be activated if tombstoned
		if info.Tombstoned {
			return types.ErrValidatorInactivated
		}

		// cannot be activated until out of inactivation
		if ctx.BlockHeader().Time.Before(info.InactiveUntil) {
			return types.ErrValidatorInactivated
		}
	}

	k.sk.Activate(ctx, consAddr)
	return nil
}
