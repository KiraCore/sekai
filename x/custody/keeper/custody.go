package keeper

import (
	"github.com/KiraCore/sekai/x/custody/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetCustodyInfoByAddress(ctx sdk.Context, address sdk.AccAddress) *types.CustodySettings {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PrefixKeyCustodyRecord))
	bz := prefixStore.Get(address)

	if bz == nil {
		return nil
	}

	info := new(types.CustodySettings)
	k.cdc.MustUnmarshal(bz, info)

	return info
}

func (k Keeper) SetCustodyRecord(ctx sdk.Context, record types.CustodyRecord) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyCustodyRecord), record.Address...)

	store.Set(key, k.cdc.MustMarshal(&record.CustodySettings))
}
