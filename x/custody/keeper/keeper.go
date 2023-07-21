package keeper

import (
	appparams "github.com/KiraCore/sekai/app/params"
	"github.com/KiraCore/sekai/x/custody/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper is for managing token module
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey sdk.StoreKey
	gk       types.CustomGovKeeper
	bk       types.BankKeeper
}

// NewKeeper returns instance of a keeper
func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryCodec, gk types.CustomGovKeeper, bk types.BankKeeper) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		gk:       gk,
		bk:       bk,
	}
}

// DefaultDenom returns the denom that is basically used for fee payment
func (k Keeper) DefaultDenom(ctx sdk.Context) string {
	return appparams.DefaultDenom
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}
