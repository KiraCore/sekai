package keeper

import (
	"github.com/KiraCore/sekai/x/layer2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetDappLeaderDenouncement(ctx sdk.Context, denouncement types.DappLeaderDenouncement) {
	bz := k.cdc.MustMarshal(&denouncement)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DappLeaderDenouncementKey(denouncement.DappName, denouncement.Leader, denouncement.Sender), bz)
}

func (k Keeper) DeleteDappLeaderDenouncement(ctx sdk.Context, name, denouncement, sender string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.DappLeaderDenouncementKey(name, denouncement, sender))
}

func (k Keeper) GetDappLeaderDenouncement(ctx sdk.Context, name, denouncement, sender string) types.DappLeaderDenouncement {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DappLeaderDenouncementKey(name, denouncement, sender))
	if bz == nil {
		return types.DappLeaderDenouncement{}
	}

	denouncementInfo := types.DappLeaderDenouncement{}
	k.cdc.MustUnmarshal(bz, &denouncementInfo)
	return denouncementInfo
}

func (k Keeper) GetDappLeaderDenouncements(ctx sdk.Context, name, denouncement string) []types.DappLeaderDenouncement {
	store := ctx.KVStore(k.storeKey)

	denouncements := []types.DappLeaderDenouncement{}
	it := sdk.KVStorePrefixIterator(store, append(append([]byte(types.PrefixDappLeaderDenouncementKey), name...), denouncement...))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		denouncement := types.DappLeaderDenouncement{}
		k.cdc.MustUnmarshal(it.Value(), &denouncement)
		denouncements = append(denouncements, denouncement)
	}
	return denouncements
}

func (k Keeper) GetAllDappLeaderDenouncements(ctx sdk.Context) []types.DappLeaderDenouncement {
	store := ctx.KVStore(k.storeKey)

	denouncements := []types.DappLeaderDenouncement{}
	it := sdk.KVStorePrefixIterator(store, []byte(types.PrefixDappLeaderDenouncementKey))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		denouncement := types.DappLeaderDenouncement{}
		k.cdc.MustUnmarshal(it.Value(), &denouncement)
		denouncements = append(denouncements, denouncement)
	}
	return denouncements
}
