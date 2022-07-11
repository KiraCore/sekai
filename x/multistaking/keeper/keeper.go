package keeper

import (
	govkeeper "github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/multistaking/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper represents the keeper that maintains the Validator Registry.
type Keeper struct {
	storeKey    sdk.StoreKey
	cdc         codec.BinaryCodec
	bankKeeper  types.BankKeeper
	tokenKeeper types.TokensKeeper
	govKeeper   govkeeper.Keeper
	sk          types.StakingKeeper
}

// NewKeeper returns new keeper.
func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryCodec, bankKeeper types.BankKeeper, tokenKeeper types.TokensKeeper, govKeeper govkeeper.Keeper, sk types.StakingKeeper) Keeper {
	return Keeper{
		storeKey:    storeKey,
		cdc:         cdc,
		bankKeeper:  bankKeeper,
		tokenKeeper: tokenKeeper,
		govKeeper:   govKeeper,
		sk:          sk,
	}
}
