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
	ak       types.AccountKeeper
	bk       types.BankKeeper
	sk       types.StakingKeeper
	gk       types.GovKeeper
	msk      types.MultiStakingKeeper
	ck       types.CollectivesKeeper
	spk      types.SpendingKeeper
	custodyk types.CustodyKeeper
}

// NewKeeper creates a recovery keeper
func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey,
	ak types.AccountKeeper,
	sk types.StakingKeeper,
	gk types.GovKeeper,
	msk types.MultiStakingKeeper,
	ck types.CollectivesKeeper,
	spk types.SpendingKeeper,
	custodyk types.CustodyKeeper,
) Keeper {

	return Keeper{
		storeKey: key,
		cdc:      cdc,
		ak:       ak,
		sk:       sk,
		gk:       gk,
		msk:      msk,
		ck:       ck,
		spk:      spk,
		custodyk: custodyk,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
