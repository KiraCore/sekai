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

	// update validator info to active status - exception will happen for jail case
	err = k.sk.Activate(ctx, validator.ValKey)
	if err != nil {
		return err
	}

	// Update validator signing info to restart uptime counters
	signInfo, found := k.GetValidatorSigningInfo(ctx, consAddr)
	if found {
		// cannot be activated until out of inactive period finish
		if ctx.BlockTime().Before(signInfo.InactiveUntil) {
			duration := signInfo.InactiveUntil.Sub(ctx.BlockTime())
			return sdkerrors.Wrap(types.ErrValidatorInactivated, fmt.Sprintf("Can NOT activate inactivate validator, jail time remaining %d seconds", duration/time.Second))
		}

		// automatically set the mischance to 0 and last_present_block to latest_block_height
		signInfo.Mischance = 0
		signInfo.MischanceConfidence = 0

		k.SetValidatorSigningInfo(ctx, consAddr, signInfo)
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
	if validator.IsJailed() {
		return sdkerrors.Wrap(types.ErrValidatorJailed, "Can NOT pause jailed validator")
	}

	// cannot be paused if not paused already
	if validator.IsInactivated() {
		return sdkerrors.Wrap(types.ErrValidatorInactivated, "Can NOT pause inactivated validator")
	}

	// cannot be paused if not paused already
	if validator.IsPaused() {
		return sdkerrors.Wrap(types.ErrValidatorPaused, "Can NOT pause already paused validator")
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
