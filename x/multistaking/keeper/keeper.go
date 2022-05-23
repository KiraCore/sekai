package keeper

import (
	"github.com/KiraCore/sekai/x/multistaking/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper represents the keeper that maintains the Validator Registry.
type Keeper struct {
	storeKey    sdk.StoreKey
	cdc         *codec.LegacyAmino
	bankKeeper  types.BankKeeper
	tokenKeeper types.TokensKeeper
}

// NewKeeper returns new keeper.
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.LegacyAmino) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}
