package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default CustomGo genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Aliases: []*TokenAlias{
			NewTokenAlias("KEX", "Kira", "", 6, []string{"ukex", "mkex"}, false),
		},
		Rates: []*TokenRate{
			NewTokenRate("ukex", sdk.NewDec(1), true, sdk.NewDecWithPrec(50, 2), sdk.OneInt(), true, false),             // 1
			NewTokenRate("ubtc", sdk.NewDec(10), true, sdk.NewDecWithPrec(25, 2), sdk.OneInt(), true, false),            // 10
			NewTokenRate("xeth", sdk.NewDecWithPrec(1, 1), true, sdk.NewDecWithPrec(10, 2), sdk.OneInt(), false, false), // 0.1
			NewTokenRate("frozen", sdk.NewDecWithPrec(1, 1), true, sdk.ZeroDec(), sdk.OneInt(), false, false),           // 0.1
		},
		TokenBlackWhites: &TokensWhiteBlack{
			Whitelisted: []string{"ukex"},
			Blacklisted: []string{"frozen"},
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
