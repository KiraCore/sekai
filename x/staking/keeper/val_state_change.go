package keeper

import (
	"errors"

	"github.com/KiraCore/sekai/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BlockValidatorUpdates calculates the ValidatorUpdates for the current block
// Called in each EndBlock
func (k Keeper) BlockValidatorUpdates(ctx sdk.Context) []abci.ValidatorUpdate {
	validatorUpdates, err := k.ApplyAndReturnValidatorSetUpdates(ctx)
	if err != nil {
		panic(err)
	}
	return validatorUpdates
}

// ApplyAndReturnValidatorSetUpdates applies and return accumulated updates to the joined validator set.
func (k Keeper) ApplyAndReturnValidatorSetUpdates(ctx sdk.Context) (updates []abci.ValidatorUpdate, err error) {
	var valUpdate []abci.ValidatorUpdate

	valSet := k.GetPendingValidatorSet(ctx)
	for _, val := range valSet {
		k.AddValidator(ctx, val)
		k.AfterValidatorCreated(ctx, val.ValKey)

		val, err = k.joinValidator(ctx, val)
		if err != nil {
			return nil, err
		}

		consPk, err := val.TmConsPubKey()
		if err != nil {
			return nil, err
		}

		valUpdate = append(valUpdate, abci.ValidatorUpdate{
			Power:  1,
			PubKey: consPk,
		})

		k.RemovePendingValidator(ctx, val)
	}

	// Remove validators from the set, paused or inactivated.
	removeVals := k.GetRemovingValidatorSet(ctx)
	for _, val := range removeVals {
		validator, err := k.GetValidator(ctx, val)
		if err != nil {
			return nil, errors.New("validator not found")
		}

		consPk, err := validator.TmConsPubKey()
		if err != nil {
			return nil, err
		}

		valUpdate = append(valUpdate, abci.ValidatorUpdate{
			Power:  0,
			PubKey: consPk,
		})
		k.RemoveRemovingValidator(ctx, validator)
	}

	// Remove validators from the set, paused or inactivated.
	reactivateVals := k.GetReactivatingValidatorSet(ctx)
	for _, val := range reactivateVals {
		validator, err := k.GetValidator(ctx, val)
		if err != nil {
			return nil, errors.New("validator not found")
		}

		consPk, err := validator.TmConsPubKey()
		if err != nil {
			return nil, err
		}

		valUpdate = append(valUpdate, abci.ValidatorUpdate{
			Power:  1,
			PubKey: consPk,
		})
		k.RemoveReactivatingValidator(ctx, validator)
	}

	return valUpdate, nil
}

// perform all the store operations for when a validator status becomes joined
func (k Keeper) joinValidator(ctx sdk.Context, validator types.Validator) (types.Validator, error) {

	// trigger hook
	consAddr := validator.GetConsAddr()
	k.AfterValidatorJoined(ctx, consAddr, validator.ValKey)

	return validator, nil
}
