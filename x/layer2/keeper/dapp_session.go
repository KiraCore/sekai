package keeper

import (
	"github.com/KiraCore/sekai/x/layer2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetDappSession(ctx sdk.Context, session types.DappSession) {
	bz := k.cdc.MustMarshal(&session)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DappSessionKey(session.DappName), bz)
}

func (k Keeper) DeleteDappSession(ctx sdk.Context, name string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.DappSessionKey(name))
}

func (k Keeper) GetDappSession(ctx sdk.Context, name string) types.DappSession {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DappSessionKey(name))
	if bz == nil {
		return types.DappSession{}
	}

	sessionInfo := types.DappSession{}
	k.cdc.MustUnmarshal(bz, &sessionInfo)
	return sessionInfo
}

func (k Keeper) GetDappSessions(ctx sdk.Context, name string) []types.DappSession {
	store := ctx.KVStore(k.storeKey)

	sessions := []types.DappSession{}
	it := sdk.KVStorePrefixIterator(store, append([]byte(types.PrefixDappSessionKey), name...))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		session := types.DappSession{}
		k.cdc.MustUnmarshal(it.Value(), &session)
		sessions = append(sessions, session)
	}
	return sessions
}

func (k Keeper) GetAllDappSessions(ctx sdk.Context) []types.DappSession {
	store := ctx.KVStore(k.storeKey)

	sessions := []types.DappSession{}
	it := sdk.KVStorePrefixIterator(store, []byte(types.PrefixDappSessionKey))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		session := types.DappSession{}
		k.cdc.MustUnmarshal(it.Value(), &session)
		sessions = append(sessions, session)
	}
	return sessions
}
