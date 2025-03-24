package keeper

import (
	"github.com/KiraCore/sekai/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetChangeCosmosEthereumRecord(ctx sdk.Context, record types.ChangeCosmosEthereumRecord) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyBridgeCosmosEthereumRecord), record.From...)

	store.Set(key, k.cdc.MustMarshal(&record))
}

func (k Keeper) SetChangeEthereumCosmosRecord(ctx sdk.Context, record types.ChangeEthereumCosmosRecord) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyBridgeEthereumCosmosRecord), record.To...)

	store.Set(key, k.cdc.MustMarshal(&record))
}

func (k Keeper) GetChangeCosmosEthereumRecord(ctx sdk.Context, address sdk.AccAddress) *types.ChangeCosmosEthereumRecord {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PrefixKeyBridgeCosmosEthereumRecord))
	bz := prefixStore.Get(address)

	if bz == nil {
		return nil
	}

	info := new(types.ChangeCosmosEthereumRecord)
	k.cdc.MustUnmarshal(bz, info)

	return info
}

func (k Keeper) GetChangeEthereumCosmosRecord(ctx sdk.Context, address sdk.AccAddress) *types.ChangeEthereumCosmosRecord {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PrefixKeyBridgeCosmosEthereumRecord))
	bz := prefixStore.Get(address)

	if bz == nil {
		return nil
	}

	info := new(types.ChangeEthereumCosmosRecord)
	k.cdc.MustUnmarshal(bz, info)

	return info
}
