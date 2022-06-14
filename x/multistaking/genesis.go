package multistaking

import (
	"github.com/KiraCore/sekai/x/multistaking/keeper"
	"github.com/KiraCore/sekai/x/multistaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the evidence module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, gs *types.GenesisState) {
	for _, pool := range gs.Pools {
		k.SetStakingPool(ctx, pool)
	}
	for _, undelegation := range gs.Undelegations {
		k.SetUndelegation(ctx, undelegation)
	}
	for _, reward := range gs.Rewards {
		delegator, err := sdk.AccAddressFromBech32(reward.Delegator)
		if err != nil {
			panic(err)
		}
		k.SetDelegatorRewards(ctx, delegator, reward.Rewards)
	}
}

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Pools:         keeper.GetAllStakingPools(ctx),
		Undelegations: keeper.GetAllUndelegations(ctx),
		Rewards:       keeper.GetAllDelegatorRewards(ctx),
	}
}
