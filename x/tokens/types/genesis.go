package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default CustomGo genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		TokenInfos: []*TokenInfo{
			NewTokenInfo("ukex", sdk.NewDec(1), true, sdk.NewDecWithPrec(50, 2), sdk.OneInt(), true, false, "KEX", "KEX", "", 6),                   // 1
			NewTokenInfo("ubtc", sdk.NewDec(10), true, sdk.NewDecWithPrec(25, 2), sdk.OneInt(), true, false, "BTC", "Bitcoin", "", 9),              // 10
			NewTokenInfo("xeth", sdk.NewDecWithPrec(1, 1), true, sdk.NewDecWithPrec(10, 2), sdk.OneInt(), false, false, "ETH", "Ethereum", "", 18), // 0.1
			NewTokenInfo("frozen", sdk.NewDecWithPrec(1, 1), true, sdk.ZeroDec(), sdk.OneInt(), false, false, "FROZEN", "FROZEN", "", 6),           // 0.1
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
