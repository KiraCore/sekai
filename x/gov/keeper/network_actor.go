package keeper

import (
	"github.com/KiraCore/sekai/x/gov/types"
	types2 "github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SaveNetworkActor(ctx sdk.Context, actor types.NetworkActor) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), NetworkActorsPrefix)

	bz := k.cdc.MustMarshalBinaryBare(&actor)
	prefixStore.Set(actor.Address.Bytes(), bz)
}

func (k Keeper) GetNetworkActorByAddress(ctx sdk.Context, address sdk.AccAddress) (types.NetworkActor, bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), NetworkActorsPrefix)

	bz := prefixStore.Get(address.Bytes())
	if bz == nil {
		return types.NetworkActor{}, false
	}

	var na types.NetworkActor
	k.cdc.MustUnmarshalBinaryBare(bz, &na)

	return na, true
}

// AddWhitelistPermission whitelist a permission to an address. It saves the actor after it.
func (k Keeper) AddWhitelistPermission(ctx sdk.Context, actor types.NetworkActor, perm types.PermValue) error {
	err := actor.Permissions.AddToWhitelist(perm)
	if err != nil {
		return err
	}

	k.SaveNetworkActor(ctx, actor)

	store := ctx.KVStore(k.storeKey)
	store.Set(WhitelistAddressPermKey(actor.Address, perm), actor.Address.Bytes())

	return nil
}

// RemoveWhitelistPermission removes the whitelisted permission that an address has
func (k Keeper) RemoveWhitelistPermission(ctx sdk.Context, actor types.NetworkActor, perm types.PermValue) error {
	err := actor.Permissions.RemoveFromWhitelist(perm)
	if err != nil {
		return err
	}

	k.SaveNetworkActor(ctx, actor)

	store := ctx.KVStore(k.storeKey)
	store.Delete(WhitelistAddressPermKey(actor.Address, perm))

	return nil
}

func (k Keeper) AssignRoleToAddress(ctx sdk.Context, addr sdk.AccAddress, role types.Role) error {
	_, found := k.GetNetworkActorByAddress(ctx, addr)
	if !found {
		return types2.ErrNetworkActorNotFound
	}

	return nil
}

// GetNetworkActorByWhitelistedPermission returns all the actors that have Perm in whitelist.
func (k Keeper) GetNetworkActorByWhitelistedPermission(ctx sdk.Context, perm types.PermValue) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, WhitelistPermKey(perm))
}

// WhitelsitAddressPermKey returns the prefix key in format <0x31 + Perm_Bytes + address_bytes>
func WhitelistAddressPermKey(address sdk.AccAddress, perm types.PermValue) []byte {
	return append(WhitelistPermKey(perm), address.Bytes()...)
}

// WhitelistPermKey returns the prefix key in format <0x31 + Perm_Bytes>
func WhitelistPermKey(perm types.PermValue) []byte {
	return append(WhitelistActorPrefix, getPermBytes(perm)...)
}

// getPermBytes returns a PermValue in bytes representation.
func getPermBytes(perm types.PermValue) []byte {
	return sdk.Uint64ToBigEndian(uint64(perm))
}
