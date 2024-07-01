package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/KiraCore/sekai/x/recovery/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the recovery store
type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec
	ak       types.AccountKeeper
	bk       types.BankKeeper
	sk       types.StakingKeeper
	gk       types.GovKeeper
	msk      types.MultiStakingKeeper
	ck       types.CollectivesKeeper
	spk      types.SpendingKeeper
	custodyk types.CustodyKeeper
	tk       types.TokensKeeper
}

// NewKeeper creates a recovery keeper
func NewKeeper(cdc codec.BinaryCodec, key storetypes.StoreKey,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	sk types.StakingKeeper,
	gk types.GovKeeper,
	msk types.MultiStakingKeeper,
	ck types.CollectivesKeeper,
	spk types.SpendingKeeper,
	custodyk types.CustodyKeeper,
	tk types.TokensKeeper,
) Keeper {

	return Keeper{
		storeKey: key,
		cdc:      cdc,
		ak:       ak,
		bk:       bk,
		sk:       sk,
		gk:       gk,
		msk:      msk,
		ck:       ck,
		spk:      spk,
		custodyk: custodyk,
		tk:       tk,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
