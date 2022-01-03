package keeper

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/KiraCore/sekai/x/gov/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func ValidateIdentityRecordKey(key string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z][_0-9a-zA-Z]*$`)
	return regex.MatchString(key)
}

func FormalizeIdentityRecordKey(key string) string {
	return strings.ToLower(key)
}

func CheckIfWithinStringArray(key string, keys []string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}

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
	// validate key
	if !ValidateIdentityRecordKey(record.Key) {
		panic("identity record key is invalid")
	}
	properties := k.GetNetworkProperties(ctx)
	uniqueKeys := strings.Split(properties.UniqueIdentityKeys, ",")
	if CheckIfWithinStringArray(record.Key, uniqueKeys) {
		addrs := k.GetAddressesByIdRecordKey(ctx, record.Key, record.Value)
		if len(addrs) == 1 && addrs[0].String() == record.Address {

		} else if len(addrs) > 0 {
			panic(fmt.Sprintf("the key %s, value %s is already registered by %s", record.Key, record.Value, addrs[0].String()))
		}
	}
	// set the key to non case-sensitive
	record.Key = FormalizeIdentityRecordKey(record.Key)

	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixIdentityRecord)
	bz := k.cdc.MustMarshal(&record)
	prefixStore.Set(sdk.Uint64ToBigEndian(record.Id), bz)

	// connect address + key to id
	store := ctx.KVStore(k.storeKey)
	addrPrefixStore := prefix.NewStore(store, types.IdentityRecordByAddressPrefix(record.Address))
	addrPrefixStore.Set([]byte(record.Key), sdk.Uint64ToBigEndian(record.Id))
}

func (k Keeper) GetIdentityRecordById(ctx sdk.Context, recordId uint64) *types.IdentityRecord {
	record := types.IdentityRecord{}

	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixIdentityRecord)
	bz := prefixStore.Get(sdk.Uint64ToBigEndian(recordId))
	if bz == nil {
		return nil
	}
	k.cdc.MustUnmarshal(bz, &record)
	return &record
}

// Get identity record id by address and key
func (k Keeper) GetIdentityRecordIdByAddressKey(ctx sdk.Context, address sdk.AccAddress, key string) uint64 {
	// validate key
	if !ValidateIdentityRecordKey(key) {
		return 0
	}
	// set the key to non case-sensitive
	key = FormalizeIdentityRecordKey(key)

	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.IdentityRecordByAddressPrefix(address.String()))
	recordIdBytes := prefixStore.Get([]byte(key))
	if recordIdBytes == nil {
		return 0
	}
	return sdk.BigEndianToUint64(recordIdBytes)
}

// DeleteIdentityRecord defines a method to delete identity record by id
func (k Keeper) DeleteIdentityRecordById(ctx sdk.Context, recordId uint64) {
	record := k.GetIdentityRecordById(ctx, recordId)
	if record == nil {
		return
	}
	recordKey := sdk.Uint64ToBigEndian(recordId)
	prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixIdentityRecord).Delete(recordKey)
	prefix.NewStore(ctx.KVStore(k.storeKey), types.IdentityRecordByAddressPrefix(record.Address)).Delete(sdk.Uint64ToBigEndian(recordId))
}

// RegisterIdentityRecord defines a method to register identity records for an address
func (k Keeper) RegisterIdentityRecords(ctx sdk.Context, address sdk.AccAddress, infos []types.IdentityInfoEntry) error {
	// validate key and set the key to non case-sensitive
	properties := k.GetNetworkProperties(ctx)
	uniqueKeys := strings.Split(properties.UniqueIdentityKeys, ",")
	for i, info := range infos {
		if !ValidateIdentityRecordKey(info.Key) {
			return sdkerrors.Wrap(types.ErrInvalidIdentityRecordKey, fmt.Sprintf("invalid key exists: key=%s", info.Key))
		}
		infos[i].Key = FormalizeIdentityRecordKey(info.Key)

		if infos[i].Key == "moniker" && len(infos[i].Info) > 32 {
			return stakingtypes.ErrInvalidMonikerLength
		}

		if CheckIfWithinStringArray(infos[i].Key, uniqueKeys) {
			addrs := k.GetAddressesByIdRecordKey(ctx, infos[i].Key, infos[i].Info)
			if len(addrs) == 1 && bytes.Equal(addrs[0], address) {

			} else if len(addrs) > 0 {
				return sdkerrors.Wrap(types.ErrKeyShouldBeUnique, fmt.Sprintf("the key %s, value %s is already registered by %s", infos[i].Key, infos[i].Info, addrs[0].String()))
			}
		}
	}

	for _, info := range infos {
		// use existing record id if it already exists
		recordId := k.GetIdentityRecordIdByAddressKey(ctx, address, info.Key)
		if recordId == 0 {
			recordId = k.GetLastIdentityRecordId(ctx) + 1
			k.SetLastIdentityRecordId(ctx, recordId)
		}
		// create or update identity record
		k.SetIdentityRecord(ctx, types.IdentityRecord{
			Id:        recordId,
			Address:   address.String(),
			Key:       info.Key,
			Value:     info.Info,
			Date:      ctx.BlockTime(),
			Verifiers: []string{},
		})
	}
	return nil
}

// DeleteIdentityRecords defines a method to delete identity records owned by an address
func (k Keeper) DeleteIdentityRecords(ctx sdk.Context, address sdk.AccAddress, keys []string) error {
	// validate key and set the key to non case-sensitive
	for i, key := range keys {
		if !ValidateIdentityRecordKey(key) {
			return sdkerrors.Wrap(types.ErrInvalidIdentityRecordKey, fmt.Sprintf("invalid key exists: key=%s", key))
		}
		keys[i] = FormalizeIdentityRecordKey(key)

		// we prevent deleting moniker field of a validator
		if key == "moniker" {
			return sdkerrors.Wrap(types.ErrMonikerDeletionNotAllowed, fmt.Sprintf("moniker field is not allowed to delete"))
		}
	}

	store := ctx.KVStore(k.storeKey)
	prefix := types.IdentityRecordByAddressPrefix(address.String())
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	keyMap := make(map[string]bool)
	for _, key := range keys {
		keyMap[key] = true
	}
	recordIds := []uint64{}
	for ; iterator.Valid(); iterator.Next() {
		key := bytes.TrimPrefix(iterator.Key(), prefix)
		if len(keys) == 0 || keyMap[string(key)] {
			// if no specific keys are provided remove all
			// invalid keys are ignored
			recordIds = append(recordIds, sdk.BigEndianToUint64(iterator.Value()))
			store.Delete(iterator.Key())
		}
	}

	for _, recordId := range recordIds {
		prevRecord := k.GetIdentityRecordById(ctx, recordId)
		if prevRecord == nil {
			return sdkerrors.Wrap(types.ErrInvalidIdentityRecordId, fmt.Sprintf("identity record with specified id does NOT exist: id=%d", recordId))
		}

		k.DeleteIdentityRecordById(ctx, recordId)
	}

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
		k.cdc.MustUnmarshal(iterator.Value(), &record)
		records = append(records, record)
	}

	return records
}

// GetIdRecordsByAddressAndKeys query identity record by address and keys
func (k Keeper) GetIdRecordsByAddressAndKeys(ctx sdk.Context, address sdk.AccAddress, keys []string) ([]types.IdentityRecord, error) {
	// validate key and set the key to non case-sensitive
	for i, key := range keys {
		if !ValidateIdentityRecordKey(key) {
			return []types.IdentityRecord{}, sdkerrors.Wrap(types.ErrInvalidIdentityRecordKey, fmt.Sprintf("invalid key exists: key=%s", key))
		}
		keys[i] = FormalizeIdentityRecordKey(key)
	}

	if len(keys) == 0 {
		return k.GetIdRecordsByAddress(ctx, address), nil
	}

	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.IdentityRecordByAddressPrefix(address.String()))

	records := []types.IdentityRecord{}
	for _, key := range keys {
		bz := prefixStore.Get([]byte(key))
		recordId := sdk.BigEndianToUint64(bz)
		record := k.GetIdentityRecordById(ctx, recordId)
		if record == nil {
			return records, sdkerrors.Wrap(types.ErrInvalidIdentityRecordId, fmt.Sprintf("identity record with specified id does NOT exist: id=%d", recordId))
		}
		records = append(records, *record)
	}
	return records, nil
}

func (k Keeper) GetAddressesByIdRecordKey(ctx sdk.Context, key, value string) []sdk.AccAddress {
	addrs := []sdk.AccAddress{}
	for _, record := range k.GetAllIdentityRecords(ctx) {
		if record.Key == key && record.Value == value {
			addr, err := sdk.AccAddressFromBech32(record.Address)
			if err != nil {
				panic(err)
			}
			addrs = append(addrs, addr)
		}
	}
	return addrs
}

// GetIdRecordsByAddress query identity record by address
func (k Keeper) GetIdRecordsByAddress(ctx sdk.Context, address sdk.AccAddress) []types.IdentityRecord {
	store := ctx.KVStore(k.storeKey)
	prefix := types.IdentityRecordByAddressPrefix(address.String())
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	records := []types.IdentityRecord{}
	for ; iterator.Valid(); iterator.Next() {
		recordId := sdk.BigEndianToUint64(iterator.Value())
		record := k.GetIdentityRecordById(ctx, recordId)
		if record == nil {
			panic(fmt.Sprintf("invalid recordId exists: recordId = %d, key=%s, value=%s", recordId, string(iterator.Key()), string(iterator.Value())))
		}
		records = append(records, *record)
	}
	return records
}

// SetIdentityRecordsVerifyRequest saves identity verify request into the store
func (k Keeper) SetIdentityRecordsVerifyRequest(ctx sdk.Context, request types.IdentityRecordsVerify) {
	requestId := request.Id
	bz := k.cdc.MustMarshal(&request)
	prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixIdRecordVerifyRequest).Set(sdk.Uint64ToBigEndian(requestId), bz)
	prefix.NewStore(
		ctx.KVStore(k.storeKey),
		types.IdRecordVerifyRequestByRequesterPrefix(request.Address),
	).Set(sdk.Uint64ToBigEndian(requestId), sdk.Uint64ToBigEndian(requestId))
	prefix.NewStore(
		ctx.KVStore(k.storeKey),
		types.IdRecordVerifyRequestByApproverPrefix(request.Verifier),
	).Set(sdk.Uint64ToBigEndian(requestId), sdk.Uint64ToBigEndian(requestId))
}

// RequestIdentityRecordsVerify defines a method to request verify request from specific verifier
func (k Keeper) RequestIdentityRecordsVerify(ctx sdk.Context, address, verifier sdk.AccAddress, recordIds []uint64, tip sdk.Coin) (uint64, error) {
	requestId := k.GetLastIdRecordVerifyRequestId(ctx) + 1

	lastRecordEditDate := time.Time{}
	for _, recordId := range recordIds {
		record := k.GetIdentityRecordById(ctx, recordId)
		if record == nil {
			return requestId, sdkerrors.Wrap(types.ErrInvalidIdentityRecordId, fmt.Sprintf("identity record with specified id does NOT exist: id=%d", recordId))
		}
		if lastRecordEditDate.Before(record.Date) {
			lastRecordEditDate = record.Date
		}
	}

	request := types.IdentityRecordsVerify{
		Id:                 requestId,
		Address:            address.String(),
		Verifier:           verifier.String(),
		RecordIds:          recordIds,
		Tip:                tip,
		LastRecordEditDate: lastRecordEditDate,
	}

	minApprovalTip := k.GetNetworkProperties(ctx).MinIdentityApprovalTip
	if sdk.NewInt(int64(minApprovalTip)).GT(tip.Amount) {
		return requestId, sdkerrors.Wrap(types.ErrInvalidApprovalTip, fmt.Sprintf("approval tip is lower than minimum tip configured by the network"))
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
	k.cdc.MustUnmarshal(bz, &request)
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
		types.IdRecordVerifyRequestByRequesterPrefix(request.Address),
	).Delete(sdk.Uint64ToBigEndian(requestId))
	prefix.NewStore(
		ctx.KVStore(k.storeKey),
		types.IdRecordVerifyRequestByApproverPrefix(request.Verifier),
	).Delete(sdk.Uint64ToBigEndian(requestId))
}

// ApproveIdentityRecords defines a method to accept or reject verification request
func (k Keeper) HandleIdentityRecordsVerifyRequest(ctx sdk.Context, verifier sdk.AccAddress, requestId uint64, approve bool) error {
	request := k.GetIdRecordsVerifyRequest(ctx, requestId)
	if request == nil {
		return sdkerrors.Wrap(types.ErrInvalidVerifyRequestId, fmt.Sprintf("specified identity record verify request does NOT exist: id=%d", requestId))
	}
	if verifier.String() != request.Verifier {
		return errors.New("verifier does not match with requested")
	}

	// send the balance regardless approve or reject
	if !request.Tip.Amount.IsZero() {
		if err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, verifier, sdk.Coins{request.Tip}); err != nil {
			return err
		}
	}

	// automatically reject if last record edit date is incorrect
	for _, recordId := range request.RecordIds {
		record := k.GetIdentityRecordById(ctx, recordId)
		if record == nil {
			return sdkerrors.Wrap(types.ErrInvalidIdentityRecordId, fmt.Sprintf("identity record with specified id does NOT exist: id=%d", recordId))
		}

		if record.Date.After(request.LastRecordEditDate) {
			approve = false
			break
		}
	}

	if approve == false {
		k.DeleteIdRecordsVerifyRequest(ctx, requestId)
		return nil
	}

	for _, recordId := range request.RecordIds {
		record := k.GetIdentityRecordById(ctx, recordId)
		if record == nil {
			return sdkerrors.Wrap(types.ErrInvalidIdentityRecordId, fmt.Sprintf("identity record with specified id does NOT exist: id=%d", recordId))
		}

		// if already exist, skip
		if CheckIfWithinStringArray(verifier.String(), record.Verifiers) {
			continue
		}
		record.Verifiers = append(record.Verifiers, verifier.String())
		k.SetIdentityRecord(ctx, *record)
	}

	k.DeleteIdRecordsVerifyRequest(ctx, requestId)
	return nil
}

// CancelIdentityRecordsVerifyRequest defines a method to cancel verification request
func (k Keeper) CancelIdentityRecordsVerifyRequest(ctx sdk.Context, executor sdk.AccAddress, requestId uint64) error {
	request := k.GetIdRecordsVerifyRequest(ctx, requestId)
	if request == nil {
		return sdkerrors.Wrap(types.ErrInvalidVerifyRequestId, fmt.Sprintf("specified identity record verify request does NOT exist: id=%d", requestId))
	}
	if executor.String() != request.Address {
		return errors.New("executor is not identity record creator")
	}

	if !request.Tip.Amount.IsZero() {
		requester, err := sdk.AccAddressFromBech32(request.Address)
		if err != nil {
			return err
		}
		if err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, requester, sdk.Coins{request.Tip}); err != nil {
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
	iterator := sdk.KVStorePrefixIterator(store, types.IdRecordVerifyRequestByRequesterPrefix(requester.String()))
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
	iterator := sdk.KVStorePrefixIterator(store, types.IdRecordVerifyRequestByApproverPrefix(requester.String()))
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
		k.cdc.MustUnmarshal(iterator.Value(), &request)
		requests = append(requests, request)
	}

	return requests
}
