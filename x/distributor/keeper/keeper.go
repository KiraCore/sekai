package keeper

import (
	appparams "github.com/KiraCore/sekai/app/params"
	"github.com/KiraCore/sekai/x/distributor/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper is for managing token module
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
	ak       types.AccountKeeper
	bk       types.BankKeeper
	tk       types.TokensKeeper
	sk       types.StakingKeeper
	gk       types.CustomGovKeeper
	mk       types.MultiStakingKeeper
	rk       types.RecoveryKeeper
}

// NewKeeper returns instance of a keeper
func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	sk types.StakingKeeper,
	gk types.CustomGovKeeper,
	mk types.MultiStakingKeeper,
	rk types.RecoveryKeeper,
	tk types.TokensKeeper,
) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		ak:       ak,
		bk:       bk,
		sk:       sk,
		gk:       gk,
		mk:       mk,
		rk:       rk,
		tk:       tk,
	}
}

// DefaultDenom returns the denom that is basically used for fee payment
func (k Keeper) DefaultDenom(ctx sdk.Context) string {
	return appparams.DefaultDenom
}
