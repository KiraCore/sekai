package types

// DefaultGenesis returns the default CustomGo genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		UbiRecords: []UBIRecord{
			{
				Name:              "ValidatorBasicRewardsPoolUBI",
				DistributionStart: 0,
				DistributionEnd:   0,
				DistributionLast:  0,
				Amount:            500_000,    // 500k KEX
				Period:            86400 * 30, // 30 days
				Pool:              "ValidatorBasicRewardsPool",
				Dynamic:           true,
			},
		},
	}
}
