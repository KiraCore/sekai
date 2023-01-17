package types

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	recoveryRecords []RecoveryRecord,
	recoveryTokens []RecoveryToken,
) *GenesisState {

	return &GenesisState{
		RecoveryRecords: recoveryRecords,
		RecoveryTokens:  recoveryTokens,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		RecoveryRecords: []RecoveryRecord{},
		RecoveryTokens:  []RecoveryToken{},
	}
}

// ValidateGenesis validates the recovery genesis parameters
func ValidateGenesis(data GenesisState) error {

	return nil
}
