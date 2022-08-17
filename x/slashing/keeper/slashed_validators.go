package keeper

import (
	"time"

	"github.com/KiraCore/sekai/x/slashing/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetSlashedValidator(ctx sdk.Context, val sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.SlashedValidatorByTimeKey(ctx.BlockTime(), val), val)
}

func (k Keeper) GetSlashedValidatorsAfter(ctx sdk.Context, timestamp time.Time) []sdk.ValAddress {
	store := ctx.KVStore(k.storeKey)
	startKey := types.SlashedValidatorByTimePrefixKey(timestamp)
	it := store.Iterator(startKey, sdk.PrefixEndBytes(types.SlashedValidatorsByTimeKeyPrefix))
	defer it.Close()

	validators := []sdk.ValAddress{}
	for ; it.Valid(); it.Next() {
		validators = append(validators, sdk.ValAddress(it.Value()))
	}
	return validators
}
