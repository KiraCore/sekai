package keeper

import (
	"cosmossdk.io/math"
	"github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	for _, coin := range amt {
		tokenInfo := k.GetTokenInfo(ctx, coin.Denom)
		if tokenInfo == nil {
			tokenInfo = &types.TokenInfo{
				Denom:        coin.Denom,
				FeeRate:      math.LegacyZeroDec(),
				FeeEnabled:   false,
				Supply:       math.ZeroInt(),
				StakeCap:     math.LegacyZeroDec(),
				StakeMin:     math.OneInt(),
				StakeEnabled: false,
				Inactive:     false,
				Symbol:       coin.Denom,
				Name:         coin.Denom,
				Icon:         "",
				Decimals:     6,
			}
		}

		tokenInfo.Supply = tokenInfo.Supply.Add(coin.Amount)
		err := k.UpsertTokenInfo(ctx, *tokenInfo)
		if err != nil {
			return err
		}
	}
	return k.bankKeeper.MintCoins(ctx, moduleName, amt)
}
