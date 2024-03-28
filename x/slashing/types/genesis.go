package types

// NewGenesisState creates a new GenesisState object
func NewGenesisState(signingInfos []SigningInfo) *GenesisState {
	return &GenesisState{
		SigningInfos: signingInfos,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		SigningInfos: []SigningInfo{},
	}
}

// ValidateGenesis validates the slashing genesis parameters
func ValidateGenesis(data GenesisState) error {

	return nil
}
