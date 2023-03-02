package keeper

import (
	"github.com/KiraCore/sekai/x/distributor/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper is for managing token module
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey sdk.StoreKey
	ak       types.AccountKeeper
	bk       types.BankKeeper
	sk       types.StakingKeeper
	gk       types.CustomGovKeeper
	mk       types.MultiStakingKeeper
	rk       types.RecoveryKeeper
}

// NewKeeper returns instance of a keeper
func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryCodec, ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper, gk types.CustomGovKeeper, mk types.MultiStakingKeeper, rk types.RecoveryKeeper) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		ak:       ak,
		bk:       bk,
		sk:       sk,
		gk:       gk,
		mk:       mk,
		rk:       rk,
	}
}

// BondDenom returns the denom that is basically used for fee payment
func (k Keeper) BondDenom(ctx sdk.Context) string {
	return "ukex"
}
