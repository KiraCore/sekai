package types

import (
	"fmt"
	"time"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params, signingInfos []SigningInfo,
) *GenesisState {

	return &GenesisState{
		Params:       params,
		SigningInfos: signingInfos,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:       DefaultParams(),
		SigningInfos: []SigningInfo{},
	}
}

// ValidateGenesis validates the slashing genesis parameters
func ValidateGenesis(data GenesisState) error {

	downtimeInactive := data.Params.DowntimeInactiveDuration
	if downtimeInactive < 1*time.Minute {
		return fmt.Errorf("downtime unblond duration must be at least 1 minute, is %s", downtimeInactive.String())
	}

	return nil
}
