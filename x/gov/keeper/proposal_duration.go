package keeper

import (
	"fmt"

	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetProposalDuration(ctx sdk.Context, proposalType string, duration uint64) error {
	// duration should be longer than minimum proposal duration
	properties := k.GetNetworkProperties(ctx)
	if duration < properties.MinimumProposalEndTime {
		return fmt.Errorf("duration should be longer than minimum proposal duration")
	}

	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.KeyPrefixProposalDuration)
	prefixStore.Set([]byte(proposalType), sdk.Uint64ToBigEndian(duration))
	return nil
}

func (k Keeper) GetProposalDuration(ctx sdk.Context, proposalType string) uint64 {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.KeyPrefixProposalDuration)
	bz := prefixStore.Get([]byte(proposalType))
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) GetAllProposalDurations(ctx sdk.Context) map[string]uint64 {
	proposalDurations := make(map[string]uint64)
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.KeyPrefixProposalDuration)
	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposalDurations[string(iterator.Key())] = sdk.BigEndianToUint64(iterator.Value())
	}
	return proposalDurations
}
