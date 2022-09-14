package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// constants
var (
	ModuleName = "basket"
	// RouterKey to be used for routing msgs
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	KeyLastBasketId        = []byte("last_basket_id")
	PrefixBasketKey        = []byte("basket_by_id")
	PrefixBasketByDenomKey = []byte("basket_by_denom")
	PrefixBasketMintByTime = []byte("basket_mint_by_time")
	PrefixBasketBurnByTime = []byte("basket_burn_by_time")
	PrefixBasketSwapByTime = []byte("basket_swap_by_time")
)

func BasketMintByTimeKey(basketId uint64, time time.Time) []byte {
	return append(append([]byte(PrefixBasketMintByTime), sdk.Uint64ToBigEndian(basketId)...), sdk.FormatTimeBytes(time)...)
}

func BasketBurnByTimeKey(basketId uint64, time time.Time) []byte {
	return append(append([]byte(PrefixBasketBurnByTime), sdk.Uint64ToBigEndian(basketId)...), sdk.FormatTimeBytes(time)...)
}

func BasketSwapByTimeKey(basketId uint64, time time.Time) []byte {
	return append(append([]byte(PrefixBasketSwapByTime), sdk.Uint64ToBigEndian(basketId)...), sdk.FormatTimeBytes(time)...)
}
