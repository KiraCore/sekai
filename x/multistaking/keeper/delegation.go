package keeper

import (
	"github.com/KiraCore/sekai/x/multistaking/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	key := append([]byte(types.KeyPrefixStakingPool), sdk.Uint64ToBigEndian(id)...)
	bz := store.Get(key)
	if bz == nil {
		return undelegation, false
	}
	k.cdc.MustUnmarshal(bz, &undelegation)
	return undelegation, true
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

func (k Keeper) IncreaseDelegatorRewards(ctx sdk.Context, delegator sdk.AccAddress, amounts sdk.Coins) {
	rewards := k.GetDelegatorRewards(ctx, delegator)
	rewards = rewards.Add(amounts...)

	store := ctx.KVStore(k.storeKey)
	key := append(types.KeyPrefixRewards, delegator...)
	store.Set(key, []byte(rewards.String()))
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

func (k Keeper) IncreasePoolRewards(ctx sdk.Context, poolId uint64, rewards sdk.Coins) {
	// TODO: loop pool delegators
	delegators := k.GetPoolDelegators(ctx, poolId)
	for _, delegator := range delegators {
		_ = delegator

		// TODO: get pool deposit tokens
		// TODO: split based on pool token's total balance and reward rate
	}

	// 	**Pseudo algorithm**

	// 1. Calculate cumulative delegator rewards `Dr` in all tokens that were deposited to the pool after the validator `commission` and his `validators_fee_share`
	// 2. For each delegator in `delegators` array or validator (if he self delegates) calculate their pool share
	//     1. For each stakeable token (as defined by Token Rates Registry) in the `balances`
	//         1. Divide deleagatoer `shares` (or `claims` if the token was compounded but the derivative itself is not stakeable) by a total of all bonded shares of a particular token `bonds[<denom>].shares` to get delegator token share `Dts` in the pool
	//         2. Multiply `Dts` by the `stake_cap` of that particular token to get delegator rewards token share `Rts = Dts  * stake_cap`
	//     2. Sum all `Rts` of a delegator to get a final delagator rewards share `Drs`
	// 3. For each reward token update multiply its amount by `Drs` and add the value to corresponding `claims`. If the record for a specific token does NOT exist, then create it.
	// 4. If `claims` amount of any stakeable token is non `0` and user set `compound` to `true` for that token then issue staking derivative to the staking pool, set the value of original claim to `0` and set the `claim` value of the derivative record to the same amount. Note that compounding should occur only if rewards equal to no less than `stake_min` and depositing newly created derivatives should happen only once - cumulatively for each token.

	// *NOTE: User holding derivatives on this account (`shares`) or owning them as `claims` as result of compounding should be treated exactly the same way while calculating the rewards that the delegator is eligible to receive - otherwise the compounding would not work.*

	// *NOTE 2: The delegators and validator section in the staking registrar are separated so that we can easily distinguish self bond as well as dispose of dust when the delegator balance drops below predefined minimum (without deleting validator address). We can optionally keep validator stake data in delegators array along all other records with a notion that we need to add all other extra fields and store them somewhere else - address, commission.*

	// 	**Pseudo algorithm**

	// 1. Calculate time period `P` that passed since the last time the block rewards were deposited to the staking pool `P = (current-block-time - timestamp)`,
	//     1. **If `P` is less or equal to `0`, or greater than `unstaking_period`, then do NOT issue any new KEX tokens to the staking pool and set the `P` to `current_block_time`. (To prevent all possible edge cases where delegators would be able to migrate the stake before validator proposes another block / prevent any possible timestamp manipulation exploits).** Notice here that this rule is essential because we do not store information when exactly the tokens were staked to the pool and in what particular amount.
	// 2. For each token in the `bonds` array calculate total amount of all bonded coins that exists to get `total_bonded_supply` by iterating over all existing staking pools and summing up `amount`.
	//     1. *NOTE: It might be good to store the `total_bonded_supply` for each token in some registry on-chain (E.g. Token Rates Registry) so that access to those variables is easy and can be queried without iterating over list of all staking pools (as there can be hundreds or thousands of those). In such case we should also be updating that registrar every time new tokens are bonded or unbonded to save a lot of processing & I/O operations time!*
	// 3. Calculate block rewards that should be credited to the staking pool for the time period `P` according to the inflation network properties
	//     1. Get inflation per second `Ips = (total_ukex_supply * inflation_rate) / inflation_period`. **Notice here that we will be inflating entire KEX supply, but only granting rewards to those who stake tokens, this way delegations will be encouraged by being more profitable the smaller total number of tokens is staked.**
	//     2. Multiply time period `P` by inflation per second `Ips` to get total block rewards `Tbr` for the time span.
	//     3. Calculate what fraction of the `Tbr` the staking pool is eligible for, by iterating over each token in `bonds`
	//         1. Divide `amount` by a `total_bonded_supply` of the specific token to get global share (`Gs = amount / total_bonded_supply`), that is percentage of all stake bonded to all validators
	//         2. Multiply global share `Gs` by a `stake_cap` of a specific token to get a rewards share (`Rs = Gs * stake_cap`), that is a percentage of all block rewards that the bonded tokens are eligible for
	//         3. Multiply rewards share `Rs` by total block rewards `Tbr` to get pool block rewards for a token in terms of `ukex` amount `Sbr = Rs * Tbr`
	//     4. Sum all `Sbr`’s together to get the pool block rewards `Pbr` amount of `ukex` that should be printed and deposited to the staking pool.
	//     5. Make absolutely sure that `Pbr` is smaller or equal to `Tbr` (Otherwise something must have been wrong with our calculations. The case where `Pbr` equails `Tbr` would occur only if all staking tokens in existence were deposited and staked to only 1 validator)
	// 4. Issue and deposit amount of `ukex` equal to `Pbr`, into our staking pool
	//     1. Grant to validator amount `Br * validator.commission`                           [ see Staking Rewards Distribution section for more details ]
	//     2. Split between delegators amount `Br * (1 - validator.commission)`     [ see Staking Rewards Distribution section for more details ]

}
