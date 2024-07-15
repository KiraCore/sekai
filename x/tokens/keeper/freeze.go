package keeper

import (
	"github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AddTokensToBlacklist(ctx sdk.Context, tokens []string) {
	tokensBlackWhites := k.GetTokenBlackWhites(ctx)
	tokensBlackWhites.Blacklisted = addTokens(tokensBlackWhites.Blacklisted, tokens)
	k.SetTokenBlackWhites(ctx, tokensBlackWhites)
}

func (k Keeper) RemoveTokensFromBlacklist(ctx sdk.Context, tokens []string) {
	tokensBlackWhites := k.GetTokenBlackWhites(ctx)
	tokensBlackWhites.Blacklisted = removeTokens(tokensBlackWhites.Blacklisted, tokens)
	k.SetTokenBlackWhites(ctx, tokensBlackWhites)
}

func (k Keeper) AddTokensToWhitelist(ctx sdk.Context, tokens []string) {
	tokensBlackWhites := k.GetTokenBlackWhites(ctx)
	tokensBlackWhites.Whitelisted = addTokens(tokensBlackWhites.Whitelisted, tokens)
	k.SetTokenBlackWhites(ctx, tokensBlackWhites)
}

func (k Keeper) RemoveTokensFromWhitelist(ctx sdk.Context, tokens []string) {
	tokensBlackWhites := k.GetTokenBlackWhites(ctx)
	tokensBlackWhites.Whitelisted = removeTokens(tokensBlackWhites.Whitelisted, tokens)
	k.SetTokenBlackWhites(ctx, tokensBlackWhites)
}

func (k Keeper) SetTokenBlackWhites(ctx sdk.Context, tokensBlackWhite types.TokensWhiteBlack) {
	store := ctx.KVStore(k.storeKey)
	store.Set(PrefixKeyTokenBlackWhite, k.cdc.MustMarshal(&tokensBlackWhite))
}

func (k Keeper) GetTokenBlackWhites(ctx sdk.Context) types.TokensWhiteBlack {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(PrefixKeyTokenBlackWhite)
	if bz == nil {
		return types.TokensWhiteBlack{}
	}

	tokensBlackWhite := types.TokensWhiteBlack{}
	k.cdc.MustUnmarshal(bz, &tokensBlackWhite)

	return tokensBlackWhite
}
