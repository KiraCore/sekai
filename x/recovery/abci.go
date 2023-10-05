package recovery

import (
	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/KiraCore/sekai/x/recovery/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {

}
