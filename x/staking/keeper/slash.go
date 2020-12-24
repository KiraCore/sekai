package keeper

import (
	customstakingtypes "github.com/KiraCore/sekai/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Activate a validator
func (k Keeper) Activate(ctx sdk.Context, consAddr sdk.ConsAddress) {
	// TODO: don't do anything for now, implement this
	// validator := k.GetValidatorByConsAddr(ctx, consAddr)
	// if !validator.Inactivated {
	// 	panic(fmt.Sprintf("cannot unjail already inactivated validator, validator: %v\n", validator))
	// }

	// validator.Inactivated = false
	// k.SetValidator(ctx, validator)
}

// Inactivate inactivate the validator
func (k Keeper) Inactivate(sdk.Context, sdk.ConsAddress) { // inactivate a validator
	// TODO: don't do anything for now, implement this
}

// Pause a validator
func (k Keeper) Pause(ctx sdk.Context, valAddress sdk.ValAddress) error {
	validator, err := k.GetValidator(ctx, valAddress)
	if err != nil {
		return err
	}

	validator.Status = customstakingtypes.Paused
	k.AddValidator(ctx, validator)

	return nil
}

// Unpause unpause the validator
func (k Keeper) Unpause(sdk.Context, sdk.ValAddress) error { // inactivate a validator
	return nil
}

// TODO: should take care of relation between Activate / Pause
// Inactivate is not possible if it's paused
// Activate is not possible if it's paused
// Pause is not possible if it's inactivated
// Unpause is not possible if it's inactivated
// Paused / Inactivated validator shouldn't participate in block generation (Previous Jailed = Paused | Inactivated)
