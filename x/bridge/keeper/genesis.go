package keeper

import (
	"github.com/KiraCore/sekai/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetBridgeAddress(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.BridgeAddressKey, []byte(address))
}

func (k Keeper) GetBridgeAddress(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.BridgeAddressKey)

	if bz == nil {
		return ""
	}

	return string(bz)
}
