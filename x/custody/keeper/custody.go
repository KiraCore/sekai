package keeper

import (
	"github.com/KiraCore/sekai/x/custody/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetCustodyInfoByAddress(ctx sdk.Context, address sdk.AccAddress) *types.CustodySettings {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PrefixKeyCustodyRecord))
	bz := prefixStore.Get(address)

	if bz == nil {
		return nil
	}

	info := new(types.CustodySettings)
	k.cdc.MustUnmarshal(bz, info)

	return info
}

func (k Keeper) SetCustodyRecord(ctx sdk.Context, record types.CustodyRecord) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyCustodyRecord), record.Address...)

	store.Set(key, k.cdc.MustMarshal(record.CustodySettings))
}

func (k Keeper) DisableCustodyRecord(ctx sdk.Context, record types.CustodyRecord) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyCustodyRecord), record.Address...)

	record.CustodySettings.CustodyEnabled = false

	store.Set(key, k.cdc.MustMarshal(record.CustodySettings))
}

func (k Keeper) DropCustodyRecord(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyCustodyRecord), address...)

	store.Delete(key)
}

func (k Keeper) SetCustodyRecordKey(ctx sdk.Context, record types.CustodyKeyRecord) {
	info := k.GetCustodyInfoByAddress(ctx, record.Address)
	info.Key = record.Key
	info.NextController = record.NextController

	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyCustodyRecord), record.Address...)

	store.Set(key, k.cdc.MustMarshal(info))
}

func (k Keeper) GetCustodyCustodiansByAddress(ctx sdk.Context, address sdk.AccAddress) *types.CustodyCustodianList {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PrefixKeyCustodyCustodians))
	bz := prefixStore.Get(address)

	if bz == nil {
		return nil
	}

	info := new(types.CustodyCustodianList)
	k.cdc.MustUnmarshal(bz, info)

	return info
}

func (k Keeper) AddToCustodyCustodians(ctx sdk.Context, record types.CustodyCustodiansRecord) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyCustodyCustodians), record.Address...)

	store.Set(key, k.cdc.MustMarshal(record.CustodyCustodians))
}

func (k Keeper) DropCustodyCustodiansByAddress(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyCustodyCustodians), address...)

	store.Delete(key)
}

func (k Keeper) GetCustodyWhiteListByAddress(ctx sdk.Context, address sdk.AccAddress) *types.CustodyWhiteList {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PrefixKeyCustodyWhiteList))
	bz := prefixStore.Get(address)

	if bz == nil {
		return nil
	}

	info := new(types.CustodyWhiteList)
	k.cdc.MustUnmarshal(bz, info)

	return info
}

func (k Keeper) AddToCustodyWhiteList(ctx sdk.Context, record types.CustodyWhiteListRecord) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyCustodyWhiteList), record.Address...)

	store.Set(key, k.cdc.MustMarshal(record.CustodyWhiteList))
}

func (k Keeper) DropCustodyWhiteListByAddress(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyCustodyWhiteList), address...)

	store.Delete(key)
}

func (k Keeper) GetCustodyLimitsByAddress(ctx sdk.Context, address sdk.AccAddress) *types.CustodyLimits {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PrefixKeyCustodyLimits))
	bz := prefixStore.Get(address)

	if bz == nil {
		return nil
	}

	info := new(types.CustodyLimits)
	k.cdc.MustUnmarshal(bz, info)

	return info
}

func (k Keeper) AddToCustodyLimits(ctx sdk.Context, record types.CustodyLimitRecord) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyCustodyLimits), record.Address...)

	store.Set(key, k.cdc.MustMarshal(record.CustodyLimits))
}

func (k Keeper) DropCustodyLimitsByAddress(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyCustodyLimits), address...)

	store.Delete(key)
}

func (k Keeper) GetCustodyLimitsStatusByAddress(ctx sdk.Context, address sdk.AccAddress) *types.CustodyStatuses {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PrefixKeyCustodyLimitsStatus))
	bz := prefixStore.Get(address)

	if bz == nil {
		return nil
	}

	info := new(types.CustodyStatuses)
	k.cdc.MustUnmarshal(bz, info)

	return info
}

func (k Keeper) AddToCustodyLimitsStatus(ctx sdk.Context, record types.CustodyLimitStatusRecord) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyCustodyLimitsStatus), record.Address...)

	store.Set(key, k.cdc.MustMarshal(record.CustodyStatuses))
}

func (k Keeper) DropCustodyLimitsStatus(ctx sdk.Context, addr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyCustodyLimitsStatus), addr...)
	store.Delete(key)
}

func (k Keeper) AddToCustodyPool(ctx sdk.Context, record types.CustodyPool) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyCustodyPool), record.Address...)

	store.Set(key, k.cdc.MustMarshal(record.Transactions))
}

func (k Keeper) DropCustodyPool(ctx sdk.Context, addr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyCustodyPool), addr...)
	store.Delete(key)
}

func (k Keeper) GetCustodyPoolByAddress(ctx sdk.Context, address sdk.AccAddress) *types.TransactionPool {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PrefixKeyCustodyPool))
	bz := prefixStore.Get(address)

	if bz == nil {
		return nil
	}

	info := new(types.TransactionPool)
	k.cdc.MustUnmarshal(bz, info)

	return info
}

func (k Keeper) GetApproveCustody(ctx sdk.Context, msg *types.MsgApproveCustodyTransaction) string {
	store := ctx.KVStore(k.storeKey)
	key1 := append([]byte(types.PrefixKeyCustodyVote), msg.FromAddress...)
	key2 := append(msg.TargetAddress, msg.Hash...)
	key := append(key1, key2...)
	bz := store.Get(key)

	if bz == nil {
		return "0"
	}

	return string(bz)
}

func (k Keeper) GetDeclineCustody(ctx sdk.Context, msg *types.MsgDeclineCustodyTransaction) string {
	store := ctx.KVStore(k.storeKey)
	key1 := append([]byte(types.PrefixKeyCustodyVote), msg.FromAddress...)
	key2 := append(msg.TargetAddress, msg.Hash...)
	key := append(key1, key2...)
	bz := store.Get(key)

	if bz == nil {
		return "0"
	}

	return string(bz)
}

func (k Keeper) ApproveCustody(ctx sdk.Context, msg *types.MsgApproveCustodyTransaction) {
	store := ctx.KVStore(k.storeKey)
	key1 := append([]byte(types.PrefixKeyCustodyVote), msg.FromAddress...)
	key2 := append(msg.TargetAddress, msg.Hash...)
	key := append(key1, key2...)

	store.Set(key, []byte("1"))
}

func (k Keeper) DeclineCustody(ctx sdk.Context, msg *types.MsgDeclineCustodyTransaction) {
	store := ctx.KVStore(k.storeKey)
	key1 := append([]byte(types.PrefixKeyCustodyVote), msg.FromAddress...)
	key2 := append(msg.TargetAddress, msg.Hash...)
	key := append(key1, key2...)

	store.Set(key, []byte("-1"))
}
