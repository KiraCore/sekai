package keeper

import (
	"github.com/KiraCore/sekai/x/recovery/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var MinHoldAmount = sdk.NewInt(1000_000)

func (k Keeper) SetRRTokenHolder(ctx sdk.Context, rrToken string, delegator sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := append(append(types.KeyPrefixRRTokenHolder, rrToken...), delegator...)
	store.Set(key, delegator)
}

func (k Keeper) RemoveRRTokenHolder(ctx sdk.Context, rrToken string, delegator sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := append(append(types.KeyPrefixRRTokenHolder, rrToken...), delegator...)
	store.Delete(key)
}

func (k Keeper) IsRRTokenHolder(ctx sdk.Context, rrToken string, delegator sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	key := append(append(types.KeyPrefixRRTokenHolder, rrToken...), delegator...)
	bz := store.Get(key)
	return bz != nil
}

func (k Keeper) GetRRTokenHolders(ctx sdk.Context, rrToken string) []sdk.AccAddress {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.KeyPrefixRRTokenHolder, rrToken...))

	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	delegators := []sdk.AccAddress{}
	for ; iterator.Valid(); iterator.Next() {
		delegators = append(delegators, sdk.AccAddress(iterator.Value()))
	}
	return delegators
}

func (k Keeper) SetRRTokenHolderRewards(ctx sdk.Context, holder sdk.AccAddress, rewards sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixRewards, holder...)
	store.Set(key, []byte(rewards.String()))
}

func (k Keeper) IncreaseRRTokenHolderRewards(ctx sdk.Context, holder sdk.AccAddress, amounts sdk.Coins) {
	rewards := k.GetRRTokenHolderRewards(ctx, holder)
	rewards = rewards.Add(amounts...)

	k.SetRRTokenHolderRewards(ctx, holder, rewards)
}

func (k Keeper) RemoveRRTokenHolderRewards(ctx sdk.Context, holder sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixRewards, holder...)
	store.Delete(key)
}

func (k Keeper) GetRRTokenHolderRewards(ctx sdk.Context, holder sdk.AccAddress) sdk.Coins {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixRewards, holder...)
	bz := store.Get(key)
	if bz == nil {
		return sdk.Coins{}
	}

	coinStr := string(bz)
	coins, err := sdk.ParseCoinsNormalized(coinStr)
	if err != nil {
		panic(err)
	}

	return coins
}

func (k Keeper) GetAllRRHolderRewards(ctx sdk.Context) []types.Rewards {
	rewards := []types.Rewards{}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixRewards)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		rewardsCoins, err := sdk.ParseCoinsNormalized(string(iterator.Value()))
		if err != nil {
			panic(err)
		}
		rewards = append(rewards, types.Rewards{
			Holder:  sdk.AccAddress(iterator.Key()).String(),
			Rewards: rewardsCoins,
		})
	}
	return rewards
}

func (k Keeper) RegisterRRTokenHolder(ctx sdk.Context, delegator sdk.AccAddress) {
	balances := k.bk.GetAllBalances(ctx, delegator)

	recoveryTokens := k.GetAllRecoveryTokens(ctx)
	for _, recoveryToken := range recoveryTokens {
		if k.IsRRTokenHolder(ctx, recoveryToken.Token, delegator) {
			continue
		}

		balance := balances.AmountOf(recoveryToken.Token)
		if balance.GTE(MinHoldAmount) {
			k.SetRRTokenHolder(ctx, recoveryToken.Token, delegator)
			break
		}
	}
}

func (k Keeper) UnregisterNotEnoughAmountHolder(ctx sdk.Context, rrToken string) {
	holders := k.GetRRTokenHolders(ctx, rrToken)

	for _, holder := range holders {
		toBeRemoved := true
		balances := k.bk.GetAllBalances(ctx, holder)
		if balances.AmountOf(rrToken).GTE(MinHoldAmount) {
			toBeRemoved = false
		}
		if toBeRemoved {
			k.RemoveRRTokenHolder(ctx, rrToken, holder)
		}
	}
}

func (k Keeper) ClaimRewards(ctx sdk.Context, delegator sdk.AccAddress) sdk.Coins {
	rewards := k.GetRRTokenHolderRewards(ctx, delegator)
	err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delegator, rewards)
	if err != nil {
		panic(err)
	}

	k.RemoveRRTokenHolderRewards(ctx, delegator)
	return rewards
}
