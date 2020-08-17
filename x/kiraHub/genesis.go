package kiraHub

import (
	sdkTypes "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/x/kiraHub/keeper"
)

type GenesisState struct{}

func DefaultGenesisState() GenesisState { return GenesisState{} }

func ValidateGenesis(genesisState GenesisState) error { return nil }

func InitializeGenesisState(context sdkTypes.Context, keeper keeper.Keeper, genesisState GenesisState) {
}
func ExportGenesis(context sdkTypes.Context, keeper keeper.Keeper) GenesisState {
	return GenesisState{}
}
