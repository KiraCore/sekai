package keeper

import (
	"fmt"

	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper represents the keeper that maintains the Validator Registry.
type Keeper struct {
	storeKey  sdk.StoreKey
	cdc       *codec.LegacyAmino
	hooks     types.StakingHooks
	govkeeper types.GovKeeper
}

// NewKeeper returns new keeper.
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.LegacyAmino, govkeeper types.GovKeeper) Keeper {
	return Keeper{storeKey: storeKey, cdc: cdc, govkeeper: govkeeper}
}

// BondDenom returns the denom that is basically used for fee payment
func (k Keeper) BondDenom(ctx sdk.Context) string {
	return "ukex"
}

// Set the validator hooks
func (k *Keeper) SetHooks(sh types.StakingHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set validator hooks twice")
	}

	k.hooks = sh

	return k
}

func (k Keeper) AddValidator(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&validator)
	store.Set(GetValidatorKey(validator.ValKey), bz)

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
	bz := k.cdc.MustMarshal(&validator)
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
	addrs := k.govkeeper.GetAddressesByIdRecordKey(ctx, "moniker", moniker)
	if len(addrs) != 1 {
		return types.Validator{}, fmt.Errorf("validator with moniker %s not found", moniker)
	}

	return k.GetValidator(ctx, sdk.ValAddress(addrs[0]))
}

func (k Keeper) getValidatorByKey(ctx sdk.Context, key []byte) (types.Validator, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)

	if bz == nil {
		return types.Validator{}, fmt.Errorf("validator not found")
	}

	var validator types.Validator
	k.cdc.MustUnmarshal(bz, &validator)

	return validator, nil
}

func (k Keeper) GetValidatorSet(ctx sdk.Context) []types.Validator {
	store := ctx.KVStore(k.storeKey)

	iter := sdk.KVStorePrefixIterator(store, ValidatorsKey)
	defer iter.Close()

	var validators []types.Validator
	for ; iter.Valid(); iter.Next() {
		var validator types.Validator
		k.cdc.MustUnmarshal(iter.Value(), &validator)
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
		k.cdc.MustUnmarshal(iter.Value(), &validator)
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

// iterate through the active validator set and perform the provided function
func (k Keeper) IterateLastValidators(ctx sdk.Context, fn func(index int64, validator types.Validator) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, LastValidatorPowerKey)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		address := AddressFromLastValidatorPowerKey(iterator.Key())

		validator, err := k.GetValidator(ctx, address)
		if err != nil {
			panic(fmt.Sprintf("validator record not found for address: %v\n", address))
		}

		stop := fn(i, validator) // XXX is this safe will the validator unexposed fields be able to get written to?
		if stop {
			break
		}
		i++
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

// MaxValidators returns the maximum number of joined validators
func (k Keeper) MaxValidators(sdk.Context) uint32 {
	// TODO: this needs to be calculated dynamically by code by looking around validator iterator?
	// This number was discussed with @asmo
	return 1000
}

// IsNetworkActive returns true if network has more than the validators required in network property
func (k Keeper) IsNetworkActive(ctx sdk.Context) bool {
	vals := k.GetValidatorSet(ctx)
	return len(vals) >= int(k.govkeeper.GetNetworkProperties(ctx).MinValidators)
}

func AddressFromLastValidatorPowerKey(key []byte) []byte {
	return key[2:] // remove prefix bytes and address length
}

// GetIdRecordsByAddress query identity records by address
func (k Keeper) GetIdRecordsByAddress(ctx sdk.Context, creator sdk.AccAddress) []govtypes.IdentityRecord {
	return k.govkeeper.GetIdRecordsByAddress(ctx, creator)
}

func (k Keeper) GetMonikerByAddress(ctx sdk.Context, addr sdk.AccAddress) (string, error) {
	records, err := k.govkeeper.GetIdRecordsByAddressAndKeys(ctx, addr, []string{"moniker"})
	if err != nil {
		return "", err
	}
	if len(records) != 1 {
		return "", fmt.Errorf("failed fetching the field moniker from identity registrar for address=%s", addr.String())
	}
	return records[0].Value, nil
}
