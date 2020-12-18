package keeper

import (
	"fmt"

	"github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper represents the keeper that maintains the Validator Registry.
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.LegacyAmino
}

// NewKeeper returns new keeper.
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.LegacyAmino) Keeper {
	return Keeper{storeKey: storeKey, cdc: cdc}
}

// BondDenom returns the denom that is basically used for fee payment
func (k Keeper) BondDenom(ctx sdk.Context) string {
	return "ukex"
}

func (k Keeper) AddValidator(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&validator)
	store.Set(GetValidatorKey(validator.ValKey), bz)

	// Save by moniker
	store.Set(GetValidatorByMonikerKey(validator.Moniker), GetValidatorKey(validator.ValKey))
	k.AddValidatorByConsAddr(ctx, validator)
}

// validator index
func (k Keeper) AddValidatorByConsAddr(ctx sdk.Context, validator types.Validator) {
	consPk := validator.GetConsAddr()

	store := ctx.KVStore(k.storeKey)
	store.Set(GetValidatorByConsAddrKey(consPk), validator.ValKey)
}

func (k Keeper) AddPendingValidator(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&validator)
	store.Set(GetPendingValidatorKey(validator.ValKey), bz)
}

func (k Keeper) RemovePendingValidator(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetPendingValidatorKey(validator.ValKey))
}

func (k Keeper) GetValidator(ctx sdk.Context, address sdk.ValAddress) (types.Validator, error) {
	return k.getValidatorByKey(ctx, GetValidatorKey(address))
}

func (k Keeper) GetValidatorByAccAddress(ctx sdk.Context, address sdk.AccAddress) (types.Validator, error) {
	return k.getValidatorByKey(ctx, GetValidatorKeyAcc(address))
}

func (k Keeper) GetValidatorByMoniker(ctx sdk.Context, moniker string) (types.Validator, error) {
	store := ctx.KVStore(k.storeKey)

	valKey := store.Get(GetValidatorByMonikerKey(moniker))
	if valKey == nil {
		return types.Validator{}, fmt.Errorf("validator with moniker %s not found", moniker)
	}

	return k.getValidatorByKey(ctx, valKey)
}

func (k Keeper) getValidatorByKey(ctx sdk.Context, key []byte) (types.Validator, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)

	if bz == nil {
		return types.Validator{}, fmt.Errorf("validator not found")
	}

	var validator types.Validator
	k.cdc.MustUnmarshalBinaryBare(bz, &validator)

	return validator, nil
}

func (k Keeper) GetValidatorSet(ctx sdk.Context) []types.Validator {
	store := ctx.KVStore(k.storeKey)

	iter := sdk.KVStorePrefixIterator(store, ValidatorsKey)
	defer iter.Close()

	var validators []types.Validator
	for ; iter.Valid(); iter.Next() {
		var validator types.Validator
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &validator)
		validators = append(validators, validator)
	}

	return validators
}

func (k Keeper) GetPendingValidatorSet(ctx sdk.Context) []types.Validator {
	store := ctx.KVStore(k.storeKey)

	iter := sdk.KVStorePrefixIterator(store, PendingValidatorQueue)
	defer iter.Close()

	var validators []types.Validator
	for ; iter.Valid(); iter.Next() {
		var validator types.Validator
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &validator)
		validators = append(validators, validator)
	}

	return validators
}

// IterateValidators iterate through validators by operator address, execute func for each validator
func (k Keeper) IterateValidators(ctx sdk.Context,
	handler func(index int64, validator *types.Validator) (stop bool)) {
	validators := k.GetValidatorSet(ctx)
	for i, val := range validators {
		if handler(int64(i), &val) {
			break
		}
	}
}

// GetValidatorByConsAddr get validator by sdk.ConsAddress
func (k Keeper) GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (types.Validator, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(GetValidatorByConsAddrKey(consAddr))
	if bz == nil {
		return types.Validator{}, fmt.Errorf("validator not found")
	}

	validator, err := k.GetValidatorByAccAddress(ctx, bz)
	if err != nil {
		return types.Validator{}, err
	}

	return validator, nil
}

// MaxValidators returns the maximum amount of bonded validators
func (k Keeper) MaxValidators(sdk.Context) uint32 {
	// TODO: don't do anything for now, implement this
	return 10
}
