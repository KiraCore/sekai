package keeper

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CheckIfWithinAddressArray(addr sdk.AccAddress, array []sdk.AccAddress) bool {
	for _, itemAddr := range array {
		if bytes.Equal(addr, itemAddr) {
			return true
		}
	}
	return false
}

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
	record := k.GetIdentityRecord(ctx, recordId)
	if record == nil {
		return
	}
	recordKey := sdk.Uint64ToBigEndian(recordId)
	prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixIdentityRecord).Delete(recordKey)
	prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixIdentityRecordByAddress).Delete(record.Address)
}

// CreateIdentityRecord defines a method to create identity record
func (k Keeper) CreateIdentityRecord(ctx sdk.Context, address sdk.AccAddress, infos []types.IdentityInfoEntry) (uint64, error) {
	if k.GetIdRecordByAddress(ctx, address) != nil {
		return 0, fmt.Errorf("identity record already registered for the address: %s", address.String())
	}
	recordId := k.GetLastIdentityRecordId(ctx) + 1
	k.SetIdentityRecord(ctx, types.IdentityRecord{
		Id:        recordId,
		Address:   address,
		Infos:     types.UnwrapInfos(infos),
		Date:      ctx.BlockTime(),
		Verifiers: []sdk.AccAddress{},
	})

	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.KeyPrefixIdentityRecordByAddress)
	prefixStore.Set(address, sdk.Uint64ToBigEndian(recordId))

	k.SetLastIdentityRecordId(ctx, recordId)
	return recordId, nil
}

// EditIdentityRecord defines a method to edit identity record, it removes all verifiers for the record
func (k Keeper) EditIdentityRecord(ctx sdk.Context, recordId uint64, address sdk.AccAddress, infos []types.IdentityInfoEntry) error {
	prevRecord := k.GetIdentityRecord(ctx, recordId)
	if prevRecord == nil {
		return fmt.Errorf("identity record with specified id does NOT exist: id=%d", recordId)
	}

	if !bytes.Equal(address, prevRecord.Address) {
		return fmt.Errorf("identity record is owned by someone else: %s", prevRecord.Address.String())
	}

	k.SetIdentityRecord(ctx, types.IdentityRecord{
		Id:        recordId,
		Address:   address,
		Infos:     types.UnwrapInfos(infos),
		Date:      ctx.BlockTime(),
		Verifiers: []sdk.AccAddress{},
	})

	return nil
}

// GetAllIdentityRecords query all identity records
func (k Keeper) GetAllIdentityRecords(ctx sdk.Context) []types.IdentityRecord {
	records := []types.IdentityRecord{}
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixIdentityRecord)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		record := types.IdentityRecord{}
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &record)
		records = append(records, record)
	}

	return records
}

// GetIdRecordByAddress query identity record by address
func (k Keeper) GetIdRecordByAddress(ctx sdk.Context, creator sdk.AccAddress) *types.IdentityRecord {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(append(types.KeyPrefixIdentityRecordByAddress, []byte(creator)...))
	if bz == nil {
		return nil
	}
	recordId := sdk.BigEndianToUint64(bz)
	return k.GetIdentityRecord(ctx, recordId)
}

// SetIdentityRecordsVerifyRequest saves identity verify request into the store
func (k Keeper) SetIdentityRecordsVerifyRequest(ctx sdk.Context, request types.IdentityRecordsVerify) {
	requestId := request.Id
	bz := k.cdc.MustMarshalBinaryBare(&request)
	prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixIdRecordVerifyRequest).Set(sdk.Uint64ToBigEndian(requestId), bz)
	prefix.NewStore(
		ctx.KVStore(k.storeKey),
		append(types.KeyPrefixIdRecordVerifyRequestByRequester, request.Address...),
	).Set(sdk.Uint64ToBigEndian(requestId), sdk.Uint64ToBigEndian(requestId))
	prefix.NewStore(
		ctx.KVStore(k.storeKey),
		append(types.KeyPrefixIdRecordVerifyRequestByApprover, request.Verifier...),
	).Set(sdk.Uint64ToBigEndian(requestId), sdk.Uint64ToBigEndian(requestId))
}

// RequestIdentityRecordsVerify defines a method to request verify request from specific verifier
func (k Keeper) RequestIdentityRecordsVerify(ctx sdk.Context, address, verifier sdk.AccAddress, recordIds []uint64, tip sdk.Coin) (uint64, error) {
	requestId := k.GetLastIdRecordVerifyRequestId(ctx) + 1

	for _, recordId := range recordIds {
		record := k.GetIdentityRecord(ctx, recordId)
		if record == nil {
			return requestId, fmt.Errorf("identity record with specified id does NOT exist: id=%d", recordId)
		}
	}

	request := types.IdentityRecordsVerify{
		Id:        requestId,
		Address:   address,
		Verifier:  verifier,
		RecordIds: recordIds,
		Tip:       tip,
	}

	k.SetIdentityRecordsVerifyRequest(ctx, request)
	k.SetLastIdRecordVerifyRequestId(ctx, requestId)

	if !tip.Amount.IsZero() {
		if err := k.bk.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, sdk.Coins{tip}); err != nil {
			return 0, err
		}
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
	request := k.GetIdRecordsVerifyRequest(ctx, requestId)
	if request == nil {
		return
	}
	prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixIdRecordVerifyRequest).Delete(sdk.Uint64ToBigEndian(requestId))
	prefix.NewStore(
		ctx.KVStore(k.storeKey),
		append(types.KeyPrefixIdRecordVerifyRequestByRequester, request.Address...),
	).Delete(sdk.Uint64ToBigEndian(requestId))
	prefix.NewStore(
		ctx.KVStore(k.storeKey),
		append(types.KeyPrefixIdRecordVerifyRequestByApprover, request.Verifier...),
	).Delete(sdk.Uint64ToBigEndian(requestId))
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

	if !request.Tip.Amount.IsZero() {
		if err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, verifier, sdk.Coins{request.Tip}); err != nil {
			return err
		}
	}

	k.DeleteIdRecordsVerifyRequest(ctx, requestId)
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

	if !request.Tip.Amount.IsZero() {
		if err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, request.Address, sdk.Coins{request.Tip}); err != nil {
			return err
		}
	}

	k.DeleteIdRecordsVerifyRequest(ctx, requestId)
	return nil
}

// GetIdRecordsVerifyRequestsByRequester query identity record verify requests by requester
func (k Keeper) GetIdRecordsVerifyRequestsByRequester(ctx sdk.Context, requester sdk.AccAddress) []types.IdentityRecordsVerify {
	requests := []types.IdentityRecordsVerify{}
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, append(types.KeyPrefixIdRecordVerifyRequestByRequester, []byte(requester)...))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		requestId := sdk.BigEndianToUint64(iterator.Value())
		request := k.GetIdRecordsVerifyRequest(ctx, requestId)
		if request == nil {
			panic(fmt.Errorf("invalid id available on requests: %d", requestId))
		}
		requests = append(requests, *request)
	}

	return requests
}

// GetIdRecordsVerifyRequestsByApprover query identity records verify requests by approver
func (k Keeper) GetIdRecordsVerifyRequestsByApprover(ctx sdk.Context, requester sdk.AccAddress) []types.IdentityRecordsVerify {
	requests := []types.IdentityRecordsVerify{}
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, append(types.KeyPrefixIdRecordVerifyRequestByApprover, []byte(requester)...))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		requestId := sdk.BigEndianToUint64(iterator.Value())
		request := k.GetIdRecordsVerifyRequest(ctx, requestId)
		if request == nil {
			panic(fmt.Errorf("invalid id available on requests: %d", requestId))
		}
		requests = append(requests, *request)
	}

	return requests
}

// GetAllIdRecordsVerifyRequests query all identity records verify requests
func (k Keeper) GetAllIdRecordsVerifyRequests(ctx sdk.Context) []types.IdentityRecordsVerify {
	requests := []types.IdentityRecordsVerify{}
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixIdRecordVerifyRequest)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		request := types.IdentityRecordsVerify{}
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &request)
		requests = append(requests, request)
	}

	return requests
}
