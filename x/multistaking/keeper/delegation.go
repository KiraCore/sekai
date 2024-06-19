package keeper

import (
	"fmt"
	"strings"

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

func (k Keeper) RemoveUndelegation(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixUndelegation, sdk.Uint64ToBigEndian(id)...)
	store.Delete(key)
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

func (k Keeper) IsPoolDelegator(ctx sdk.Context, poolId uint64, delegator sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	key := append(append(types.KeyPrefixPoolDelegator, sdk.Uint64ToBigEndian(poolId)...), delegator...)
	bz := store.Get(key)
	return bz != nil
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
	k.UnregisterNotEnoughStakeDelegator(ctx, pool)

	delegators := k.GetPoolDelegators(ctx, pool.Id)
	for _, shareToken := range pool.TotalShareTokens {
		nativeDenom := types.GetNativeDenom(pool.Id, shareToken.Denom)
		rate := k.tokenKeeper.GetTokenInfo(ctx, nativeDenom)
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
				sdk.NewCoin(reward.Denom, sdk.NewDecFromInt(reward.Amount).Mul(rate.StakeCap).RoundInt()),
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
		properties := k.govKeeper.GetNetworkProperties(ctx)
		if compoundInfo.LastExecBlock+properties.AutocompoundIntervalNumBlocks > uint64(ctx.BlockHeight()) {
			continue
		}
		autoCompoundRewards := sdk.Coins{}
		if compoundInfo.AllDenom {
			autoCompoundRewards = rewards
			k.RemoveDelegatorRewards(ctx, delegator)
		} else {
			for _, reward := range rewards {
				rate := k.tokenKeeper.GetTokenInfo(ctx, reward.Denom)
				if rate.StakeEnabled && reward.Amount.GTE(rate.StakeMin) && isWithinArray(reward.Denom, compoundInfo.CompoundDenoms) {
					autoCompoundRewards = autoCompoundRewards.Add(reward)
				}
			}
			if !autoCompoundRewards.IsZero() {
				k.SetDelegatorRewards(ctx, delegator, rewards.Sub(autoCompoundRewards...))
			}
		}
		if !autoCompoundRewards.IsZero() {
			err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, delegator, autoCompoundRewards)
			if err != nil {
				panic(err)
			}
			err = k.Delegate(ctx, &types.MsgDelegate{
				DelegatorAddress: delegator.String(),
				ValidatorAddress: pool.Validator,
				Amounts:          autoCompoundRewards,
			})
			if err != nil {
				panic(err)
			}
			compoundInfo.LastExecBlock = uint64(ctx.BlockHeight())
			k.SetCompoundInfo(ctx, compoundInfo)
		}
	}
}

func (k Keeper) GetMinDelegatorWithValue(ctx sdk.Context, pool types.StakingPool) (sdk.AccAddress, sdk.Int) {
	poolDelegators := k.GetPoolDelegators(ctx, pool.Id)
	minDelegationValue := sdk.ZeroInt()
	minDelegator := sdk.AccAddress{}
	valAddr, err := sdk.ValAddressFromBech32(pool.Validator)
	if err != nil {
		return minDelegator, minDelegationValue
	}
	for _, delegator := range poolDelegators {
		if delegator.String() == sdk.AccAddress(valAddr).String() {
			continue
		}
		delegationValue := k.GetPoolDelegationValue(ctx, pool, delegator)
		if minDelegationValue.IsZero() || minDelegationValue.GT(delegationValue) {
			minDelegationValue = delegationValue
			minDelegator = delegator
		}
	}
	return minDelegator, minDelegationValue
}

func (k Keeper) Delegate(ctx sdk.Context, msg *types.MsgDelegate) error {
	// check if validator is active
	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return err
	}

	validator, err := k.sk.GetValidator(ctx, valAddr)
	if err != nil {
		return err
	}

	if !validator.IsActive() {
		return types.ErrNotActiveValidator
	}

	pool, found := k.GetStakingPoolByValidator(ctx, msg.ValidatorAddress)
	if !found {
		return types.ErrStakingPoolNotFound
	}

	if pool.Slashed.IsPositive() {
		return types.ErrActionNotSupportedForSlashedPool
	}

	delegator, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return err
	}

	// if it is going to create a new delegator
	if !k.IsPoolDelegator(ctx, pool.Id, delegator) {
		properties := k.govKeeper.GetNetworkProperties(ctx)
		poolDelegators := k.GetPoolDelegators(ctx, pool.Id)
		if len(poolDelegators) >= int(properties.MaxDelegators) {
			minDelegator, minDelegationValue := k.GetMinDelegatorWithValue(ctx, pool)

			// if it exceeds 10x of min delegation remove previous pool delegator
			delegatorValue := k.GetPoolDelegationValue(ctx, pool, delegator)
			newDelegatorValue := delegatorValue.Add(k.GetCoinsValue(ctx, msg.Amounts))
			if minDelegationValue.IsPositive() &&
				newDelegatorValue.GTE(minDelegationValue.Mul(sdk.NewInt(int64(properties.MinDelegationPushout)))) {
				k.RemovePoolDelegator(ctx, pool.Id, minDelegator)
			} else {
				return types.ErrMaxDelegatorsReached
			}
		}
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, delegator, types.ModuleName, msg.Amounts)
	if err != nil {
		return err
	}

	for _, amount := range msg.Amounts {
		rate := k.tokenKeeper.GetTokenInfo(ctx, amount.Denom)
		if !rate.StakeEnabled {
			return types.ErrNotAllowedStakingToken
		}
		if amount.Amount.LT(rate.StakeMin) {
			return types.ErrDenomStakingMinTokensNotReached
		}
	}

	pool.TotalStakingTokens = sdk.Coins(pool.TotalStakingTokens).Add(msg.Amounts...)
	poolCoins := types.GetPoolCoins(pool, msg.Amounts)
	pool.TotalShareTokens = sdk.Coins(pool.TotalShareTokens).Add(poolCoins...)
	k.SetStakingPool(ctx, pool)

	err = k.tokenKeeper.MintCoins(ctx, minttypes.ModuleName, poolCoins)
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

func (k Keeper) Undelegate(ctx sdk.Context, msg *types.MsgUndelegate) error {
	pool, found := k.GetStakingPoolByValidator(ctx, msg.ValidatorAddress)
	if !found {
		return types.ErrStakingPoolNotFound
	}

	delegator, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return err
	}

	poolCoins := types.GetPoolCoins(pool, msg.Amounts)

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, delegator, types.ModuleName, poolCoins)
	if err != nil {
		return err
	}
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, poolCoins)
	if err != nil {
		return err
	}

	if !sdk.Coins(pool.TotalStakingTokens).IsAllGTE(msg.Amounts) {
		return types.ErrInsufficientTotalStakingTokens
	}

	pool.TotalStakingTokens = sdk.Coins(pool.TotalStakingTokens).Sub(msg.Amounts...)
	pool.TotalShareTokens = sdk.Coins(pool.TotalShareTokens).Sub(poolCoins...)
	k.SetStakingPool(ctx, pool)

	lastUndelegationId := k.GetLastUndelegationId(ctx) + 1
	k.SetLastUndelegationId(ctx, lastUndelegationId)
	properties := k.govKeeper.GetNetworkProperties(ctx)
	k.SetUndelegation(ctx, types.Undelegation{
		Id:         lastUndelegationId,
		Address:    msg.DelegatorAddress,
		ValAddress: msg.ValidatorAddress,
		Expiry:     uint64(ctx.BlockTime().Unix()) + properties.UnstakingPeriod,
		Amount:     msg.Amounts,
	})

	balances := k.bankKeeper.GetAllBalances(ctx, delegator)
	prefix := fmt.Sprintf("v%d_", pool.Id)
	if !strings.Contains(balances.String(), prefix) {
		k.RemovePoolDelegator(ctx, pool.Id, delegator)
	}
	return nil
}

func (k Keeper) GetPoolDelegationValue(ctx sdk.Context, pool types.StakingPool, delegator sdk.AccAddress) sdk.Int {
	delegationValue := sdk.ZeroInt()
	balances := k.bankKeeper.GetAllBalances(ctx, delegator)
	for _, stakingToken := range pool.TotalStakingTokens {
		rate := k.tokenKeeper.GetTokenInfo(ctx, stakingToken.Denom)
		if rate == nil {
			continue
		}
		shareToken := types.GetShareDenom(pool.Id, stakingToken.Denom)
		balance := balances.AmountOf(shareToken)
		delegationValue = delegationValue.Add(sdk.NewDecFromInt(balance).Mul(rate.FeeRate).RoundInt())
	}
	return delegationValue
}

func (k Keeper) GetCoinsValue(ctx sdk.Context, coins sdk.Coins) sdk.Int {
	delegationValue := sdk.ZeroInt()
	for _, coin := range coins {
		rate := k.tokenKeeper.GetTokenInfo(ctx, coin.Denom)
		if rate == nil {
			continue
		}
		delegationValue = delegationValue.Add(sdk.NewDecFromInt(coin.Amount).Mul(rate.FeeRate).RoundInt())
	}
	return delegationValue
}

func (k Keeper) RegisterDelegator(ctx sdk.Context, delegator sdk.AccAddress) {
	balances := k.bankKeeper.GetAllBalances(ctx, delegator)

	pools := k.GetAllStakingPools(ctx)
	for _, pool := range pools {
		if k.IsPoolDelegator(ctx, pool.Id, delegator) {
			continue
		}

		properties := k.govKeeper.GetNetworkProperties(ctx)
		poolDelegators := k.GetPoolDelegators(ctx, pool.Id)
		if len(poolDelegators) >= int(properties.MaxDelegators) {
			minDelegator, minDelegationValue := k.GetMinDelegatorWithValue(ctx, pool)

			// if it exceeds 10x of min delegation remove previous pool delegator
			delegatorValue := k.GetPoolDelegationValue(ctx, pool, delegator)
			if minDelegationValue.IsPositive() &&
				delegatorValue.GTE(minDelegationValue.Mul(sdk.NewInt(int64(properties.MinDelegationPushout)))) {
				k.RemovePoolDelegator(ctx, pool.Id, minDelegator)
			} else {
				continue
			}
		}

		for _, stakingToken := range pool.TotalStakingTokens {
			rate := k.tokenKeeper.GetTokenInfo(ctx, stakingToken.Denom)
			shareToken := types.GetShareDenom(pool.Id, stakingToken.Denom)
			balance := balances.AmountOf(shareToken)
			if balance.GTE(rate.StakeMin) {
				k.SetPoolDelegator(ctx, pool.Id, delegator)
				break
			}
		}
	}
}

func (k Keeper) UnregisterNotEnoughStakeDelegator(ctx sdk.Context, pool types.StakingPool) {
	delegators := k.GetPoolDelegators(ctx, pool.Id)

	// distribute rewards allocated for the staked denom to delegators
	for _, delegator := range delegators {
		toBeRemoved := true
		balances := k.bankKeeper.GetAllBalances(ctx, delegator)
		for _, shareToken := range pool.TotalShareTokens {
			nativeDenom := types.GetNativeDenom(pool.Id, shareToken.Denom)
			rate := k.tokenKeeper.GetTokenInfo(ctx, nativeDenom)
			if rate == nil {
				continue
			}

			// total share token amount validation
			if shareToken.Amount.IsZero() {
				continue
			}

			balance := balances.AmountOf(shareToken.Denom)
			if balance.GTE(rate.StakeMin) {
				toBeRemoved = false
				break
			}
		}
		if toBeRemoved {
			k.RemovePoolDelegator(ctx, pool.Id, delegator)
		}
	}
}

func (k Keeper) ClaimRewards(ctx sdk.Context, delegator sdk.AccAddress) sdk.Coins {
	rewards := k.GetDelegatorRewards(ctx, delegator)
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, delegator, rewards)
	if err != nil {
		panic(err)
	}

	k.RemoveDelegatorRewards(ctx, delegator)
	return rewards
}

func (k Keeper) ClaimRewardsFromModule(ctx sdk.Context, moduleName string) sdk.Coins {
	delegator := authtypes.NewModuleAddress(moduleName)
	rewards := k.GetDelegatorRewards(ctx, delegator)
	if rewards.IsAllPositive() {
		err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, authtypes.FeeCollectorName, moduleName, rewards)
		if err != nil {
			panic(err)
		}

		k.RemoveDelegatorRewards(ctx, delegator)
	}
	return rewards
}
