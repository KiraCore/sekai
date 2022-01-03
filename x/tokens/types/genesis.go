package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default CustomGo genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Aliases: []*TokenAlias{
			NewTokenAlias("KEX", "Kira", "", 6, []string{"ukex", "mkex"}),
		},
		Rates: []*TokenRate{
			NewTokenRate("ukex", sdk.NewDec(1), true),              // 1
			NewTokenRate("ubtc", sdk.NewDec(10), true),             // 10
			NewTokenRate("xeth", sdk.NewDecWithPrec(1, 1), true),   // 0.1
			NewTokenRate("frozen", sdk.NewDecWithPrec(1, 1), true), // 0.1
		},
		TokenBlackWhites: &TokensWhiteBlack{
			Whitelisted: []string{"ukex"},
			Blacklisted: []string{"frozen"},
		},
	}
}
