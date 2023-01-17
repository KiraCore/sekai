package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/KiraCore/sekai/x/recovery/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the recovery store
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      codec.BinaryCodec
	sk       types.StakingKeeper
	gk       types.GovKeeper
	msk      types.MultiStakingKeeper
}

// NewKeeper creates a recovery keeper
func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, sk types.StakingKeeper) Keeper {

	return Keeper{
		storeKey: key,
		cdc:      cdc,
		sk:       sk,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
