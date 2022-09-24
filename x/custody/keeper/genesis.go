package keeper

import (
	"github.com/KiraCore/sekai/x/custody/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetMaxCustodyBufferSize(ctx sdk.Context, size uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CustodyBufferSizeKey, sdk.Uint64ToBigEndian(size))
}

func (k Keeper) SetMaxCustodyTxSize(ctx sdk.Context, size uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CustodyTxSizeKey, sdk.Uint64ToBigEndian(size))
}

func (k Keeper) GetMaxCustodyBufferSize(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CustodyBufferSizeKey)

	if bz == nil {
		return 1
	}

	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) GetMaxCustodyTxSize(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CustodyTxSizeKey)

	if bz == nil {
		return 1
	}

	return sdk.BigEndianToUint64(bz)
}
