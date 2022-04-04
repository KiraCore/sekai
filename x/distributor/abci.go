package distributor

import (
	"time"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/KiraCore/sekai/x/distributor/keeper"
	"github.com/KiraCore/sekai/x/distributor/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker sets the proposer for determining distributor during endblock
// and distribute rewards for the previous block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	// determine the total power signing the block
	var previousTotalPower, sumPreviousPrecommitPower int64
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		previousTotalPower += voteInfo.Validator.Power
		if voteInfo.SignedLastBlock {
			sumPreviousPrecommitPower += voteInfo.Validator.Power
		}
	}

	// TODO this is Tendermint-dependent
	// ref https://github.com/cosmos/cosmos-sdk/issues/3095
	if ctx.BlockHeight() > 1 {
		previousProposer := k.GetPreviousProposerConsAddr(ctx)
		k.AllocateTokens(ctx, sumPreviousPrecommitPower, previousTotalPower, previousProposer, req.LastCommitInfo.GetVotes())
	}

	for _, bondedVote := range req.LastCommitInfo.GetVotes() {
		k.SetValidatorVote(ctx, bondedVote.Validator.Address, ctx.BlockHeight())
	}

	// remove votes older than snap period
	snapPeriod := k.GetSnapPeriod(ctx)
	allVotes := k.GetAllValidatorVotes(ctx)
	for _, vote := range allVotes {
		if vote.Height+snapPeriod < ctx.BlockHeight() {
			consAddr, err := sdk.ConsAddressFromBech32(vote.ConsAddr)
			if err != nil {
				panic(err)
			}
			k.DeleteValidatorVote(ctx, consAddr, vote.Height)
		}
	}

	// record the proposer for when we payout on the next block
	consAddr := sdk.ConsAddress(req.Header.ProposerAddress)
	k.SetPreviousProposerConsAddr(ctx, consAddr)
}
