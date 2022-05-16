package keeper

import (
	"github.com/KiraCore/sekai/x/multistaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetLastUndelegationId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyLastUndelegationId)
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLastUndelegationId(ctx sdk.Context, id uint64) {
	idBz := sdk.Uint64ToBigEndian(id)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyLastUndelegationId, idBz)
}

func (k Keeper) GetUndelegationById(ctx sdk.Context, id uint64) (undelegation types.Undelegation, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.KeyPrefixStakingPool), sdk.Uint64ToBigEndian(id)...)
	bz := store.Get(key)
	if bz == nil {
		return undelegation, false
	}
	k.cdc.MustUnmarshal(bz, &undelegation)
	return undelegation, true
}

func (k Keeper) SetUndelegation(ctx sdk.Context, undelegation types.Undelegation) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.KeyPrefixUndelegation), sdk.Uint64ToBigEndian(undelegation.Id)...)
	store.Set(key, k.cdc.MustMarshal(&undelegation))
}
