package keeper

import (
	"github.com/KiraCore/sekai/x/recovery/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) GetRecoveryRecord(ctx sdk.Context, address string) (types.RecoveryRecord, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RecoveryRecordKey(address))
	if bz == nil {
		return types.RecoveryRecord{}, sdkerrors.Wrapf(types.ErrRecoveryRecordDoesNotExist, "RecoveryRecord: %s does not exist", address)
	}
	record := types.RecoveryRecord{}
	k.cdc.MustUnmarshal(bz, &record)
	return record, nil
}

func (k Keeper) SetRecoveryRecord(ctx sdk.Context, record types.RecoveryRecord) {
	bz := k.cdc.MustMarshal(&record)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.RecoveryRecordKey(record.Address), bz)
}

func (k Keeper) DeleteRecoveryRecord(ctx sdk.Context, record types.RecoveryRecord) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.RecoveryRecordKey(record.Address))
}

func (k Keeper) GetAllRecoveryRecords(ctx sdk.Context) []types.RecoveryRecord {
	store := ctx.KVStore(k.storeKey)

	records := []types.RecoveryRecord{}
	it := sdk.KVStorePrefixIterator(store, types.RecoveryRecordKeyPrefix)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		record := types.RecoveryRecord{}
		k.cdc.MustUnmarshal(it.Value(), &record)
		records = append(records, record)
	}
	return records
}

func (k Keeper) GetRecoveryToken(ctx sdk.Context, address string) (types.RecoveryToken, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RecoveryTokenKey(address))
	if bz == nil {
		return types.RecoveryToken{}, sdkerrors.Wrapf(types.ErrRecoveryTokenDoesNotExist, "RecoveryToken: %s does not exist", address)
	}
	recovery := types.RecoveryToken{}
	k.cdc.MustUnmarshal(bz, &recovery)
	return recovery, nil
}

func (k Keeper) SetRecoveryToken(ctx sdk.Context, recovery types.RecoveryToken) {
	bz := k.cdc.MustMarshal(&recovery)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.RecoveryTokenKey(recovery.Address), bz)
}

func (k Keeper) DeleteRecoveryToken(ctx sdk.Context, recovery types.RecoveryToken) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.RecoveryTokenKey(recovery.Address))
}

func (k Keeper) GetAllRecoveryTokens(ctx sdk.Context) []types.RecoveryToken {
	store := ctx.KVStore(k.storeKey)

	recoveries := []types.RecoveryToken{}
	it := sdk.KVStorePrefixIterator(store, types.RecoveryRecordKeyPrefix)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		recovery := types.RecoveryToken{}
		k.cdc.MustUnmarshal(it.Value(), &recovery)
		recoveries = append(recoveries, recovery)
	}
	return recoveries
}
