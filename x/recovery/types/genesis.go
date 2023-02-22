package types

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	recoveryRecords []RecoveryRecord,
	recoveryTokens []RecoveryToken,
	rewards []Rewards,
	rotations []Rotation,
) *GenesisState {
	return &GenesisState{
		RecoveryRecords: recoveryRecords,
		RecoveryTokens:  recoveryTokens,
		Rewards:         rewards,
		Rotations:       rotations,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		RecoveryRecords: []RecoveryRecord{},
		RecoveryTokens:  []RecoveryToken{},
		Rewards:         []Rewards{},
	}
}

// ValidateGenesis validates the recovery genesis parameters
func ValidateGenesis(data GenesisState) error {

	return nil
}
