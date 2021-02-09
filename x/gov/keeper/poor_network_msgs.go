package keeper

import (
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SavePoorNetworkMsgs store poor network messages by gov or by genesis
func (k Keeper) SavePoorNetworkMsgs(ctx sdk.Context, allows *types.AllowedMessages) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(allows)
	store.Set(PoorNetworkMsgsPrefix, bz)
}

// GetPoorNetworkMsgs returns poor network messages stored inside keeper
func (k Keeper) GetPoorNetworkMsgs(ctx sdk.Context) (*types.AllowedMessages, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(PoorNetworkMsgsPrefix)
	if bz == nil {
		return &types.AllowedMessages{}, false
	}

	var am types.AllowedMessages
	k.cdc.MustUnmarshalBinaryBare(bz, &am)

	return &am, true
}
