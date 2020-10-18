package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// store prefixes
var (
	PrefixKeyTokenAlias = []byte("token_alias_registry")
	PrefixKeyDenomToken = []byte("denom_token_registry")
	PrefixKeyTokenRate  = []byte("token_rate_registry")
)

// Keeper is for managing token module
type Keeper struct {
	cdc      codec.BinaryMarshaler
	storeKey sdk.StoreKey
}

// NewKeeper returns instance of a keeper
func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryMarshaler) Keeper {
	return Keeper{cdc: cdc, storeKey: storeKey}
}
