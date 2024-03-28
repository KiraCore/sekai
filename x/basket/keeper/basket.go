package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/KiraCore/sekai/x/basket/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		return types.Basket{}, errorsmod.Wrapf(types.ErrBasketDoesNotExist, "basket: %d does not exist", id)
	}
	basket := types.Basket{}
	k.cdc.MustUnmarshal(bz, &basket)
	return basket, nil
}

func (k Keeper) GetBasketByDenom(ctx sdk.Context, denom string) (types.Basket, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(append(types.PrefixBasketByDenomKey, denom...))
	if bz == nil {
		return types.Basket{}, errorsmod.Wrapf(types.ErrBasketDoesNotExist, "basket: %s does not exist", denom)
	}
	id := sdk.BigEndianToUint64(bz)
	return k.GetBasketById(ctx, id)
}

func (k Keeper) SetBasket(ctx sdk.Context, basket types.Basket) {
	idBz := sdk.Uint64ToBigEndian(basket.Id)
	bz := k.cdc.MustMarshal(&basket)
	store := ctx.KVStore(k.storeKey)
	store.Set(append(types.PrefixBasketKey, idBz...), bz)
	store.Set(append(types.PrefixBasketByDenomKey, basket.GetBasketDenom()...), idBz)
}

func (k Keeper) DeleteBasket(ctx sdk.Context, basket types.Basket) {
	idBz := sdk.Uint64ToBigEndian(basket.Id)
	store := ctx.KVStore(k.storeKey)
	store.Delete(append(types.PrefixBasketKey, idBz...))
	store.Delete(append(types.PrefixBasketByDenomKey, basket.GetBasketDenom()...))
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

func (k Keeper) CreateBasket(ctx sdk.Context, basket types.Basket) error {
	// create a new basket id
	basketId := k.GetLastBasketId(ctx) + 1
	k.SetLastBasketId(ctx, basketId)
	basket.Id = basketId
	basket.Surplus = sdk.Coins{} // surplus is zero at initial

	if len(basket.Tokens) == 0 {
		return types.ErrEmptyUnderlyingTokens
	}
	usedDenom := make(map[string]bool)
	for index, token := range basket.Tokens {
		// ensure tokens amount is zero
		basket.Tokens[index].Amount = sdk.ZeroInt()
		// validate denom for the token
		if err := sdk.ValidateDenom(token.Denom); err != nil {
			return err
		}
		if token.Weight.IsZero() {
			return types.ErrTokenWeightShouldNotBeZero
		}
		if usedDenom[token.Denom] {
			return types.ErrDuplicateDenomExistsOnTokens
		}
		usedDenom[token.Denom] = true
	}

	k.SetBasket(ctx, basket)
	return nil
}

func (k Keeper) EditBasket(ctx sdk.Context, basket types.Basket) error {
	oldBasket, err := k.GetBasketById(ctx, basket.Id)
	if err != nil {
		return err
	}

	// TODO: what happens if suffix change?
	// TODO: what happens if a basket token is removed?

	// use previous surplus
	basket.Surplus = oldBasket.Surplus

	prevAmounts := make(map[string]sdk.Int)
	for _, token := range oldBasket.Tokens {
		prevAmounts[token.Denom] = token.Amount
	}

	if len(basket.Tokens) == 0 {
		return types.ErrEmptyUnderlyingTokens
	}

	usedDenom := make(map[string]bool)
	rates, _ := basket.RatesAndIndexes()
	basketDenomSupplyEst := sdk.ZeroDec()
	for index, token := range basket.Tokens {
		// ensure tokens amount is derivated from previous by denom
		if !prevAmounts[token.Denom].IsNil() {
			basket.Tokens[index].Amount = prevAmounts[token.Denom]
		} else {
			basket.Tokens[index].Amount = sdk.ZeroInt()
		}

		// validate denom for the token
		if err := sdk.ValidateDenom(token.Denom); err != nil {
			return err
		}
		if token.Weight.IsZero() {
			return types.ErrTokenWeightShouldNotBeZero
		}
		if usedDenom[token.Denom] {
			return types.ErrDuplicateDenomExistsOnTokens
		}
		usedDenom[token.Denom] = true
		basketDenomSupplyEst = basketDenomSupplyEst.
			Add(basket.Tokens[index].Amount.ToLegacyDec().Mul(rates[token.Denom]))
	}

	supply := k.bk.GetSupply(ctx, basket.GetBasketDenom())
	if supply.Amount.GT(basketDenomSupplyEst.TruncateInt()) {
		return types.ErrBasketDenomSupplyTooBig
	}

	k.SetBasket(ctx, basket)
	return nil
}
