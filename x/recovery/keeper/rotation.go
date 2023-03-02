package keeper

import (
	"github.com/KiraCore/sekai/x/recovery/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetRotationHistory(ctx sdk.Context, rotation types.Rotation) {
	bz := k.cdc.MustMarshal(&rotation)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.RotationHistoryKey(rotation.Address), bz)
}

func (k Keeper) DeleteRotationHistory(ctx sdk.Context, rotation types.Rotation) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.RotationHistoryKey(rotation.Address))
}

func (k Keeper) GetRotationHistory(ctx sdk.Context, address string) types.Rotation {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RotationHistoryKey(address))
	if bz == nil {
		return types.Rotation{}
	}

	rotation := types.Rotation{}
	k.cdc.MustUnmarshal(bz, &rotation)
	return rotation
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
