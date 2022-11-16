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

func (k Keeper) GetAllCollectives(ctx sdk.Context) []types.Collective {
	store := ctx.KVStore(k.storeKey)

	collectives := []types.Collective{}
	it := sdk.KVStorePrefixIterator(store, types.PrefixCollectiveKey)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		collective := types.Collective{}
		k.cdc.MustUnmarshal(it.Value(), &collective)
		collectives = append(collectives, collective)
	}
	return collectives
}

func (k Keeper) SetCollective(ctx sdk.Context, collective types.Collective) {
	bz := k.cdc.MustMarshal(&collective)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CollectiveKey(collective.Name), bz)
}

func (k Keeper) DeleteCollective(ctx sdk.Context, name string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.CollectiveKey(name))
}

func (k Keeper) SetCollectiveContributer(ctx sdk.Context, cc types.CollectiveContributor) {
	bz := k.cdc.MustMarshal(&cc)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CollectiveContributerKey(cc.Name, cc.Address), bz)
}

func (k Keeper) DeleteCollectiveContributer(ctx sdk.Context, name, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.CollectiveContributerKey(name, address))
}

func (k Keeper) GetCollectiveContributer(ctx sdk.Context, name string, contributer string) types.CollectiveContributor {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CollectiveContributerKey(name, contributer))
	if bz == nil {
		return types.CollectiveContributor{}
	}

	cc := types.CollectiveContributor{}
	k.cdc.MustUnmarshal(bz, &cc)
	return cc
}

func (k Keeper) GetAllCollectiveContributers(ctx sdk.Context, name string) []types.CollectiveContributor {
	store := ctx.KVStore(k.storeKey)

	cclist := []types.CollectiveContributor{}
	it := sdk.KVStorePrefixIterator(store, append([]byte(types.PrefixCollectiveContributerKey), name...))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		cc := types.CollectiveContributor{}
		k.cdc.MustUnmarshal(it.Value(), &cc)
		cclist = append(cclist, cc)
	}
	return cclist
}

func (k Keeper) SendDonation(ctx sdk.Context, name string, account sdk.AccAddress, coins sdk.Coins) error {
	collective := k.GetCollective(ctx, name)

	donations := sdk.Coins(collective.Donations)
	if donations.IsAllGTE(coins) {
		collective.Donations = donations.Sub(coins)
	}
	k.SetCollective(ctx, collective)

	err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, account, coins)
	if err != nil {
		return err
	}
	return nil
}
