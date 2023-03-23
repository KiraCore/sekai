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

func (k Keeper) CreateNewSession(ctx sdk.Context, name string, prevLeader string) {
	leader := ""
	executors := k.GetDappExecutors(ctx, name)
	if len(executors) > 0 {
		leader = executors[0].Interx
		for index, executor := range executors {
			if executor.Interx == prevLeader {
				leader = executors[(index+1)%len(executors)].Interx
			}
		}
	}
	k.SetDappSession(ctx, types.DappSession{
		DappName:   name,
		Leader:     leader,
		Start:      uint64(ctx.BlockTime().Unix()),
		StatusHash: "",
		Status:     types.Unscheduled,
		// Unscheduled
		// Scheduled
		// Ongoing
	})

	// TODO: remove operators exiting when session ends
	// k.keeper.DeleteDappOperator(ctx, msg.DappName, msg.Sender)
	// TODO: If the operator leaving the dApp was a verifier then
	// as the result of the exit tx his LP tokens bond should be returned once the record is deleted.
	// The bond can only be claimed if and only if the status didn’t change to jailed in the meantime.
}
