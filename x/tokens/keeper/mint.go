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
			k.UpsertTokenInfo(ctx, types.TokenInfo{
				Denom:       coin.Denom,
				FeeRate:     math.LegacyZeroDec(),
				FeePayments: false,
				StakeCap:    math.LegacyZeroDec(),
				StakeMin:    math.OneInt(),
				StakeToken:  false,
				Invalidated: false,
				Symbol:      coin.Denom,
				Name:        coin.Denom,
				Icon:        "",
				Decimals:    6,
			})
		}
	}
	return k.bankKeeper.MintCoins(ctx, moduleName, amt)
}
