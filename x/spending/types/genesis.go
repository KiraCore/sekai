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
				Name:          "ValidatorBasicRewardsPool",
				ClaimStart:    0,
				ClaimEnd:      0,
				Token:         "ukex",
				Rates:         sdk.DecCoins{sdk.NewDecCoin("ukex", sdk.NewInt(385))}, // 1k KEX per month per validator
				VoteQuorum:    33,
				VotePeriod:    300, // 300s
				VoteEnactment: 300, // 300s
				Owners:        &PermInfo{OwnerRoles: []uint64{govtypes.RoleValidator}},
				Beneficiaries: &WeightedPermInfo{
					Roles: []WeightedRole{{govtypes.RoleValidator, 1}},
				},
				Balances: sdk.Coins{},
			},
		},
	}
}
