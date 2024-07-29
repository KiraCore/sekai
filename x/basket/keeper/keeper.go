package keeper

import (
	"github.com/KiraCore/sekai/x/basket/types"
	govkeeper "github.com/KiraCore/sekai/x/gov/keeper"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
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
	gk       govkeeper.Keeper
	tk       types.TokensKeeper
	mk       types.MultiStakingKeeper
}

// NewKeeper returns instance of a keeper
func NewKeeper(storeKey storetypes.StoreKey, cdc codec.BinaryCodec, ak types.AccountKeeper, bk types.BankKeeper, gk govkeeper.Keeper, tk types.TokensKeeper, mk types.MultiStakingKeeper) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		ak:       ak,
		bk:       bk,
		gk:       gk,
		tk:       tk,
		mk:       mk,
	}
}

func (k Keeper) CheckIfAllowedPermission(ctx sdk.Context, addr sdk.AccAddress, permValue govtypes.PermValue) bool {
	return govkeeper.CheckIfAllowedPermission(ctx, k.gk, addr, govtypes.PermHandleBasketEmergency)
}
