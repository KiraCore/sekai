package keeper

import (
	"time"

	"github.com/KiraCore/sekai/x/basket/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) RegisterMintAction(ctx sdk.Context, basketId uint64, amount sdk.Int) {
	mintedAmount := k.GetMintAmount(ctx, basketId, ctx.BlockTime())
	mintedAmount = mintedAmount.Add(amount)
	k.SetMintAmount(ctx, ctx.BlockTime(), basketId, mintedAmount)
}

func (k Keeper) RegisterBurnAction(ctx sdk.Context, basketId uint64, amount sdk.Int) {
	burnedAmount := k.GetBurnAmount(ctx, basketId, ctx.BlockTime())
	burnedAmount = burnedAmount.Add(amount)
	k.SetBurnAmount(ctx, ctx.BlockTime(), basketId, burnedAmount)
}

func (k Keeper) RegisterSwapAction(ctx sdk.Context, basketId uint64, amount sdk.Int) {
	swapedAmount := k.GetSwapAmount(ctx, basketId, ctx.BlockTime())
	swapedAmount = swapedAmount.Add(amount)
	k.SetSwapAmount(ctx, ctx.BlockTime(), basketId, swapedAmount)
}

func (k Keeper) SetMintAmount(ctx sdk.Context, time time.Time, basketId uint64, amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	amountByTime := types.AmountAtTime{
		BasketId: basketId,
		Amount:   amount,
		Time:     uint64(time.Unix()),
	}
	bz := k.cdc.MustMarshal(&amountByTime)
	store.Set(types.BasketMintByTimeKey(basketId, time), bz)
}

func (k Keeper) SetBurnAmount(ctx sdk.Context, time time.Time, basketId uint64, amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	amountByTime := types.AmountAtTime{
		BasketId: basketId,
		Amount:   amount,
		Time:     uint64(time.Unix()),
	}
	bz := k.cdc.MustMarshal(&amountByTime)
	store.Set(types.BasketBurnByTimeKey(basketId, time), bz)
}

func (k Keeper) SetSwapAmount(ctx sdk.Context, time time.Time, basketId uint64, amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	amountByTime := types.AmountAtTime{
		BasketId: basketId,
		Amount:   amount,
		Time:     uint64(time.Unix()),
	}
	bz := k.cdc.MustMarshal(&amountByTime)
	store.Set(types.BasketSwapByTimeKey(basketId, time), bz)
}

func (k Keeper) GetMintAmount(ctx sdk.Context, basketId uint64, time time.Time) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.BasketMintByTimeKey(basketId, time))
	if bz == nil {
		return sdk.ZeroInt()
	}
	amountByTime := types.AmountAtTime{}
	k.cdc.MustUnmarshal(bz, &amountByTime)
	return amountByTime.Amount
}

func (k Keeper) GetBurnAmount(ctx sdk.Context, basketId uint64, time time.Time) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.BasketBurnByTimeKey(basketId, time))
	if bz == nil {
		return sdk.ZeroInt()
	}
	amountByTime := types.AmountAtTime{}
	k.cdc.MustUnmarshal(bz, &amountByTime)
	return amountByTime.Amount
}

func (k Keeper) GetSwapAmount(ctx sdk.Context, basketId uint64, time time.Time) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.BasketSwapByTimeKey(basketId, time))
	if bz == nil {
		return sdk.ZeroInt()
	}
	amountByTime := types.AmountAtTime{}
	k.cdc.MustUnmarshal(bz, &amountByTime)
	return amountByTime.Amount
}

func (k Keeper) GetLimitsPeriodMintAmount(ctx sdk.Context, basketId uint64, limitsPeriod uint64) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	startTime := ctx.BlockTime().Add(-time.Second * time.Duration(limitsPeriod))
	iterator := store.Iterator(
		types.BasketMintByTimeKey(basketId, startTime),
		sdk.PrefixEndBytes(append(types.PrefixBasketMintByTime, sdk.Uint64ToBigEndian(basketId)...)),
	)

	defer iterator.Close()

	totalAmount := sdk.ZeroInt()
	for ; iterator.Valid(); iterator.Next() {
		amountByTime := types.AmountAtTime{}
		k.cdc.MustUnmarshal(iterator.Value(), &amountByTime)
		totalAmount = totalAmount.Add(amountByTime.Amount)
	}
	return totalAmount
}

func (k Keeper) GetLimitsPeriodBurnAmount(ctx sdk.Context, basketId uint64, limitsPeriod uint64) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	startTime := ctx.BlockTime().Add(-time.Second * time.Duration(limitsPeriod))
	iterator := store.Iterator(
		types.BasketBurnByTimeKey(basketId, startTime),
		sdk.PrefixEndBytes(append(types.PrefixBasketBurnByTime, sdk.Uint64ToBigEndian(basketId)...)),
	)

	defer iterator.Close()

	totalAmount := sdk.ZeroInt()
	for ; iterator.Valid(); iterator.Next() {
		amountByTime := types.AmountAtTime{}
		k.cdc.MustUnmarshal(iterator.Value(), &amountByTime)
		totalAmount = totalAmount.Add(amountByTime.Amount)
	}
	return totalAmount
}

func (k Keeper) GetLimitsPeriodSwapAmount(ctx sdk.Context, basketId uint64, limitsPeriod uint64) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	startTime := ctx.BlockTime().Add(-time.Second * time.Duration(limitsPeriod))
	iterator := store.Iterator(
		types.BasketSwapByTimeKey(basketId, startTime),
		sdk.PrefixEndBytes(append(types.PrefixBasketSwapByTime, sdk.Uint64ToBigEndian(basketId)...)),
	)

	defer iterator.Close()

	totalAmount := sdk.ZeroInt()
	for ; iterator.Valid(); iterator.Next() {
		amountByTime := types.AmountAtTime{}
		k.cdc.MustUnmarshal(iterator.Value(), &amountByTime)
		totalAmount = totalAmount.Add(amountByTime.Amount)
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

func (k Keeper) GetAllMintAmounts(ctx sdk.Context) []types.AmountAtTime {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixBasketMintByTime)
	iterator := prefixStore.Iterator(nil, nil)

	defer iterator.Close()

	historicalMints := []types.AmountAtTime{}
	for ; iterator.Valid(); iterator.Next() {
		amount := types.AmountAtTime{}
		k.cdc.MustUnmarshal(iterator.Value(), &amount)
		historicalMints = append(historicalMints, amount)
	}
	return historicalMints
}

func (k Keeper) GetAllBurnAmounts(ctx sdk.Context) []types.AmountAtTime {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixBasketBurnByTime)
	iterator := prefixStore.Iterator(nil, nil)

	defer iterator.Close()

	historicalBurns := []types.AmountAtTime{}
	for ; iterator.Valid(); iterator.Next() {
		amount := types.AmountAtTime{}
		k.cdc.MustUnmarshal(iterator.Value(), &amount)
		historicalBurns = append(historicalBurns, amount)
	}
	return historicalBurns
}

func (k Keeper) GetAllSwapAmounts(ctx sdk.Context) []types.AmountAtTime {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixBasketSwapByTime)
	iterator := prefixStore.Iterator(nil, nil)

	defer iterator.Close()

	historicalSwaps := []types.AmountAtTime{}
	for ; iterator.Valid(); iterator.Next() {
		amount := types.AmountAtTime{}
		k.cdc.MustUnmarshal(iterator.Value(), &amount)
		historicalSwaps = append(historicalSwaps, amount)
	}
	return historicalSwaps
}
