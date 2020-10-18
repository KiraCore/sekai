package types

// DefaultGenesis returns the default CustomGo genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Aliases: map[string]*TokenAlias{
			"KEX": NewTokenAlias(0, 0, []VoteType{0, 1}, "KEX", "Kira", "", 6, []string{"ukex", "mkex"}, ProposalStatus_active),
		},
		Rates: map[string]*TokenRate{
			"ukex": NewTokenRate("ukex", 1, true),
			"ubtc": NewTokenRate("ubtc", 0.00001, true),
			"xeth": NewTokenRate("xeth", 0.0001, true),
		},
	}
}
