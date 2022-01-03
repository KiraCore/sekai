package keeper

import (
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SavePoorNetworkMessages store poor network messages by gov or by genesis
func (k Keeper) SavePoorNetworkMessages(ctx sdk.Context, allows *types.AllowedMessages) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(allows)
	store.Set(PoorNetworkMessagesPrefix, bz)
}

// GetPoorNetworkMessages returns poor network messages stored inside keeper
func (k Keeper) GetPoorNetworkMessages(ctx sdk.Context) *types.AllowedMessages {
	var am types.AllowedMessages
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(PoorNetworkMessagesPrefix)
	if bz == nil {
		return &am
	}

	k.cdc.MustUnmarshal(bz, &am)

	return &am
}
