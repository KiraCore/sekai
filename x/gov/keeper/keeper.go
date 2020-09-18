package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	KeyPrefixPermissionsRegistry       = []byte("perm_registry")
	KeyPrefixCouncilorIdentityRegistry = []byte("council_registry")
	KeyPrefixActors                    = []byte("network_actors")
)

type Keeper struct {
	cdc      codec.BinaryMarshaler
	storeKey sdk.StoreKey
}

func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryMarshaler) Keeper {
	return Keeper{cdc: cdc, storeKey: storeKey}
}
