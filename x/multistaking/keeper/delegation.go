package keeper

import (
	"github.com/KiraCore/sekai/x/multistaking/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (k Keeper) GetLastUndelegationId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyLastUndelegationId)
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLastUndelegationId(ctx sdk.Context, id uint64) {
	idBz := sdk.Uint64ToBigEndian(id)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyLastUndelegationId, idBz)
}

func (k Keeper) GetUndelegationById(ctx sdk.Context, id uint64) (undelegation types.Undelegation, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.KeyPrefixUndelegation), sdk.Uint64ToBigEndian(id)...)
	bz := store.Get(key)
	if bz == nil {
		return undelegation, false
	}
	k.cdc.MustUnmarshal(bz, &undelegation)
	return undelegation, true
}

func (k Keeper) GetAllUndelegations(ctx sdk.Context) []types.Undelegation {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixUndelegation)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	undelegations := []types.Undelegation{}
	for ; iterator.Valid(); iterator.Next() {
		undelegation := types.Undelegation{}
		k.cdc.MustUnmarshal(iterator.Value(), &undelegation)
		undelegations = append(undelegations, undelegation)
	}
	return undelegations
}

func (k Keeper) SetUndelegation(ctx sdk.Context, undelegation types.Undelegation) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixUndelegation, sdk.Uint64ToBigEndian(undelegation.Id)...)
	store.Set(key, k.cdc.MustMarshal(&undelegation))
}

func (k Keeper) SetPoolDelegator(ctx sdk.Context, poolId uint64, delegator sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := append(append(types.KeyPrefixPoolDelegator, sdk.Uint64ToBigEndian(poolId)...), delegator...)
	store.Set(key, delegator)
}

func (k Keeper) RemovePoolDelegator(ctx sdk.Context, poolId uint64, delegator sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := append(append(types.KeyPrefixPoolDelegator, sdk.Uint64ToBigEndian(poolId)...), delegator...)
	store.Delete(key)
}

func (k Keeper) GetPoolDelegators(ctx sdk.Context, poolId uint64) []sdk.AccAddress {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.KeyPrefixPoolDelegator, sdk.Uint64ToBigEndian(poolId)...))

	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	delegators := []sdk.AccAddress{}
	for ; iterator.Valid(); iterator.Next() {
		delegators = append(delegators, sdk.AccAddress(iterator.Value()))
	}
	return delegators
}

func (k Keeper) SetDelegatorRewards(ctx sdk.Context, delegator sdk.AccAddress, rewards sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixRewards, delegator...)
	store.Set(key, []byte(rewards.String()))
}

func (k Keeper) IncreaseDelegatorRewards(ctx sdk.Context, delegator sdk.AccAddress, amounts sdk.Coins) {
	rewards := k.GetDelegatorRewards(ctx, delegator)
	rewards = rewards.Add(amounts...)

	k.SetDelegatorRewards(ctx, delegator, rewards)
}

func (k Keeper) RemoveDelegatorRewards(ctx sdk.Context, delegator sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixRewards, delegator...)
	store.Delete(key)
}

func (k Keeper) GetDelegatorRewards(ctx sdk.Context, delegator sdk.AccAddress) sdk.Coins {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixRewards, delegator...)
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

func (k Keeper) GetAllDelegatorRewards(ctx sdk.Context) []types.Rewards {
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
			Delegator: sdk.AccAddress(iterator.Key()).String(),
			Rewards:   rewardsCoins,
		})
	}
	return rewards
}

func isWithinArray(s string, arr []string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}

func (k Keeper) IncreasePoolRewards(ctx sdk.Context, pool types.StakingPool, rewards sdk.Coins) {
	delegators := k.GetPoolDelegators(ctx, pool.Id)
	for _, shareToken := range pool.TotalShareTokens {
		nativeDenom := getNativeDenom(pool.Id, shareToken.Denom)
		rate := k.tokenKeeper.GetTokenRate(ctx, nativeDenom)
		if rate == nil {
			continue
		}

		if rate.StakeCap.IsZero() {
			continue
		}

		// total share token amount validation
		if shareToken.Amount.IsZero() {
			continue
		}

		// rewards allocated for the denom
		denomAllocation := sdk.Coins{}
		for _, reward := range rewards {
			denomAllocation = denomAllocation.Add(
				sdk.NewCoin(reward.Denom, reward.Amount.ToDec().Mul(rate.StakeCap).RoundInt()),
			)
		}

		// distribute rewards allocated for the staked denom to delegators
		for _, delegator := range delegators {
			balances := k.bankKeeper.GetAllBalances(ctx, delegator)
			balance := balances.AmountOf(shareToken.Denom)

			delegatorRewards := sdk.Coins{}
			for _, reward := range denomAllocation {
				delegatorRewards = delegatorRewards.Add(sdk.NewCoin(reward.Denom, reward.Amount.Mul(balance).Quo(shareToken.Amount)))
			}

			k.IncreaseDelegatorRewards(ctx, delegator, delegatorRewards)
		}
	}

	for _, delegator := range delegators {
		// autocompound rewards
		rewards := k.GetDelegatorRewards(ctx, delegator)
		compoundInfo := k.GetCompoundInfoByAddress(ctx, delegator.String())
		autoCompoundRewards := sdk.Coins{}
		if compoundInfo.AllDenom {
			autoCompoundRewards = rewards
			k.RemoveDelegatorRewards(ctx, delegator)
		} else {
			for _, reward := range rewards {
				rate := k.tokenKeeper.GetTokenRate(ctx, reward.Denom)
				if rate.StakeToken && reward.Amount.GTE(rate.StakeMin) && isWithinArray(reward.Denom, compoundInfo.CompoundDenoms) {
					autoCompoundRewards = autoCompoundRewards.Add(reward)
				}
			}
			if !autoCompoundRewards.IsZero() {
				err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, delegator, autoCompoundRewards)
				if err != nil {
					panic(err)
				}
				k.SetDelegatorRewards(ctx, delegator, rewards.Sub(autoCompoundRewards))
			}
		}
		if !autoCompoundRewards.IsZero() {
			msgServer := NewMsgServerImpl(k, k.bankKeeper, k.govKeeper, k.sk)
			_, err := msgServer.Delegate(sdk.WrapSDKContext(ctx), &types.MsgDelegate{
				DelegatorAddress: delegator.String(),
				ValidatorAddress: pool.Validator,
				Amounts:          autoCompoundRewards,
			})
			if err != nil {
				panic(err)
			}
		}
	}
}

func (k Keeper) Delegate(ctx sdk.Context, msg *types.MsgDelegate) error {
	pool, found := k.GetStakingPoolByValidator(ctx, msg.ValidatorAddress)
	if !found {
		return types.ErrStakingPoolNotFound
	}

	delegator, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, delegator, types.ModuleName, msg.Amounts)
	if err != nil {
		return err
	}

	for _, amount := range msg.Amounts {
		rate := k.tokenKeeper.GetTokenRate(ctx, amount.Denom)
		if !rate.StakeToken {
			return types.ErrNotAllowedStakingToken
		}
		if amount.Amount.LT(rate.StakeMin) {
			return types.ErrDenomStakingMinTokensNotReached
		}
	}

	pool.TotalStakingTokens = sdk.Coins(pool.TotalStakingTokens).Add(msg.Amounts...)
	poolCoins := getPoolCoins(pool.Id, msg.Amounts)
	pool.TotalShareTokens = sdk.Coins(pool.TotalShareTokens).Add(poolCoins...)
	k.SetStakingPool(ctx, pool)

	// TODO: should check the ratio between poolCoins and coins
	err = k.bankKeeper.MintCoins(ctx, minttypes.ModuleName, poolCoins)
	if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delegator, poolCoins)
	if err != nil {
		return err
	}

	k.SetPoolDelegator(ctx, pool.Id, delegator)
	return nil
}
