package keeper

import (
	"github.com/KiraCore/sekai/x/layer2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetDappSessionApproval(ctx sdk.Context, session types.DappSessionApproval) {
	bz := k.cdc.MustMarshal(&session)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DappSessionApprovalKey(session.DappName, session.Approver), bz)
}

func (k Keeper) DeleteDappSessionApproval(ctx sdk.Context, name, session string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.DappSessionApprovalKey(name, session))
}

func (k Keeper) GetDappSessionApproval(ctx sdk.Context, name, verifier string) types.DappSessionApproval {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DappSessionApprovalKey(name, verifier))
	if bz == nil {
		return types.DappSessionApproval{}
	}

	sessionInfo := types.DappSessionApproval{}
	k.cdc.MustUnmarshal(bz, &sessionInfo)
	return sessionInfo
}

func (k Keeper) GetDappSessionApprovals(ctx sdk.Context, name string) []types.DappSessionApproval {
	store := ctx.KVStore(k.storeKey)

	sessions := []types.DappSessionApproval{}
	it := sdk.KVStorePrefixIterator(store, append([]byte(types.PrefixDappSessionApprovalKey), name...))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		session := types.DappSessionApproval{}
		k.cdc.MustUnmarshal(it.Value(), &session)
		sessions = append(sessions, session)
	}
	return sessions
}

func (k Keeper) GetAllDappSessionApprovals(ctx sdk.Context) []types.DappSessionApproval {
	store := ctx.KVStore(k.storeKey)

	sessions := []types.DappSessionApproval{}
	it := sdk.KVStorePrefixIterator(store, []byte(types.PrefixDappSessionApprovalKey))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		session := types.DappSessionApproval{}
		k.cdc.MustUnmarshal(it.Value(), &session)
		sessions = append(sessions, session)
	}
	return sessions
}
