package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default CustomGo genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Aliases: map[string]*TokenAlias{
			"KEX": NewTokenAlias(0, 0, []VoteType{0, 1}, "KEX", "Kira", "", 6, []string{"ukex", "mkex"}, ProposalStatus_active),
		},
		Rates: map[string]*TokenRate{
			"ukex": NewTokenRate("ukex", sdk.NewDec(1), true),            // 1
			"ubtc": NewTokenRate("ubtc", sdk.NewDec(10), true),           // 10
			"xeth": NewTokenRate("xeth", sdk.NewDecWithPrec(1, 1), true), // 0.1
		},
	}
}
