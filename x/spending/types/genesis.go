package types

import (
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default CustomGo genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Pools: []SpendingPool{
			{
				Name:          "ValidatorRewardPool",
				ClaimStart:    0,
				ClaimEnd:      0,
				Token:         "ukex",
				Rate:          sdk.NewDec(100),
				VoteQuorum:    33,
				VotePeriod:    300, // 300s
				VoteEnactment: 300, // 300s
				Owners:        &PermInfo{OwnerRoles: []uint64{govtypes.RoleValidator}},
				Beneficiaries: &PermInfo{OwnerRoles: []uint64{govtypes.RoleValidator}},
				Balance:       sdk.ZeroInt(),
			},
		},
	}
}
