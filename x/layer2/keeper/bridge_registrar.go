package keeper

import (
	"github.com/KiraCore/sekai/x/layer2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetBridgeRegistrarHelper(ctx sdk.Context, helper types.BridgeRegistrarHelper) {
	bz := k.cdc.MustMarshal(&helper)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.BridgeRegistrarHelperKey, bz)
}

func (k Keeper) GetBridgeRegistrarHelper(ctx sdk.Context) types.BridgeRegistrarHelper {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.BridgeRegistrarHelperKey)
	if bz == nil {
		return types.BridgeRegistrarHelper{}
	}

	helper := types.BridgeRegistrarHelper{}
	k.cdc.MustUnmarshal(bz, &helper)
	return helper
}

func (k Keeper) SetBridgeAccount(ctx sdk.Context, account types.BridgeAccount) {
	bz := k.cdc.MustMarshal(&account)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.BridgeAccountKey(account.Index), bz)
}

func (k Keeper) GetBridgeAccount(ctx sdk.Context, index uint64) types.BridgeAccount {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.BridgeAccountKey(index))
	if bz == nil {
		return types.BridgeAccount{}
	}

	account := types.BridgeAccount{}
	k.cdc.MustUnmarshal(bz, &account)
	return account
}

func (k Keeper) GetBridgeAccounts(ctx sdk.Context) []types.BridgeAccount {
	store := ctx.KVStore(k.storeKey)

	accounts := []types.BridgeAccount{}
	it := sdk.KVStorePrefixIterator(store, []byte(types.PrefixBridgeAccountKey))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		account := types.BridgeAccount{}
		k.cdc.MustUnmarshal(it.Value(), &account)
		accounts = append(accounts, account)
	}
	return accounts
}

func (k Keeper) SetBridgeToken(ctx sdk.Context, token types.BridgeToken) {
	bz := k.cdc.MustMarshal(&token)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.BridgeTokenKey(token.Index), bz)
}

func (k Keeper) GetBridgeToken(ctx sdk.Context, index uint64) types.BridgeToken {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.BridgeTokenKey(index))
	if bz == nil {
		return types.BridgeToken{}
	}

	token := types.BridgeToken{}
	k.cdc.MustUnmarshal(bz, &token)
	return token
}

func (k Keeper) GetBridgeTokens(ctx sdk.Context) []types.BridgeToken {
	store := ctx.KVStore(k.storeKey)

	tokens := []types.BridgeToken{}
	it := sdk.KVStorePrefixIterator(store, []byte(types.PrefixBridgeTokenKey))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		token := types.BridgeToken{}
		k.cdc.MustUnmarshal(it.Value(), &token)
		tokens = append(tokens, token)
	}
	return tokens
}

func (k Keeper) SetXAM(ctx sdk.Context, xam types.XAM) {
	bz := k.cdc.MustMarshal(&xam)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.XAMKey(xam.Res.Xid), bz)
}

func (k Keeper) GetXAM(ctx sdk.Context, xid uint64) types.XAM {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.XAMKey(xid))
	if bz == nil {
		return types.XAM{}
	}

	xam := types.XAM{}
	k.cdc.MustUnmarshal(bz, &xam)
	return xam
}

func (k Keeper) GetXAMs(ctx sdk.Context) []types.XAM {
	store := ctx.KVStore(k.storeKey)

	xams := []types.XAM{}
	it := sdk.KVStorePrefixIterator(store, []byte(types.PrefixXAMKey))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		xam := types.XAM{}
		k.cdc.MustUnmarshal(it.Value(), &xam)
		xams = append(xams, xam)
	}
	return xams
}
