package keeper

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/tokens/types"
)

// GetTokenAlias returns a token alias
func (k Keeper) GetTokenAlias(ctx sdk.Context, symbol string) *types.TokenAlias {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), PrefixKeyTokenAlias)
	bz := prefixStore.Get([]byte(symbol))
	if bz == nil {
		return nil
	}

	alias := new(types.TokenAlias)
	k.cdc.MustUnmarshal(bz, alias)

	return alias
}

// ListTokenAlias returns all list of token alias
func (k Keeper) ListTokenAlias(ctx sdk.Context) []*types.TokenAlias {
	var tokenAliases []*types.TokenAlias

	// get iterator for token aliases
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, PrefixKeyTokenAlias)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		symbol := strings.TrimPrefix(string(iterator.Key()), string(PrefixKeyTokenAlias))
		tokenAlias := k.GetTokenAlias(ctx, symbol)
		if tokenAlias != nil {
			tokenAliases = append(tokenAliases, tokenAlias)
		}
	}
	return tokenAliases
}

// GetTokenAliasesByDenom returns all list of token alias
func (k Keeper) GetTokenAliasesByDenom(ctx sdk.Context, denoms []string) map[string]*types.TokenAlias {
	// get iterator for token aliases
	store := ctx.KVStore(k.storeKey)
	tokenAliasesMap := make(map[string]*types.TokenAlias)

	for _, denom := range denoms {
		denomTokenStoreID := append([]byte(PrefixKeyDenomToken), []byte(denom)...)

		if store.Has(denomTokenStoreID) {
			symbol := string(store.Get(denomTokenStoreID))
			tokenAlias := k.GetTokenAlias(ctx, symbol)
			tokenAliasesMap[denom] = tokenAlias
		}
	}
	return tokenAliasesMap
}

// UpsertTokenAlias upsert a token alias to the registry
func (k Keeper) UpsertTokenAlias(ctx sdk.Context, alias types.TokenAlias) error {
	store := ctx.KVStore(k.storeKey)
	// we use symbol of TokenAlias as an ID inside KVStore storage
	tokenAliasStoreID := append([]byte(PrefixKeyTokenAlias), []byte(alias.Symbol)...)

	for _, denom := range alias.Denoms {
		denomTokenStoreID := append([]byte(PrefixKeyDenomToken), []byte(denom)...)

		if store.Has(denomTokenStoreID) {
			symbol := string(store.Get(denomTokenStoreID))
			if symbol != alias.Symbol {
				return fmt.Errorf("%s denom is already registered for %s token alias", denom, symbol)
			}
		}
		store.Set(denomTokenStoreID, []byte(alias.Symbol))
	}
	store.Set(tokenAliasStoreID, k.cdc.MustMarshal(&alias))
	return nil
}

// DeleteTokenAlias delete token alias by symbol
func (k Keeper) DeleteTokenAlias(ctx sdk.Context, symbol string) error {
	store := ctx.KVStore(k.storeKey)
	// we use symbol of TokenAlias as an ID inside KVStore storage
	tokenAliasStoreID := append([]byte(PrefixKeyTokenAlias), []byte(symbol)...)

	if !store.Has(tokenAliasStoreID) {
		return fmt.Errorf("no alias is available for %s symbol", symbol)
	}

	var alias types.TokenAlias
	bz := store.Get(tokenAliasStoreID)
	k.cdc.MustUnmarshal(bz, &alias)

	for _, denom := range alias.Denoms {
		denomTokenStoreID := append([]byte(PrefixKeyDenomToken), []byte(denom)...)
		store.Delete(denomTokenStoreID)
	}

	store.Delete(tokenAliasStoreID)
	return nil
}
