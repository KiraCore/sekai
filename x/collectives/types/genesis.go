package types

// DefaultGenesis returns the default CustomGo genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Collectives:  []Collective{},
		Contributers: []CollectiveContributor{},
	}
}
