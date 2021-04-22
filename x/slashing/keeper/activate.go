package keeper

import (
	"fmt"
	"time"

	"github.com/KiraCore/sekai/x/slashing/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Activate calls the staking Activate function to activate a validator if the
// inactivated period has concluded
func (k Keeper) Activate(ctx sdk.Context, validatorAddr sdk.ValAddress) error {
	validator, err := k.sk.GetValidator(ctx, validatorAddr)
	if err != nil {
		return types.ErrNoValidatorForAddress
	}

	// cannot be activated if not inactivated
	if !validator.IsInactivated() {
		return sdkerrors.Wrap(types.ErrValidatorNotInactivated, "Can NOT activate NOT inactivated validator")
	}

	consAddr := validator.GetConsAddr()

	// If the validator has a ValidatorSigningInfo object that signals that the
	// validator was joined and so we must check that the validator is not tombstoned
	// and can be activated at the current block.
	//
	// A validator that is inactivated but has no ValidatorSigningInfo object signals
	// that the validator was never joined and must've been inactivated due to falling
	// below their minimum self-delegation. The validator can activate at any point
	// assuming they've now joined above their minimum self-delegation.
	info, found := k.GetValidatorSigningInfo(ctx, consAddr)
	if found {
		// cannot be activated if tombstoned
		if info.Tombstoned {
			return sdkerrors.Wrap(types.ErrValidatorInactivated, "Can NOT activate tombstoned validator, governance proposal required")
		}

		// cannot be activated until out of inactive period finish
		if ctx.BlockTime().Before(info.InactiveUntil) {
			duration := info.InactiveUntil.Sub(ctx.BlockTime())
			return sdkerrors.Wrap(types.ErrValidatorInactivated, fmt.Sprintf("Can NOT activate inactivate validator, jail time remaining %d seconds", duration/time.Second))
		}

		// automatically set the mischance to 0 and last_present_block to latest_block_height
		info.Mischance = 0
		info.LastPresentBlock = ctx.BlockHeight()
		k.SetValidatorSigningInfo(ctx, consAddr, info)
	}

	err = k.sk.Activate(ctx, validator.ValKey)
	if err != nil {
		return err
	}

	return nil
}

// Pause calls the staking Pause function to pause a validator if validator is not paused / inactivated
func (k Keeper) Pause(ctx sdk.Context, validatorAddr sdk.ValAddress) error {
	validator, err := k.sk.GetValidator(ctx, validatorAddr)
	if err != nil {
		return types.ErrNoValidatorForAddress
	}

	// cannot be paused if not paused already
	if validator.IsPaused() {
		return sdkerrors.Wrap(types.ErrValidatorPaused, "Can NOT pause already paused validator")
	}

	// cannot be paused if not paused already
	if validator.IsInactivated() {
		return sdkerrors.Wrap(types.ErrValidatorInactivated, "Can NOT pause inactivated validator")
	}

	k.sk.Pause(ctx, validator.ValKey)
	return nil
}

// Unpause calls the staking Unpause function to unpause a validator if validator is paused
func (k Keeper) Unpause(ctx sdk.Context, validatorAddr sdk.ValAddress) error {
	validator, err := k.sk.GetValidator(ctx, validatorAddr)
	if err != nil {
		return types.ErrNoValidatorForAddress
	}

	// cannot be unpaused if not paused
	if !validator.IsPaused() {
		return sdkerrors.Wrap(types.ErrValidatorNotPaused, "Can NOT pause inactivated validator")
	}

	k.sk.Unpause(ctx, validator.ValKey)
	return nil
}
