package keeper

import (
	"fmt"

	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
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
	feesCollected := k.bk.GetAllBalances(ctx, feeCollector.GetAddress())

	// transfer collected fees to the distribution module account
	err := k.bk.SendCoinsFromModuleToModule(ctx, authtypes.FeeCollectorName, types.ModuleName, feesCollected)
	if err != nil {
		panic(err)
	}

	// pay previous proposer
	proposerValidator, err := k.sk.GetValidatorByConsAddr(ctx, previousProposer)

	if err == nil {
		k.AllocateTokensToValidator(ctx, proposerValidator, feesCollected)
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

	// TODO: give rest of the tokens not distributed to community pool
	// allocate community funding
	// feePool := k.GetFeePool(ctx)
	// feePool.CommunityPool = feePool.CommunityPool.Add(remaining...)
	// k.SetFeePool(ctx, feePool)
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
