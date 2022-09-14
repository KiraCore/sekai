package keeper

import (
	"time"

	"github.com/KiraCore/sekai/x/basket/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) RegisterMintAction(ctx sdk.Context, basketId uint64, amount sdk.Int) {
	mintedAmount := k.GetMintAmount(ctx, basketId, ctx.BlockTime())
	k.SetMintAmount(ctx, ctx.BlockTime(), basketId, mintedAmount.Add(amount))
}

func (k Keeper) RegisterBurnAction(ctx sdk.Context, basketId uint64, amount sdk.Int) {
	burnedAmount := k.GetBurnAmount(ctx, basketId, ctx.BlockTime())
	k.SetBurnAmount(ctx, ctx.BlockTime(), basketId, burnedAmount.Add(amount))
}

func (k Keeper) RegisterSwapAction(ctx sdk.Context, basketId uint64, amount sdk.Int) {
	swapedAmount := k.GetSwapAmount(ctx, basketId, ctx.BlockTime())
	k.SetMintAmount(ctx, ctx.BlockTime(), basketId, swapedAmount.Add(amount))
}

func (k Keeper) SetMintAmount(ctx sdk.Context, time time.Time, basketId uint64, amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := []byte(amount.String())
	store.Set(types.BasketMintByTimeKey(basketId, time), bz)
}

func (k Keeper) SetBurnAmount(ctx sdk.Context, time time.Time, basketId uint64, amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := []byte(amount.String())
	store.Set(types.BasketBurnByTimeKey(basketId, time), bz)
}

func (k Keeper) SetSwapAmount(ctx sdk.Context, time time.Time, basketId uint64, amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := []byte(amount.String())
	store.Set(types.BasketSwapByTimeKey(basketId, time), bz)
}

func (k Keeper) GetMintAmount(ctx sdk.Context, basketId uint64, time time.Time) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.BasketMintByTimeKey(basketId, time))
	if bz == nil {
		return sdk.ZeroInt()
	}
	amountStr := string(bz)
	amount, ok := sdk.NewIntFromString(amountStr)
	if !ok {
		panic("invalid mint amount recorded")
	}
	return amount
}

func (k Keeper) GetBurnAmount(ctx sdk.Context, basketId uint64, time time.Time) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.BasketBurnByTimeKey(basketId, time))
	if bz == nil {
		return sdk.ZeroInt()
	}
	amountStr := string(bz)
	amount, ok := sdk.NewIntFromString(amountStr)
	if !ok {
		panic("invalid burn amount recorded")
	}
	return amount
}

func (k Keeper) GetSwapAmount(ctx sdk.Context, basketId uint64, time time.Time) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.BasketSwapByTimeKey(basketId, time))
	if bz == nil {
		return sdk.ZeroInt()
	}
	amountStr := string(bz)
	amount, ok := sdk.NewIntFromString(amountStr)
	if !ok {
		panic("invalid swap amount recorded")
	}
	return amount
}

func (k Keeper) GetLimitsPeriodMintAmount(ctx sdk.Context, basketId uint64, limitsPeriod uint64) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	startTime := ctx.BlockTime().Add(-time.Second * time.Duration(limitsPeriod))
	iterator := store.Iterator(
		sdk.PrefixEndBytes(append(types.PrefixBasketMintByTime, sdk.Uint64ToBigEndian(basketId)...)),
		types.BasketMintByTimeKey(basketId, startTime),
	)

	defer iterator.Close()

	totalAmount := sdk.ZeroInt()
	for ; iterator.Valid(); iterator.Next() {
		amountStr := string(iterator.Value())
		amount, ok := sdk.NewIntFromString(amountStr)
		if !ok {
			panic("invalid burn amount recorded")
		}
		totalAmount = totalAmount.Add(amount)
	}
	return totalAmount
}

func (k Keeper) GetLimitsPeriodBurnAmount(ctx sdk.Context, basketId uint64, limitsPeriod uint64) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	startTime := ctx.BlockTime().Add(-time.Second * time.Duration(limitsPeriod))
	iterator := store.Iterator(
		sdk.PrefixEndBytes(append(types.PrefixBasketBurnByTime, sdk.Uint64ToBigEndian(basketId)...)),
		types.BasketBurnByTimeKey(basketId, startTime),
	)

	defer iterator.Close()

	totalAmount := sdk.ZeroInt()
	for ; iterator.Valid(); iterator.Next() {
		amountStr := string(iterator.Value())
		amount, ok := sdk.NewIntFromString(amountStr)
		if !ok {
			panic("invalid mint amount recorded")
		}
		totalAmount = totalAmount.Add(amount)
	}
	return totalAmount
}

func (k Keeper) GetLimitsPeriodSwapAmount(ctx sdk.Context, basketId uint64, limitsPeriod uint64) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	startTime := ctx.BlockTime().Add(-time.Second * time.Duration(limitsPeriod))
	iterator := store.Iterator(
		sdk.PrefixEndBytes(append(types.PrefixBasketSwapByTime, sdk.Uint64ToBigEndian(basketId)...)),
		types.BasketSwapByTimeKey(basketId, startTime),
	)

	defer iterator.Close()

	totalAmount := sdk.ZeroInt()
	for ; iterator.Valid(); iterator.Next() {
		amountStr := string(iterator.Value())
		amount, ok := sdk.NewIntFromString(amountStr)
		if !ok {
			panic("invalid swap amount recorded")
		}
		totalAmount = totalAmount.Add(amount)
	}
	return totalAmount
}

func (k Keeper) ClearOldMintAmounts(ctx sdk.Context, basketId uint64, limitsPeriod uint64) {
	store := ctx.KVStore(k.storeKey)
	startTime := ctx.BlockTime().Add(-time.Second * time.Duration(limitsPeriod))
	iterator := store.Iterator(
		append(types.PrefixBasketMintByTime, sdk.Uint64ToBigEndian(basketId)...),
		types.BasketMintByTimeKey(basketId, startTime))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

func (k Keeper) ClearOldBurnAmounts(ctx sdk.Context, basketId uint64, limitsPeriod uint64) {
	store := ctx.KVStore(k.storeKey)
	startTime := ctx.BlockTime().Add(-time.Second * time.Duration(limitsPeriod))
	iterator := store.Iterator(
		append(types.PrefixBasketBurnByTime, sdk.Uint64ToBigEndian(basketId)...),
		types.BasketBurnByTimeKey(basketId, startTime))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

func (k Keeper) ClearOldSwapAmounts(ctx sdk.Context, basketId uint64, limitsPeriod uint64) {
	store := ctx.KVStore(k.storeKey)
	startTime := ctx.BlockTime().Add(-time.Second * time.Duration(limitsPeriod))
	iterator := store.Iterator(
		append(types.PrefixBasketSwapByTime, sdk.Uint64ToBigEndian(basketId)...),
		types.BasketSwapByTimeKey(basketId, startTime))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}
