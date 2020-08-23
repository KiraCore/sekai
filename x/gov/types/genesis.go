package types

// DefaultGenesis returns the default CustomGo genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Permissions: map[uint64]Permissions{
			uint64(RoleValidator): NewPermissions([]PermValue{PermClaimGovernanceSeat}, nil),
			uint64(RoleCouncilor): {},
			uint64(RoleGovLeader): {},
		},
	}
}
