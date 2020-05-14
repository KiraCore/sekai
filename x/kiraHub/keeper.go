package kiraHub

import (
	"github.com/KiraCore/cosmos-sdk/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
)

// Keeper of the createorderboook store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
}






