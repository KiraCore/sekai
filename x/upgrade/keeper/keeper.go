package keeper

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper is for managing upgrade module
type Keeper struct {
	cdc      codec.BinaryMarshaler
	storeKey sdk.StoreKey
}

// NewKeeper constructs an upgrade Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryMarshaler) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

func (am Keeper) ExportGenesisRollbackState(clientCtx sdk.Context, marshaler codec.JSONMarshaler) json.RawMessage {
	return nil
}
