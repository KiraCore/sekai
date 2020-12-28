package keeper

import (
	"github.com/KiraCore/sekai/x/staking/types"
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
	k.addReactivatingValidator(ctx, validator)

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
	k.addRemovingValidator(ctx, validator)

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
	k.addRemovingValidator(ctx, validator)

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
	k.addReactivatingValidator(ctx, validator)

	return nil
}

// GetRemovingValidatorSet returns the keys of the validators that needs to be removed from
// the set.
func (k Keeper) GetRemovingValidatorSet(ctx sdk.Context) [][]byte {
	store := ctx.KVStore(k.storeKey)

	iter := sdk.KVStorePrefixIterator(store, RemovingValidatorQueue)
	defer iter.Close()

	var validatorKeys [][]byte
	for ; iter.Valid(); iter.Next() {
		validatorKeys = append(validatorKeys, iter.Value())
	}

	return validatorKeys
}

// GetReactivatingValidatorSet returns the keys of the validators need to be reactivated, this
// is used in the Enblock function to reactivate those validators who have been unpaused
// or activated.
func (k Keeper) GetReactivatingValidatorSet(ctx sdk.Context) [][]byte {
	store := ctx.KVStore(k.storeKey)

	iter := sdk.KVStorePrefixIterator(store, ReactivatingValidatorQueue)
	defer iter.Close()

	var validatorKeys [][]byte
	for ; iter.Valid(); iter.Next() {
		validatorKeys = append(validatorKeys, iter.Value())
	}

	return validatorKeys
}

func (k Keeper) addRemovingValidator(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetRemovingValidatorKey(validator.ValKey), validator.ValKey)
}

func (k Keeper) RemoveRemovingValidator(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetRemovingValidatorKey(validator.ValKey))
}

func (k Keeper) addReactivatingValidator(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetReactivatingValidatorKey(validator.ValKey), validator.ValKey)
}
