package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper is for managing upgrade module
type Keeper struct {
	cdc      codec.BinaryMarshaler
	storeKey sdk.StoreKey
}
