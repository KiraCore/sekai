package multistaking

import (
	"github.com/KiraCore/sekai/x/multistaking/keeper"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker called every block, update validator set
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
