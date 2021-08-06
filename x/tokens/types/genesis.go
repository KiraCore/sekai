package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default CustomGo genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Aliases: map[string]*TokenAlias{
			"KEX": NewTokenAlias("KEX", "Kira", "", 6, []string{"ukex", "mkex"}),
		},
		Rates: map[string]*TokenRate{
			"ukex":   NewTokenRate("ukex", sdk.NewDec(1), true),              // 1
			"ubtc":   NewTokenRate("ubtc", sdk.NewDec(10), true),             // 10
			"xeth":   NewTokenRate("xeth", sdk.NewDecWithPrec(1, 1), true),   // 0.1
			"frozen": NewTokenRate("frozen", sdk.NewDecWithPrec(1, 1), true), // 0.1
		},
		TokenBlackWhites: &TokensWhiteBlack{
			Whitelisted: []string{"ukex"},
			Blacklisted: []string{"frozen"},
		},
	}
}
