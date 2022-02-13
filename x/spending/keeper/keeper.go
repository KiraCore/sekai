package keeper

import (
	"github.com/KiraCore/sekai/x/spending/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper is for managing token module
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey sdk.StoreKey

	bk types.BankKeeper
	gk types.CustomGovKeeper
}

// NewKeeper returns instance of a keeper
func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryCodec, bk types.BankKeeper, gk types.CustomGovKeeper) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		bk:       bk,
		gk:       gk,
	}
}

func (k Keeper) IsAllowedAddress(ctx sdk.Context, address sdk.AccAddress, permInfo types.PermInfo) bool {
	for _, owner := range permInfo.OwnerAccounts {
		if owner == address.String() {
			return true
		}
	}

	actor, found := k.gk.GetNetworkActorByAddress(ctx, address)
	if !found {
		return false
	}

	flags := make(map[uint64]bool)
	for _, role := range permInfo.OwnerRoles {
		flags[role] = true
	}

	for _, role := range actor.Roles {
		if flags[role] {
			return true
		}
	}
	return false
}
