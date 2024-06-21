package types

import (
	"encoding/json"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default CustomGo genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		TokenInfos: []TokenInfo{
			NewTokenInfo("ukex", "adr20", sdk.NewDec(1), true, math.ZeroInt(), math.ZeroInt(), sdk.NewDecWithPrec(50, 2), sdk.OneInt(), true, false, "KEX", "KEX", "", 6, "", "", "", 0, math.ZeroInt(), "", false, "", ""),                   // 1
			NewTokenInfo("ubtc", "adr20", sdk.NewDec(10), true, math.ZeroInt(), math.ZeroInt(), sdk.NewDecWithPrec(25, 2), sdk.OneInt(), true, false, "BTC", "Bitcoin", "", 9, "", "", "", 0, math.ZeroInt(), "", false, "", ""),              // 10
			NewTokenInfo("xeth", "adr20", sdk.NewDecWithPrec(1, 1), true, math.ZeroInt(), math.ZeroInt(), sdk.NewDecWithPrec(10, 2), sdk.OneInt(), false, false, "ETH", "Ethereum", "", 18, "", "", "", 0, math.ZeroInt(), "", false, "", ""), // 0.1
			NewTokenInfo("frozen", "adr20", sdk.NewDecWithPrec(1, 1), true, math.ZeroInt(), math.ZeroInt(), sdk.ZeroDec(), sdk.OneInt(), false, false, "FROZEN", "FROZEN", "", 6, "", "", "", 0, math.ZeroInt(), "", false, "", ""),           // 0.1
		},
		TokenBlackWhites: TokensWhiteBlack{
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
