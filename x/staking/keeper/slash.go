package keeper

import (
	"github.com/KiraCore/sekai/x/staking/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Activate a validator
func (k Keeper) Activate(ctx sdk.Context, valAddress sdk.ValAddress) error {
	validator, err := k.GetValidator(ctx, valAddress)
	if err != nil {
		return err
	}

	if validator.IsPaused() {
		return sdkerrors.Wrap(stakingtypes.ErrValidatorPaused, "Can NOT activate paused validator, you must unpause")
	}

	if validator.IsJailed() {
		return sdkerrors.Wrap(stakingtypes.ErrValidatorJailed, "Can NOT activate jailed validator, you must unjail via proposal")
	}

	if validator.IsActive() {
		return sdkerrors.Wrap(stakingtypes.ErrValidatorActive, "Can NOT activate already active validator")
	}

	k.setStatusToValidator(ctx, validator, stakingtypes.Active)
	k.addReactivatingValidator(ctx, validator)
	k.RemoveRemovingValidator(ctx, validator)

	return nil
}

// ResetWholeValidatorRank reset whole validators' status, rank and streak
func (k Keeper) ResetWholeValidatorRank(ctx sdk.Context) {
	// TODO: is it correct to use this iterator @Jonathan?
	k.IterateValidators(ctx, func(index int64, validator *types.Validator) (stop bool) {
		validator.Status = stakingtypes.Active
		validator.Rank = 0
		validator.Streak = 0
		k.AddValidator(ctx, *validator)
		return false
	})
}

// Inactivate inactivate the validator
func (k Keeper) Inactivate(ctx sdk.Context, valAddress sdk.ValAddress) error { // inactivate a validator
	validator, err := k.GetValidator(ctx, valAddress)
	if err != nil {
		return err
	}

	if validator.IsPaused() {
		return stakingtypes.ErrValidatorPaused
	}

	networkProperties := k.govkeeper.GetNetworkProperties(ctx)
	validator.Status = stakingtypes.Inactive
	validator.Rank = validator.Rank * int64(100-networkProperties.InactiveRankDecreasePercent) / 100
	validator.Streak = 0

	k.AddValidator(ctx, validator)
	k.addRemovingValidator(ctx, validator)
	k.RemoveReactivatingValidator(ctx, validator)

	return nil
}

// HandleValidatorSignature manage rank and streak by block miss / sign result
func (k Keeper) HandleValidatorSignature(ctx sdk.Context, valAddress sdk.ValAddress, missed bool, mischance int64) error {
	validator, err := k.GetValidator(ctx, valAddress)
	if err != nil {
		return err
	}
	networkProperties := k.govkeeper.GetNetworkProperties(ctx)
	if missed {
		if mischance > 0 { // it means mischance confidence is set, we update streak and rank properties
			// set validator streak by 0 and decrease rank by X
			validator.Streak = 0
			validator.Rank -= int64(networkProperties.MischanceRankDecreaseAmount)
			if validator.Rank < 0 {
				validator.Rank = 0
			}
		}
	} else {
		// increase streak and reset rank if streak is higher than rank
		validator.Streak++
		if validator.Streak > validator.Rank {
			validator.Rank = validator.Streak
		}
	}
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
		return stakingtypes.ErrValidatorInactive
	}

	k.setStatusToValidator(ctx, validator, stakingtypes.Paused)
	k.addRemovingValidator(ctx, validator)
	k.RemoveReactivatingValidator(ctx, validator)

	return nil
}

// Unpause unpause the validator
func (k Keeper) Unpause(ctx sdk.Context, valAddress sdk.ValAddress) error { // inactivate a validator
	validator, err := k.GetValidator(ctx, valAddress)
	if err != nil {
		return err
	}

	if validator.IsInactivated() {
		return stakingtypes.ErrValidatorInactive
	}

	k.setStatusToValidator(ctx, validator, stakingtypes.Active)
	k.addReactivatingValidator(ctx, validator)
	k.RemoveRemovingValidator(ctx, validator)

	return nil
}

// Jail a validator
func (k Keeper) Jail(ctx sdk.Context, valAddress sdk.ValAddress) error {
	validator, err := k.GetValidator(ctx, valAddress)
	if err != nil {
		return err
	}

	k.setStatusToValidator(ctx, validator, stakingtypes.Jailed)
	k.addRemovingValidator(ctx, validator)
	k.setJailValidatorInfo(ctx, validator)
	k.RemoveReactivatingValidator(ctx, validator)

	return nil
}

// GetValidatorJailInfo returns information about a jailed validor, found is false
// if there is no validator, so a validator that is not jailed should return false.
func (k Keeper) GetValidatorJailInfo(ctx sdk.Context, valAddress sdk.ValAddress) (stakingtypes.ValidatorJailInfo, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetValidatorJailInfoKey(valAddress))
	if bz == nil {
		return stakingtypes.ValidatorJailInfo{}, false
	}

	var info stakingtypes.ValidatorJailInfo
	k.cdc.MustUnmarshal(bz, &info)

	return info, true
}

// TODO: add proposal to jail a validator by governance
// TODO: Would be good to merge signInfo with validator to reduce complexity
// TODO: Same evidence shouldn't be used again for another jail - is it already supported?
// When unjail proposal come, all the evidence info should be disregarded at the time - that's older.

// Unjail a validator
func (k Keeper) Unjail(ctx sdk.Context, valAddress sdk.ValAddress) error {
	validator, err := k.GetValidator(ctx, valAddress)
	if err != nil {
		return err
	}

	// Unjail move validator status to Inactive from Jailed
	// User can activate it after counting activation time - uptime counter is reset at that time
	k.setStatusToValidator(ctx, validator, stakingtypes.Inactive)
	k.removeJailValidatorInfo(ctx, validator)

	return nil
}

func (k Keeper) setStatusToValidator(ctx sdk.Context, validator stakingtypes.Validator, status stakingtypes.ValidatorStatus) {
	validator.Status = status
	k.AddValidator(ctx, validator)
}

func (k Keeper) setJailValidatorInfo(ctx sdk.Context, validator stakingtypes.Validator) {
	jailInfo := stakingtypes.ValidatorJailInfo{
		Time: ctx.BlockTime(),
	}

	bz := k.cdc.MustMarshal(jailInfo)

	store := ctx.KVStore(k.storeKey)
	store.Set(GetValidatorJailInfoKey(validator.ValKey), bz)
}

func (k Keeper) removeJailValidatorInfo(ctx sdk.Context, validator stakingtypes.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetValidatorJailInfoKey(validator.ValKey))
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

func (k Keeper) RemoveReactivatingValidator(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetReactivatingValidatorKey(validator.ValKey))
}
