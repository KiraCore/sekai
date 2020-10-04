package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/tokens/types"
)

// UpsertTokenAlias upsert a token alias to the registry
func (k Keeper) UpsertTokenAlias(ctx sdk.Context, alias types.TokenAlias) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), KeyPrefixPermissionsRegistry)

	prefixStore.Set(types.RoleToKey(role), k.cdc.MustMarshalBinaryBare(permissions))
}

// GetTokenAlias returns a token alias
func (k Keeper) GetTokenAlias(ctx sdk.Context, role types.Role) *types.Permissions {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), KeyPrefixPermissionsRegistry)
	bz := prefixStore.Get(types.RoleToKey(role))
	if bz == nil {
		return nil
	}

	perm := new(types.Permissions)
	k.cdc.MustUnmarshalBinaryBare(bz, perm)

	return perm
}
