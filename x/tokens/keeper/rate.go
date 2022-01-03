package keeper

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/tokens/types"
)

// GetTokenRate returns a token rate
func (k Keeper) GetTokenRate(ctx sdk.Context, denom string) *types.TokenRate {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), PrefixKeyTokenRate)
	bz := prefixStore.Get([]byte(denom))
	if bz == nil {
		return nil
	}

	rate := new(types.TokenRate)
	k.cdc.MustUnmarshal(bz, rate)

	return rate
}

// ListTokenRate returns all list of token rate
func (k Keeper) ListTokenRate(ctx sdk.Context) []*types.TokenRate {
	var tokenRates []*types.TokenRate

	// get iterator for token rates
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, PrefixKeyTokenRate)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		denom := strings.TrimPrefix(string(iterator.Key()), string(PrefixKeyTokenRate))
		tokenRate := k.GetTokenRate(ctx, denom)
		if tokenRate != nil {
			tokenRates = append(tokenRates, tokenRate)
		}
	}
	return tokenRates
}

// GetTokenRatesByDenom returns all list of token rate
func (k Keeper) GetTokenRatesByDenom(ctx sdk.Context, denoms []string) map[string]*types.TokenRate {
	// get iterator for token aliases
	store := ctx.KVStore(k.storeKey)
	tokenRatesMap := make(map[string]*types.TokenRate)

	for _, denom := range denoms {
		denomTokenStoreID := append([]byte(PrefixKeyTokenRate), []byte(denom)...)

		if store.Has(denomTokenStoreID) {
			tokenRate := k.GetTokenRate(ctx, denom)
			tokenRatesMap[denom] = tokenRate
		}
	}
	return tokenRatesMap
}

// UpsertTokenRate upsert a token rate to the registry
func (k Keeper) UpsertTokenRate(ctx sdk.Context, rate types.TokenRate) error {
	store := ctx.KVStore(k.storeKey)
	// we use denom of TokenRate as an ID inside KVStore storage
	tokenRateStoreID := append([]byte(PrefixKeyTokenRate), []byte(rate.Denom)...)
	if rate.Denom == k.BondDenom(ctx) && store.Has(tokenRateStoreID) {
		return errors.New("bond denom rate is read-only")
	}
	store.Set(tokenRateStoreID, k.cdc.MustMarshal(&rate))
	return nil
}

// DeleteTokenRate delete token denom by denom
func (k Keeper) DeleteTokenRate(ctx sdk.Context, denom string) error {
	store := ctx.KVStore(k.storeKey)
	// we use symbol of DeleteTokenRate as an ID inside KVStore storage
	tokenRateStoreID := append([]byte(PrefixKeyTokenRate), []byte(denom)...)

	if !store.Has(tokenRateStoreID) {
		return fmt.Errorf("no rate registry is available for %s denom", denom)
	}

	store.Delete(tokenRateStoreID)
	return nil
}
