package recovery

import (
	"github.com/KiraCore/sekai/x/recovery/keeper"
	"github.com/KiraCore/sekai/x/recovery/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, stakingKeeper types.StakingKeeper, data *types.GenesisState) {
	for _, record := range data.RecoveryRecords {
		keeper.SetRecoveryRecord(ctx, record)
	}

	for _, token := range data.RecoveryTokens {
		keeper.SetRecoveryToken(ctx, token)
	}

	for _, rewards := range data.Rewards {
		holder := sdk.MustAccAddressFromBech32(rewards.Holder)
		keeper.SetRRTokenHolderRewards(ctx, holder, rewards.Rewards)
	}

	for _, rotation := range data.Rotations {
		keeper.SetRotationHistory(ctx, rotation)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) (data *types.GenesisState) {
	records := keeper.GetAllRecoveryRecords(ctx)
	tokens := keeper.GetAllRecoveryTokens(ctx)
	rewards := keeper.GetAllRRHolderRewards(ctx)
	rotations := keeper.GetAllRotationHistory(ctx)
	return types.NewGenesisState(records, tokens, rewards, rotations)
}
