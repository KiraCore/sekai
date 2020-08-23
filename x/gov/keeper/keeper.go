package keeper

import (
	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	cdc      codec.BinaryMarshaler
	storeKey sdk.StoreKey
}

func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryMarshaler) Keeper {
	return Keeper{cdc: cdc, storeKey: storeKey}
}

func (k Keeper) SetPermissionsForRole(ctx sdk.Context, role types.Role, permissions types.Permissions) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPermissionsRegistry)

	prefixStore.Set(types.RoleToKey(role), k.cdc.MustMarshalBinaryBare(&permissions))
}

func (k Keeper) GetPermissionsForRole(ctx sdk.Context, councilor types.Role) types.Permissions {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPermissionsRegistry)
	bz := prefixStore.Get(types.RoleToKey(councilor))

	var perm types.Permissions
	k.cdc.MustUnmarshalBinaryBare(bz, &perm)

	return perm
}
