package types

import (
	sekaitypes "github.com/KiraCore/sekai/types"
)

// DefaultGenesis returns the default CustomGo genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Version:     sekaitypes.SekaiVersion,
		CurrentPlan: nil,
		NextPlan:    nil,
	}
}
