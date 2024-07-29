package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/tokens/types"
)

// GetTokenInfo returns a token info
func (k Keeper) GetTokenInfo(ctx sdk.Context, denom string) *types.TokenInfo {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), PrefixKeyTokenInfo)
	bz := prefixStore.Get([]byte(denom))
	if bz == nil {
		return nil
	}

	info := new(types.TokenInfo)
	k.cdc.MustUnmarshal(bz, info)

	return info
}

// GetAllTokenInfos returns all list of token info
func (k Keeper) GetAllTokenInfos(ctx sdk.Context) []types.TokenInfo {
	var tokenInfos []types.TokenInfo

	// get iterator for token infos
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, PrefixKeyTokenInfo)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		info := types.TokenInfo{}
		k.cdc.MustUnmarshal(iterator.Value(), &info)
		tokenInfos = append(tokenInfos, info)
	}
	return tokenInfos
}

// GetTokenInfosByDenom returns all list of token info
func (k Keeper) GetTokenInfosByDenom(ctx sdk.Context, denoms []string) map[string]types.TokenInfoResponse {
	tokenInfosMap := make(map[string]types.TokenInfoResponse)

	for _, denom := range denoms {
		tokenInfo := k.GetTokenInfo(ctx, denom)
		supply := k.bankKeeper.GetSupply(ctx, denom)
		tokenInfosMap[denom] = types.TokenInfoResponse{
			Data:   tokenInfo,
			Supply: supply,
		}
	}
	return tokenInfosMap
}

// UpsertTokenInfo upsert a token info to the registry
func (k Keeper) UpsertTokenInfo(ctx sdk.Context, info types.TokenInfo) error {
	store := ctx.KVStore(k.storeKey)
	// we use denom of TokenInfo as an ID inside KVStore storage
	tokenInfoStoreID := append([]byte(PrefixKeyTokenInfo), []byte(info.Denom)...)
	// if info.Denom == k.DefaultDenom(ctx) && store.Has(tokenInfoStoreID) {
	// 	return types.ErrBondDenomIsReadOnly
	// }

	if !info.SupplyCap.IsNil() && info.SupplyCap.IsPositive() && info.Supply.GT(info.SupplyCap) {
		return types.ErrCannotExceedTokenCap
	}

	store.Set(tokenInfoStoreID, k.cdc.MustMarshal(&info))

	totalRewardsCap := sdk.ZeroDec()
	infos := k.GetAllTokenInfos(ctx)
	for _, info := range infos {
		totalRewardsCap = totalRewardsCap.Add(info.StakeCap)
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
	tokenInfoStoreID := append([]byte(PrefixKeyTokenInfo), []byte(denom)...)

	if !store.Has(tokenInfoStoreID) {
		return fmt.Errorf("no token info registry is available for %s denom", denom)
	}

	store.Delete(tokenInfoStoreID)
	return nil
}
