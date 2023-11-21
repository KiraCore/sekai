package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetPoolPrefix(poolID uint64) string {
	return fmt.Sprintf("v%d/", poolID)
}
func GetPoolCoins(pool StakingPool, coins sdk.Coins) sdk.Coins {
	prefix := GetPoolPrefix(pool.Id)
	poolCoins := sdk.Coins{}
	for _, coin := range coins {
		poolCoins = poolCoins.Add(sdk.NewCoin(prefix+coin.Denom, sdk.NewDecFromInt(coin.Amount).Mul(sdk.OneDec().Sub(pool.Slashed)).RoundInt()))
	}
	return poolCoins
}
func GetShareDenom(poolID uint64, denom string) string {
	prefix := GetPoolPrefix(poolID)
	return prefix + denom
}
func GetNativeDenom(poolID uint64, denom string) string {
	return strings.TrimPrefix(denom, GetPoolPrefix(poolID))
}
