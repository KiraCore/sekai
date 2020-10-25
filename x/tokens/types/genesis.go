package types

// DefaultGenesis returns the default CustomGo genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Aliases: map[string]*TokenAlias{
			"KEX": NewTokenAlias(0, 0, []VoteType{0, 1}, "KEX", "Kira", "", 6, []string{"ukex", "mkex"}, ProposalStatus_active),
		},
		Rates: map[string]*TokenRate{
			"ukex": NewTokenRate("ukex", 1000000000, true),
			"ubtc": NewTokenRate("ubtc", 10000000000, true),
			"xeth": NewTokenRate("xeth", 100000, true),
		},
	}
}
