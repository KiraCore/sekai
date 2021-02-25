package slashing

import (
	"fmt"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/KiraCore/sekai/x/slashing/keeper"
	"github.com/KiraCore/sekai/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	fmt.Println("abci.BeginBlocker", req.LastCommitInfo.GetVotes())

	// Iterate over all the validators which *should* have signed this block
	// store whether or not they have actually signed it and slash/unbond any
	// which have missed too many blocks in a row (downtime slashing)
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		fmt.Println("abci.BeginBlockerxx1", sdk.ConsAddress(voteInfo.Validator.Address).String())
		fmt.Println("abci.BeginBlockerxx2", sdk.ValAddress(voteInfo.Validator.Address).String())
		fmt.Println("abci.BeginBlockerxx3", sdk.AccAddress(voteInfo.Validator.Address).String())
		k.HandleValidatorSignature(ctx, voteInfo.Validator.Address, voteInfo.Validator.Power, voteInfo.SignedLastBlock)
	}
}
