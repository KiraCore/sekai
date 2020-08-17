package keeper

import (
	"github.com/KiraCore/cosmos-sdk/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
)

// Keeper represents the keeper that maintains the Validator Registry.
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

// NewKeeper returns new keeper.
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{storeKey: storeKey, cdc: cdc}
}
