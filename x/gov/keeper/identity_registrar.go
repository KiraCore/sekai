package keeper

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetLastIdentityRecordId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyLastIdentityRecordId)
	if bz == nil {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLastIdentityRecordId(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyLastIdentityRecordId, sdk.Uint64ToBigEndian(id))
}

func (k Keeper) GetLastIdRecordVerifyRequestId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyLastIdRecordVerifyRequestId)
	if bz == nil {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

func CheckIfWithinAddressArray(addr sdk.AccAddress, array []sdk.AccAddress) bool {
	for _, itemAddr := range array {
		if bytes.Equal(addr, itemAddr) {
			return true
		}
	}
	return false
}

func (k Keeper) SetLastIdRecordVerifyRequestId(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyLastIdRecordVerifyRequestId, sdk.Uint64ToBigEndian(id))
}

// SetIdentityRecord defines a method to set identity record
func (k Keeper) SetIdentityRecord(ctx sdk.Context, record types.IdentityRecord) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixIdentityRecord)
	bz := k.cdc.MustMarshalBinaryBare(&record)
	prefixStore.Set(sdk.Uint64ToBigEndian(record.Id), bz)
}

func (k Keeper) GetIdentityRecord(ctx sdk.Context, recordId uint64) *types.IdentityRecord {
	record := types.IdentityRecord{}

	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixIdentityRecord)
	bz := prefixStore.Get(sdk.Uint64ToBigEndian(recordId))
	if bz == nil {
		return nil
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &record)
	return &record
}

// DeleteIdentityRecord defines a method to delete identity record by id
func (k Keeper) DeleteIdentityRecord(ctx sdk.Context, recordId uint64) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixIdentityRecord)
	prefixStore.Delete(sdk.Uint64ToBigEndian(recordId))
}

// CreateIdentityRecord defines a method to create identity record
func (k Keeper) CreateIdentityRecord(ctx sdk.Context, address sdk.AccAddress, infos map[string]string, date time.Time) uint64 {
	recordId := k.GetLastIdentityRecordId(ctx) + 1

	k.SetIdentityRecord(ctx, types.IdentityRecord{
		Id:        recordId,
		Address:   address,
		Infos:     infos,
		Date:      date,
		Verifiers: []sdk.AccAddress{},
	})

	k.SetLastIdentityRecordId(ctx, recordId)
	return recordId
}

// EditIdentityRecord defines a method to edit identity record, it removes all verifiers for the record
func (k Keeper) EditIdentityRecord(ctx sdk.Context, recordId uint64, address sdk.AccAddress, infos map[string]string, date time.Time) error {
	prevRecord := k.GetIdentityRecord(ctx, recordId)
	if prevRecord == nil {
		return fmt.Errorf("identity record with specified id does NOT exist: id=%d", recordId)
	}

	k.SetIdentityRecord(ctx, types.IdentityRecord{
		Id:        recordId,
		Address:   address,
		Infos:     infos,
		Date:      date,
		Verifiers: []sdk.AccAddress{},
	})

	return nil
}

// GetAllIdentityRecords query all identity records
func (k Keeper) GetAllIdentityRecords(ctx sdk.Context) []types.IdentityRecord {
	records := []types.IdentityRecord{}
	// get iterator for token aliases
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixIdentityRecord)

	for ; iterator.Valid(); iterator.Next() {
		record := types.IdentityRecord{}
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &record)
		records = append(records, record)
	}

	return records
}

// RequestIdentityRecordsVerify defines a method to request verify request from specific verifier
func (k Keeper) RequestIdentityRecordsVerify(ctx sdk.Context, address, verifier sdk.AccAddress, recordIds []uint64, tip sdk.Coin) (uint64, error) {
	requestId := k.GetLastIdRecordVerifyRequestId(ctx) + 1
	request := types.IdentityRecordsVerify{
		Id:        requestId,
		Address:   address,
		Verifier:  verifier,
		RecordIds: recordIds,
		Tip:       tip,
	}
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixIdRecordVerifyRequest)
	bz := k.cdc.MustMarshalBinaryBare(&request)
	prefixStore.Set(sdk.Uint64ToBigEndian(requestId), bz)
	k.SetLastIdRecordVerifyRequestId(ctx, requestId)

	if err := k.bk.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, sdk.Coins{tip}); err != nil {
		return 0, err
	}
	return requestId, nil
}

// GetIdRecordsVerifyRequest defines a method to get an identity records verify request by id
func (k Keeper) GetIdRecordsVerifyRequest(ctx sdk.Context, requestId uint64) *types.IdentityRecordsVerify {
	request := types.IdentityRecordsVerify{}

	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixIdRecordVerifyRequest)
	bz := prefixStore.Get(sdk.Uint64ToBigEndian(requestId))
	if bz == nil {
		return nil
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &request)
	return &request
}

// DeleteIdRecordsVerifyRequest defines a method to get an identity records verify request by id
func (k Keeper) DeleteIdRecordsVerifyRequest(ctx sdk.Context, requestId uint64) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixIdRecordVerifyRequest)
	prefixStore.Delete(sdk.Uint64ToBigEndian(requestId))
}

// ApproveIdentityRecords defines a method to accept verification request
func (k Keeper) ApproveIdentityRecords(ctx sdk.Context, verifier sdk.AccAddress, requestId uint64) error {
	request := k.GetIdRecordsVerifyRequest(ctx, requestId)
	if request == nil {
		return fmt.Errorf("specified identity record verify request does NOT exist: id=%d", requestId)
	}
	if !bytes.Equal(verifier, request.Verifier) {
		return errors.New("verifier does not match with requested")
	}

	for _, recordId := range request.RecordIds {
		record := k.GetIdentityRecord(ctx, recordId)
		if record == nil {
			return fmt.Errorf("identity record with specified id does NOT exist: id=%d", recordId)
		}
		// if already exist, skip
		if CheckIfWithinAddressArray(verifier, record.Verifiers) {
			continue
		}
		record.Verifiers = append(record.Verifiers, verifier)
		k.SetIdentityRecord(ctx, *record)
	}

	if err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, verifier, sdk.Coins{request.Tip}); err != nil {
		return err
	}

	return nil
}

// CancelIdentityRecordsVerifyRequest defines a method to cancel verification request
func (k Keeper) CancelIdentityRecordsVerifyRequest(ctx sdk.Context, executor sdk.AccAddress, requestId uint64) error {
	request := k.GetIdRecordsVerifyRequest(ctx, requestId)
	if request == nil {
		return fmt.Errorf("specified identity record verify request does NOT exist: id=%d", requestId)
	}
	if !bytes.Equal(executor, request.Address) {
		return errors.New("executor is not identity record creator")
	}

	k.DeleteIdRecordsVerifyRequest(ctx, requestId)
	return nil
}
