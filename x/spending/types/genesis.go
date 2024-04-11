package types

import (
	"encoding/json"

	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/codec"
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
				Rates:         sdk.DecCoins{sdk.NewDecCoin("ukex", sdk.NewInt(385))}, // 1k KEX per month per validator
				VoteQuorum:    sdk.NewDecWithPrec(33, 2),
				VotePeriod:    300, // 300s
				VoteEnactment: 300, // 300s
				Owners:        &PermInfo{OwnerRoles: []uint64{govtypes.RoleValidator}},
				Beneficiaries: &WeightedPermInfo{
					Roles: []WeightedRole{{govtypes.RoleValidator, sdk.OneDec()}},
				},
				Balances: sdk.Coins{},
			},
		},
	}
}

// GetGenesisStateFromAppState returns x/auth GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return genesisState
}
