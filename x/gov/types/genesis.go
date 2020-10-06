package types

// DefaultGenesis returns the default CustomGo genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Permissions: map[uint64]*Permissions{
			uint64(RoleSudo):      NewPermissions([]PermValue{PermSetPermissions, PermClaimCouncilor, PermClaimValidator}, nil),
			uint64(RoleValidator): NewPermissions([]PermValue{PermClaimValidator}, nil),
		},
		StartingProposalId: 1,
		NetworkProperties: &NetworkProperties{
			MinTxFee: 1,
			MaxTxFee: 10000,
		},
		ExecutionFees: []*ExecutionFee{
			{
				Name:              "Claim Validator Seat",
				TransactionType:   "A",
				ExecutionFee:      10,
				FailureFee:        1,
				Timeout:           10,
				DefaultParameters: 0,
			},
			{
				Name:              "Claim Governance Seat",
				TransactionType:   "A",
				ExecutionFee:      10,
				FailureFee:        1,
				Timeout:           10,
				DefaultParameters: 0,
			},
			{
				Name:              "Claim Proposal Type X",
				TransactionType:   "B",
				ExecutionFee:      10,
				FailureFee:        1,
				Timeout:           10,
				DefaultParameters: 0,
			},
			{
				Name:              "Vote Proposal Type X",
				TransactionType:   "B",
				ExecutionFee:      10,
				FailureFee:        1,
				Timeout:           10,
				DefaultParameters: 0,
			},
			{
				Name:              "Submit Proposal Type X",
				TransactionType:   "B",
				ExecutionFee:      10,
				FailureFee:        1,
				Timeout:           10,
				DefaultParameters: 0,
			},
			{
				Name:              "Veto Proposal Type X",
				TransactionType:   "B",
				ExecutionFee:      10,
				FailureFee:        1,
				Timeout:           10,
				DefaultParameters: 0,
			},
		},
	}
}
