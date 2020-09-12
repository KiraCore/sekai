package types

// DefaultGenesis returns the default CustomGo genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Permissions: map[uint64]*Permissions{
			uint64(RoleSudo):      NewPermissions([]PermValue{PermSetPermissions, PermClaimGovernance, PermClaimValidator}, nil),
			uint64(RoleValidator): NewPermissions([]PermValue{PermClaimValidator}, nil),
		},
	}
}
