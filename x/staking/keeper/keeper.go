package keeper

import (
	"github.com/KiraCore/cosmos-sdk/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/x/staking/types"
)

// Keeper represents the keeper that maintains the Validator Registry.
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

// NewKeeper returns new keeper.
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{storeKey: storeKey, cdc: cdc}
}

func (k Keeper) AddValidator(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(validator)
	store.Set(types.GetValidatorKey(validator.ValKey), bz)
}

func (k Keeper) GetValidator(ctx sdk.Context, address sdk.ValAddress) types.Validator {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorKey(address))

	var validator types.Validator
	k.cdc.MustUnmarshalBinaryBare(bz, &validator)

	return validator
}

func (k Keeper) GetValidatorSet(ctx sdk.Context) []types.Validator {
	store := ctx.KVStore(k.storeKey)

	iter := store.Iterator(types.ValidatorsKey, nil)

	var validators []types.Validator
	for ; iter.Valid(); iter.Next() {
		var validator types.Validator
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &validator)
		validators = append(validators, validator)
	}

	return validators
}
