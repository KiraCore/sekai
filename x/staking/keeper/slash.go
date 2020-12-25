package keeper

import (
	customstakingtypes "github.com/KiraCore/sekai/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Activate a validator
func (k Keeper) Activate(ctx sdk.Context, valAddress sdk.ValAddress) error {
	validator, err := k.GetValidator(ctx, valAddress)
	if err != nil {
		return err
	}

	if validator.IsPaused() {
		return customstakingtypes.ErrValidatorPaused
	}

	validator.Status = customstakingtypes.Active
	k.AddValidator(ctx, validator)

	return nil
}

// Inactivate inactivate the validator
func (k Keeper) Inactivate(ctx sdk.Context, valAddress sdk.ValAddress) error { // inactivate a validator
	validator, err := k.GetValidator(ctx, valAddress)
	if err != nil {
		return err
	}

	if validator.IsPaused() {
		return customstakingtypes.ErrValidatorPaused
	}

	validator.Status = customstakingtypes.Inactive
	k.AddValidator(ctx, validator)

	return nil
}

// Pause a validator
func (k Keeper) Pause(ctx sdk.Context, valAddress sdk.ValAddress) error {
	validator, err := k.GetValidator(ctx, valAddress)
	if err != nil {
		return err
	}

	if validator.IsInactivated() {
		return customstakingtypes.ErrValidatorInactive
	}

	validator.Status = customstakingtypes.Paused
	k.AddValidator(ctx, validator)

	return nil
}

// Unpause unpause the validator
func (k Keeper) Unpause(ctx sdk.Context, valAddress sdk.ValAddress) error { // inactivate a validator
	validator, err := k.GetValidator(ctx, valAddress)
	if err != nil {
		return err
	}

	if validator.IsInactivated() {
		return customstakingtypes.ErrValidatorInactive
	}

	validator.Status = customstakingtypes.Active
	k.AddValidator(ctx, validator)

	return nil
}

// TODO: should take care of relation between Activate / Pause
// Inactivate is not possible if it's paused
// Activate is not possible if it's paused
// Pause is not possible if it's inactivated
// Unpause is not possible if it's inactivated
// Paused / Inactivated validator shouldn't participate in block generation (Previous Jailed = Paused | Inactivated)
