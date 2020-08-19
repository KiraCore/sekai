package kiraHub

import (
	"github.com/KiraCore/sekai/x/kiraHub/keeper"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct{}

func DefaultGenesisState() GenesisState { return GenesisState{} }

func ValidateGenesis(genesisState GenesisState) error { return nil }

func InitializeGenesisState(context sdkTypes.Context, keeper keeper.Keeper, genesisState GenesisState) {
}
func ExportGenesis(context sdkTypes.Context, keeper keeper.Keeper) GenesisState {
	return GenesisState{}
}
