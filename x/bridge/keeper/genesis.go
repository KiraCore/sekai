package keeper

import (
	"github.com/KiraCore/sekai/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetBridgeAddress(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.BridgeAddressKey, []byte(address))
}

func (k Keeper) SetCosmosEthereumExchangeRate(ctx sdk.Context, rate int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.BridgeCosmosEthereumExchangeRateKey, sdk.Uint64ToBigEndian(uint64(rate)))
}

func (k Keeper) SetEthereumCosmosExchangeRate(ctx sdk.Context, rate int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.BridgeEthereumCosmosExchangeRateKey, sdk.Uint64ToBigEndian(uint64(rate)))
}

func (k Keeper) GetBridgeAddress(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.BridgeAddressKey)

	if bz == nil {
		return ""
	}

	return string(bz)
}

func (k Keeper) GetCosmosEthereumExchangeRate(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.BridgeCosmosEthereumExchangeRateKey)

	if bz == nil {
		return 1
	}

	return int64(sdk.BigEndianToUint64(bz))
}

func (k Keeper) GetEthereumCosmosExchangeRate(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.BridgeCosmosEthereumExchangeRateKey)

	if bz == nil {
		return 1
	}

	return int64(sdk.BigEndianToUint64(bz))
}
