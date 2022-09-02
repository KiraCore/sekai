package basket

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/KiraCore/sekai/x/basket/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker sets the proposer for determining distributor during endblock
// and distribute rewards for the previous block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
}
