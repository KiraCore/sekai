package types

// DefaultGenesis returns the default CustomGo genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		CurrentPlan: nil,
		NextPlan:    nil,
	}
}
