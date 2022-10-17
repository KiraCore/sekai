package keeper

import (
	"github.com/KiraCore/sekai/x/basket/types"
	govkeeper "github.com/KiraCore/sekai/x/gov/keeper"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
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
	gk       govkeeper.Keeper
	mk       types.MultiStakingKeeper
}

// NewKeeper returns instance of a keeper
func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryCodec, ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper, gk govkeeper.Keeper, mk types.MultiStakingKeeper) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		ak:       ak,
		bk:       bk,
		sk:       sk,
		gk:       gk,
		mk:       mk,
	}
}

// BondDenom returns the denom that is basically used for fee payment
func (k Keeper) BondDenom(ctx sdk.Context) string {
	return "ukex"
}

func (k Keeper) CheckIfAllowedPermission(ctx sdk.Context, addr sdk.AccAddress, permValue govtypes.PermValue) bool {
	return govkeeper.CheckIfAllowedPermission(ctx, k.gk, addr, govtypes.PermHandleBasketEmergency)
}
