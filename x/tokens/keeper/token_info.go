package keeper

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/tokens/types"
)

// GetTokenInfo returns a token rate
func (k Keeper) GetTokenInfo(ctx sdk.Context, denom string) *types.TokenInfo {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), PrefixKeyTokenInfo)
	bz := prefixStore.Get([]byte(denom))
	if bz == nil {
		return nil
	}

	rate := new(types.TokenInfo)
	k.cdc.MustUnmarshal(bz, rate)

	return rate
}

// GetAllTokenInfos returns all list of token rate
func (k Keeper) GetAllTokenInfos(ctx sdk.Context) []types.TokenInfo {
	var tokenRates []types.TokenInfo

	// get iterator for token rates
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, PrefixKeyTokenInfo)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		info := types.TokenInfo{}
		k.cdc.MustUnmarshal(iterator.Value(), &info)
		tokenRates = append(tokenRates, info)
	}
	return tokenRates
}

// GetTokenInfosByDenom returns all list of token rate
func (k Keeper) GetTokenInfosByDenom(ctx sdk.Context, denoms []string) map[string]types.TokenInfoResponse {
	tokenRatesMap := make(map[string]types.TokenInfoResponse)

	for _, denom := range denoms {
		tokenRate := k.GetTokenInfo(ctx, denom)
		supply := k.bankKeeper.GetSupply(ctx, denom)
		tokenRatesMap[denom] = types.TokenInfoResponse{
			Data:   tokenRate,
			Supply: supply,
		}
	}
	return tokenRatesMap
}

// UpsertTokenInfo upsert a token rate to the registry
func (k Keeper) UpsertTokenInfo(ctx sdk.Context, rate types.TokenInfo) error {
	store := ctx.KVStore(k.storeKey)
	// we use denom of TokenInfo as an ID inside KVStore storage
	tokenRateStoreID := append([]byte(PrefixKeyTokenInfo), []byte(rate.Denom)...)
	if rate.Denom == k.DefaultDenom(ctx) && store.Has(tokenRateStoreID) {
		return errors.New("bond denom rate is read-only")
	}

	store.Set(tokenRateStoreID, k.cdc.MustMarshal(&rate))

	totalRewardsCap := sdk.ZeroDec()
	rates := k.GetAllTokenInfos(ctx)
	for _, rate := range rates {
		totalRewardsCap = totalRewardsCap.Add(rate.StakeCap)
	}
	if totalRewardsCap.GT(sdk.OneDec()) {
		return types.ErrTotalRewardsCapExceeds100Percent
	}
	return nil
}

// DeleteTokenInfo delete token denom by denom
func (k Keeper) DeleteTokenInfo(ctx sdk.Context, denom string) error {
	store := ctx.KVStore(k.storeKey)
	// we use symbol of DeleteTokenInfo as an ID inside KVStore storage
	tokenRateStoreID := append([]byte(PrefixKeyTokenInfo), []byte(denom)...)

	if !store.Has(tokenRateStoreID) {
		return fmt.Errorf("no rate registry is available for %s denom", denom)
	}

	store.Delete(tokenRateStoreID)
	return nil
}
