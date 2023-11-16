package basket

import (
	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/KiraCore/sekai/x/basket/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker sets the proposer for determining distributor during endblock
// and distribute rewards for the previous block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	baskets := k.GetAllBaskets(ctx)
	for _, basket := range baskets {
		k.ClearOldMintAmounts(ctx, basket.Id, basket.LimitsPeriod)
		k.ClearOldBurnAmounts(ctx, basket.Id, basket.LimitsPeriod)
		k.ClearOldSwapAmounts(ctx, basket.Id, basket.LimitsPeriod)
	}
}
