package keeper

import (
	"fmt"

	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// AllocateTokens handles distribution of the collected fees
// bondedVotes is a list of (validator address, validator voted on last block flag) for all
// validators in the bonded set.
func (k Keeper) AllocateTokens(
	ctx sdk.Context, sumPreviousPrecommitPower, totalPreviousPower int64,
	previousProposer sdk.ConsAddress, bondedVotes []abci.VoteInfo,
) {

	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)
	feeCollector := k.ak.GetModuleAccount(ctx, authtypes.FeeCollectorName)
	feesAccBalance := k.bk.GetAllBalances(ctx, feeCollector.GetAddress())
	feesTreasury := k.GetFeesTreasury(ctx)

	// mint inflated tokens
	totalSupply := k.bk.GetSupply(ctx, k.BondDenom(ctx))
	properties := k.gk.GetNetworkProperties(ctx)
	inflationRewards := totalSupply.Amount.Mul(sdk.NewInt(int64(properties.InflationRate))).Quo(sdk.NewInt(int64(properties.InflationPeriod)))
	inflationCoin := sdk.NewCoin(totalSupply.Denom, inflationRewards)

	if inflationRewards.IsPositive() {
		err := k.bk.MintCoins(ctx, minttypes.ModuleName, sdk.Coins{inflationCoin})
		if err != nil {
			panic(err)
		}
		err = k.bk.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, authtypes.FeeCollectorName, sdk.Coins{inflationCoin})
		if err != nil {
			panic(err)
		}
	}

	// combine fees and inflated tokens for rewards allocation
	feesCollected := sdk.Coins{}
	if feesAccBalance.IsAllGTE(feesTreasury) {
		feesCollected = feesAccBalance.Sub(feesTreasury)
	}

	validatorsFeeShare := k.gk.GetNetworkProperties(ctx).ValidatorsFeeShare
	if validatorsFeeShare.GT(sdk.OneDec()) {
		validatorsFeeShare = sdk.OneDec()
	}

	// pay previous proposer
	proposerValidator, err := k.sk.GetValidatorByConsAddr(ctx, previousProposer)
	if err == nil {
		// calculate reward based on historical bonded votes of the validator
		votes := k.GetValidatorVotes(ctx, previousProposer)
		power := int64(len(votes))
		snapPeriod := k.GetSnapPeriod(ctx)
		validatorRewards := sdk.Coins{}
		poolRewards := sdk.Coins{}

		// add fee rewards for validator
		for _, r := range feesCollected {
			cutAmount := r.Amount.Mul(sdk.NewInt(power)).Quo(sdk.NewInt(snapPeriod))
			valReward := cutAmount.ToDec().Mul(validatorsFeeShare).RoundInt()
			if valReward.IsPositive() {
				validatorRewards = validatorRewards.Add(sdk.NewCoin(r.Denom, valReward))
			}
			poolReward := cutAmount.Sub(valReward)
			if poolReward.IsPositive() {
				poolRewards = poolRewards.Add(sdk.NewCoin(r.Denom, poolReward))
			}
		}

		// add block inflation rewards for validator
		cutInflationRewards := inflationRewards.Mul(sdk.NewInt(power)).Quo(sdk.NewInt(snapPeriod))
		inflationCommissionReward := cutInflationRewards.ToDec().Mul(proposerValidator.Commission).RoundInt()
		validatorRewards = validatorRewards.Add(sdk.NewCoin(totalSupply.Denom, inflationCommissionReward))
		inflationPoolReward := cutInflationRewards.Sub(inflationCommissionReward)
		poolRewards = poolRewards.Add(sdk.NewCoin(totalSupply.Denom, inflationPoolReward))

		if !validatorRewards.Empty() {
			k.AllocateTokensToValidator(ctx, proposerValidator, validatorRewards)
		}

		if !poolRewards.Empty() {
			pool, found := k.mk.GetStakingPoolByValidator(ctx, proposerValidator.ValKey.String())
			if found {
				k.mk.IncreasePoolRewards(ctx, pool, poolRewards)
			}
		}
	} else {
		// previous proposer can be unknown if say, the unbonding period is 1 block, so
		// e.g. a validator undelegates at block X, it's removed entirely by
		// block X+1's endblock, then X+2 we need to refer to the previous
		// proposer for X+1, but we've forgotten about them.
		fmt.Println(fmt.Sprintf(
			"WARNING: Attempt to allocate proposer rewards to unknown proposer %s. "+
				"This should happen only if the proposer unbonded completely within a single block, "+
				"which generally should not happen except in exceptional circumstances (or fuzz testing). "+
				"We recommend you investigate immediately.",
			previousProposer.String()))
	}

	// give rest of the tokens to community pool
	remainings := k.bk.GetAllBalances(ctx, feeCollector.GetAddress())
	k.SetFeesTreasury(ctx, remainings)
}

// AllocateTokensToValidator allocate tokens to a particular validator, splitting according to commission
func (k Keeper) AllocateTokensToValidator(ctx sdk.Context, val stakingtypes.Validator, tokens sdk.Coins) {
	// send coins from fee pool to validator account
	err := k.bk.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, sdk.AccAddress(val.GetValKey()), tokens)
	if err != nil {
		panic(err)
	}
}

// GetPreviousProposerConsAddr returns the proposer consensus address for the
// current block.
func (k Keeper) GetPreviousProposerConsAddr(ctx sdk.Context) sdk.ConsAddress {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ProposerKey)
	if bz == nil {
		panic("previous proposer not set")
	}

	return sdk.ConsAddress(bz)
}

// set the proposer public key for this block
func (k Keeper) SetPreviousProposerConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ProposerKey, consAddr)
}

func (k Keeper) GetValidatorPerformance(ctx sdk.Context, valAddr sdk.ValAddress) (int64, error) {
	validator, err := k.sk.GetValidator(ctx, valAddr)
	if err != nil {
		return 0, err
	}
	return int64(len(k.GetValidatorVotes(ctx, validator.GetConsAddr()))), nil
}
