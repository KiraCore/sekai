package keeper

import (
	"github.com/KiraCore/sekai/x/recovery/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetRotationHistory(ctx sdk.Context, recovery types.Rotation) {
	bz := k.cdc.MustMarshal(&recovery)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.RotationHistoryKey(recovery.Address), bz)
}

func (k Keeper) DeleteRotationHistory(ctx sdk.Context, recovery types.Rotation) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.RotationHistoryKey(recovery.Address))
}

func (k Keeper) GetAllRotationHistory(ctx sdk.Context) []types.Rotation {
	store := ctx.KVStore(k.storeKey)

	rotations := []types.Rotation{}
	it := sdk.KVStorePrefixIterator(store, types.RotationHistoryKeyPrefix)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		rotation := types.Rotation{}
		k.cdc.MustUnmarshal(it.Value(), &rotation)
		rotations = append(rotations, rotation)
	}
	return rotations
}
