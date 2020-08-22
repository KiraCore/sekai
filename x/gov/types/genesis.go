package types

// DefaultGenesis returns the default CustomGo genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Permissions: map[string]Permissions{
			string(RoleValidator): NewPermissions([]PermValue{PermClaimGovernanceSeat}, nil),
			string(RoleCouncilor): {},
			string(RoleGovLeader): {},
		},
	}
}
