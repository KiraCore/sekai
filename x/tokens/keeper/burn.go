package keeper

import (
	"github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	for _, coin := range amt {
		tokenInfo := k.GetTokenInfo(ctx, coin.Denom)
		if tokenInfo == nil {
			return types.ErrTokenNotRegistered
		}
		tokenInfo.Supply = tokenInfo.Supply.Sub(coin.Amount)
		err := k.UpsertTokenInfo(ctx, *tokenInfo)
		if err != nil {
			return err
		}
	}
	return k.bankKeeper.BurnCoins(ctx, moduleName, amt)
}
