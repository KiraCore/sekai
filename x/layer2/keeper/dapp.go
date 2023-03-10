package keeper

import (
	"github.com/KiraCore/sekai/x/layer2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetDapp(ctx sdk.Context, name string) types.Dapp {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DappKey(name))
	if bz == nil {
		return types.Dapp{}
	}

	dapp := types.Dapp{}
	k.cdc.MustUnmarshal(bz, &dapp)
	return dapp
}

func (k Keeper) GetAllDapps(ctx sdk.Context) []types.Dapp {
	store := ctx.KVStore(k.storeKey)

	dapps := []types.Dapp{}
	it := sdk.KVStorePrefixIterator(store, types.KeyPrefixDapp)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		dapp := types.Dapp{}
		k.cdc.MustUnmarshal(it.Value(), &dapp)
		dapps = append(dapps, dapp)
	}
	return dapps
}

func (k Keeper) SetDapp(ctx sdk.Context, dapp types.Dapp) {
	bz := k.cdc.MustMarshal(&dapp)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DappKey(dapp.Name), bz)
}

func (k Keeper) DeleteDapp(ctx sdk.Context, name string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.DappKey(name))
}

func (k Keeper) SetUserDappBond(ctx sdk.Context, bond types.UserDappBond) {
	bz := k.cdc.MustMarshal(&bond)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.UserDappBondKey(bond.DappName, bond.User), bz)
}

func (k Keeper) DeleteUserDappBond(ctx sdk.Context, name, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.UserDappBondKey(name, address))
}

func (k Keeper) GetUserDappBond(ctx sdk.Context, name string, user string) types.UserDappBond {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.UserDappBondKey(name, user))
	if bz == nil {
		return types.UserDappBond{}
	}

	bond := types.UserDappBond{}
	k.cdc.MustUnmarshal(bz, &bond)
	return bond
}

func (k Keeper) GetUserDappBonds(ctx sdk.Context, name string) []types.UserDappBond {
	store := ctx.KVStore(k.storeKey)

	bondlist := []types.UserDappBond{}
	it := sdk.KVStorePrefixIterator(store, append([]byte(types.PrefixUserDappBondKey), name...))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		bond := types.UserDappBond{}
		k.cdc.MustUnmarshal(it.Value(), &bond)
		bondlist = append(bondlist, bond)
	}
	return bondlist
}

func (k Keeper) GetAllUserDappBonds(ctx sdk.Context) []types.UserDappBond {
	store := ctx.KVStore(k.storeKey)

	bondlist := []types.UserDappBond{}
	it := sdk.KVStorePrefixIterator(store, []byte(types.PrefixUserDappBondKey))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		bond := types.UserDappBond{}
		k.cdc.MustUnmarshal(it.Value(), &bond)
		bondlist = append(bondlist, bond)
	}
	return bondlist
}

func (k Keeper) ExecuteDappRemove(ctx sdk.Context, dapp types.Dapp) error {
	for _, userBond := range k.GetUserDappBonds(ctx, dapp.Name) {
		// send tokens back to user
		addr := sdk.MustAccAddressFromBech32(userBond.User)
		err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.Coins{userBond.Bond})
		if err != nil {
			return err
		}
		k.DeleteUserDappBond(ctx, dapp.Name, userBond.User)
	}
	k.DeleteDapp(ctx, dapp.Name)
	return nil
}
