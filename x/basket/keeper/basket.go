package keeper

import (
	"github.com/KiraCore/sekai/x/basket/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) GetLastBasketId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyLastBasketId)
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLastBasketId(ctx sdk.Context, id uint64) {
	idBz := sdk.Uint64ToBigEndian(id)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyLastBasketId, idBz)
}

func (k Keeper) GetBasketById(ctx sdk.Context, id uint64) (types.Basket, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(append(types.PrefixBasketKey, sdk.Uint64ToBigEndian(id)...))
	if bz == nil {
		return types.Basket{}, sdkerrors.Wrapf(types.ErrBasketDoesNotExist, "basket: %d does not exist", id)
	}
	basket := types.Basket{}
	k.cdc.MustUnmarshal(bz, &basket)
	return basket, nil
}

func (k Keeper) SetBasket(ctx sdk.Context, basket types.Basket) {
	idBz := sdk.Uint64ToBigEndian(basket.Id)
	bz := k.cdc.MustMarshal(&basket)
	store := ctx.KVStore(k.storeKey)
	store.Set(append(types.PrefixBasketKey, idBz...), bz)
}

func (k Keeper) GetAllBaskets(ctx sdk.Context) []types.Basket {
	store := ctx.KVStore(k.storeKey)

	baskets := []types.Basket{}
	it := sdk.KVStorePrefixIterator(store, types.PrefixBasketKey)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		basket := types.Basket{}
		k.cdc.MustUnmarshal(it.Value(), &basket)
		baskets = append(baskets, basket)
	}
	return baskets
}
