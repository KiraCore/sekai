package keeper

import (
	"github.com/KiraCore/sekai/x/ethereum/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetRelay(ctx sdk.Context, record *types.MsgRelay) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyRelay), record.Address...)

	store.Set(key, k.cdc.MustMarshal(record))
}

func (k Keeper) GetRelayByAddress(ctx sdk.Context, address sdk.AccAddress) *types.MsgRelay {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PrefixKeyRelay))
	bz := prefixStore.Get(address)

	if bz == nil {
		return nil
	}

	info := new(types.MsgRelay)
	k.cdc.MustUnmarshal(bz, info)

	return info
}
