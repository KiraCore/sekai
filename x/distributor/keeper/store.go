package keeper

import (
	"bytes"

	"github.com/KiraCore/sekai/x/distributor/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetFeesTreasury(ctx sdk.Context, coins sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.FeesTreasuryKey, []byte(coins.String()))
}

func (k Keeper) GetFeesTreasury(ctx sdk.Context) sdk.Coins {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.FeesTreasuryKey)
	if bz == nil {
		return sdk.Coins{}
	}
	coins, err := sdk.ParseCoinsNormalized(string(bz))
	if err != nil {
		panic(err)
	}
	return coins
}

func (k Keeper) SetSnapPeriod(ctx sdk.Context, period int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.SnapPeriodKey, sdk.Uint64ToBigEndian(uint64(period)))
}

func (k Keeper) GetSnapPeriod(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.SnapPeriodKey)
	if bz == nil {
		return 1
	}
	return int64(sdk.BigEndianToUint64(bz))
}

func (k Keeper) SetValidatorVote(ctx sdk.Context, consAddr sdk.ConsAddress, height int64) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixKeyValidatorVote)
	bz := sdk.Uint64ToBigEndian(uint64(height))
	prefixStore.Set(append(consAddr, bz...), bz)
}

func (k Keeper) DeleteValidatorVote(ctx sdk.Context, consAddr sdk.ConsAddress, height int64) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixKeyValidatorVote)
	bz := sdk.Uint64ToBigEndian(uint64(height))
	prefixStore.Delete(append(consAddr, bz...))
}

func (k Keeper) GetValidatorVotes(ctx sdk.Context, consAddr sdk.ConsAddress) []int64 {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.PrefixKeyValidatorVote, []byte(consAddr)...))

	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	voteHeights := []int64{}
	for ; iterator.Valid(); iterator.Next() {
		voteHeights = append(voteHeights, int64(sdk.BigEndianToUint64(iterator.Value())))
	}
	return voteHeights
}

func (k Keeper) GetAllValidatorVotes(ctx sdk.Context) []types.ValidatorVote {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrefixKeyValidatorVote)

	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	valVotes := []types.ValidatorVote{}
	for ; iterator.Valid(); iterator.Next() {
		consAddr := bytes.TrimSuffix(iterator.Key(), iterator.Value())
		valVotes = append(valVotes, types.ValidatorVote{
			ConsAddr: sdk.ConsAddress(consAddr).String(),
			Height:   int64(sdk.BigEndianToUint64(iterator.Value())),
		})
	}
	return valVotes
}
