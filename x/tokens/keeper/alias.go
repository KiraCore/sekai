package keeper

import (
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
	k.cdc.MustUnmarshalBinaryBare(bz, alias)

	return alias
}

// ListTokenAlias returns all list of token alias
func (k Keeper) ListTokenAlias(ctx sdk.Context) []types.TokenAlias {

}

// UpsertTokenAlias upsert a token alias to the registry
func (k Keeper) UpsertTokenAlias(ctx sdk.Context, alias types.TokenAlias) error {
	store := ctx.KVStore(k.storeKey)
	// we use symbol of TokenAlias as an ID inside KVStore storage
	tokenAliasStoreID := append([]byte(PrefixKeyTokenAlias), []byte(alias.Symbol)...)

	store.Set(tokenAliasStoreID, k.cdc.MustMarshalBinaryBare(alias))
	return nil
}

// DeleteTokenAlias delete token alias by symbol
func (k Keeper) DeleteTokenAlias(ctx sdk.Context, symbol string) error {
	store := ctx.KVStore(k.storeKey)
	// we use symbol of TokenAlias as an ID inside KVStore storage
	tokenAliasStoreID := append([]byte(PrefixKeyTokenAlias), []byte(symbol)...)

	store.Delete(tokenAliasStoreID)
	return nil
}
