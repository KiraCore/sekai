package keeper

import (
	"github.com/KiraCore/sekai/x/collectives/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetCollective(ctx sdk.Context, name string) types.Collective {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CollectiveKey(name))
	if bz == nil {
		return types.Collective{}
	}

	collective := types.Collective{}
	k.cdc.MustUnmarshal(bz, &collective)
	return collective
}

func (k Keeper) SetCollective(ctx sdk.Context, collective types.Collective) {
	bz := k.cdc.MustMarshal(&collective)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CollectiveKey(collective.Name), bz)
}
